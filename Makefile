include makefiles/common.mk
include makefiles/windows.mk


# Run commands
msi: ${WINDOWS_OUTPUT}/${NAME}.msi

clean:
	rm -rf output

test:
	go test ./...

output/:
	@mkdir $@

# Linux
output/${NAME}: ${GO_FILES} | output/
	@GOOS=linux GOARCH=amd64 go build -o $@ . 
