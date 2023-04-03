package app

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type initState struct {
	logger        logger.Logger
	writer        screenWriter
	rootState     state
	user          *domain.User
	currentBuffer string
	nextState     state
}

const username = "-u"
const clinicDomain = "-d"

func (initState) name() string {
	return "init"
}

func (init *initState) initialize() {
	init.logger.Debugf(
		"entering init. available states %v",
		init.rootState.name(),
	)
	init.nextState = init
	init.writer.newLine()
	init.writer.writeString("")
}

func (init *initState) handleKeyinput(character rune, key keyboard.Key) state {
	keyHandler(key, &init.currentBuffer, init.triggerAutocomplete, init.submit)

	if character != 0 {
		init.currentBuffer += string(character)
	}

	init.writer.writeString(init.currentBuffer)
	return init.nextState
}

func (init *initState) triggerAutocomplete() {
}

func (init *initState) submit() {
	flags := parseFlags(init.currentBuffer)
	init.logger.Debugf("submitting query flags: %v", flags)

	missingFlags := map[string]string{}
	if init.user.Auth.Domain == "" {
		missingFlags[clinicDomain] = ""
	}
	if init.user.Auth.Username == "" {
		missingFlags[username] = ""
	}

	if domain, exists := flags[clinicDomain]; exists {
		delete(missingFlags, clinicDomain)
		init.user.Auth.Domain = domain
	}
	if providedUserName, exists := flags[username]; exists {
		delete(missingFlags, username)
		init.user.Auth.Username = providedUserName
	}

	if len(missingFlags) != 0 {
		missingParameters := []string{}
		if _, exists := missingFlags[clinicDomain]; exists {
			missingParameters = append(missingParameters, "clinic domain")
		}
		if _, exists := missingFlags[username]; exists {
			missingParameters = append(missingParameters, "username")
		}
		init.writer.writeString(fmt.Sprintf("missing user data %v", missingParameters))
	}

	err := init.user.SaveUserFile()
	if err != nil {
		init.logger.Infof("error writing userfile: %v", err)
	}

	init.nextState = init.rootState
}

func (init *initState) shutdown() {
	init.currentBuffer = ""
}
