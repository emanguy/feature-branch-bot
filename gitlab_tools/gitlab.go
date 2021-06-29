package gitlab_tools

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

func FetchMergeRequestsWithTag(client *gitlab.Client, repositoryPathWithNamespace, tag string) ([]*gitlab.MergeRequest, error) {
	mrFetchOptions := &gitlab.ListProjectMergeRequestsOptions{
		Labels: gitlab.Labels{tag},
	}
	mergeRequests, _, mrFetchErr := client.MergeRequests.ListProjectMergeRequests(repositoryPathWithNamespace, mrFetchOptions)
	if mrFetchErr != nil {
		return nil, fmt.Errorf("failed to list merge requests: %w", mrFetchErr)
	}

	// Filter out closed MRs
	nonClosedMergeRequests := make([]*gitlab.MergeRequest, 0, len(mergeRequests))
	for _, mergeRequest := range mergeRequests {
		if mergeRequest.ClosedAt == nil {
			nonClosedMergeRequests = append(nonClosedMergeRequests, mergeRequest)
		}
	}
	// TODO verify this filtering works

	return nonClosedMergeRequests, nil
}
