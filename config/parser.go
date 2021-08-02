package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseConfig(fileName string) (BotConfig, error) {
	var configuration BotConfig

	fileData, fileReadErr := ioutil.ReadFile(fileName)
	if fileReadErr != nil {
		return BotConfig{}, fmt.Errorf("failed to read the configuration file at %v: %w", fileName, fileReadErr)
	}

	if parseErr := json.Unmarshal(fileData, &configuration); parseErr != nil {
		return BotConfig{}, fmt.Errorf("parsing of configuration file failed: %w", parseErr)
	}

	return configuration, nil
}

func DetermineSSHCredentialsForProject(vcsServer VCSServer, vcsProject VCSProject) (SSHCredentials, error) {
	if vcsProject.SSHCreds.IsPresent {
		return vcsProject.SSHCreds, nil
	}
	if vcsServer.SSHCreds.IsPresent {
		return vcsServer.SSHCreds, nil
	}

	return SSHCredentials{}, fmt.Errorf("no ssh credentials are present between project and server")
}

func DetermineSyncTagForProject(serverTag, projectTag string) (string, error) {
	if len(projectTag) > 0 {
		return projectTag, nil
	}
	if len(serverTag) > 0 {
		return serverTag, nil
	}

	return "", fmt.Errorf("no synchronization tag is present between project and server")
}
