package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/tywil04/slavartdl/divolt"
	"github.com/tywil04/slavartdl/downloader"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SlavartdlUI struct {
	ctx     context.Context
	session *divolt.Session
}

func Init() *SlavartdlUI {
	return &SlavartdlUI{}
}

func (s *SlavartdlUI) Startup(ctx context.Context) {
	s.ctx = ctx
	s.session = &divolt.Session{}
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
	return divolt.SlavartAllowedHosts
}

func (s *SlavartdlUI) Login(email, password string) bool {
	err := s.session.AuthenticateWithCredentials(email, password)
	if err != nil {
		runtime.LogError(s.ctx, err.Error())
		return false
	}
	return true
}

func (s *SlavartdlUI) DownloadUrl(url, outputDir string, quality, timeout, cooldown int, skipUnzip, ignoreCover, ignoreSubdirs bool) bool {
	timeoutTime := time.Duration(timeout) * time.Second
	// cooldownDuration := time.Duration(cooldown) * time.Second

	status, err := s.session.SlavartGetBotStatus()
	if status == divolt.SlavartBotStatusOffline {
		runtime.LogError(s.ctx, "bot is offline")
		return false
	}

	message, err := s.session.SlavartSendDownloadCommand(url, quality)
	if err != nil {
		runtime.LogError(s.ctx, err.Error())
		return false
	}

	musicUrl, err := s.session.SlavartGetUploadUrl(message.Id, url, timeoutTime)
	if err != nil {
		runtime.LogError(s.ctx, err.Error())
		return false
	}

	buffer, bytesWritten, err := downloader.DownloadFile(musicUrl)
	if err != nil {
		runtime.LogError(s.ctx, err.Error())
		return false
	}

	if !skipUnzip {
		err := downloader.Unzip(buffer, bytesWritten, outputDir, ignoreSubdirs, ignoreCover)
		if err != nil {
			runtime.LogError(s.ctx, err.Error())
			return false
		}
	} else {
		outputPath := outputDir + string(os.PathSeparator) + filepath.Clean("slavart-"+time.Now().String()) + ".zip"
		err := downloader.CopyFile(buffer, outputPath)
		if err != nil {
			runtime.LogError(s.ctx, err.Error())
			return false
		}
	}

	return true
}
