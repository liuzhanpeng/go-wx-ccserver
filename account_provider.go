package wxccserver

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// AccountProvider 公众号账号信息提供器
type AccountProvider interface {
	// 返回公众号账号列表
	Obtain() ([]*Account, error)
}

// TomlAccountProvider 基于toml格式文件的公众号账号提供器
type TomlAccountProvider struct {
	filename string // 配置文件
}

// NewTomlAccountProvider 创建公众号账号提供器
func NewTomlAccountProvider(filename string) *TomlAccountProvider {
	return &TomlAccountProvider{
		filename: filename,
	}
}

// Obtain 返回公众号账号列表
func (p *TomlAccountProvider) Obtain() ([]*Account, error) {
	viper.SetConfigFile(p.filename)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "加载公众号账号配置错误")
	}

	var cfg struct {
		Accounts []*Account
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "加载公众号账号配置错误")
	}

	return cfg.Accounts, nil
}
