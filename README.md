# Feature Branch Bot

This bot is designed to keep long-running feature branches up-to-date with their target branches.
The bot searches for open merge requests with a label specified by the bot's configuration file,
clones the referenced branch, performs a merge from the destination branch into the source branch,
and pushes the result if it's successful. In the event of a merge conflict, the bot will post a comment
on the merge request saying that a conflict occurred, and will ask humans to resolve the conflict so that it can continue
to keep the branch up-to-date.

The configuration file is very flexible, and you can read about its format [here](./README_CONFIG.md).

## State of the bot

The bot has been fully re-written from the ground up in Kotlin, as the Go libraries for Git have been found to be
relatively immature in comparison. The old Go implementation can be found on the `old-go-implementation` branch for
anyone interested.

The bot can be run from source using the included `gradlew` and `gradlew.bat` files with the gradle command
`./gradlew run`. Note that Java 11 is required to run the feature branch bot.

The following tasks still need to get done before this project is fully released:
 - [ ] Update the Gradle file to build a JAR file so the code can be run without the source code or build chain
 - [ ] Add a Dockerfile so the bot can be run as part of CI or automation without any dependency requirements
 - [ ] Implement a CI pipeline to build the Dockerfile and JAR file automatically for future releases

In the future, I plan to add GitHub support but that will have to happen after the other "initial release" tasks
are completed. Pull requests and contributions are welcome!
