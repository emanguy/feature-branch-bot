package me.erittenhouse.featurebranchbot.git

import org.eclipse.jgit.api.Git
import org.eclipse.jgit.lib.ProgressMonitor
import org.eclipse.jgit.transport.CredentialItem
import org.eclipse.jgit.transport.CredentialsProvider
import org.eclipse.jgit.transport.SshTransport
import org.eclipse.jgit.transport.URIish
import java.io.File

fun cloneRepository(cloneURL: String, cloneDir: String, credentials: Credentials, liveProgress: Boolean): Git {
    return Git.cloneRepository()
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
        }
        .call()
}