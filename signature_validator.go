package wxccserver

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"strconv"
	"time"
)

var (
	ErrTimeout = errors.New("签名超时")
	ErrInvalid = errors.New("无效签名")
)

// SignatureValidator 签名验证器接口
type SignatureValidator interface {
	Validate(appID string, timestamp int64, signature string) error
}

// HMacSignatureValidator 基于HMac算法(md5)的签名验证组件
type HMacSignatureValidator struct {
	key     string // 签名密钥
	timeout int    // 签名超时时间(单位:秒)
}

// NewHMacSignatureValidator 创建基于HMac算法的签名验证组件
func NewHMacSignatureValidator(key string, timeout int) *HMacSignatureValidator {
	return &HMacSignatureValidator{
		key:     key,
		timeout: timeout,
	}
}

// Validate 验证签名合法性
func (validator *HMacSignatureValidator) Validate(appID string, timestamp int64, signature string) error {
	if time.Now().Add(-time.Duration(validator.timeout)*time.Second).Unix() > timestamp {
		return ErrTimeout
	}

	mac := hmac.New(md5.New, []byte(validator.key))
	mac.Write([]byte(appID + strconv.FormatInt(timestamp, 10)))

	s, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return ErrInvalid
	}

	if hmac.Equal(mac.Sum(nil), []byte(s)) {
		return nil
	}

	return ErrInvalid
}
