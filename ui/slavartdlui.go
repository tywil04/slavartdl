package main

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/tywil04/slavartdl/slavart"
)

type SlavartdlUI struct {
	ctx context.Context

	sessionToken string
}

func Init() *SlavartdlUI {
	return &SlavartdlUI{}
}

func (s *SlavartdlUI) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *SlavartdlUI) Shutdown(ctx context.Context) {
}

func (s *SlavartdlUI) OpenDirectoryDialog(title string) string {
	opts := runtime.OpenDialogOptions{
		Title:                      title,
		CanCreateDirectories:       true,
		ResolvesAliases:            true,
		TreatPackagesAsDirectories: false,
	}

	file, err := runtime.OpenDirectoryDialog(s.ctx, opts)
	if err != nil {
		runtime.LogFatal(s.ctx, err.Error())
	}

	return file
}

func (s *SlavartdlUI) OpenFileDialog(title string) string {
	opts := runtime.OpenDialogOptions{
		Title:                      title,
		CanCreateDirectories:       false,
		ResolvesAliases:            true,
		TreatPackagesAsDirectories: false,
	}

	file, err := runtime.OpenFileDialog(s.ctx, opts)
	if err != nil {
		runtime.LogFatal(s.ctx, err.Error())
	}

	return file
}

func (s *SlavartdlUI) SaveFileDialog(defaultFilename string) string {
	opts := runtime.SaveDialogOptions{
		CanCreateDirectories:       false,
		TreatPackagesAsDirectories: false,
	}

	if defaultFilename != "" {
		opts.DefaultFilename = defaultFilename
	}

	file, err := runtime.SaveFileDialog(s.ctx, opts)
	if err != nil {
		runtime.LogFatal(s.ctx, err.Error())
	}

	return file
}

func (s *SlavartdlUI) GetAllowedHosts() []string {
	return slavart.AllowedHosts
}

func (s *SlavartdlUI) Login(email, password string) bool {
	var err error
	s.sessionToken, err = slavart.GetSessionTokenFromCredentials(email, password)
	if err != nil {
		runtime.LogError(s.ctx, err.Error())
		return false
	}
	return true
}

func (s *SlavartdlUI) DownloadUrl(url, outputDir string, quality, timeout, cooldown int, skipUnzip, ignoreCover, ignoreSubdirs bool) {
	timeoutTime := time.Now().Add(time.Duration(timeout) * time.Second)
	cooldownDuration := time.Duration(cooldown) * time.Second

	slavart.DownloadUrl(
		url,
		s.sessionToken,
		"silent",
		quality,
		timeoutTime,
		cooldownDuration,
		outputDir,
		skipUnzip,
		ignoreCover,
		ignoreSubdirs,
	)
}
