package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SlavartdlUI struct {
	ctx context.Context
}

func Init() *SlavartdlUI {
	return &SlavartdlUI{}
}

func (s *SlavartdlUI) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *SlavartdlUI) Shutdown(ctx context.Context) {
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
