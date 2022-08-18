// 接口类型发送消息，支持群聊和单聊，但是不支持at某人 ，需要有企业内部应用或企业内部机器人
// 单聊 https://open.dingtalk.com/document/group/chatbots-send-one-on-one-chat-messages-in-batches
// 群聊 https://open.dingtalk.com/document/group/the-robot-sends-a-group-message

package ding

import (
	"bytes"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	oToMessageBatchSendUrl = "https://api.dingtalk.com/v1.0/robot/oToMessages/batchSend"
	groupMessageSendUrl    = "https://api.dingtalk.com/v1.0/robot/groupMessages/send"
	cache                  *bigcache.BigCache
)

// OtOClient 单聊客户端
type OtOClient struct {
	*IClient
}

// GroupClient 群聊客户端，好像没有办法 @某人
type GroupClient struct {
	*IClient
}

// IClient 接口方式发送消息，支持群聊和单聊
type IClient struct {
	url       string
	RobotCode string `json:"robotCode"`
	AppKeySecret
}

// OtOMessageBody 发送单聊post body
type OtOMessageBody struct {
	*RobotCodeMsgKeyParam
	UserIds []string `json:"userIds"`
}

func (o *OtOMessageBody) String() string {
	marshal, _ := json.Marshal(o)
	return string(marshal)
}

// GroupMessageBody 发送群聊post body
type GroupMessageBody struct {
	*RobotCodeMsgKeyParam
	// 开放的群id。需要是加密的。可以从钉钉post来的请求的：conversationId 字段获取
	OpenConversationId string `json:"openConversationId"`
}

func (g *GroupMessageBody) String() string {
	marshal, _ := json.Marshal(g)
	return string(marshal)
}

type RobotCodeMsgKeyParam struct {
	// 参考： https://open.dingtalk.com/document/group/message-types-and-data-format?spm=ding_open_doc.document.0.0.68116771ZgXHJA
	//消息模板参数。
	MsgParam string `json:"msgParam"`
	//消息模板key。
	MsgKey string `json:"msgKey"`
	//企业内部开发-机器人：填写企业自建应用的appKey。
	//第三方企业机器人：填写第三方企业应用绑定机器人的robotCode。
	RobotCode string `json:"robotCode"`
}

func initCache() {
	// 钉钉默认的accessToken 有效期为7200秒（2小时）
	var err error
	cache, err = bigcache.NewBigCache(bigcache.DefaultConfig(7000 * time.Second))
	if err != nil {
		log.Fatalf("init BigCache failed: %s\n", err)
	}
}

// NewOtOClient 创建单聊客户端
func NewOtOClient(robotCode string, appKey, appSecret string) *OtOClient {
	initCache()
	return &OtOClient{
		IClient: &IClient{
			url:       oToMessageBatchSendUrl,
			RobotCode: robotCode,
			AppKeySecret: AppKeySecret{
				AppKey:    appKey,
				AppSecret: appSecret,
			},
		},
	}
}

// NewGroupClient 创建群聊客户端
func NewGroupClient(robotCode string, appKey, appSecret string) *GroupClient {
	initCache()
	return &GroupClient{
		IClient: &IClient{
			url:       groupMessageSendUrl,
			RobotCode: robotCode,
			AppKeySecret: AppKeySecret{
				AppKey:    appKey,
				AppSecret: appSecret,
			},
		},
	}
}

// 通过接口的方式发送钉钉消息
func (c *IClient) sendDingInterfaceMsg(url string, msg any) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msgByte))
	if err != nil {
		return err
	}
	accessToken, err := GetAccessToken(AppKeySecret{AppKey: c.AppKey, AppSecret: c.AppSecret})
	if err != nil {
		return err
	}
	client := http.Client{Timeout: 2 * time.Second}
	request.Header.Set("x-acs-dingtalk-access-token", accessToken)
	request.Header.Set("Content-Type", ContentTypeJson)

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("发送钉钉接口消息后，收到钉钉的回复: %v\n", string(respByte))
	return nil
}

func (c *IClient) createRobotCodeMessageKeyParam(msgKey, msgParam string) *RobotCodeMsgKeyParam {
	return &RobotCodeMsgKeyParam{
		RobotCode: c.RobotCode,
		MsgKey:    msgKey,
		MsgParam:  msgParam,
	}
}

func (c *IClient) createOtOMessageBody(msgKey, msgParam string, userIds []string) *OtOMessageBody {
	msg := &OtOMessageBody{
		UserIds:              userIds,
		RobotCodeMsgKeyParam: c.createRobotCodeMessageKeyParam(msgKey, msgParam),
	}
	return msg
}

func (c *IClient) createGroupMessageBody(openConversationId, msgKey, msgParam string) *GroupMessageBody {
	msg := &GroupMessageBody{
		OpenConversationId:   openConversationId,
		RobotCodeMsgKeyParam: c.createRobotCodeMessageKeyParam(msgKey, msgParam),
	}
	return msg
}

// SendTextMsgWithUserIds 发送单聊文本消息给userIds这些用户，可以从postReq.senderStaffId 获取
func (o *OtOClient) SendTextMsgWithUserIds(content string, userIds []string) error {
	msg := o.createOtOMessageBody(IMsgKeyText, Text{Content: content}.String(), userIds)
	return o.sendDingInterfaceMsg(o.url, msg)
}

// SendMarkdownMsgWithUserIds 发送单聊markdown消息给userIds这些用户，可以从postReq.senderStaffId 获取
func (o *OtOClient) SendMarkdownMsgWithUserIds(title, text string, userIds []string) error {
	msg := o.createOtOMessageBody(IMsgKeyMarkdown, (&Markdown{Title: title, Text: text}).String(), userIds)
	return o.sendDingInterfaceMsg(o.url, msg)
}

// SendImageMsg 发送单聊图片消息给userIds这些用户，可以从postReq.senderStaffId 获取
func (o *OtOClient) SendImageMsg(photoURL string, userIds []string) error {
	msg := o.createOtOMessageBody(IMsgKeyImage, (&Image{PhotoURL: photoURL}).String(), userIds)
	return o.sendDingInterfaceMsg(o.url, msg)
}

// SendLinkMsg 发送单聊Link链接消息给userIds这些用户，可以从postReq.senderStaffId 获取
func (o *OtOClient) SendLinkMsg(title, text, picUrl, messageUrl string, userIds []string) error {
	msg := o.createOtOMessageBody(IMsgKeyLink, NewLink(title, text, picUrl, messageUrl).String(), userIds)
	return o.sendDingInterfaceMsg(o.url, msg)
}

func (o *OtOClient) SendActionCardMsg(title, text, singleTitle, singleURL string, userIds []string) error {
	msg := o.createOtOMessageBody(IMsgKeyActionCard, NewEntiretyActionCard(title, text, singleTitle, singleURL).String(), userIds)
	return o.sendDingInterfaceMsg(o.url, msg)
}

// SendTextMsg 发送群聊文本消息给conversationId这个群，可以从postReq.conversationId 获取
func (g *GroupClient) SendTextMsg(content, conversationId string) error {
	msg := g.createGroupMessageBody(conversationId, IMsgKeyText, (&Text{Content: content}).String())
	return g.sendDingInterfaceMsg(g.url, msg)
}

// SendMarkdownMsg 发送群聊markdown消息给conversationId这个群，可以从postReq.conversationId 获取
func (g *GroupClient) SendMarkdownMsg(title, text, conversationId string) error {
	return g.sendDingInterfaceMsg(g.url, g.createGroupMessageBody(conversationId, IMsgKeyMarkdown, (&Markdown{Title: title, Text: text}).String()))
}

// SendImageMsg 发送群聊图片消息给conversationId群，可以从postReq.conversationId 获取
func (g *GroupClient) SendImageMsg(photoURL, conversationId string) error {
	msg := g.createGroupMessageBody(conversationId, IMsgKeyImage, (&Image{PhotoURL: photoURL}).String())
	return g.sendDingInterfaceMsg(g.url, msg)
}

// SendLinkMsg 发送群聊Link链接消息给conversationId群，可以从postReq.conversationId 获取
func (g *GroupClient) SendLinkMsg(title, text, picUrl, messageUrl, conversationId string) error {
	msg := g.createGroupMessageBody(conversationId, IMsgKeyLink, NewLink(title, text, picUrl, messageUrl).String())
	return g.sendDingInterfaceMsg(g.url, msg)
}
func (g *GroupClient) SendActionCardMsg(title, text, singleTitle, singleURL, conversationId string) error {
	msg := g.createGroupMessageBody(conversationId, IMsgKeyActionCard, NewEntiretyActionCard(title, text, singleTitle, singleURL).String())
	return g.sendDingInterfaceMsg(g.url, msg)
}
