package gitrepo

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
)

// CloneRepository clones a git repository from url into destinationDir via the passed SSH credentials. This function will optionally show live
// progress if interactiveProgress is true.
func CloneRepository(url, destinationDir string, credentials SSHCredentials, interactiveProgress bool) (*git.Repository, error) {
	var alreadyNewlined bool
	cloneOpts := git.CloneOptions{
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafe,
			ProgressCallback: func(_ string, completed, total uint) git.ErrorCode {
				if interactiveProgress {
					fmt.Printf("\rChecking out repository: %v/%v complete", completed, total)
				}

				return git.ErrorCodeOK
			},
		},
		FetchOptions: &git.FetchOptions{
			RemoteCallbacks: git.RemoteCallbacks{
				TransferProgressCallback: func(progress git.TransferProgress) git.ErrorCode {
					if !alreadyNewlined && interactiveProgress {
						fmt.Printf("\rDownloading repository: %v/%v complete", progress.ReceivedObjects, progress.TotalObjects)
						if progress.ReceivedObjects == progress.TotalObjects {
							fmt.Println()
							alreadyNewlined = true
						}
					}
					return git.ErrorCodeOK
				},
				CredentialsCallback: func(string, string, git.CredentialType) (*git.Credential, error) {
					return git.NewCredentialSSHKeyFromMemory("git", credentials.SSHPubKey, credentials.SSHPrivKey, "")
				},
				CertificateCheckCallback: func(*git.Certificate, bool, string) git.ErrorCode {
					return git.ErrorCodeOK
				},
			},
		},
		CheckoutBranch: "master",
	}

	repo, err := git.Clone(url, destinationDir, &cloneOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repo: %w", err)
	}

	if interactiveProgress {
		// Need to newline because the interactive status will just be repeating on the same line
		fmt.Println()
	}

	return repo, nil
}
