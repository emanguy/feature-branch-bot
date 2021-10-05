package me.erittenhouse.featurebranchbot.config

import kotlinx.serialization.modules.SerializersModule
import kotlinx.serialization.modules.polymorphic

val configSerializationModule = SerializersModule {
    polymorphic(CredentialSource::class) {
        subclass(FileCredentialSource::class, FileCredentialSource.serializer())
        subclass(EnvVarCredentialSource::class, EnvVarCredentialSource.serializer())
        subclass(InlineCredentialSource::class, InlineCredentialSource.serializer())
    }
}