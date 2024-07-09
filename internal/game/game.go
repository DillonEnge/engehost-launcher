package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type StateStore struct {
	store map[string]interface{}
}

func (s *StateStore) GetState(key string) (interface{}, error) {
	v, ok := s.store[key]
	if !ok {
		return nil, fmt.Errorf("state not found with key: %s", key)
	}

	return v, nil
}

func (s *StateStore) SetState(key string, val interface{}) {
	s.store[key] = val
}

func NewStateStore() *StateStore {
	return &StateStore{
		store: make(map[string]interface{}),
	}
}

type Game struct {
	txtRender   *etxt.Renderer
	drawables   []Drawable
	sharedState *StateStore
}

type Drawable interface {
	Draw(*ebiten.Image)
	Update(*Game) error
}

func NewGame(r *etxt.Renderer, d []Drawable, ss *StateStore) *Game {
	return &Game{
		txtRender:   r,
		drawables:   d,
		sharedState: ss,
	}
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{42, 42, 42, 0})
	for _, v := range g.drawables {
		v.Draw(screen)
	}
}

func (g *Game) Update() error {
	for _, v := range g.drawables {
		if err := v.Update(g); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	panic("unused")
}

func (g *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	canvasWidth := math.Ceil(logicWinWidth * scale)
	canvasHeight := math.Ceil(logicWinHeight * scale)
	return canvasWidth, canvasHeight
}
