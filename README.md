# Feature Branch Bot

This bot is designed to keep long-running feature branches up-to-date with their target branches.
The bot searches for open merge requests with a label specified by the bot's configuration file,
clones the referenced branch, performs a merge from the destination branch into the source branch,
and pushes the result if it's successful. In the event of a merge conflict, the bot will post a comment
on the merge request saying that a conflict occurred, and will ask humans to resolve the conflict so that it can continue
to keep the branch up-to-date.

The configuration file is very flexible, and you can read about its format [here](./README_CONFIG.md).

## State of the bot

The bot is known to run fine on MacOS, but I plan to automate the building of the software and release a docker
image so the bot can easily be used in CI or other forms of automation. The strategy for building on Linux (which will
be the primary target) is still TBD at the moment. I also plan to add local development instructions to this README
at a future date.

If you want to take a shot at running locally on a Mac, just install `libgit2` and `pkg-config` from homebrew, then set up
a config file and run `go run .` in the root directory.
