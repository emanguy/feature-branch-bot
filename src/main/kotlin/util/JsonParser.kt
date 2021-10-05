package me.erittenhouse.featurebranchbot.util

import kotlinx.serialization.json.Json
import me.erittenhouse.featurebranchbot.config.configSerializationModule

val serializer = Json {
    serializersModule = configSerializationModule
}
