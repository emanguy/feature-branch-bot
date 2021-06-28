package gitrepo

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
)

var stdCheckoutOptions = &git.CheckoutOptions{
	Strategy: git.CheckoutSafe,
}

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

// SwitchToBranch looks up a branch and switches to it, returning the commit the branch points
// to for easy git reset or branch operations if necessary
func SwitchToBranch(repo *git.Repository, branch string) (*git.Commit, error) {
	foundBranch, lookupErr := repo.LookupBranch(branch, git.BranchAll)
	if lookupErr != nil {
		return nil, fmt.Errorf("failed to find branch %v: %w", branch, lookupErr)
	}
	if foundBranch == nil {
		return nil, fmt.Errorf("no branch in the repo matches the name %v", branch)
	}
	defer foundBranch.Free()

	branchCommit, commitFindErr := repo.LookupCommit(foundBranch.Target())
	if commitFindErr != nil {
		return nil, fmt.Errorf("couldn't find commit for branch %v: %w", branch, branchCommit)
	}
	commitTree, treeLookupErr := repo.LookupTree(branchCommit.TreeId())
	if treeLookupErr != nil {
		return nil, fmt.Errorf("failed to retrieve file info for branch %v: %w", branch, treeLookupErr)
	}
	defer commitTree.Free()

	checkoutErr := repo.CheckoutTree(commitTree, stdCheckoutOptions)
	if checkoutErr != nil {
		branchCommit.Free()
		return nil, fmt.Errorf("failed to check out code from resolved branch %v's tree: %w", branch, checkoutErr)
	}
	headSetErr := repo.SetHead("refs/heads/" + branch)
	if headSetErr != nil {
		branchCommit.Free()
		return nil, fmt.Errorf("failed to set repo to newly checked out branch %v: %w", branch, headSetErr)
	}

	return branchCommit, nil
}

// MergeBranchToCurrent merges the specified branch to the current HEAD
func MergeBranchToCurrent(repo *git.Repository, branchToMerge string) error {
	branchResult, branchLookupErr := repo.LookupBranch(branchToMerge, git.BranchAll)
	if branchLookupErr != nil {
		return fmt.Errorf("failed to find branch to merge into current: %v: %w", branchToMerge, branchLookupErr)
	}
	if branchResult == nil {
		return fmt.Errorf("branch to merge not found: %v", branchToMerge)
	}
	defer branchResult.Free()

	annotatedBranchToMergeCommit, commitLookupErr := repo.LookupAnnotatedCommit(branchResult.Target())
	if commitLookupErr != nil {
		return fmt.Errorf("could not find commit referenced by branch to merge: %v: %w", branchToMerge, commitLookupErr)
	}
	defer annotatedBranchToMergeCommit.Free()

	mergeOpts := &git.MergeOptions{
		TreeFlags: git.MergeTreeFailOnConflict | git.MergeTreeFindRenames,
	}
	mergeErr := repo.Merge([]*git.AnnotatedCommit{annotatedBranchToMergeCommit}, mergeOpts, stdCheckoutOptions)
	if mergeErr != nil {
		return fmt.Errorf("merge of branch %v failed: %w", branchToMerge, mergeErr)
	}

	return nil
}

// ResetRepo does a hard reset to the specified commit
func ResetRepo(repo *git.Repository, branchCommit *git.Commit) error {
	resetErr := repo.ResetToCommit(branchCommit, git.ResetHard, stdCheckoutOptions)
	if resetErr != nil {
		return fmt.Errorf("failed to hard reset: %w", resetErr)
	}

	return nil
}

// PushChanges pushes currentBranch's changes to the "origin" remote, which should be the default
func PushChanges(repo *git.Repository, currentBranch string) error {
	originRemote, remoteLookupErr := repo.Remotes.Lookup("origin")
	if remoteLookupErr != nil {
		return fmt.Errorf("could not get \"origin\" remote. failed to push: %w", remoteLookupErr)
	}
	defer originRemote.Free()

	pushErr := originRemote.Push([]string{"refs/heads/" + currentBranch}, &git.PushOptions{})
	if pushErr != nil {
		return fmt.Errorf("failed to push updated branch %v: %w", currentBranch, pushErr)
	}

	return nil
}
