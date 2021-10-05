package me.erittenhouse.featurebranchbot.config

import kotlinx.serialization.Serializable

@Serializable
data class Configuration(
    val servers: List<VCSServer>,
    val interactiveProgress: Boolean = false,
    val showErrorStacktraces: Boolean = false,
)

@Serializable
data class VCSServer(
    val type: String,
    val baseURL: String,
    val apiToken: CredentialSource,
    val syncTag: String? = null,
    val sshCreds: SSHCredentials? = null,
    val projectsToSync: List<VCSProject>,
)

@Serializable
data class VCSProject(
    val pathWithNamespace: String,
    val sshCloneURL: String,
    val syncTag: String? = null,
    val mainBranchName: String = "master",
    val sshCredentials: SSHCredentials? = null,
)

@Serializable
data class SSHCredentials(
    val publicKey: CredentialSource,
    val privateKey: CredentialSource,
)
