package main

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
)

const sshPublicKey = `FILL ME`
const sshPrivateKey = `FILL ME`
const repoSSHURL = `FILL ME`

func main() {
	fmt.Println("Hello world! Cloning portal.")

	var alreadyNewlined bool
	cloneOpts := git.CloneOptions{
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafe,
			ProgressCallback: func(_ string, completed, total uint) git.ErrorCode {
				fmt.Printf("\rChecking out repository: %v/%v complete", completed, total)

				return git.ErrorCodeOK
			},
		},
		FetchOptions: &git.FetchOptions{
			RemoteCallbacks: git.RemoteCallbacks{
				TransferProgressCallback: func(progress git.TransferProgress) git.ErrorCode {
					if !alreadyNewlined {
						fmt.Printf("\rDownloading repository: %v/%v complete", progress.ReceivedObjects, progress.TotalObjects)
						if progress.ReceivedObjects == progress.TotalObjects {
							fmt.Println()
							alreadyNewlined = true
						}
					}
					return git.ErrorCodeOK
				},
				CredentialsCallback: func(string, string, git.CredentialType) (*git.Credential, error) {
					return git.NewCredentialSSHKeyFromMemory("git", sshPublicKey, sshPrivateKey, "")
				},
				CertificateCheckCallback: func(*git.Certificate, bool, string) git.ErrorCode {
					return git.ErrorCodeOK
				},
			},
		},
		CheckoutBranch: "master",
	}
	repo, err := git.Clone(repoSSHURL, "./cloned-repo/", &cloneOpts)
	if err != nil {
		fmt.Printf("Failed to clone.", err)
		return
	}
	defer repo.Free()

	// Need to newline because the
	fmt.Println()

	branchIterator, iterCreateErr := repo.NewBranchIterator(git.BranchAll)
	if iterCreateErr != nil {
		fmt.Println("Failed to create branch iterator. ", iterCreateErr)
		return
	}
	defer branchIterator.Free()

	iterationErr := branchIterator.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		branchName, nameReadErr := branch.Name()
		if nameReadErr != nil {
			fmt.Println("Couldn't read branch name.")
			return nil
		}

		fmt.Println("Found branch: ", branchName)
		return nil
	})
	if iterationErr != nil {
		fmt.Println("Branch iteration failed.", iterationErr)
		return
	}
}
