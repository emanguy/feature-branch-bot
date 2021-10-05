package me.erittenhouse.featurebranchbot

import me.erittenhouse.featurebranchbot.config.VCSProject
import me.erittenhouse.featurebranchbot.config.VCSServer
import me.erittenhouse.featurebranchbot.git.Credentials

fun SyncRepository(server: VCSServer, project: VCSProject, credentials: Credentials, syncTag: String, mainBranchName: String, liveProgress: Boolean) {
    // Fetch merge requests for project

    // Clone, replacing slashes in name with underscores

    // Sync each detected merge request

    // Delete local repo
}