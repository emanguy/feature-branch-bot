package me.erittenhouse.featurebranchbot.config

fun <T> determineProjectValue(serverValue: T?, projectValue: T?): T? {
    if (projectValue != null) return projectValue
    return serverValue
}
