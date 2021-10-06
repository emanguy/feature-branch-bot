package me.erittenhouse.featurebranchbot.gitlab

import org.gitlab4j.api.GitLabApi
import org.gitlab4j.api.models.MergeRequest
import org.gitlab4j.api.models.MergeRequestFilter

fun GitLabApi.fetchOpenMergeRequestsWithLabel(repositoryNamespaceWithPath: String, tag: String): List<MergeRequest> {
    val project = this.projectApi.getProject(repositoryNamespaceWithPath)
    val matchingMergeRequests = this.mergeRequestApi.getMergeRequests(MergeRequestFilter().apply {
        projectId = project.id
        labels = listOf(tag)
    })

    return matchingMergeRequests.filter { it.closedAt == null }
}