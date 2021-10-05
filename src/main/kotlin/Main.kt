package me.erittenhouse.featurebranchbot

import kotlinx.serialization.decodeFromString
import me.erittenhouse.featurebranchbot.config.Configuration
import me.erittenhouse.featurebranchbot.config.determineProjectValue
import me.erittenhouse.featurebranchbot.gitlab.fetchOpenMergeRequestsWithTag
import me.erittenhouse.featurebranchbot.util.printStackTraceIfEnabled
import me.erittenhouse.featurebranchbot.util.serializer
import org.gitlab4j.api.GitLabApi
import java.io.File

fun main(args: Array<String>) {
    val configFileName = args.getOrNull(0) ?: "bot-config.json"
    val programConfig = try {
        serializer.decodeFromString<Configuration>(File(configFileName).readText())
    } catch(e: Exception) {
        println("Error: failed to read configuration file.")
        e.printStackTrace()
        return
    }

    for (server in programConfig.servers) {
        println("Syncing merge requests on server ${server.baseURL}...")
        val gitLabApi = try {
            GitLabApi(server.baseURL, server.apiToken.retrieveCredential())
        } catch (e: Exception) {
            println("Error: failed to connect to GitLab API for server ${server.baseURL}.")
            e.printStackTraceIfEnabled(programConfig)
            continue
        }

        for (project in server.projectsToSync) {
            println("Syncing merge requests for project ${project.pathWithNamespace}...")
            val syncTag = determineProjectValue(server.syncTag, project.syncTag)
            if (syncTag == null) {
                println("Error: failed to determine sync tag for project ${project.pathWithNamespace}. " +
                        "It must be provided for either the server or project or both.")
                continue
            }

            val mergeRequestsToSync = try {
                gitLabApi.fetchOpenMergeRequestsWithTag(project.pathWithNamespace, syncTag)
            } catch (e: Exception) {
                println("Error: failed to list merge requests in ${project.pathWithNamespace} with the label $syncTag." +
                        "Please check the provided credentials.")
                e.printStackTraceIfEnabled(programConfig)
                continue
            }

            mergeRequestsToSync.forEach(::println)
        }
    }
}