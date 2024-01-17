GO_FILES=$(shell find . -type f -name "*.go")
NAME=jane_cli
ICON=${NAME}.ico
TAG=$(shell git describe --tags | sed 's|^v||' | sed 's|\(\.*\)-.*|\1|')
COMPILE_COMMAND=go build -ldflags="-X main.version=${TAG}"
