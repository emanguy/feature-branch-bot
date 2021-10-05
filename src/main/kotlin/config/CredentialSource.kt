package me.erittenhouse.featurebranchbot.config

import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import java.io.File
import kotlinx.serialization.Serializable

interface CredentialSource {
    fun retrieveCredential(): String
}

@Serializable
@SerialName("FILE")
data class FileCredentialSource(val location: String) : CredentialSource {
    override fun retrieveCredential(): String = File(location).readText()
}

@Serializable
@SerialName("ENVIRONMENT")
data class EnvVarCredentialSource(val envVar: String) : CredentialSource {
    override fun retrieveCredential(): String = System.getenv(envVar)
}

@Serializable
@SerialName("INLINE")
data class InlineCredentialSource(val value: String) : CredentialSource {
    override fun retrieveCredential(): String = value
}

