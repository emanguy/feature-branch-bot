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
    var workDone: Int = 0
    var totalProgress: Int = 0

    private fun printProgress(currentProgress: Int, isUpdate: Boolean = true) {
        if (currentTaskTitle.isEmpty()) return
        if (isUpdate) {
            print("\r")
        } else {
            println()
        }

        print("$currentTaskTitle: $currentProgress")
        if (totalProgress > 0) {
            print("/$totalProgress")
        }
    }

    override fun start(totalTasks: Int) {
        // Do nothing
    }

    override fun beginTask(title: String?, totalWork: Int) {
        currentTaskTitle = title ?: "Unnamed Task"
        workDone = 0
        totalProgress = totalWork
        printProgress(0, isUpdate = false)
    }

    override fun update(completed: Int) {
        workDone += completed
        printProgress(workDone)
    }

    override fun endTask() {}

    override fun isCancelled(): Boolean = false
}