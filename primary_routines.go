package main

import (
	"feature-branch-bot/gitlab_tools"
	"feature-branch-bot/gitrepo"
	"fmt"
	git "github.com/libgit2/git2go/v31"
	"github.com/xanzy/go-gitlab"
	"os"
	"strings"
)

func SyncRepository(glClient *gitlab.Client, projectPathWithNamespace, triggerTag, cloneURL string, gitCreds gitrepo.SSHCredentials) error {
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
	clonedRepo, repoCloneErr := gitrepo.CloneRepository(cloneURL, outputDir, gitCreds, false)
	if repoCloneErr != nil {
		return repoCloneErr
	}
	defer clonedRepo.Free()

	for _, mergeRequestToSync := range mrsToSync {
		fmt.Printf("Syncing merge request !%v with base branch...\n", mergeRequestToSync.IID)
		mergeErr := SyncMR(glClient, *mergeRequestToSync, clonedRepo)
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

func SyncMR(glClient *gitlab.Client, mergeRequest gitlab.MergeRequest, repo *git.Repository) error {
	panic("implement me")
}
