package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// CredentialSource is an interface for retrieving a credential via various strategies.
type CredentialSource interface {
	RetrieveCredential() (string, error)
}

// CredentialSourceIdentifier contains a key that is used to initially unmarshal a CredentialSource. It is used
// to determine which implementation of CredentialSource to unmarshal.
type CredentialSourceIdentifier struct {
	Type string `json:"type"`
}

// FileCredentialSource retrieves credentials from an existing file on disk.
type FileCredentialSource struct {
	Location string `json:"location"`
}

// RetrieveCredential reads the file designated by the FileCredentialSource and returns its content as a string.
func (fcs FileCredentialSource) RetrieveCredential() (string, error) {
	fileContent, fileReadErr := ioutil.ReadFile(fcs.Location)
	if fileReadErr != nil {
		return "", fmt.Errorf("failed to read credential at %v: %w", fcs.Location, fileReadErr)
	}

	return string(fileContent), nil
}

// EnvironmentCredentialSource retrieves credentials from a specified environment variable.
type EnvironmentCredentialSource struct {
	VariableName string `json:"envVar"`
}

// RetrieveCredential retrieves the content of the environment variable specified by the EnvironmentCredentialSource
// and returns it as a string.
func (ecs EnvironmentCredentialSource) RetrieveCredential() (string, error) {
	envVarContent, envVarExists := os.LookupEnv(ecs.VariableName)
	if !envVarExists {
		return "", fmt.Errorf("failed to read credential from %v, the variable did not exist", ecs.VariableName)
	}

	return envVarContent, nil
}

// InlineCredentialSource contains the content of a credential directly in JSON. This method is insecure, as it
// stores sensitive information directly inside the configuration file, which is presumably kept in version control.
type InlineCredentialSource struct {
	Value string `json:"value"`
}

// RetrieveCredential returns the literal credential contained in the InlineCredentialSource.
func (ics InlineCredentialSource) RetrieveCredential() (string, error) {
	return ics.Value, nil
}

// unmarshalCredentialSource parses the JSON representation of a CredentialSource into one of its implementations.
// It intentionally has the same format as json.Unmarshal().
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
