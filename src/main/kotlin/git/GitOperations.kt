package me.erittenhouse.featurebranchbot.git

import org.eclipse.jgit.api.CreateBranchCommand
import org.eclipse.jgit.api.Git
import org.eclipse.jgit.api.MergeResult
import org.eclipse.jgit.api.ResetCommand
import org.eclipse.jgit.lib.ProgressMonitor
import org.eclipse.jgit.transport.CredentialItem
import org.eclipse.jgit.transport.CredentialsProvider
import org.eclipse.jgit.transport.SshTransport
import org.eclipse.jgit.transport.URIish
import java.io.File

fun cloneRepository(cloneURL: String, cloneDir: String, credentials: Credentials, liveProgress: Boolean): Git {
    val repo = Git.cloneRepository()
        .setURI(cloneURL)
        .setDirectory(File(cloneDir))
        .setTransportConfigCallback { transport ->
            val convertedTransport = transport as? SshTransport ?:
                throw Error("Did not receive SSH transport. Did you provide an SSH URL?")
            convertedTransport.sshSessionFactory = SSHSessionFactory(credentials)
        }.let { cloneCommand ->
            if (liveProgress) {
                cloneCommand.setProgressMonitor(BotProgressMonitor())
            } else {
                cloneCommand.setProgressMonitor(NoOpMonitor())
            }
        }.call()
    repo.repository.config.setString("user", null, "name", "Feature-Branch Bot")
    repo.repository.config.setString("user", null, "email", "noreply@featurebranchbot.net")

    return repo
}

fun checkoutBranch(repo: Git, branch: String) {
    repo.checkout()
        .setName(branch)
        .setStartPoint("origin/$branch")
        .setCreateBranch(true)
        .setUpstreamMode(CreateBranchCommand.SetupUpstreamMode.TRACK)
        .call()
}

fun mergeBranchToCurrent(repo: Git, branchToMerge: String): MergeResult.MergeStatus {
    val targetBranchCommit = repo.repository.resolve("origin/${branchToMerge}")
    val mergeResult = repo.merge()
        .include(targetBranchCommit)
        .setCommit(true)
        .setMessage("Automatic merge from $branchToMerge via Feature Branch Bot")
        .call()

    return mergeResult.mergeStatus
}

fun hardReset(repo: Git) {
    repo.reset()
        .setMode(ResetCommand.ResetType.HARD)
        .call()
}

fun pushChanges(repo: Git, credentials: Credentials, liveProgress: Boolean) {
    repo.push()
        .setTransportConfigCallback { transport ->
            val convertedTransport = transport as? SshTransport ?:
                throw Error("Did not receive SSH transport. Did you provide an SSH URL?")
            convertedTransport.sshSessionFactory = SSHSessionFactory(credentials)
        }.let { pushCommand ->
            if (liveProgress) {
                pushCommand.setProgressMonitor(BotProgressMonitor())
            } else {
                pushCommand.setProgressMonitor(NoOpMonitor())
            }
        }.call()
}
