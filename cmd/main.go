package main

import (
	"context"
	"fmt"
	"image/color"
	"image/png"
	"log/slog"
	"net/http"

	"github.com/DillonEnge/keizai-launcher/internal/fonts"
	"github.com/DillonEnge/keizai-launcher/internal/game"
	"github.com/DillonEnge/keizai-launcher/internal/requests"
	"github.com/DillonEnge/keizai-launcher/internal/sysio"
	"github.com/DillonEnge/keizai-launcher/internal/ui/button"
	"github.com/DillonEnge/keizai-launcher/internal/ui/drawer"
	"github.com/DillonEnge/keizai-launcher/internal/ui/label"
	"github.com/DillonEnge/keizai-launcher/internal/ui/panel"
	"github.com/google/go-github/v62/github"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Engehost Launcher")

	client := requests.Client{URL: "https://game-registry.engehost.net"}

	games, err := client.GetGames()
	if err != nil {
		panic(err)
	}

	gClient := github.NewClient(nil)

	release, resp, err := gClient.Repositories.GetLatestRelease(context.Background(), games[0].RepoOwner, games[0].RepoName)
	if err != nil {
		panic(fmt.Sprintf("err: %+v, resp: %+v", err, resp))
	}

	for _, v := range release.Assets {
		url := v.BrowserDownloadURL
		if url == nil {
			panic("failed to find release url")
		}
	}

	t, err := newTxtRenderer()
	if err != nil {
		slog.Error("failed to create txt renderer", "err", err)
		return
	}

	sio, err := sysio.NewSysio()
	if err != nil {
		panic(err)
	}

	ss := game.NewStateStore()

	ss.SetState("game", games[0])

	checkGameButton := button.NewButton(
		.28, .9,
		.12, .07,
		36,
		color.RGBA{32, 96, 246, 255},
		color.White,
		"",
		t,
	)

	checkGameButton.AddHandler(button.HANDLER_ON_MOUNT, func(b *button.Button) error {
		s, err := ss.GetState("game")
		if err != nil {
			return err
		}
		g, ok := s.(requests.Game)
		if !ok {
			return fmt.Errorf("failed to convert state to Game")
		}
		ok, err = sio.CheckForGame(g)
		if err != nil {
			return err
		}
		if ok {
			latest, err := sio.CheckLatest(gClient, g)
			if err != nil {
				return err
			}
			if latest {
				b.SetState(button.STATE_PLAY)
			} else {
				b.SetState(button.STATE_UPDATE)
			}
		} else {
			b.SetState(button.STATE_INSTALL)
		}
		return nil
	})
	checkGameButton.AddHandler(button.HANDLER_ON_CLICK, func(b *button.Button) error {
		s, err := ss.GetState("game")
		if err != nil {
			return err
		}
		g, ok := s.(requests.Game)
		if !ok {
			return fmt.Errorf("failed to convert state to Game")
		}

		switch b.GetState() {
		case button.STATE_PLAY:
			path, err := sio.GetInstallDirPath()
			if err != nil {
				return err
			}

			if err = sio.ExecuteGame(path, g); err != nil {
				return err
			}
		case button.STATE_INSTALL:
			fallthrough
		case button.STATE_UPDATE:
			fp, err := sio.DownloadLatestRelease(gClient, g)
			if err != nil {
				return err
			}
			if err = sio.InstallLatestRelease(fp, g); err != nil {
				return err
			}
		default:
		}

		ok, err = sio.CheckForGame(g)
		if err != nil {
			return err
		}
		if ok {
			// checkLatest("MacOS")
			b.SetState(button.STATE_PLAY)
		} else {
			b.SetState(button.STATE_INSTALL)
		}
		return nil
	})

	options := make([]drawer.Option, 0)

	for _, v := range games {
		var img *ebiten.Image

		if v.IconURL != "" {
			img, err = ebitenImageFromURL(v.IconURL)
			if err != nil {
				panic(err)
			}
		}

		options = append(options, drawer.NewOption(v.Name, img))
	}
	gamesDrawer := drawer.NewDrawer(
		0, 0.2,
		.25, .8,
		0.07,
		32,
		options,
		color.RGBA{42, 42, 42, 255},
		color.White,
		t,
	)

	gamesDrawer.AddHandler(drawer.HANDLER_ON_CLICK, func(d *drawer.Drawer) error {
		for _, v := range games {
			if v.Name == d.GetSelection().GetText() {
				ss.SetState("game", v)
				break
			}
		}
		return nil
	})

	d := []game.Drawable{
		panel.NewPanel(0, 0, 1, 1, color.RGBA{42, 42, 42, 255}),
		gamesDrawer,
		panel.NewPanel(0.25, 0, .75, 1, color.RGBA{32, 32, 32, 255}),
		checkGameButton,
		label.NewLabel(
			0.1, 0.1,
			36,
			color.White,
			"Engehost Games",
			t,
		),
	}

	g := game.NewGame(t, d, ss)

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("closing game", "err", err)
		return
	}
}

func newTxtRenderer() (*etxt.Renderer, error) {
	robotoFont := fonts.F

	fontLib := etxt.NewFontLibrary()
	_, err := fontLib.ParseFontBytes(robotoFont)
	// _, _, err = fontLib.ParseDirFonts("assets/fonts")
	if err != nil {
		return nil, err
	}

	// create a new text renderer and configure it
	txtRenderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	txtRenderer.SetFont(fontLib.GetFont("Roboto"))
	txtRenderer.SetAlign(etxt.Bottom, etxt.Left)
	txtRenderer.SetSizePx(12)

	return txtRenderer, nil
}

func ebitenImageFromURL(url string) (*ebiten.Image, error) {
	resp, err := http.Get(url)
	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
