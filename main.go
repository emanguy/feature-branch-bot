package main

import (
	"feature-branch-bot/gitrepo"
	"fmt"
	"github.com/xanzy/go-gitlab"
)

const sshPublicKey = `FILL ME`
const sshPrivateKey = `FILL ME`
const repoSSHURL = `FILL ME`
const projectPathWithNamespace = `FILL ME`
const serverBaseURL = `FILL ME`
const gitlabAPIToken = `FILL ME`
const keepUpToDateTag = `FILL ME`

func main() {
	fmt.Println("Hello world! Cloning portal.")

	sshCredentials := gitrepo.SSHCredentials{
		SSHPubKey:  sshPublicKey,
		SSHPrivKey: sshPrivateKey,
	}

	glClient, clientCreateErr := gitlab.NewClient(gitlabAPIToken, gitlab.WithBaseURL(serverBaseURL))
	if clientCreateErr != nil {
		fmt.Println("Failed to connect to GitLab:", clientCreateErr)
		return
	}

	repoSyncErr := SyncRepository(glClient, projectPathWithNamespace, keepUpToDateTag, repoSSHURL, sshCredentials, false)
	if repoSyncErr != nil {
		fmt.Println("Failed to sync the requested repository: ", repoSyncErr)
	} else {
		fmt.Println("Repo sync successful.")
	}
}
