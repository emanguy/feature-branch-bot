package config

import (
	"fmt"
	"io/ioutil"
	"os"
)

type CredentialSource interface {
	RetrieveCredential() (string, error)
}

type CredentialSourceIdentifier struct {
	Type string
}

type FileCredentialSource struct {
	location string
}

func (fcs FileCredentialSource) RetrieveCredential() (string, error) {
	fileContent, fileReadErr := ioutil.ReadFile(fcs.location)
	if fileReadErr != nil {
		return "", fmt.Errorf("failed to read ssh credential at %v: %w", fcs.location, fileReadErr)
	}

	return string(fileContent), nil
}

type EnvironmentCredentialSource struct {
	variableName string
}

func (ecs EnvironmentCredentialSource) RetrieveCredential() (string, error) {
	envVarContent, envVarExists := os.LookupEnv(ecs.variableName)
	if !envVarExists {
		return "", fmt.Errorf("failed to read ssh credential from %v, the variable did not exist", ecs.variableName)
	}

	return envVarContent, nil
}

type InlineCredentialSource struct {
	value string
}

func (ics InlineCredentialSource) RetrieveCredential() (string, error) {
	return ics.value, nil
}
