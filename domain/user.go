package domain

import (
	"fmt"
	"io"
	"os"

	"github.com/Bekreth/jane_cli/logger"
	"gopkg.in/yaml.v3"
)

type User struct {
	filePath string
	logger   logger.Logger
	Auth     Auth `yaml:"auth"`
}

func NewUser(logger logger.Logger, filePath string) (User, error) {
	output := User{
		filePath: filePath,
		logger:   logger,
	}

	userFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("The provided user file at %v is missing.  Creating one now\n", filePath)
		err = output.SaveUserFile()
		if err != nil {
			fmt.Printf("Unable to create new user file at %v: %v", filePath, err)
			return output, err
		}
		return output, nil
	}

	bytes, err := io.ReadAll(userFile)
	if err != nil {
		fmt.Printf("Unable to read from user file at %v\n", userFile)
		return output, err
	}

	err = yaml.Unmarshal(bytes, &output)
	if err != nil {
		logger.Infof("failed to read user file %v: %v", filePath, err)
		newFilePath := filePath + ".swap"
		fmt.Printf(
			"Unable to parse user file at %v.  Moving it to %v and creating new user file\n",
			filePath,
			newFilePath,
		)
		newFile, err := os.Create(newFilePath)
		if err != nil {
			logger.Infof("unable to copy existing user file to %v: %v", newFilePath, err)
			fmt.Printf("Unable to copy existing user file to %v\n", newFilePath)
			return output, err
		}
		_, err = newFile.Write(bytes)
		if err != nil {
			logger.Infof("unable to copy existing user file to %v: %v", newFilePath, err)
			fmt.Printf("Unable to copy existing user file to %v\n", newFilePath)
			os.Remove(newFilePath)
			return output, err
		}
	}

	return output, nil
}

func (user User) PostCheck() {
	message := "User data isn't fully initialized.  Missing data for %v, run the following commands: \n%v\n"
	missingData := []string{}
	commands := []string{}

	shouldPrint := false

	if user.Auth.Domain == "" {
		missingData = append(missingData, "clinic domain")
		commands = append(commands, "* init -d ${clinicDomain}")
		shouldPrint = true
	}
	if user.Auth.Username == "" {
		missingData = append(missingData, "username")
		commands = append(commands, "* init -u ${username}")
		shouldPrint = true
	}
	if user.Auth.AuthCookie == "" {
		missingData = append(missingData, "authentication token")
		commands = append(commands, "* auth -p ${password}")
		shouldPrint = true
	}

	if shouldPrint {
		fmt.Printf(message, missingData, commands)
	}
}

func (user User) SaveUserFile() error {
	userFile, err := os.OpenFile(user.filePath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		user.logger.Infof("the user file %v doesn't exist, creating it", user.filePath)
		userFile, err = os.Create(user.filePath)
		if err != nil {
			user.logger.Infof("unable to create user file at %v: %v", user.filePath, err)
			return err
		}
	}

	bytes, err := yaml.Marshal(user)
	if err != nil {
		return err
	}

	_, err = userFile.Write(bytes)
	if err != nil {
		user.logger.Infof("unable to write date to user file %v", user.filePath)
		return err
	}

	return nil
}
