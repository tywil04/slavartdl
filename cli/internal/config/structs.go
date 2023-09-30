package config

type ConfigCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfigDownloadCmdIgnore struct {
	Cover   bool `json:"cover"`
	SubDirs bool `json:"subdirs"`
}

type ConfigDownloadCmdSkip struct {
	Unzip bool `json:"unzip"`
}

type ConfigDownloadCmd struct {
	UseDiscord bool                     `json:"usediscord"`
	OutputDir  string                   `json:"outputdir"`
	LogLevel   string                   `json:"loglevel"`
	Quality    int                      `json:"quality"`
	Timeout    int                      `json:"timeout"`
	Cooldown   int                      `json:"cooldown"`
	Ignore     *ConfigDownloadCmdIgnore `json:"ignore"`
	Skip       *ConfigDownloadCmdSkip   `json:"skip"`
}

type Config struct {
	DivoltSessionTokens     []string            `json:"divoltsessiontokens"`
	DivoltLoginCredentials  []*ConfigCredential `json:"divoltlogincredentials"`
	DiscordSessionTokens    []string            `json:"discordsessiontokens"`
	DiscordLoginCredentials []*ConfigCredential `json:"discordlogincredentials"`
	DownloadCmd             *ConfigDownloadCmd  `json:"downloadcmd"`
}
