package me.erittenhouse.featurebranchbot.git

import org.eclipse.jgit.lib.ProgressMonitor

class NoOpMonitor : ProgressMonitor {
    override fun start(totalTasks: Int) {}
    override fun beginTask(title: String?, totalWork: Int) {}
    override fun update(completed: Int) {}
    override fun endTask() {}
    override fun isCancelled(): Boolean = false

}

class BotProgressMonitor : ProgressMonitor {
    var currentTaskTitle: String = ""
    var totalProgress: Int = 0

    override fun start(totalTasks: Int) {
        // Do nothing
    }

    override fun beginTask(title: String?, totalWork: Int) {
        currentTaskTitle = title ?: "Unnamed Task"
        totalProgress = totalWork
        println("$currentTaskTitle: 0/$totalWork")
    }

    override fun update(completed: Int) {
        println("\r$currentTaskTitle: $completed/$totalProgress")
    }

    override fun endTask() {
        // Move to a new line
        println()
    }

    override fun isCancelled(): Boolean = false
}