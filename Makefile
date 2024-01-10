GO_FILES=$(shell find . -type f -name "*.go")
WINDOWS_OUTPUT=output/windows
MSI_PACKAGER=windows_msi.xml

# Run commands

msi: ${WINDOWS_OUTPUT}/jane_cli.msi

clean:
	rm -rf output

test:
	go test ./...

output/:
	@mkdir $@

# Windows

${WINDOWS_OUTPUT}/: | output/
	@mkdir $@

${WINDOWS_OUTPUT}/etc/: | ${WINDOWS_OUTPUT}/
	@mkdir $@

${WINDOWS_OUTPUT}/var/: | ${WINDOWS_OUTPUT}/
	@mkdir $@

${WINDOWS_OUTPUT}/jane_cli.exe: ${GO_FILES}| ${WINDOWS_OUTPUT}/
	@GOOS=windows GOARCH=amd64 go build  -o $@ .

${WINDOWS_OUTPUT}/jane_cli.wxs: packaging/windows_msi.xml | ${WINDOWS_OUTPUT}/
	@cp $< $@

${WINDOWS_OUTPUT}/etc/config.yaml: etc/config.yaml | ${WINDOWS_OUTPUT}/etc/
	@cp $< $@

${WINDOWS_OUTPUT}/var/user.yaml: etc/user.template.yaml | ${WINDOWS_OUTPUT}/var/
	@cp $< $@

${WINDOWS_OUTPUT}/var/log: | ${WINDOWS_OUTPUT}/var/
	@touch $@

${WINDOWS_OUTPUT}/jane_cli.wixobj: ${WINDOWS_OUTPUT}/jane_cli.wxs \
	${WINDOWS_OUTPUT}/jane_cli.exe \
	${WINDOWS_OUTPUT}/etc/config.yaml \
	${WINDOWS_OUTPUT}/var/user.yaml \
	${WINDOWS_OUTPUT}/var/log
	docker run --rm -v $(shell pwd)/output/windows:/wix dactiv/wix candle jane_cli.wxs

${WINDOWS_OUTPUT}/jane_cli.msi: ${WINDOWS_OUTPUT}/jane_cli.wixobj
	docker run --rm -v $(shell pwd)/output/windows:/wix dactiv/wix light ${<F} -sval

# Linux
output/jane_cli: ${GO_FILES} | output/
	@GOOS=linux GOARCH=amd64 go build -o $@ . 
