package initialize

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

const username = "-u"
const clinicDomain = "-d"

type initState struct {
	logger    logger.Logger
	writer    terminal.ScreenWriter
	user      *domain.User
	rootState terminal.State

	currentBuffer string
	nextState     terminal.State
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	user *domain.User,
	rootState terminal.State,
) terminal.State {
	return &initState{
		logger:    logger,
		writer:    writer,
		user:      user,
		rootState: rootState,
	}
}

func (initState) Name() string {
	return "init"
}

func (init *initState) Initialize() {
	init.logger.Debugf(
		"entering init. available states %v",
		init.rootState.Name(),
	)
	init.nextState = init
	init.writer.NewLine()
	init.writer.WriteString("")
}

func (init *initState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, &init.currentBuffer, init.triggerAutocomplete, init.submit)

	if character != 0 {
		init.currentBuffer += string(character)
	}

	init.writer.WriteString(init.currentBuffer)
	return init.nextState
}

func (init *initState) triggerAutocomplete() {
}

func (init *initState) submit() {
	flags := terminal.ParseFlags(init.currentBuffer)
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
		init.writer.WriteString(fmt.Sprintf("missing user data %v", missingParameters))
	}

	err := init.user.SaveUserFile()
	if err != nil {
		init.logger.Infof("error writing userfile: %v", err)
	}

	init.nextState = init.rootState
}
