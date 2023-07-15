# Slavart Downloader
This is a simple tool written in Go(lang) to download music using the SlavArt Divolt server. This tool was inspired by [https://github.com/D0otDo0t/slavolt-scraper](https://github.com/D0otDo0t/slavolt-scraper), however I choose to write my own tool because I noticed there were inefficiencies in how the download link was collected.

I created this tool for educational purposes and I do not condone any form of piracy.

You can find pre-build versions in the releases section.

## Usage
Just run the executable for your system. There are prebuild executables for Windows, Linux and MacOS (darwin). Linux and MacOS have arm compatible versions.

On Windows don't double click the executable, it will open a command prompt and close instead open a command prompt and navigate to the directory containing the `slavartdl.exe`. You then need to run this executable in the command prompt to successfully run the command.

## Config
Session tokens are stored in a local config file (use `slavartdl config` to find the location). You do not need to manually edit the config, you can use the commands show below. The session tokens are stored in plaintext due to the simplicity of this program, this means anyone who has access to your file system can use your revolt account(s). Dont use your main account for this, I am not liable for your account getting hacked or stolen.

You can have multiple session tokens that will randomly get used per request.

### Getting session tokens to add to config
I recommend following the guide from [D0otDo0t](https://github.com/D0otDo0t/slavolt-scraper).

## Commands
- `<value>`: 1 argument needed for value
- `...<value>`: as many arguments needed for value (0 to inf)
- `[value]`: optional value

Download from services with flags
```
slavartdl download <url> [flags]

flags:
-o, --output-directory:         (required) the directory to save the files to

-q, --quality:                  the quality of music to download, omit (or -1) for best quality available
                                (1: 128kbps MP3/AAC, 2: 320kbps MP3/AAC, 3: 16bit 44.1kHz, 4: 24bit ≤96kHz, 5: 24bit ≤192kHz)

-s, --timeout-duration-seconds: how long it takes to search for a link before it gives up in seconds
                                (this combines with --timeout-duration-minutes)

-m, --timeout-duration-minutes: how long it takes to search for a link before it gives up in minutes
                                (this combines with --timeout-duration-seconds)

-c, --ignore-cover:             whether to ignore cover.jpg when unzipping downloaded music

-d, --ignore-subdirectories:    whether to ignore subdirectories when unzipping downloaded music

-h, --help:                     help command
```

Display the location of the config file
```
slavartdl config [flags]

flags:
-h, --help: help command
```

Add one or multiple tokens to config
```
slavartdl config add tokens ...<tokens> [flags]

flags:
-h, --help: help command
```

List all tokens stored in config with there `tokenindex`
```
slavartdl config list tokens [flags]

flags:
-h, --help: help command
```

Remove token using index show in `slavartdl list tokens`
```
slavartdl config remove tokens ...<tokenindexes> [flags]

flags:
-h, --help: help command
```

## Building
To build this application you need to have Go(lang) installed. Once installed run the following commands:
```
git clone https://github.com/tywil04/slavartdl
cd slavartdl
go build -o slavartdl main.go
```
