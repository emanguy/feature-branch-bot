package config

import (
	"encoding/json"
	"fmt"
)

type BotConfig struct {
	InteractiveProgress bool
	Servers             []VCSServer
}

type VCSServer struct {
	Type     string
	BaseURL  string
	APIToken string
	SyncTag  string
	SSHCreds ParsableCredentials
}

type ParsableCredentials struct {
	PublicKey  CredentialSource
	PrivateKey CredentialSource
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

func (pc *ParsableCredentials) UnmarshalJSON(bytes []byte) error {
	type IntermediateParsableCredentials struct {
		PublicKey  json.RawMessage `json:"publicKey"`
		PrivateKey json.RawMessage `json:"privateKey"`
	}
	var intermediateRepresentation IntermediateParsableCredentials
	if initialParseErr := json.Unmarshal(bytes, &intermediateRepresentation); initialParseErr != nil {
		return fmt.Errorf("failed to read ssh credentials: %w", initialParseErr)
	}

	if pkParseErr := unmarshalCredentialSource(intermediateRepresentation.PublicKey, &pc.PublicKey); pkParseErr != nil {
		return fmt.Errorf("failed to read public key information from ssh credentials: %w", pkParseErr)
	}
	if privkParseErr := unmarshalCredentialSource(intermediateRepresentation.PrivateKey, &pc.PrivateKey); privkParseErr != nil {
		return fmt.Errorf("failed to read private key information from ssh credentials: %w", privkParseErr)
	}

	return nil
}
