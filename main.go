package main

import (
	"feature-branch-bot/config"
	"feature-branch-bot/gitrepo"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"
)

func main() {
	configFileName := "./bot-config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}

	botConfiguration, readErr := config.ParseConfig(configFileName)
	if readErr != nil {
		fmt.Println("Bot failed, could not read configuration.")
		fmt.Println(readErr)
	}

	fmt.Println("Configuration read successfully. Starting branch sync on tagged merge/pull requests.")

	for _, vcsServer := range botConfiguration.Servers {
		fmt.Println("Synchronizing projects on server:", vcsServer.BaseURL)

		serverAPIToken, tokenFetchErr := vcsServer.APIToken.RetrieveCredential()
		if tokenFetchErr != nil {
			fmt.Println("Failed to read API token:", tokenFetchErr)
		}
		glClient, clientCreateErr := gitlab.NewClient(serverAPIToken, gitlab.WithBaseURL(vcsServer.BaseURL))
		if clientCreateErr != nil {
			fmt.Println("Failed to connect to GitLab:", clientCreateErr)
			continue
		}

		for _, vcsProject := range vcsServer.ProjectsToSync {
			fmt.Println("Now synchronizing project:", vcsProject.PathWithNamespace)

			credsForProject, credsSelectErr := config.DetermineSSHCredentialsForProject(vcsServer, vcsProject)
			if credsSelectErr != nil {
				fmt.Println("Couldn't determine appropriate SSH credentials: ", credsSelectErr)
				continue
			}
			sshPublicKey, pubkeyFetchErr := credsForProject.PublicKey.RetrieveCredential()
			if pubkeyFetchErr != nil {
				fmt.Println("Failed to read public key:", pubkeyFetchErr)
				continue
			}
			sshPrivateKey, privkeyFetchErr := credsForProject.PrivateKey.RetrieveCredential()
			if privkeyFetchErr != nil {
				fmt.Println("Failed to read private key:", privkeyFetchErr)
				continue
			}

			sshCredentials := gitrepo.SSHCredentials{
				SSHPubKey:  sshPublicKey,
				SSHPrivKey: sshPrivateKey,
			}

			keepUpToDateTag, tagCheckErr := config.DetermineSyncTagForProject(vcsServer.SyncTag, vcsProject.SyncTag)
			if tagCheckErr != nil {
				fmt.Println("Could not determine the tag to check for on merge/pull requests:", tagCheckErr)
				continue
			}

			repoSyncErr := SyncRepository(glClient, vcsProject.PathWithNamespace, keepUpToDateTag, vcsProject.SSHCloneURL, sshCredentials, botConfiguration.InteractiveProgress)
			if repoSyncErr != nil {
				fmt.Println("Failed to sync the requested repository: ", repoSyncErr)
			} else {
				fmt.Println("Repo sync successful.")
			}
		}
	}
}
