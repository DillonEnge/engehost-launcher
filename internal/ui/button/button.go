package button

import (
	"fmt"
	"image/color"

	"github.com/DillonEnge/keizai-launcher/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

type HandlerType int

type Handler func(b *Button) error

type Handlers map[HandlerType]Handler

const (
	HANDLER_ON_CLICK HandlerType = iota
	HANDLER_ON_MOUNT
)

const (
	STATE_INSTALL ButtonState = iota
	STATE_PLAY
	STATE_UPDATE
)

type ButtonState int

type Button struct {
	primaryColor color.Color
	textColor    color.Color
	textSize     int
	text         string
	x            float32
	y            float32
	width        float32
	height       float32
	txtRenderer  *etxt.Renderer
	hovered      bool
	clicked      bool
	handlers     Handlers
	state        ButtonState
}

func NewButton(
	x, y, width, height float32, textSize int,
	primaryColor, textColor color.Color,
	text string,
	t *etxt.Renderer,
) *Button {
	return &Button{
		width:        width,
		height:       height,
		x:            x,
		y:            y,
		primaryColor: primaryColor,
		textColor:    textColor,
		textSize:     textSize,
		text:         text,
		txtRenderer:  t,
		handlers:     Handlers{},
	}
}

func (b *Button) Update(g *game.Game) error {
	if f, ok := b.handlers[HANDLER_ON_MOUNT]; ok {
		f(b)
		delete(b.handlers, HANDLER_ON_MOUNT)
	}

	lX, lY := g.LayoutF(1280, 720)
	x, y := float32(lX), float32(lY)

	mouseX, mouseY := ebiten.CursorPosition()
	mX, mY := float32(mouseX), float32(mouseY)

	switch b.state {
	case STATE_PLAY:
		b.text = "Play"
	case STATE_INSTALL:
		b.text = "Install"
	case STATE_UPDATE:
		b.text = "Update"
	default:
	}

	if mX > b.x*x && mX < (b.x*x)+(b.width*x) && mY > b.y*y && mY < (b.y*y)+(b.height*y) {
		b.hovered = true

		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			b.clicked = true
		} else {
			b.clicked = false
		}
	} else {
		b.hovered = false
		b.clicked = false
	}

	if b.hovered && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		f, ok := b.handlers[HANDLER_ON_CLICK]
		if !ok {
			return fmt.Errorf("failed to find ON_CLICK handler")
		}

		if err := f(b); err != nil {
			return err
		}
	}

	return nil
}

func subUInt8(n uint8, subn int) uint8 {
	if int(n)-subn < 0 {
		return 0
	}

	return n - uint8(subn)
}

func (b *Button) Draw(screen *ebiten.Image) {
	tx := b.x * float32(screen.Bounds().Dx())
	ty := b.y * float32(screen.Bounds().Dy())
	tw := b.width * float32(screen.Bounds().Dx())
	th := b.height * float32(screen.Bounds().Dy())

	c := b.primaryColor
	r, g, bl, _ := c.RGBA()
	cr, cg, cbl := uint8(r), uint8(g), uint8(bl)

	if b.hovered {
		c = color.RGBA{subUInt8(cr, 20), subUInt8(cg, 20), subUInt8(cbl, 20), 255}
	}
	if b.clicked {
		c = color.RGBA{subUInt8(cr, 40), subUInt8(cg, 40), subUInt8(cbl, 40), 255}
	}

	vector.DrawFilledRect(
		screen,
		tx, ty,
		tw, th,
		c,
		false,
	)

	b.txtRenderer.SetColor(b.textColor)
	b.txtRenderer.SetTarget(screen)
	b.txtRenderer.SetSizePx(b.textSize)
	b.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	b.txtRenderer.Draw(
		b.text,
		int(tx+(tw/2)),
		int(ty+(th/2)),
	)
}

func (b *Button) AddHandler(key HandlerType, h Handler) {
	b.handlers[key] = h
}

func (b *Button) SetText(t string) {
	b.text = t
}

func (b *Button) SetState(s ButtonState) {
	b.state = s
}

func (b *Button) GetState() ButtonState {
	return b.state
}
