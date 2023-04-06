package client

type Config struct {
	UserFilePath string `yaml:"userFilePath"`
}

var DefaultConfig = Config{}
