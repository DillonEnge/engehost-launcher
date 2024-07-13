//go:build windows

package sysio

import (
	"errors"
	"runtime"

	"github.com/DillonEnge/keizai-launcher/internal/requests"
	"github.com/google/go-github/v62/github"
)

var (
	ErrUnsupportedSystem = errors.New("unsupported system")
)

type Adapter interface {
	GetInstallDirPath() (string, error)
	GetHomeDirPath() (string, error)
	DownloadLatestRelease(client *github.Client, g requests.Game) (*string, error)
	InstallLatestRelease(filePath *string, g requests.Game) error
	CheckForGame(g requests.Game) (bool, error)
	CheckLatest(client *github.Client, g requests.Game) (bool, error)
	GetVersion(appPath string, g requests.Game) (*string, error)
	GetExecutableName(appPath string, g requests.Game) (*string, error)
	ExecuteGame(appPath string, g requests.Game) error
}

func NewSysio() (Adapter, error) {
	switch runtime.GOOS {
	case "windows":
		return &WindowsAdapter{}, nil
	default:
		return nil, ErrUnsupportedSystem
	}
}
