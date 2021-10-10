package me.erittenhouse.featurebranchbot

import me.erittenhouse.featurebranchbot.config.VCSProject
import me.erittenhouse.featurebranchbot.git.*
import me.erittenhouse.featurebranchbot.gitlab.fetchOpenMergeRequestsWithLabel
import me.erittenhouse.featurebranchbot.util.printStackTraceIfEnabled
import org.eclipse.jgit.api.Git
import org.eclipse.jgit.api.MergeResult
import org.gitlab4j.api.GitLabApi
import org.gitlab4j.api.models.MergeRequest
import java.io.File

/**
 * Synchronizes labeled merge requests for the given repository. Returns true if any merge request sync succeeded
 * for the given repo.
 */
fun syncRepository(gitlabApi: GitLabApi, project: VCSProject, triggerLabel: String, credentials: Credentials, liveProgress: Boolean, stackTracesEnabled: Boolean): Boolean {
    // Fetch merge requests for project
    val mergeRequestsToSync = try {
        gitlabApi.fetchOpenMergeRequestsWithLabel(project.pathWithNamespace, triggerLabel)
    } catch (e: Exception) {
        throw Exception("Failed to list merge requests in ${project.pathWithNamespace} with the label $triggerLabel." +
                "Please check the provided credentials.", e)
    }

    // Clone, replacing slashes in name with underscores
    println("Cloning repository from project ${project.pathWithNamespace}.")
    val outputDir = project.pathWithNamespace.replace("/", "_")
    val clonedRepo = try {
        cloneRepository(project.sshCloneURL, outputDir, credentials, liveProgress)
    } catch (e: Exception) {
        throw Exception("Failed to clone repository for project: ${project.pathWithNamespace}.", e)
    }

    // Sync each detected merge request
    var anyMergeSucceeded = false
    for (mergeRequest in mergeRequestsToSync) {
        println("Synchronizing MR !${mergeRequest.iid} (${mergeRequest.sourceBranch} -> ${mergeRequest.targetBranch})...")
        try {
            syncMR(gitlabApi, clonedRepo, mergeRequest, credentials, liveProgress, stackTracesEnabled)
            anyMergeSucceeded = true
        } catch (e: Exception) {
            println("Error: failed to sync mr !${mergeRequest.iid}: ${e.message}")
            e.printStackTraceIfEnabled(stackTracesEnabled)
        }
    }

    // Delete local repo
    try {
        File(outputDir).deleteRecursively()
    } catch (e: Exception) {
        throw Exception("Failed to remove cloned project ${project.pathWithNamespace} at $outputDir.", e)
    }

    return anyMergeSucceeded
}

fun syncMR(
    gitlabApi: GitLabApi,
    repo: Git,
    mergeRequest: MergeRequest,
    credentials: Credentials,
    liveProgress: Boolean,
    stackTracesEnabled: Boolean,
) {
    // Switch to branch
    checkoutBranch(repo, mergeRequest.sourceBranch)

    // Perform Merge
    val mergeStatus = mergeBranchToCurrent(repo, mergeRequest.targetBranch)

    // If merge fails, make comment on MR and hard reset
    if (!mergeStatus.isSuccessful) {
        println("Error: merge of MR !${mergeRequest.iid} failed with status $mergeStatus. Hard resetting.")
        if (mergeStatus == MergeResult.MergeStatus.CONFLICTING) {
            try {
                gitlabApi.notesApi.createMergeRequestNote(
                    mergeRequest.projectId,
                    mergeRequest.iid,
                    ":warning:  Error: automatic merge failed due to merge conflict. Please merge manually.",
                )
            } catch(e: Exception) {
                println("Error: failed to alert GitLab users of merge conflict: ${e.message}")
                e.printStackTraceIfEnabled(stackTracesEnabled)
            }
        }

        hardReset(repo)
        throw Exception("Merge could not be completed due to conflict.")
    }

    println("Success: merge for MR !${mergeRequest.iid} completed with status $mergeStatus. Pushing result.")

    // Only push if there was actually anything new to merge
    if (mergeStatus != MergeResult.MergeStatus.ALREADY_UP_TO_DATE) {
        pushChanges(repo, credentials, liveProgress)
    }
}
