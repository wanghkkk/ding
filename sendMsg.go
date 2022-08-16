package ding

// SendMsgWithUserIds 发送消息， 如果是单聊则给userIds，如果是群聊则@userIds
type SendMsgWithUserIds interface {
	// SendTextMsgWithUserIds 发送文本消息， 如果是单聊则给userIds，如果是群聊则@userIds
	SendTextMsgWithUserIds(content string, userIds []string) error
	// SendMarkdownMsgWithUserIds 发送markdown消息， 如果是单聊则给userIds，如果是群聊则@userIds
	SendMarkdownMsgWithUserIds(title, text string, userIds []string) error
}
