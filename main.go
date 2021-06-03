package main

import (
	"feature-branch-bot/gitlab_tools"
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

	// repo, err := gitrepo.CloneRepository(repoSSHURL, "./cloned-repo/", gitrepo.SSHCredentials{
	// 	SSHPubKey:  sshPublicKey,
	// 	SSHPrivKey: sshPrivateKey,
	// }, true)
	// if err != nil {
	// 	fmt.Println("Repo clone failed:", err)
	// 	return
	// }
	// defer repo.Free()

	glClient, clientCreateErr := gitlab.NewClient(gitlabAPIToken, gitlab.WithBaseURL(serverBaseURL))
	if clientCreateErr != nil {
		fmt.Println("Failed to connect to GitLab:", clientCreateErr)
		return
	}

	mergeRequests, fetchErr := gitlab_tools.FetchMergeRequestsWithTag(glClient, projectPathWithNamespace, keepUpToDateTag)
	if fetchErr != nil {
		fmt.Printf("Failed to get merge requests with the tag: %v. Error: %v\n", keepUpToDateTag, fetchErr)
		return
	}

	fmt.Println(len(mergeRequests))
}
