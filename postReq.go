package ding

import "fmt"

var (
	// Dan 单聊
	Dan = "1"
	// Qun 群聊
	Qun = "2"
)

// PostReq  钉钉发来的post请求
type PostReq struct {
	//加密的会话ID。即群ID，可用于通过接口给群发送消息
	ConversationId string `json:"conversationId"`
	// 1单聊   2群聊
	ConversationType string `json:"conversationType"`
	// 群聊时才有的会话标题。
	ConversationTitle string `json:"conversationTitle,omitempty"`
	// 被@人的信息。
	AtUsers []AtUser `json:"atUsers"`
	// 加密的会话ID。
	ChatbotCorpId string `json:"chatbotCorpId,omitempty"`
	// 加密的机器人ID。
	ChatbotUserId string `json:"chatbotUserId"`
	// 加密的消息ID。
	MsgId string `json:"msgId"`
	// 消息类型，一般为text
	Msgtype string `json:"msgtype"`
	// 发送者昵称。
	SenderNick string `json:"senderNick"`
	//加密的发送者ID。
	SenderId string `json:"senderId"`
	// 企业内部群中@该机器人的成员userid。该字段在机器人发布线上版本后，才会返回。
	SenderStaffId string `json:"senderStaffId,omitempty"`
	// 企业内部群有的发送者当前群的企业corpId。
	SenderCorpId string `json:"senderCorpId,omitempty"`
	// 是否为管理员。
	IsAdmin bool `json:"isAdmin,omitempty"`
	// 当前会话的Webhook地址过期时间。
	SessionWebhookExpiredTime int64 `json:"sessionWebhookExpiredTime"`
	//消息的时间戳，单位ms。
	CreateAt int64 `json:"createAt"`
	//是否在@列表中。
	IsInAtList bool `json:"isInAtList,omitempty"`
	// 当前会话的Webhook地址。
	SessionWebhook string `json:"sessionWebhook"`
	// 消息
	Text Text `json:"text"`
	// 机器人code， 一般为normal
	RobotCode string `json:"robotCode"`
}

// AtUser at用户
type AtUser struct {
	// 加密的发送者ID。
	DingtalkId string `json:"dingtalkId"`
	// 当前企业内部群中员工userid值。
	StaffId string `json:"staffId"`
}

//String 优雅打印postReq
func (p *PostReq) String() string {
	return fmt.Sprintf(`{
    "conversationId": %q,
    "atUsers": [
        {
            "dingtalkId": %q,
            "staffId": %q
        }
    ],
    "chatbotUserId": %q,
    "msgId": %q,
    "senderNick": %q,
    "isAdmin": %t,
	"senderStaffId": %q,
	"senderCorpId": %q,
    "sessionWebhookExpiredTime": %d,
    "createAt": %d,
    "conversationType": %q,
    "senderId": %q,
    "conversationTitle": %q,
    "isInAtList": %t,
    "sessionWebhook": %q,
    "text": {
        "content": %q
    },
    "robotCode": %q,
    "msgtype": %q
}`,
		p.ConversationId,
		p.AtUsers[0].DingtalkId,
		p.AtUsers[0].StaffId,
		p.ChatbotUserId,
		p.MsgId,
		p.SenderNick,
		p.IsAdmin,
		p.SenderStaffId,
		p.SenderCorpId,
		p.SessionWebhookExpiredTime,
		p.CreateAt,
		p.ConversationType,
		p.SenderId,
		p.ConversationTitle,
		p.IsInAtList,
		p.SessionWebhook,
		p.Text.Content,
		p.RobotCode,
		p.Msgtype,
	)
}
