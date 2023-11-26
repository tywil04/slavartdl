# SlavartDL
[![Go Report Card](https://goreportcard.com/badge/github.com/tywil04/slavartdl)](https://goreportcard.com/report/github.com/tywil04/slavartdl)

This repo contains a command line tool, and a desktop application (WIP) which allows you to download music from the SlavArt Divolt server or the Pixeldrain Discord server. 

This tool was originally inspired by [slavart-scraper](https://github.com/D0otDo0t/slavolt-scraper), however I choose to write my own tool because I noticed there were inefficiencies in how the download link was collected.

You can find pre-built versions in the releases section.

I created this tool for educational purposes and I do not condone any form of piracy.


## Usage
Just run the executable for your system. There are pre-built executables for Windows (x86), Linux (Arm64 and x86) and MacOS (Arm64 and x86).

On all systems, download the correct executable from the releases section.

On Windows don't double click the executable, it will open a command prompt and close. Instead open a command prompt and navigate to the directory containing the `slavartdl.exe`. You then need to run this executable in the command prompt to with the correct arguments and flags.

On Linux and MacOS open a terminal, navigate to the directory containing the slavartdl executable. Then set the executable flag using `chmod +x slavartdl`. Now you just need to run the executable with the correct arguments and flags.


### Commands
To find out how to run the commands use `--help` on any command (or subcommand). 

Heres a brief structure of the commands:
```shell
slavartdl download [urls...] # downloads music from supported urls.

slavartdl version # lists the current version of the slavartdl executable.

slavartdl update  # updates slavartdl to the latest version.

slavartdl config add divoltTokens [tokens...]             # adds one to infinite divolt session tokens to config.
slavartdl config add divoltCredential [email] [password]  # adds one divolt credential to config. a credential is an email and password.
slavartdl config add discordTokens [tokens...]            # adds one to infinite discord session tokens to config.
slavartdl config add discordCredential [email] [password] # adds one discord credential to config. a credential is an email and password.

slavartdl config list divoltTokens       # lists all divolt session tokens with their id.
slavartdl config list divoltCredentials  # lists all divolt credentials with their id. a credential is an email and password.
slavartdl config list discordTokens      # lists all discord session tokens with their id.
slavartdl config list discordCredentials # lists all discord credentials with their id. a credential is an email and password.

slavartdl config remove divoltTokens [token ids...]             # removes one to infinite divolt session tokens from config using their id.
slavartdl config remove divoltCredentials [credential ids...]   # removes one to infinite divolt credentials from config using their id.
slavartdl config remove discordTokens [token ids...]            # removes one to infinite discord session tokens from config using their id.
slavartdl config remove discordCredentials [credential ids...]  # removes one to infinite discord credentials from config using their id.
```


## Config
Session tokens are stored in a local config file (use the `config` command to find the location). You do not need to manually edit the config, you can use the commands show below. The session tokens are stored in plaintext due to the simplicity of this program, this means anyone who has access to your file system can use your revolt account(s). Dont use your main account for this, I am not liable for your account getting hacked or stolen.

You can have multiple session tokens that will randomly get used per request.


### Structure
```
{
  "discordlogincredentials": [
    {
      "email": string,
      "password": string
    }
    ...
  ],
  "discordsessiontokens": [
    string 
    ...
  ],
  "divoltlogincredentials": [
    {
      "email": string,
      "password": string
    }
    ...
  ],
  "divoltsessiontokens": [
    string 
    ...
  ],
  "downloadcmd": {
    "ignore": {
      "cover": bool,
      "subdirs": bool
    },
    "outputdir": string,
    "loglevel": string,
    "quality": int,
    "timeout": int,
    "cooldown": int,
    "useDiscord": bool
  }
}
```
Notes on config:
- You will require either at least one session token in `divoltsessiontokens` or at least one dictionary with an email and password in `divoltlogincredentials`. Both of these allow this tool to access your Divolt account which is required for this bot to function. Please don't use your main account!
- The structure of `downloadcmd.timeout` has changed, its now just integer instead of a dictionary containing seconds and minutes as integers. This is because the timeout flag has also changed to only be seconds. The additional minutes flag/config was redundant so it was removed. Only old versions of slavartdl will have this incorrect config, it will have to be changed manually.
- `downloadcmd.outputdir` must be an absolute file path, not relative.

#### Getting Divolt session tokens to add to config
If you want to get your session token, its easy. Just note that once you get your session token, for it to remain active you must close the divolt tab and do not logout.

Follows these steps to get a session token:
- Log in to Divolt.
- Open browser DevTools.
- Navigate to the network tab and then select the 'Fetch/XHR' filter.
- Find and select a request with the domain of `api.divolt.xyz`.
- Select the header tab and copy the value of the `X-Session-token` from the response headers. If there is no `X-Session-token` select another request.


#### Getting Discord session tokens to add to config
If you want to get your session token, its easy. Just note that once you get your session token, for it to remain active you must close the discord tab and do not logout.

Follows these steps to get a session token:
- Log in to Discord.
- Open browser DevTools.
- Navigate to the network tab and then select the 'Fetch/XHR' filter.
- Find and select a request named `login` with a request url of `https://discord.com/api/v9/auth/login`.
- Select the response tab and copy the value of `token`.


## Building the command line tool
To build the command line application you need to have Go(lang) installed. Once installed run the following commands:
```
git clone https://github.com/tywil04/slavartdl
cd slavartdl
go build -o slavartdl cli/main.go
```


### Similar Tools
- [`Limestone`](https://github.com/dxbednarczyk/limestone)
- [`slavart-scraper`](https://github.com/D0otDo0t/slavolt-scraper)


#### Comparison
|  | `limestone` | `slavartdl` | `slavolt-scraper` |
|--|--|--|--|
| Command-line | ✅ | ✅  | ❌ |
| Terminal UI | ✅ | ❌ | ❌ |
| Graphical UI | ❌ | ✅[^1] | ❌ |
| Pixeldrain Discord Support | ❌ | ✅ | ❌
| Slavart Website Support | ✅ | ❌ | ❌ |
| Standard input support | ❌ | ✅ | ❌ |
| Unzip downloaded tracks | ❌ | ✅ | ❌ |
| Credential storage | ✅ | ✅ | ✅ |
| Session token storage | ✅ | ✅ | ❌ |
| Self-update | ❌ | ✅ | ❌ |
| Websocket to avoid large REST payloads | ✅ | ✅[^2] | ❌

The idea for this comparison was inspired from [`limestone`](https://github.com/dxbednarczyk/limestone#comparison). I have only added to it to better reflect the features that each tool offers!

[^1]: The graphical UI is still a work in progress
[^2]: As of version v1.1.16, versions before used large REST payloads
