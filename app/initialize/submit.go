package initialize

import (
	"github.com/Bekreth/jane_cli/app/terminal"
)

func (state *initState) submit() {
	flags := terminal.ParseFlags(state.currentBuffer)
	state.logger.Debugf("submitting query flags: %v", flags)

	if _, exists := flags["help"]; exists {
		state.printHelp()
		state.currentBuffer = ""
		return
	}

	missingFlags := map[string]string{}
	if state.user.Auth.Domain == "" {
		missingFlags[clinicDomain] = ""
	}
	if state.user.Auth.Username == "" {
		missingFlags[username] = ""
	}

	if domain, exists := flags[clinicDomain]; exists {
		delete(missingFlags, clinicDomain)
		state.user.Auth.Domain = domain
	}
	if providedUserName, exists := flags[username]; exists {
		delete(missingFlags, username)
		state.user.Auth.Username = providedUserName
	}

	if len(missingFlags) != 0 {
		missingParameters := []string{}
		if _, exists := missingFlags[clinicDomain]; exists {
			missingParameters = append(missingParameters, "clinic domain")
		}
		if _, exists := missingFlags[username]; exists {
			missingParameters = append(missingParameters, "username")
		}
		state.writer.WriteStringf("missing user data %v", missingParameters)
	}

	err := state.user.SaveUserFile()
	if err != nil {
		state.logger.Infof("error writing userfile: %v", err)
	}

	state.nextState = state.rootState
}
