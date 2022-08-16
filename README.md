
## golang 发送钉钉消息

### 支持钉钉两种方式

### 方式1： 钉钉接口方式

- 适用于钉钉企业内部应用或钉钉企业内部机器人
- 支持群聊和单聊
- 群聊的text和markdown消息支持at某人

### 方式2： 钉钉webhook方式 【推荐】

- 适用于钉钉群自定义机器人
  - 也支持钉钉企业内部机器人，使用：NewWhClientUseSessionWebhook， 传入钉钉企业内部机器人发来的post请求中的sessionWebhookUrl
  - 此时能够支持群聊和单聊
- 使用钉钉群自定义机器人只支持群消息
- 群聊的text和markdown消息支持at某人
- 支持钉钉安全设置，加签，暂不支持关键字