# Feature Branch Bot Configuration

See the [sample configuration](./sample-config.json) for a starting point.

The bot looks for `bot-config.json` as the default configuration file name, though
you can provide an alternate path for the configuration as the first argument of the executable.

If you're using the feature branch bot with docker, you can also provide the bot config inline via the environment
variable `BOT_CONFIG_JSON`.

## Top Level
* `interactiveProgress` [Boolean, optional] - True if live progress should be rendered while cloning a repo. Good for local runs, bad for automation. Defaults to false.
* `showErrorStacktraces` [Boolean, optional] - True if stacktraces should be shown when errors occur during bot operations. 
   By default, the messages are only shown. Recommended for debugging or issue reporting purposes. Defaults to false.
* `servers` [List of VCS Server] - The list of VCS servers to synchronize feature branches for.

## VCS Server
* `type` [String] - The type of VCS server we're syncing branches for. Helps determine the API to use for looking up Pull/Merge Requests. 
   Currently, the only supported value is `GITLAB`, though more will be added in the future.
* `baseURL` [String] - The base URL of the VCS server. Used to determine where to point the API so private servers are supported.
* `apiToken` [Credential Source] - The API token to use for accessing the VCS server's API.
  * For GitLab servers, this token should only require the `api` scope. The user this token was generated for should have access to read merge requests
    and push commits on the targeted GitLab server for the projects listed under `projectsToSync`.
* `syncLabel` [String, optional] - The default label to look for on merge requests for the bot to know which ones to sync. 
  This will be used for projects which do not specify a label. Not necessary if it is specified on all projects.
* `sshCreds` [SSH Credentials, optional] - The default SSH credentials to use if not specified on a project. 
  Not necessary if it is specified on all projects.
* `projectsToSync` [List of VCS Project] - The projects on this server whose feature branches should be synced with the target branch.

## SSH Credentials
*Note that SSH keys must be in PEM format, i.e. generated via `ssh-keygen -m pem` or the bot will fail to authenticate when cloning/pushing.
This is a limitation of the pure-Java SSH library used on this project, however that makes it platform-agnostic.*
* `publicKey` [Credential Source] - Information on where to find the public key for SSH cloning.
* `privateKey` [Credential Source] - Information on where to find the private key for SSH cloning.

## Credential Source
* `type` [String] - The type of credential source, or how to retrieve the credential. Value should be one of `FILE`, `ENVIRONMENT`, or `INLINE`.
* `INLINE` Fields:
  * `value` [String] - The literal content of a sensitive credential. Good for local development, but not recommended for automation.
* `ENVIRONMENT` Fields:
  * `envVar` [String] - The environment variable to pull the credential from.
* `FILE` Fields:
  * `location` - The location of the file where the bot should read the credential. Can be relative to the bot's working directory or an absolute path.
    Make sure the bot has read permissions on the file in question!

## VCS Project
* `pathWithNamespace` [String] - The project's path, with included namespace. For example, this repository's value might be `emanguy/feature-branch-bot`.
* `sshCloneURL` [String] - The SSH url where the bot can clone the code from for this project.
* `syncLabel` [String, optional] - The label to look for on merge/pull requests for this project. Other merge/pull requests will be ignored. 
  This value is more specific than the server-level value, and will override it if it is specified here.
* `mainBranchName` [String, optional] - The name of the primary branch in the repo, used during the initial clone. Defaults to "master" if not specified.
* `sshCreds` [SSH Credentials, optional] - SSH credentials to use for this specific project during clone/push operations.
  This will override the server-level value if it is provided here.
