WINDOWS_OUTPUT=output/windows
MSI_OUTPUT=${WINDOWS_OUTPUT}/${TAG}
PREVIOUS_OUTPUT=${WINDOWS_OUTPUT}/${PREVIOUS}
PATCH=${NAME}_patch
MSI_PACKAGER=windows_msi.xml

${WINDOWS_OUTPUT}/: | output/
	@mkdir $@

${MSI_OUTPUT}/: | ${WINDOWS_OUTPUT}/
	@mkdir $@

${MSI_OUTPUT}/etc/: | ${MSI_OUTPUT}/
	@mkdir $@

${MSI_OUTPUT}/var/: | ${MSI_OUTPUT}/
	@mkdir $@

${MSI_OUTPUT}/${NAME}.exe: ${GO_FILES}| ${MSI_OUTPUT}/
	@GOOS=windows GOARCH=amd64 ${COMPILE_COMMAND} -o $@ .

${MSI_OUTPUT}/etc/config.yaml: etc/config.yaml | ${MSI_OUTPUT}/etc/
	@cp $< $@

${MSI_OUTPUT}/var/user.yaml: etc/user.template.yaml | ${MSI_OUTPUT}/var/
	@cp $< $@

${MSI_OUTPUT}/var/log: | ${MSI_OUTPUT}/var/
	@touch $@

${MSI_OUTPUT}/${ICON}: packaging/${ICON} | ${MSI_OUTPUT}/var/
	@cp $< $@


# Making the MSI

${MSI_OUTPUT}/${NAME}_${TAG}.wxs: packaging/windows_msi.xml | ${MSI_OUTPUT}/
	@cat $< | sed 's|$$VERSION|${TAG}|' > $@

${MSI_OUTPUT}/${NAME}_${TAG}.wixobj ${MSI_OUTPUT}/${NAME}_${TAG}.wixpdb: \
	${MSI_OUTPUT}/${NAME}_${TAG}.wxs \
	${MSI_OUTPUT}/${ICON} \
	${MSI_OUTPUT}/${NAME}.exe \
	${MSI_OUTPUT}/etc/config.yaml \
	${MSI_OUTPUT}/var/user.yaml \
	${MSI_OUTPUT}/var/log
	@echo Building wixobj
	@docker run --rm \
		-v $(shell pwd)/${MSI_OUTPUT}:/wix \
		dactiv/wix candle ${<F}

${MSI_OUTPUT}/${NAME}_${TAG}.msi: ${MSI_OUTPUT}/${NAME}_${TAG}.wixobj
	@echo Building msi
	@docker run --rm \
		-v $(shell pwd)/${MSI_OUTPUT}:/wix \
		dactiv/wix light \
		${<F} -sval -out ${@F}


# Making the MSP

${MSI_OUTPUT}/${PATCH}_${TAG}.wxs: packaging/windows_msp.xml | ${MSI_OUTPUT}/
	@cat $< | sed 's|$$VERSION|${TAG}|' > $@

${MSI_OUTPUT}/${PATCH}_${TAG}.wixmst: \
	${MSI_OUTPUT}/${NAME}_${TAG}.wixpdb \
	${PREVIOUS_OUTPUT}/${NAME}_${PREVIOUS}.wixpdb
	@echo Building msp
	@docker run --rm \
		-v $(shell pwd)/${WINDOWS_OUTPUT}:/wix \
		dactiv/wix torch \
		-p -xi ${TAG}/${NAME}_${TAG}.wixpdb  ${PREVIOUS}/${NAME}_${TAG}.wixpdb \
		-out ${TAG}/${@F}

${MSI_OUTPUT}/${PATCH}_${TAG}.wixobj: ${MSI_OUTPUT}/${PATCH}_${TAG}.wxs
	@echo Building wixobj for msp
	@docker run --rm \
		-v $(shell pwd)/${MSI_OUTPUT}:/wix \
		dactiv/wix candle ${<F}

${MSI_OUTPUT}/${PATCH}_${TAG}.wixmsp: ${MSI_OUTPUT}/${PATCH}_${TAG}.wixobj
	@echo Building wixmsp
	@docker run --rm \
		-v $(shell pwd)/${MSI_OUTPUT}:/wix \
		dactiv/wix light \
		${<F} -sval -out ${@F}

${MSI_OUTPUT}/${NAME}_${TAG}.msp: \
	${MSI_OUTPUT}/${PATCH}_${TAG}.wixmsp \
	${MSI_OUTPUT}/${PATCH}_${TAG}.wixmst
	@echo Building msp
	@docker run --rm \
		-v $(shell pwd)/${MSI_OUTPUT}:/wix \
		dactiv/wix pyro \
		${PATCH}_${TAG}.wixmsp -out ${@F} -t Sample ${PATCH}_${TAG}.wixmst
