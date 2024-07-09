package panel

import (
	"image/color"

	"github.com/DillonEnge/keizai-launcher/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Panel struct {
	x       float32
	y       float32
	width   float32
	height  float32
	bgColor color.Color
}

func NewPanel(
	x, y, width, height float32,
	bgColor color.Color,
) *Panel {
	return &Panel{
		x:       x,
		y:       y,
		width:   width,
		height:  height,
		bgColor: bgColor,
	}
}

func (p *Panel) Update(g *game.Game) error {
	return nil
}

func (p *Panel) Draw(screen *ebiten.Image) {
	tx := p.x * float32(screen.Bounds().Dx())
	ty := p.y * float32(screen.Bounds().Dy())
	tw := p.width * float32(screen.Bounds().Dx())
	th := p.height * float32(screen.Bounds().Dy())

	vector.DrawFilledRect(
		screen,
		tx, ty,
		tw, th,
		p.bgColor,
		false,
	)
}
