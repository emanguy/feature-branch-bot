package main

import (
	"feature-branch-bot/gitlab_tools"
	"feature-branch-bot/gitrepo"
	"fmt"
	git "github.com/libgit2/git2go/v31"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/multierr"
	"os"
	"strings"
)

func SyncRepository(glClient *gitlab.Client, projectPathWithNamespace, triggerTag, cloneURL string, gitCreds gitrepo.SSHCredentials, liveProgress bool) error {
	fmt.Printf("Syncing repository %v...\n", projectPathWithNamespace)
	mrsToSync, mrFetchErr := gitlab_tools.FetchMergeRequestsWithTag(glClient, projectPathWithNamespace, triggerTag)
	if mrFetchErr != nil {
		return mrFetchErr
	}
	if len(mrsToSync) == 0 {
		fmt.Printf("No merge requests to sync in %v.\n", projectPathWithNamespace)
		return nil
	}

	fmt.Printf("%v: %v merge requests to sync. Cloning repo.\n", projectPathWithNamespace, len(mrsToSync))

	outputDir := strings.Replace(projectPathWithNamespace, "/", "_", -1)
	clonedRepo, repoCloneErr := gitrepo.CloneRepository(cloneURL, outputDir, gitCreds, liveProgress)
	if repoCloneErr != nil {
		return repoCloneErr
	}
	defer clonedRepo.Free()

	for _, mergeRequestToSync := range mrsToSync {
		fmt.Printf("Syncing merge request !%v with base branch...\n", mergeRequestToSync.IID)
		mergeErr := SyncMR(glClient, *mergeRequestToSync, clonedRepo, gitCreds, liveProgress)
		if mergeErr != nil {
			fmt.Printf("Sync failure for MR !%v: %v\n", mergeRequestToSync.IID, mergeErr)
		}
	}

	fmt.Printf("Cleaning up local files for %v...\n", projectPathWithNamespace)
	deleteRepoErr := os.RemoveAll(outputDir)
	if deleteRepoErr != nil {
		return fmt.Errorf("failed to delete local repo files after sync of %v: %w", projectPathWithNamespace, deleteRepoErr)
	}

	fmt.Printf("Feature branch sync complete for %v.\n", projectPathWithNamespace)
	return nil
}

func SyncMR(glClient *gitlab.Client, mergeRequest gitlab.MergeRequest, repo *git.Repository, gitCreds gitrepo.SSHCredentials, liveProgress bool) error {
	currentBranch := mergeRequest.SourceBranch
	targetBranch := mergeRequest.TargetBranch

	branchCommit, branchSwitchErr := gitrepo.SwitchToBranch(repo, currentBranch)
	if branchSwitchErr != nil {
		return branchSwitchErr
	}
	defer branchCommit.Free()

	fmt.Printf("Now merging branch %v into %v for MR !%v...\n", targetBranch, currentBranch, mergeRequest.IID)
	mergeErr := gitrepo.MergeBranches(repo, targetBranch)
	if mergeErr != nil {
		fmt.Printf("Merge failed, possible conflict. Notifying authors for MR !%v and hard resetting...\n", mergeRequest.IID)
		comment := ":warning:  Error: automatic merge failed due to merge conflict. Please merge manually."
		_, _, commentErr := glClient.Notes.CreateMergeRequestNote(mergeRequest.ProjectID, mergeRequest.IID, &gitlab.CreateMergeRequestNoteOptions{
			Body: &comment,
		})
		// If we fail to make the comment, just combine the comment error with the merge error
		if commentErr != nil {
			combinedErr := multierr.Combine(commentErr, mergeErr)
			mergeErr = fmt.Errorf("merge failed, failed to notify users of conflict: %w", combinedErr)
		}

		// Now hard reset
		resetErr := gitrepo.ResetRepo(repo, branchCommit)
		if resetErr != nil {
			mergeErr = multierr.Append(mergeErr, resetErr)
		}

		return mergeErr
	}

	fmt.Printf("Pushing updated branch for MR !%v...\n")
	pushErr := gitrepo.PushChanges(repo, currentBranch, gitCreds, liveProgress)
	if pushErr != nil {
		return pushErr
	}

	return nil
}
