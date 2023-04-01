package logger

type Config struct {
	Debugger bool   `yaml:"debugger"`
	Output   string `yaml:"output"`
}
