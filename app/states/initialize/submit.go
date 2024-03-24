package initialize

import "fmt"

func (state *initState) Submit(flags map[string]string) bool {
	if _, exists := flags[".."]; exists {
		state.nextState = state.rootState
		return true
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
		//TODO: Test this
		state.buffer.AddString(fmt.Sprintf(
			"missing user data %v",
			missingParameters,
		))
	}

	err := state.user.SaveUserFile()
	if err != nil {
		state.logger.Infof("error writing userfile: %v", err)
	}
	return true
}
