# Feature Branch Bot

This bot is designed to keep long-running feature branches up-to-date with their target branches.
The bot searches for open merge requests with a label specified by the bot's configuration file,
clones the referenced branch, performs a merge from the destination branch into the source branch,
and pushes the result if it's successful. In the event of a merge conflict, the bot will post a comment
on the merge request saying that a conflict occurred, and will ask humans to resolve the conflict so that it can continue
to keep the branch up-to-date.

## Using the bot

The bot requires a configuration file to run. It is very flexible, and you can read 
about its format [here](./README_CONFIG.md).

The bot can be run from the prepackaged JAR file included with every release. Just
check the "releases" page to download it. Note that Java 11 is required to run the feature branch bot.

Alternatively, if you plan to use the bot in CI, you can run it via the Docker image
published in the "packages" section.

## Old Go implementation

The bot has been fully re-written from the ground up in Kotlin, as the Go libraries for Git have been found to be
relatively immature in comparison. The old Go implementation can be found on the `old-go-implementation` branch for
anyone interested.
