package main

import (
	"github.com/markel1974/gowolfenstein/pixels"
	"image/color"
)

type MiniMap struct {
	world        *World
	p            *pixels.PictureRGBA
	defaultColor color.RGBA
	x            float64
	y            float64
}

func NewMiniMap(world *World) *MiniMap {
	return &MiniMap{
		world:        world,
		x:            float64((-world.GetWidth() / 2) - 5),
		y:            float64((world.GetHeight() / 2) + 5),
		p:            pixels.NewPictureRGBA(pixels.R(float64(0), float64(0), float64(world.GetWidth()), float64(world.GetHeight()))),
		defaultColor: color.RGBA{R: 255, A: 255},
	}
}

func (m *MiniMap) Update(posX float64, posY float64) {
	for x, row := range m.world.GetWorld() {
		for y, data := range row {
			var r, g, b, a uint8
			/*
				var c color.RGBA
				if data != nil {
					c.R = uint8(data.GetTextureId())
					c.G = uint8(128)
					c.B = uint8(data.GetTextureId())
					c.A = 6
				}
				m.p.SetColor(x, y, c)
			*/
			if data != nil {
				r = uint8(data.GetTextureId())
				g = 128
				b = uint8(data.GetTextureId())
				a = 6
			}
			m.p.SetRGBA(x, y, r, g, b, a)
		}
	}

	id, avail := m.world.GetWorldAvailable()
	for x, row := range avail {
		for y, data := range row {
			if data.id == id {
				var r, g, b, a uint8
				switch data.kind {
				case "wall":
					r = uint8(128)
					g = uint8(128)
					b = uint8(127)
					a = 6
				default:
					r = uint8(255)
					a = 255
				}

				m.p.SetRGBA(x, y, r, g, b, a)
			}
		}
	}

	m.p.SetColor(int(posX-1), int(posY), m.defaultColor)
	m.p.SetColor(int(posX), int(posY-1), m.defaultColor)
	m.p.SetColor(int(posX), int(posY), m.defaultColor)
	m.p.SetColor(int(posX+1), int(posY), m.defaultColor)
	m.p.SetColor(int(posX), int(posY+1), m.defaultColor)
	/*
		if as.active {
			i.Set(as.X, as.Y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		} else {
			i.Set(as.X, as.Y, color.RGBA{R: 64, G: 64, B: 64, A: 255})
		}
	*/
}

func (m *MiniMap) GetPicture() pixels.IPicture {
	return m.p
}
