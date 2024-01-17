WINDOWS_OUTPUT=output/windows
MSI_PACKAGER=windows_msi.xml

${WINDOWS_OUTPUT}/: | output/
	@mkdir $@

${WINDOWS_OUTPUT}/etc/: | ${WINDOWS_OUTPUT}/
	@mkdir $@

${WINDOWS_OUTPUT}/var/: | ${WINDOWS_OUTPUT}/
	@mkdir $@

${WINDOWS_OUTPUT}/${NAME}.exe: ${GO_FILES}| ${WINDOWS_OUTPUT}/
	@GOOS=windows GOARCH=amd64 ${COMPILE_COMMAND} -o $@ .

${WINDOWS_OUTPUT}/${NAME}.wxs: packaging/windows_msi.xml | ${WINDOWS_OUTPUT}/
	@cat $< | sed 's|$$VERSION|${TAG}|' > $@

${WINDOWS_OUTPUT}/etc/config.yaml: etc/config.yaml | ${WINDOWS_OUTPUT}/etc/
	@cp $< $@

${WINDOWS_OUTPUT}/var/user.yaml: etc/user.template.yaml | ${WINDOWS_OUTPUT}/var/
	@cp $< $@

${WINDOWS_OUTPUT}/var/log: | ${WINDOWS_OUTPUT}/var/
	@touch $@

${WINDOWS_OUTPUT}/${ICON}: packaging/${ICON} | ${WINDOWS_OUTPUT}/var/
	@cp $< $@

${WINDOWS_OUTPUT}/${NAME}.wixobj: ${WINDOWS_OUTPUT}/${NAME}.wxs \
	${WINDOWS_OUTPUT}/${ICON} \
	${WINDOWS_OUTPUT}/${NAME}.exe \
	${WINDOWS_OUTPUT}/etc/config.yaml \
	${WINDOWS_OUTPUT}/var/user.yaml \
	${WINDOWS_OUTPUT}/var/log
	@echo Building wixobj
	@docker run --rm -v $(shell pwd)/output/windows:/wix dactiv/wix candle ${NAME}.wxs

${WINDOWS_OUTPUT}/${NAME}_${TAG}.msi: ${WINDOWS_OUTPUT}/${NAME}.wixobj
	@echo Building msi
	@docker run --rm \
		-v $(shell pwd)/output/windows:/wix \
		dactiv/wix light \
		${<F} -sval -out ${@F}

