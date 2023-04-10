GO_FILES=$(shell find . -type f -name "*.go")

output/:
	@mkdir $@

output/jane_cli.exe: ${GO_FILES}| output/
	@GOOS=windows GOARCH=amd64 go build  -o $@ .

output/jane_cli: ${GO_FILES} | output/
	@GOOS=linux GOARCH=amd64 go build -o $@ . 

clean:
	rm -rf output

build: output/jane_cli output/jane_cli.exe
