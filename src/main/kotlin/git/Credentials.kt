package me.erittenhouse.featurebranchbot.git

import org.eclipse.jgit.transport.CredentialsProvider
import org.eclipse.jgit.transport.RemoteSession
import org.eclipse.jgit.transport.SshSessionFactory
import org.eclipse.jgit.transport.URIish
import org.eclipse.jgit.util.FS

data class Credentials(val sshPublicKey: String, val sshPrivateKey: String)

class SSHSessionFactory(private val credentials: Credentials) : SshSessionFactory() {
    override fun getSession(uri: URIish?, credentialsProvider: CredentialsProvider?, fs: FS?, tms: Int): RemoteSession {
        TODO("Not yet implemented")
    }

    override fun getType(): String {
        TODO("Not yet implemented")
    }

}