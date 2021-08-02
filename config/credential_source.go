package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type CredentialSource interface {
	RetrieveCredential() (string, error)
}

type CredentialSourceIdentifier struct {
	Type string `json:"type"`
}

type FileCredentialSource struct {
	Location string `json:"location"`
}

func (fcs FileCredentialSource) RetrieveCredential() (string, error) {
	fileContent, fileReadErr := ioutil.ReadFile(fcs.Location)
	if fileReadErr != nil {
		return "", fmt.Errorf("failed to read ssh credential at %v: %w", fcs.Location, fileReadErr)
	}

	return string(fileContent), nil
}

type EnvironmentCredentialSource struct {
	VariableName string `json:"envVar"`
}

func (ecs EnvironmentCredentialSource) RetrieveCredential() (string, error) {
	envVarContent, envVarExists := os.LookupEnv(ecs.VariableName)
	if !envVarExists {
		return "", fmt.Errorf("failed to read ssh credential from %v, the variable did not exist", ecs.VariableName)
	}

	return envVarContent, nil
}

type InlineCredentialSource struct {
	Value string `json:"value"`
}

func (ics InlineCredentialSource) RetrieveCredential() (string, error) {
	return ics.Value, nil
}

func unmarshalCredentialSource(data []byte, sourceTarget *CredentialSource) error {
	var credentialType CredentialSourceIdentifier
	if err := json.Unmarshal(data, &credentialType); err != nil {
		return err
	}

	switch credentialType.Type {
	case "FILE":
		var credentialSource FileCredentialSource
		if err := json.Unmarshal(data, &credentialSource); err != nil {
			return err
		}
		*sourceTarget = credentialSource
		return nil
	case "ENVIRONMENT":
		var credentialSource EnvironmentCredentialSource
		if err := json.Unmarshal(data, &credentialSource); err != nil {
			return err
		}
		*sourceTarget = credentialSource
		return nil
	case "INLINE":
		var credentialSource InlineCredentialSource
		if err := json.Unmarshal(data, &credentialSource); err != nil {
			return err
		}
		*sourceTarget = credentialSource
		return nil
	}

	return fmt.Errorf("%v is not a valid credential source type", credentialType.Type)
}
