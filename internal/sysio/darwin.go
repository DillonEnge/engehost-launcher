package sysio

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"slices"
	"strings"

	"github.com/DillonEnge/keizai-launcher/internal/requests"
	"github.com/google/go-github/v62/github"
	"github.com/walle/targz"
	"howett.net/plist"
)

const (
	DEFAULT_INSTALL_DIR_MAC = "/Library/Application Support/Engehost/"
)

type DarwinAdapter struct{}

var _ Adapter = (*DarwinAdapter)(nil)

func (d *DarwinAdapter) GetInstallDirPath() (string, error) {
	p, err := d.GetHomeDirPath()
	if err != nil {
		return "", err
	}

	return p + DEFAULT_INSTALL_DIR_MAC, nil
}

func (d *DarwinAdapter) GetHomeDirPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func (d *DarwinAdapter) DownloadLatestRelease(client *github.Client, g requests.Game) (*string, error) {
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

func (d *DarwinAdapter) InstallLatestRelease(filePath *string, g requests.Game) error {
	path, err := d.GetInstallDirPath()
	if err != nil {
		return err
	}

	err = targz.Extract(*filePath, path)
	if err != nil {
		return err
	}

	exeName, err := d.GetExecutableName(path, g)
	if err != nil {
		return err
	}

	err = os.Chmod(path+g.Name+".app/Contents/MacOS/"+*exeName, 0755)
	if err != nil {
		return err
	}

	err = os.Remove(*filePath)
	if err != nil {
		return err
	}

	return nil
}

func (d *DarwinAdapter) CheckForGame(g requests.Game) (bool, error) {
	path, err := d.GetInstallDirPath()
	if err != nil {
		return false, err
	}

	dir, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	if slices.ContainsFunc(dir, func(de fs.DirEntry) bool {
		return de.Name() == fmt.Sprintf("%s.app", g.Name)
	}) {
		return true, nil
	}

	return false, nil
}

func (d *DarwinAdapter) CheckLatest(client *github.Client, g requests.Game) (bool, error) {
	path, err := d.GetInstallDirPath()
	if err != nil {
		return false, err
	}

	ver, err := d.GetVersion(path, g)
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

func (d *DarwinAdapter) GetVersion(appPath string, g requests.Game) (*string, error) {
	f, err := os.Open(appPath + g.Name + ".app/Contents/Info.plist")
	if err != nil {
		return nil, err
	}

	var info struct {
		CFBundleVersion string `plist:"CFBundleVersion"`
	}

	if err = plist.NewDecoder(f).Decode(&info); err != nil {
		return nil, err
	}

	ver := fmt.Sprintf("v%s", info.CFBundleVersion)

	return &ver, nil
}

func (d *DarwinAdapter) GetExecutableName(appPath string, g requests.Game) (*string, error) {
	f, err := os.Open(appPath + g.Name + ".app/Contents/Info.plist")
	if err != nil {
		return nil, err
	}

	var info struct {
		CFBundleExecutable string `plist:"CFBundleExecutable"`
	}

	if err = plist.NewDecoder(f).Decode(&info); err != nil {
		return nil, err
	}

	return &info.CFBundleExecutable, nil
}

func (d *DarwinAdapter) ExecuteGame(appPath string, g requests.Game) error {
	exeName, err := d.GetExecutableName(appPath, g)
	if err != nil {
		return err
	}

	cmd := exec.Command(appPath + g.Name + ".app/Contents/MacOS/" + *exeName)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
