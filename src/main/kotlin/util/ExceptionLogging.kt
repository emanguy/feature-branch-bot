package me.erittenhouse.featurebranchbot.util

import me.erittenhouse.featurebranchbot.config.Configuration

fun Exception.printStackTraceIfEnabled(config: Configuration) {
    if (config.showErrorStacktraces) this.printStackTrace()
}