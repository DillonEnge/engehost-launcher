package label

import (
	"image/color"

	"github.com/DillonEnge/keizai-launcher/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type Label struct {
	textColor   color.Color
	textSize    int
	text        string
	x           float32
	y           float32
	txtRenderer *etxt.Renderer
}

func NewLabel(
	x, y float32,
	textSize int,
	textColor color.Color,
	text string,
	t *etxt.Renderer,
) *Label {
	return &Label{
		x:           float32(x),
		y:           float32(y),
		textColor:   textColor,
		textSize:    textSize,
		text:        text,
		txtRenderer: t,
	}
}

func (l *Label) Update(_ *game.Game) error {
	return nil
}

func (l *Label) Draw(screen *ebiten.Image) {
	tx := l.x * float32(screen.Bounds().Dx())
	ty := l.y * float32(screen.Bounds().Dy())

	l.txtRenderer.SetColor(l.textColor)
	l.txtRenderer.SetTarget(screen)
	l.txtRenderer.SetSizePx(l.textSize)
	l.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	l.txtRenderer.Draw(l.text, int(tx), int(ty))
}
