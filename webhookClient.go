// webhook 方式发送消息，仅支持群聊，支持at某人，不支持单聊，无需企业内部应用，直接创建群自定义机器人即可。
// 也支持通过企业内部机器人postReq 发来的SessionWebhook地址来发送消息

package ding

// 参考： https://github.com/wanghuiyt/ding/blob/main/ding.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ContentTypeJson = "application/json; charset=utf-8"
)

// WhClient 钉钉webhook客户端, webhook方式只支持在群聊会话。
type WhClient struct {
	// 通过企业内部机器人postReq 发来的SessionWebhook地址(有过期时间)来发送/回复消息。和下面的AccessToken二选一
	SessionWebhookUrl string
	// webhook url 中access_token=XXX, 填写这里的XXX
	AccessToken string
	// 钉钉安全设置，加签。
	// 参考： https://developers.dingtalk.com/document/robots/customize-robot-security-settings
	Secret string
	// 钉钉安全设置，关键字
	// todo 还未实现
	KeyWorld string
}

// NewWhClientWithoutSecret 创建钉钉客户端，不用密钥
func NewWhClientWithoutSecret(accessToken string) *WhClient {
	return &WhClient{
		AccessToken: accessToken,
	}
}

// NewWhClientWithSecret 创建钉钉客户端，使用密钥
func NewWhClientWithSecret(aToken, secret string) *WhClient {
	return &WhClient{
		AccessToken: aToken,
		Secret:      secret,
	}
}

// NewWhClientUseSessionWebhook 通过企业内部机器人postReq 发来的SessionWebhook地址来发送消息
func NewWhClientUseSessionWebhook(sessionWebhookUrl string) *WhClient {
	return &WhClient{SessionWebhookUrl: sessionWebhookUrl}
}

// GetUrl 获取 给钉钉发送post的url， 根据是否有安全设置会有不同的url
func (c *WhClient) GetUrl() string {
	if c.SessionWebhookUrl != "" {
		return c.SessionWebhookUrl
	} else {
		wh := "https://oapi.dingtalk.com/robot/send?access_token=" + c.AccessToken
		if len(c.Secret) > 0 {
			timestame := time.Now().Unix()
			data := fmt.Sprintf("%d\n%s", timestame, c.Secret)
			sign := GetSign(data, c.Secret)
			return fmt.Sprintf("%s&timestamp=%d&sign=%s", wh, timestame, sign)
		}
		return wh
	}
}

// sendDingWebhookMsg 发送钉钉webhook post 请求，即发送消息。msg为message.go里定义的
func sendDingWebhookMsg(url string, msg any) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, ContentTypeJson, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("发送webhook消息，收到钉钉的回复: %v\n", respByte)
	return nil
}

// SendTextMsgWithUserIds 发送文本消息，群聊, @userIds
func (c *WhClient) SendTextMsgWithUserIds(content string, userIds []string) error {
	msg := NewWhTextMsgWithAtUserIds(content, userIds...)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendTextMsgWithUserMobile 发送文本消息，群聊, @mobile
func (c *WhClient) SendTextMsgWithUserMobile(content string, mobiles []string) error {
	msg := NewWhTextMsgWithAtMobiles(content, mobiles...)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendTextMsgWithAtAll 发送文本消息，群聊, @all
func (c *WhClient) SendTextMsgWithAtAll(content string) error {
	msg := NewWhTextMsgWithAtAll(content)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendTextMsg 发送文本消息，群聊
func (c *WhClient) SendTextMsg(content string) error {
	msg := NewWhTextMsg(content)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendMarkdownMsgWithUserIds 发送markdown消息，群聊，@userIds
func (c *WhClient) SendMarkdownMsgWithUserIds(title, text string, userIds []string) error {
	msg := NewWhMarkdownMsgWithAtUserIds(title, text, userIds...)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendMarkdownMsgWithUserMobile 发送markdown消息，群聊，@mobile
func (c *WhClient) SendMarkdownMsgWithUserMobile(title, text string, mobiles []string) error {
	msg := NewWhMarkdownMsgWithAtMobiles(title, text, mobiles...)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendMarkdownMsgWithAtAll 发送markdown消息，群聊，@all
func (c *WhClient) SendMarkdownMsgWithAtAll(title, text string) error {
	msg := NewWhMarkdownMsgWithAtAll(title, text)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendMarkdownMsg 发送markdown消息，群聊
func (c *WhClient) SendMarkdownMsg(title, text string) error {
	msg := NewWhMarkdownMsg(title, text)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendLinkMsg 发送link链接消息，群聊，这个不能@某人
func (c *WhClient) SendLinkMsg(title, text, messageUrl, picUrl string) error {
	msg := NewWhLinkMsg(text, title, picUrl, messageUrl)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

// SendEntiretyActionCardMsg 发送整体跳转actionCard 消息
func (c *WhClient) SendEntiretyActionCardMsg(title, text, singleTitle, singleURL string) error {
	msg := NewWhEntiretyActionCardMsg(title, text, singleTitle, singleURL)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

func (c *WhClient) SendIndependentActionCardMsg(title, text string, btns []*Btn) error {
	msg := NewWhIndependentActionCardMsg(title, text, btns)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

func (c *WhClient) SendIndependentActionCardMsgWithBtnOrientation(title, text, btnOrientation string, btns []*Btn) error {
	msg := NewWhIndependentActionCardMsgWithBtnOrientation(title, text, btnOrientation, btns)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}

func (c *WhClient) SendWhFeedCardMsg(links []*Link) error {
	msg := NewWhFeedCardMsg(links)
	return sendDingWebhookMsg(c.GetUrl(), msg)
}
