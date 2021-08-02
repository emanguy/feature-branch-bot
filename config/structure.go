package config

import (
	"encoding/json"
	"fmt"
)

// BotConfig defines the JSON format for configuring the bot.
type BotConfig struct {
	InteractiveProgress bool        `json:"interactiveProgress"`
	Servers             []VCSServer `json:"servers"`
}

// staticVCSServerFields contains non-dynamic fields for a VCSServer.
// This allows these fields to be shared between VCSServer and the intermediate
// representation used when reifying interface fields.
type staticVCSServerFields struct {
	Type           string         `json:"type"`
	BaseURL        string         `json:"baseUrl"`
	SyncTag        string         `json:"syncTag"`
	SSHCreds       SSHCredentials `json:"sshCreds"`
	ProjectsToSync []VCSProject   `json:"projectsToSync"`
}

// VCSServer represents a single version control server containing projects whose merge requests
// need synchronization.
type VCSServer struct {
	staticVCSServerFields
	APIToken CredentialSource
}

// UnmarshalJSON deserializes a VCSServer from JSON, reifying its interface fields based on the JSON content.
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

// SSHCredentials contains the public and private key necessary for cloning a repository. The IsPresent field
// is set to true if the credentials were actually provided.
type SSHCredentials struct {
	IsPresent  bool
	PublicKey  CredentialSource
	PrivateKey CredentialSource
}

// UnmarshalJSON deserializes SSHCredentials from JSON, using the content to reify interface fields.
func (creds *SSHCredentials) UnmarshalJSON(bytes []byte) error {
	type IntermediateSSHCredentials struct {
		PublicKey  json.RawMessage `json:"publicKey"`
		PrivateKey json.RawMessage `json:"privateKey"`
	}
	var intermediateRepresentation IntermediateSSHCredentials
	if initialParseErr := json.Unmarshal(bytes, &intermediateRepresentation); initialParseErr != nil {
		return fmt.Errorf("failed to read ssh credentials: %w", initialParseErr)
	}

	if pkParseErr := unmarshalCredentialSource(intermediateRepresentation.PublicKey, &creds.PublicKey); pkParseErr != nil {
		return fmt.Errorf("failed to read public key information from ssh credentials: %w", pkParseErr)
	}
	if privkParseErr := unmarshalCredentialSource(intermediateRepresentation.PrivateKey, &creds.PrivateKey); privkParseErr != nil {
		return fmt.Errorf("failed to read private key information from ssh credentials: %w", privkParseErr)
	}

	creds.IsPresent = true

	return nil
}

// VCSProject represents a project on a server which should be searched for merge requests in need of synchronization.
type VCSProject struct {
	PathWithNamespace string         `json:"pathWithNamespace"`
	SSHCloneURL       string         `json:"sshCloneUrl"`
	SyncTag           string         `json:"syncTag"`
	MainBranchName    string         `json:"mainBranchName"`
	SSHCreds          SSHCredentials `json:"sshCreds"`
}

// MainBranch returns the main branch name for a VCSProject or returns "master" if no value was provided.
func (vcp VCSProject) MainBranch() string {
	if len(vcp.MainBranchName) == 0 {
		return "master"
	}
	return vcp.MainBranchName
}
