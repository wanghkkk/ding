package ding

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
)

// GetSign 获取前面， 传入数据data，和密钥secret
func GetSign(data, secret string) string {
    m := hmac.New(sha256.New, []byte(secret))
    m.Write([]byte(data))
    return base64.StdEncoding.EncodeToString(m.Sum(nil))
}

// GetDingSign 参考： https://open.dingtalk.com/document/robots/verify-valid-requests-for-intra-enterprise-robot-group-chat
func GetDingSign(timestamp, secret string) string {
    data := timestamp + "\n" + secret
    return GetSign(data, secret)
}
