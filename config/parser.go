package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ParseConfig reads the JSON configuration file for the bot and returns the parsed information.
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

// DetermineSSHCredentialsForProject determines the SSH Credentials to use for a given project given the
// values provided for the server and individual project. The project value takes priority over the server
// value, if it exists.
func DetermineSSHCredentialsForProject(vcsServer VCSServer, vcsProject VCSProject) (SSHCredentials, error) {
	if vcsProject.SSHCreds.IsPresent {
		return vcsProject.SSHCreds, nil
	}
	if vcsServer.SSHCreds.IsPresent {
		return vcsServer.SSHCreds, nil
	}

	return SSHCredentials{}, fmt.Errorf("no ssh credentials are present between project and server")
}

// DetermineSyncTagForProject determines which label should be used for discovering merge requests on a project
// that should be synchronized with their destination branches. The project value takes precedence over the server
// value if it is provided.
func DetermineSyncTagForProject(serverTag, projectTag string) (string, error) {
	if len(projectTag) > 0 {
		return projectTag, nil
	}
	if len(serverTag) > 0 {
		return serverTag, nil
	}

	return "", fmt.Errorf("no synchronization tag is present between project and server")
}
