package me.erittenhouse.featurebranchbot.git

import org.eclipse.jgit.api.Git
import org.eclipse.jgit.transport.CredentialItem
import org.eclipse.jgit.transport.CredentialsProvider
import org.eclipse.jgit.transport.SshTransport
import org.eclipse.jgit.transport.URIish

fun cloneRepository(cloneURL: String, credentials: Credentials) {
    Git.cloneRepository()
        .setURI(cloneURL)
        .setTransportConfigCallback { transport ->
            (transport as SshTransport).sshSessionFactory
        }
}