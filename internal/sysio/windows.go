//go:build windows

package sysio

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"github.com/DillonEnge/keizai-launcher/internal/requests"
	"github.com/bi-zone/go-fileversion"
	"github.com/google/go-github/v62/github"
	"github.com/walle/targz"
)

const (
	DEFAULT_INSTALL_DIR_WINDOWS = "\\Program Files\\Engehost"
)

type WindowsAdapter struct{}

var _ Adapter = (*WindowsAdapter)(nil)

func (w *WindowsAdapter) GetInstallDirPath() (string, error) {
	p, err := w.GetHomeDirPath()
	if err != nil {
		return "", err
	}

	return p + DEFAULT_INSTALL_DIR_WINDOWS, nil
}

func (w *WindowsAdapter) GetHomeDirPath() (string, error) {
	return "C:", nil
}

func (w *WindowsAdapter) DownloadLatestRelease(client *github.Client, g requests.Game) (*string, error) {
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), g.RepoOwner, g.RepoName)
	if err != nil {
		return nil, err
	}

	var asset *github.ReleaseAsset
	for _, v := range release.Assets {
		if strings.Contains(*v.Name, fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)) {
			asset = v
		}
	}

	if asset == nil {
		return nil, fmt.Errorf("failed to find asset with GOOS and GOARCH")
	}

	out, err := os.Create(*asset.Name)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	resp, err := http.Get(*asset.BrowserDownloadURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	return asset.Name, nil
}

func (w *WindowsAdapter) InstallLatestRelease(filePath *string, g requests.Game) error {
	path, err := w.GetInstallDirPath()
	if err != nil {
		return err
	}

	err = targz.Extract(*filePath, path)
	if err != nil {
		return err
	}

	err = os.Chmod(path+"\\"+strings.ToLower(g.Name)+"\\"+g.Name+".exe", 0755)
	if err != nil {
		return err
	}

	err = os.Remove(*filePath)
	if err != nil {
		return err
	}

	return nil
}

func (w *WindowsAdapter) CheckForGame(g requests.Game) (bool, error) {
	path, err := w.GetInstallDirPath()
	if err != nil {
		return false, err
	}

	dir, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	if slices.ContainsFunc(dir, func(de fs.DirEntry) bool {
		return de.Name() == strings.ToLower(g.Name)
	}) {
		return true, nil
	}

	return false, nil
}

func (w *WindowsAdapter) CheckLatest(client *github.Client, g requests.Game) (bool, error) {
	path, err := w.GetInstallDirPath()
	if err != nil {
		return false, err
	}

	ver, err := w.GetVersion(path, g)
	if err != nil {
		return false, err
	}

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), g.RepoOwner, g.RepoName)
	if err != nil {
		return false, err
	}

	if *release.Name != *ver {
		return false, nil
	}

	return true, nil
}

func (w *WindowsAdapter) GetVersion(appPath string, g requests.Game) (*string, error) {
	f, err := fileversion.New(appPath + "\\" + strings.ToLower(g.Name) + "\\" + g.Name + ".exe")
	if err != nil {
		return nil, err
	}

	v, ok := strings.CutSuffix(f.FixedInfo().FileVersion.String(), ".0")
	if !ok {
		return nil, fmt.Errorf("failed to find trailing .0 in .exe fileversion")
	}

	ver := fmt.Sprintf("v%s", v)

	return &ver, nil
}

func (w *WindowsAdapter) GetExecutableName(appPath string, g requests.Game) (*string, error) {
	return &g.Name, nil
}

func (w *WindowsAdapter) ExecuteGame(appPath string, g requests.Game) error {
	cmd := exec.Command(appPath + "\\" + strings.ToLower(g.Name) + "\\" + g.Name + ".exe")
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
