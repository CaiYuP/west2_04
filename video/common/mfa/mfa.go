package mfa

import (
	"encoding/base64"
	"fmt"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// 1. 生成 TOTP secret
func GenerateSecret(username string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "west2", // 显示在 App 里的发行者名称
		AccountName: username,
	})
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}

// 2. 按 otpauth 规范拼接 URI
func BuildOtpAuthURL(secret, account, issuer string) string {
	// 典型格式：
	// otpauth://totp/{issuer}:{account}?secret={secret}&issuer={issuer}&algorithm=SHA1&digits=6&period=30
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		issuer, account, secret, issuer,
	)
}

// 3 + 4. 生成二维码 PNG 并转为 data:image/png;base64,...
func GenerateMFAQRCodeDataURL(otpURL string) (string, error) {
	// 使用 go-qrcode 直接生成 PNG 字节
	pngBytes, err := qrcode.Encode(otpURL, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// Base64 编码
	b64 := base64.StdEncoding.EncodeToString(pngBytes)

	// 拼 data URL
	dataURL := "data:image/png;base64," + b64
	return dataURL, nil
}
