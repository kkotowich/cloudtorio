# cloudtorio
By harnessing the version control powers of Git and the file storage capacity of GitHub, this application allows for asynchronous factorio multiplayer action.

## Generate API Personal Access Token
https://github.com/settings/tokens

Generate an API access token from github, granting it scopes for repo access.

## Configuration
When you run the app for the first time you be prompted for config values. This will create a `config.json` file in the same location as the executable is launched. 

| key          | description                                                 | example  |
| ------------ |-------------------------------------------------------------| -----|
| apiKey       | api personal access token                                   |  |
| repoUrl      | the url for the save game repository                        | https://github.com/kkotowich/cloudtorio-save |
| saveGameName | the name of the save game file without the file extention   | cloudtorio-save    |
| saveGamePath | the absolute path to the save game file on local disk       | C:\\Users\\kkotowich\\AppData\\Roaming\\Factorio\\saves

**Currently, `saveGamePath` configuration setting must be entered manually**

## Running the App
Run ./cloudtorio.exe

download: Downloads the save game file from the github repo. This will write the file to disk and overwrite the save game data.
upload: Uploads the local save game file to the guthub repo. This will create a commit with the new save game data.
edit config: deletes the `config.json` file and re-runs the config input propmpts. **Be careful not to lose your API Token or else you will have to generate a new one.**
exit: exits the application

## Create Executable
```shell
go build
```
This will generate an executable that you can use to run the application on your OS. For example, windows will generate `cloudtorio.exe`

## TODO
1. unit tests
1. refactor code to not be ugly
1. fix save game path key
1. macos support
1. allow multiple configs to be saved, swapped