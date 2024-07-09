package view

import (
	"github.com/DillonEnge/keizai-launcher/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type View struct {
	children []game.Drawable
}

func NewView(children ...game.Drawable) *View {
	return &View{
		children: children,
	}
}

func (v *View) AddChild(d ...game.Drawable) {
	for _, d := range d {
		v.children = append(v.children, d)
	}
}

func (v *View) Update(g *game.Game) error {
	for _, c := range v.children {
		if err := c.Update(g); err != nil {
			return err
		}
	}

	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	for _, c := range v.children {
		c.Draw(screen)
	}
}
