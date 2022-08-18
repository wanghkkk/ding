package ding

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 通过接口发送机器人消息类型 支持批量发送单聊消息和向群内发消息
var (
	IMsgKeyText        = "sampleText"
	IMsgKeyMarkdown    = "sampleMarkdown"
	IMsgKeyImage       = "sampleImageMsg"
	IMsgKeyLink        = "sampleLink"
	IMsgKeyActionCard  = "sampleActionCard"
	IMsgKeyActionCard2 = "sampleActionCard2"
	IMsgKeyActionCard3 = "sampleActionCard3"
	IMsgKeyActionCard4 = "sampleActionCard4"
	IMsgKeyActionCard5 = "sampleActionCard5"
	IMsgKeyActionCard6 = "sampleActionCard6"

	// Debug 开启debug，默认不开启，开启后会输出钉钉返回的消息
	Debug = false
)

// OpenDebug 开启debug模式，开启后会输出钉钉返回的消息
func OpenDebug() {
	Debug = true
}

// Image 图片类型，接口方式发送机器人消息专有
type Image struct {
	PhotoURL string `json:"photoURL"`
}

func (i *Image) String() string {
	return fmt.Sprintf(`{"photoURL":%q}`, i.PhotoURL)
}

// 消息类型参考： https://open.dingtalk.com/document/group/message-types-and-data-format
// 通过webhook发送机器人消息的类型
var (
	WhMsgTypeFeedCard   = "feedCard"
	WhMsgTypeActionCard = "actionCard"
	WhMsgTypeMarkdown   = "markdown"
	WhMsgTypeText       = "text"
	WhMsgTypeLink       = "link"
)

// WhTextMsg Text 文本消息
type WhTextMsg struct {
	At   At   `json:"at"`
	Text Text `json:"text"`
	// 消息类型，必须提供
	MsgType string `json:"msgtype"`
}

// At @ 谁
type At struct {
	// 被@人的手机号。可选
	//说明 消息内容content中要带上"@手机号"，跟atMobiles参数结合使用，才有@效果
	AtMobiles []string `json:"atMobiles,omitempty"`
	// 被@人的用户userid。可选
	//说明 消息内容content中要带上"@userId"，跟atUserIds参数结合使用，才有@效果。
	AtUserIds []string `json:"atUserIds,omitempty"`
	// @所有人是true，否则为false。可选
	IsAtAll bool `json:"isAtAll,omitempty"`
}

// Text 消息
type Text struct {
	// 消息文本。
	Content string `json:"content"`
}

func (t Text) String() string {
	return fmt.Sprintf(`{"content": %q}`, t.Content)
}

func NewWhTextMsg(content string) *WhTextMsg {
	return &WhTextMsg{
		MsgType: WhMsgTypeText,
		Text: Text{
			Content: content,
		},
	}
}

func NewWhTextMsgWithAtMobiles(content string, mobiles ...string) *WhTextMsg {
	textMsg := NewWhTextMsg(content)
	textMsg.At.AtMobiles = mobiles
	return textMsg
}

func NewWhTextMsgWithAtUserIds(content string, userIds ...string) *WhTextMsg {
	textMsg := NewWhTextMsg(content)
	textMsg.At.AtUserIds = userIds
	return textMsg
}

func NewWhTextMsgWithAtAll(content string) *WhTextMsg {
	textMsg := NewWhTextMsg(content)
	textMsg.At.IsAtAll = true
	return textMsg
}

// WhLinkMsg 链接消息
type WhLinkMsg struct {
	MsgType string `json:"msgtype"`
	Link    Link   `json:"link"`
}

// Link Link链接消息体
type Link struct {
	//必填项
	//消息内容。如果太长只会部分展示。
	//当做feedCard时，无需这个字段
	Text string `json:"text"`
	// 消息标题。必填项
	Title string `json:"title"`
	// 图片URL。可选
	PicUrl string `json:"picUrl,omitempty"`
	//必填项
	//点击消息跳转的URL.
	//移动端 在钉钉客户端内打开。
	//PC端 默认外部浏览器打开。
	//希望在侧边栏打开，请参考消息链接说明 https://open.dingtalk.com/document/app/message-link-description?spm=ding_open_doc.document.0.0.316448e0uHvnQD#section-7w8-4c2-9az
	MessageUrl string `json:"messageUrl"`
}

func (l *Link) String() string {
	return fmt.Sprintf(`{"title":%q,"text":%q,"picUrl":%q,"messageUrl":%q}`, l.Title, l.Text, l.PicUrl, l.MessageUrl)
}

func NewLink(title, text, picUrl, messageUrl string) *Link {
	msg := NewLinkForFeedCard(title, picUrl, messageUrl)
	msg.Text = text
	return msg
}

func NewLinkForFeedCard(title, picUrl, messageUrl string) *Link {
	return &Link{
		Title:      title,
		PicUrl:     picUrl,
		MessageUrl: messageUrl,
	}
}

func NewWhLinkMsg(text, title, picUrl, MessageUrl string) *WhLinkMsg {
	return &WhLinkMsg{
		MsgType: WhMsgTypeLink,
		Link: Link{
			Text:       text,
			Title:      title,
			PicUrl:     picUrl,
			MessageUrl: MessageUrl,
		},
	}
}

// WhMarkdownMsg markdown消息
type WhMarkdownMsg struct {
	MsgType  string   `json:"msgtype"`
	At       At       `json:"at"`
	MarkDown Markdown `json:"markdown"`
}

// Markdown markdown 消息体
type Markdown struct {
	//必填项
	//首屏会话透出的展示内容。
	Title string `json:"title"`
	// 必填项
	//markdown格式的消息内容。
	Text string `json:"text"`
}

func (m *Markdown) String() string {
	return fmt.Sprintf(`{"title": %q, "text": %q}`, m.Title, m.Text)
}

func NewWhMarkdownMsg(title, text string) *WhMarkdownMsg {
	return &WhMarkdownMsg{
		MsgType: WhMsgTypeMarkdown,
		MarkDown: Markdown{
			Title: title,
			Text:  text,
		},
	}
}

func NewWhMarkdownMsgWithAtMobiles(title, text string, mobiles ...string) *WhMarkdownMsg {
	markdownMsg := NewWhMarkdownMsg(title, text)
	markdownMsg.At.AtMobiles = mobiles
	return markdownMsg
}

func NewWhMarkdownMsgWithAtAll(title, text string) *WhMarkdownMsg {
	markdownMsg := NewWhMarkdownMsg(title, text)
	markdownMsg.At.IsAtAll = true
	return markdownMsg
}

func NewWhMarkdownMsgWithAtUserIds(title, text string, userIds ...string) *WhMarkdownMsg {
	// 被@人的用户userid。
	//说明 消息内容content中要带上"@userId"，跟atUserIds参数结合使用，才有@效果。
	a := make([]string, len(userIds))
	for i := 0; i < len(userIds); i++ {
		a[i] = fmt.Sprintf("@%s", userIds[i])
	}
	text = fmt.Sprintf("%s %s", text, strings.Join(a, " "))
	markdownMsg := NewWhMarkdownMsg(title, text)
	markdownMsg.At.AtUserIds = userIds
	return markdownMsg
}

// WhEntiretyActionCardMsg  整体跳转actionCard消息
type WhEntiretyActionCardMsg struct {
	MsgType    string              `json:"msgtype"`
	ActionCard *EntiretyActionCard `json:"actionCard"`
}

// EntiretyActionCard 整体跳转actionCard消息体
type EntiretyActionCard struct {
	//必填项
	// 首屏会话透出的展示内容。
	Title string `json:"title"`
	//必填项
	// markdown格式的消息内容。
	Text string `json:"text"`
	//必填项
	// 单个按钮的标题。
	SingleTitle string `json:"singleTitle"`
	//必填项
	// 单个按钮的跳转链接。
	//移动端  在钉钉客户端内打开。
	//PC端  默认侧边栏打开。
	//希望在外部浏览器打开，请参考消息链接说明 https://open.dingtalk.com/document/app/message-link-description?spm=ding_open_doc.document.0.0.316448e0uHvnQD#section-7w8-4c2-9az
	SingleURL string `json:"singleURL"`
}

func NewEntiretyActionCard(title, text, singleTitle, singleURL string) *EntiretyActionCard {
	return &EntiretyActionCard{
		Title:       title,
		Text:        text,
		SingleURL:   singleURL,
		SingleTitle: singleTitle,
	}
}

func (a *EntiretyActionCard) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func NewWhEntiretyActionCardMsg(title, text, singleTitle, singleURL string) *WhEntiretyActionCardMsg {
	return &WhEntiretyActionCardMsg{
		MsgType:    WhMsgTypeActionCard,
		ActionCard: NewEntiretyActionCard(title, text, singleTitle, singleURL),
	}
}

// WhIndependentActionCardMsg 独立跳转actionCard消息
type WhIndependentActionCardMsg struct {
	MsgType    string                `json:"msgtype"`
	ActionCard IndependentActionCard `json:"actionCard"`
}

// IndependentActionCard 独立跳转actionCard消息体
type IndependentActionCard struct {
	//必填项
	// 首屏会话透出的展示内容。
	Title string `json:"title"`
	//必填项
	// markdown格式的消息内容。
	Text string `json:"text"`
	// 可选项
	// 按钮排列顺序。
	//0：按钮竖直排列
	//1：按钮横向排列
	BtnOrientation string `json:"btnOrientation,omitempty"`
	//必填项
	// 按钮。
	Btns []*Btn `json:"btns"`
}

func NewWhIndependentActionCardMsg(title, text string, btns []*Btn) *WhIndependentActionCardMsg {
	return &WhIndependentActionCardMsg{
		MsgType: WhMsgTypeActionCard,
		ActionCard: IndependentActionCard{
			Title: title,
			Text:  text,
			Btns:  btns,
		},
	}
}

func NewWhIndependentActionCardMsgWithBtnOrientation(title, text, btnOrientation string, btns []*Btn) *WhIndependentActionCardMsg {
	msg := NewWhIndependentActionCardMsg(title, text, btns)
	msg.ActionCard.BtnOrientation = btnOrientation
	return msg
}

// Btn 按钮。
type Btn struct {
	//必填项
	//按钮标题。
	Title string `json:"title"`
	//必填项
	//点击按钮触发的URL。
	//移动端    在钉钉客户端内打开。
	//PC端    默认侧边栏打开。
	//希望在外部浏览器打开，请参考消息链接说明 https://open.dingtalk.com/document/app/message-link-description?spm=ding_open_doc.document.0.0.316448e0uHvnQD#section-7w8-4c2-9az
	ActionURL string `json:"actionURL"`
}

func NewBtn(title, actionURL string) *Btn {
	return &Btn{
		Title:     title,
		ActionURL: actionURL,
	}
}

// WhFeedCardMsg FeedCard消息
type WhFeedCardMsg struct {
	MsgType  string     `json:"msgtype"`
	FeedCard WhFeedCard `json:"feedCard"`
}

// WhFeedCard FeedCard消息体
type WhFeedCard struct {
	Links []*Link `json:"links"`
}

func NewWhFeedCardMsg(links []*Link) *WhFeedCardMsg {
	return &WhFeedCardMsg{
		MsgType: WhMsgTypeFeedCard,
		FeedCard: WhFeedCard{
			Links: links,
		},
	}
}
