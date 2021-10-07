package me.erittenhouse.featurebranchbot

import me.erittenhouse.featurebranchbot.config.VCSProject
import me.erittenhouse.featurebranchbot.git.Credentials
import me.erittenhouse.featurebranchbot.git.cloneRepository
import me.erittenhouse.featurebranchbot.gitlab.fetchOpenMergeRequestsWithLabel
import org.eclipse.jgit.api.Git
import org.gitlab4j.api.GitLabApi
import org.gitlab4j.api.models.MergeRequest
import java.io.File

fun syncRepository(gitlabApi: GitLabApi, project: VCSProject, triggerLabel: String, credentials: Credentials, liveProgress: Boolean) {
    // Fetch merge requests for project
    val mergeRequestsToSync = try {
        gitlabApi.fetchOpenMergeRequestsWithLabel(project.pathWithNamespace, triggerLabel)
    } catch (e: Exception) {
        throw Error("Failed to list merge requests in ${project.pathWithNamespace} with the label $triggerLabel." +
                "Please check the provided credentials.", e)
    }

    // Clone, replacing slashes in name with underscores
    println("Cloning repository from project ${project.pathWithNamespace}.")
    val outputDir = project.pathWithNamespace.replace("/", "_")
    val clonedRepo = try {
        cloneRepository(project.sshCloneURL, outputDir, credentials, liveProgress)
    } catch (e: Exception) {
        throw Error("Failed to clone repository for project: ${project.pathWithNamespace}.", e)
    }

    // Sync each detected merge request

    // Delete local repo
    try {
        File(outputDir).deleteRecursively()
    } catch (e: Exception) {
        throw Error("Failed to remove cloned project ${project.pathWithNamespace} at $outputDir.", e)
    }
}

fun syncMR(gitlabApi: GitLabApi, repo: Git, mergeRequest: MergeRequest, credentials: Credentials, liveProgress: Boolean) {
    // Switch to branch

    // Perform Merge

    // If merge fails, make comment on MR and hard reset
}