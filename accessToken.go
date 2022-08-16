package ding

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
)

var (
    AccessTokenUrl = "https://api.dingtalk.com/v1.0/oauth2/accessToken"
)

// AppKeySecret 企业内部应用的appKey, appSecret
type AppKeySecret struct {
    // 已创建的企业内部应用的AppKey。
    AppKey string `json:"appKey"`
    // 已创建的企业内部应用的AppSecret。
    AppSecret string `json:"appSecret"`
}

// AccessToken 获取到的钉钉 accessToken和过期时间，需要缓存下来，避免重复获取
type AccessToken struct {
    // 生成的accessToken。
    Token string `json:"accessToken"`
    // accessToken的过期时间，单位秒。
    ExpireIn int64 `json:"expireIn"`
}

// getAccessTokenFromDing 从钉钉获取企业内部应用的accessToken
// 参考： https://open.dingtalk.com/document/orgapp-server/obtain-the-access_token-of-an-internal-app
func getAccessTokenFromDing(aks AppKeySecret) (*AccessToken, error) {
    ks := AppKeySecret{
        AppKey:    aks.AppKey,
        AppSecret: aks.AppSecret,
    }
    ksByte, err := json.Marshal(ks)
    if err != nil {
        return nil, err
    }
    resp, err := http.Post(AccessTokenUrl, ContentTypeJson, bytes.NewBuffer(ksByte))
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    datByte, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var dat AccessToken

    err = json.Unmarshal(datByte, &dat)
    if err != nil {
        return nil, err
    }

    return &dat, nil
}

// GetAccessToken 获取access token 先从cache，否则在从钉钉
func GetAccessToken(aks AppKeySecret) (string, error) {
    datKey := "DingAccessToken"
    tokenByte, err := cache.Get(datKey)
    if err != nil {
        // 缓存里没有，就请求钉钉获取
        at, err := getAccessTokenFromDing(aks)
        // 从钉钉也没获取到就没办法，返回错误了
        if err != nil {
            return "", err
        }
        err = cache.Set(datKey, []byte(at.Token))
        if err != nil {
            // 写入cache失败，并不是大问题，大不了再从钉钉获取，先让调用者能用再说
            log.Println("set ding access token to cache failed: " + err.Error())
        }
        return at.Token, nil
    }

    return string(tokenByte), nil
}
