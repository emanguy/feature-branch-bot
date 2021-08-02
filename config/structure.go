package config

import (
	"encoding/json"
	"fmt"
)

type BotConfig struct {
	InteractiveProgress bool        `json:"interactiveProgress"`
	Servers             []VCSServer `json:"servers"`
}

type staticVCSServerFields struct {
	Type           string         `json:"type"`
	BaseURL        string         `json:"baseUrl"`
	SyncTag        string         `json:"syncTag"`
	SSHCreds       SSHCredentials `json:"sshCreds"`
	ProjectsToSync []VCSProject   `json:"projectsToSync"`
}

type VCSServer struct {
	staticVCSServerFields
	APIToken CredentialSource
}

func (svr *VCSServer) UnmarshalJSON(bytes []byte) error {
	type IntermediateVCSServer struct {
		staticVCSServerFields
		APIToken json.RawMessage `json:"apiToken"`
	}
	var intermediateRepresentation IntermediateVCSServer
	if parseErr := json.Unmarshal(bytes, &intermediateRepresentation); parseErr != nil {
		return fmt.Errorf("failed to read vcs server: %w", parseErr)
	}

	svr.staticVCSServerFields = intermediateRepresentation.staticVCSServerFields

	if tokenParseErr := unmarshalCredentialSource(intermediateRepresentation.APIToken, &svr.APIToken); tokenParseErr != nil {
		return fmt.Errorf("failed to read api token for server %v: %w", intermediateRepresentation.BaseURL, tokenParseErr)
	}
	return nil
}

type SSHCredentials struct {
	IsPresent  bool
	PublicKey  CredentialSource
	PrivateKey CredentialSource
}

func (pc *SSHCredentials) UnmarshalJSON(bytes []byte) error {
	type IntermediateSSHCredentials struct {
		PublicKey  json.RawMessage `json:"publicKey"`
		PrivateKey json.RawMessage `json:"privateKey"`
	}
	var intermediateRepresentation IntermediateSSHCredentials
	if initialParseErr := json.Unmarshal(bytes, &intermediateRepresentation); initialParseErr != nil {
		return fmt.Errorf("failed to read ssh credentials: %w", initialParseErr)
	}

	if pkParseErr := unmarshalCredentialSource(intermediateRepresentation.PublicKey, &pc.PublicKey); pkParseErr != nil {
		return fmt.Errorf("failed to read public key information from ssh credentials: %w", pkParseErr)
	}
	if privkParseErr := unmarshalCredentialSource(intermediateRepresentation.PrivateKey, &pc.PrivateKey); privkParseErr != nil {
		return fmt.Errorf("failed to read private key information from ssh credentials: %w", privkParseErr)
	}

	pc.IsPresent = true

	return nil
}

type VCSProject struct {
	PathWithNamespace string         `json:"pathWithNamespace"`
	SSHCloneURL       string         `json:"sshCloneUrl"`
	SyncTag           string         `json:"syncTag"`
	SSHCreds          SSHCredentials `json:"sshCreds"`
}
