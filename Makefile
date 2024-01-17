include makefiles/common.mk
include makefiles/windows.mk


# Run commands
msi: ${MSI_OUTPUT}/${NAME}_${TAG}.msi
msp: ${MSI_OUTPUT}/${NAME}_${TAG}.msp

clean:
	rm -rf output

test:
	@go test ./...

output/:
	@mkdir $@

# Linux
output/${NAME}: ${GO_FILES} | output/
	@GOOS=linux GOARCH=amd64 ${COMPILE_COMMAND} -o $@ . 
