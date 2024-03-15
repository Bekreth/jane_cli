module github.com/Bekreth/jane_cli

go 1.21

toolchain go1.21.0

require (
	github.com/bekreth/screen_reader_terminal v0.0.0-20240212023332-7b8acfac6fde
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kopoli/go-terminal-size v0.0.0-20170219200355-5c97524c8b54 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace github.com/bekreth/screen_reader_terminal => ../screen_reader_terminal
