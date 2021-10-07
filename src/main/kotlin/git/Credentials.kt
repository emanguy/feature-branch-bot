package me.erittenhouse.featurebranchbot.git

import com.jcraft.jsch.JSch
import org.eclipse.jgit.transport.*
import org.eclipse.jgit.util.FS

data class Credentials(val sshPublicKey: String, val sshPrivateKey: String)

class SSHSessionFactory(private val credentials: Credentials) : JschConfigSessionFactory() {
    override fun createDefaultJSch(fs: FS?): JSch {
        val defaultJsch = super.createDefaultJSch(fs)
        // NOTE: this only works for PEM-encoded keypairs (ssh-keygen -m pem
        defaultJsch.addIdentity(
            "userProvidedKey",
            credentials.sshPrivateKey.toByteArray(),
            credentials.sshPublicKey.toByteArray(),
            ByteArray(0),
        )
        return defaultJsch
    }
}