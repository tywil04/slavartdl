# Slavart Downloader
This is a simple tool written in Go(lang) to download using the slavart divolt server. 

## Config
Session tokens are stored in a local config file (use `slavartdl config` to find the location). You do not need to manually edit the config, you can use the commands show below. The session tokens are stored in plaintext due to the simplicity of this program, this means anyone who has access to your file system can use your revolt account(s). Dont use your main account for this, I am not liable for your account getting hacked or stolen.

## Commands
Download from service with flags
```
slavartdl -> download <url> [flags]

flags:
-o, --output-directory: (required) the directory to save the files to
-q, quality: the quality of music to download, omit (or -1) for best quality available (1: 128kbps MP3/AAC, 2: 320kbps MP3/AAC, 3: 16bit 44.1kHz, 4: 24bit ≤96kHz, 5: 24bit ≤192kHz)
-h, --help: help command
```

Display the location of the config file
```
slavartdl -> config [flags]

flags:
-h, --help: help command
```

Modify the config
```
slavartdl -> config -> add -> tokens ...tokens [flags]
slavartdl -> config -> list [flags]
slavartdl -> config -> remove -> tokens ...tokenIndexes [flags]

flags:
-h, --help: help command
```