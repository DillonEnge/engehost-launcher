package drawer

import (
	"image/color"

	"github.com/DillonEnge/keizai-launcher/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

type HandlerType int

type Handler func(d *Drawer) error

type Handlers map[HandlerType]Handler

const (
	HANDLER_ON_CLICK HandlerType = iota
	HANDLER_ON_MOUNT
)

type Option struct {
	text    string
	image   *ebiten.Image
	hovered bool
}

func NewOption(text string, img *ebiten.Image) Option {
	return Option{
		text:  text,
		image: img,
	}
}

func (o Option) GetText() string {
	return o.text
}

type Drawer struct {
	primaryColor color.Color
	options      []Option
	x            float32
	y            float32
	width        float32
	height       float32
	selection    int
	textSize     int
	textColor    color.Color
	optionHeight float32
	handlers     Handlers
	txtRenderer  *etxt.Renderer
}

func NewDrawer(
	x, y, width, height, optionHeight float32, textSize int,
	options []Option,
	primaryColor, textColor color.Color,
	t *etxt.Renderer,
) *Drawer {
	return &Drawer{
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		textSize:     textSize,
		textColor:    textColor,
		selection:    0,
		optionHeight: optionHeight,
		options:      options,
		primaryColor: primaryColor,
		handlers:     Handlers{},
		txtRenderer:  t,
	}
}

func (d *Drawer) Update(g *game.Game) error {
	if f, ok := d.handlers[HANDLER_ON_MOUNT]; ok {
		f(d)
		delete(d.handlers, HANDLER_ON_MOUNT)
	}

	lX, lY := g.LayoutF(1280, 720)
	x, y := float32(lX), float32(lY)
	tx := d.x * x
	ty := d.y * y
	tw := d.width * x
	toh := d.optionHeight * y

	mouseX, mouseY := ebiten.CursorPosition()
	mX, mY := float32(mouseX), float32(mouseY)

	for i, v := range d.options {
		if mX > tx && mX < tx+tw && mY > ty+(toh*float32(i)) && mY < ty+(toh*float32(i))+toh {
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
				d.selection = i
				err := d.handlers[HANDLER_ON_CLICK](d)
				if err != nil {
					return err
				}
			}
			v.hovered = true
		} else {
			v.hovered = false
		}
		d.options[i] = v
	}
	return nil
}

func (d *Drawer) Draw(screen *ebiten.Image) {
	tx := d.x * float32(screen.Bounds().Dx())
	ty := d.y * float32(screen.Bounds().Dy())
	tw := d.width * float32(screen.Bounds().Dx())
	th := d.height * float32(screen.Bounds().Dy())
	toh := d.optionHeight * float32(screen.Bounds().Dy())

	vector.DrawFilledRect(
		screen,
		tx, ty,
		tw, th,
		d.primaryColor,
		false,
	)

	for i, v := range d.options {
		var c color.Color
		c = d.primaryColor

		if v.hovered {
			c = subRGBA(d.primaryColor, 3)
		}
		if d.selection == i {
			c = subRGBA(d.primaryColor, 6)
		}

		vector.DrawFilledRect(
			screen,
			tx, ty+(float32(i)*toh),
			tw, toh,
			c,
			false,
		)

		do := &ebiten.DrawImageOptions{}

		relImageSize := 0.75

		do.GeoM.Scale(float64((toh*float32(relImageSize))/float32(v.image.Bounds().Dx())), float64((toh*float32(relImageSize))/float32(v.image.Bounds().Dy())))
		do.GeoM.Translate(float64(tx+(tw/15)), float64(ty+(float32(i)*toh)+(toh*float32((1-relImageSize)/2))))

		screen.DrawImage(v.image, do)

		d.txtRenderer.SetColor(d.textColor)
		d.txtRenderer.SetTarget(screen)
		d.txtRenderer.SetSizePx(d.textSize)
		d.txtRenderer.SetAlign(etxt.YCenter, etxt.Left)
		d.txtRenderer.Draw(v.text, int(tx+(tw/4)), int(ty+(float32(i)*toh)+(toh/2)))
	}
}

func (d *Drawer) AddHandler(key HandlerType, h Handler) {
	d.handlers[key] = h
}

func (d *Drawer) GetSelection() Option {
	return d.options[d.selection]
}

func subUInt8(n uint8, subn int) uint8 {
	if int(n)-subn < 0 {
		return 0
	}

	return n - uint8(subn)
}

func subRGBA(c color.Color, subn int) color.Color {
	r, g, bl, _ := c.RGBA()
	cr, cg, cbl := uint8(r), uint8(g), uint8(bl)

	return color.RGBA{subUInt8(cr, subn), subUInt8(cg, subn), subUInt8(cbl, subn), 255}
}
