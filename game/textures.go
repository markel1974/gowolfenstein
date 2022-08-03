package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
)

type ColorData struct {
	color       color.RGBA
	x           int
	y           int
	minDist     float64
	minDistComp float64
}

func NewColorData(minDist int, c color.RGBA, x int, y int) *ColorData {
	cd := &ColorData{
		minDist: float64(minDist),
		color:   c,
		x:       x,
		y:       y,
	}
	cd.minDistComp = 1 / cd.minDist
	return cd
}

func (c *ColorData) Get(dist float64, light float64) color.RGBA {
	if dist <= c.minDist {
		if light <= 1 {
			return c.color
		}
		r := float64(c.color.R) * light
		if r > 255 {
			r = 255
		}
		g := float64(c.color.G) * light
		if g > 255 {
			g = 255
		}
		b := float64(c.color.B) * light
		if b > 255 {
			b = 255
		}
		return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: c.color.A}
	}
	dist *= c.minDistComp
	r := (float64(c.color.R) / dist) * light
	if r > 255 {
		r = 255
	}
	g := (float64(c.color.G) / dist) * light
	if g > 255 {
		g = 255
	}
	b := (float64(c.color.B) / dist) * light
	if b > 255 {
		b = 255
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: c.color.A}
}

type Texture struct {
	width     int
	ratioWith float64
	data      [][]*ColorData
	cache     []*ColorData
}

func NewTexture(max int) *Texture {
	data := make([][]*ColorData, max)
	for x := 0; x < max; x++ {
		data[x] = make([]*ColorData, max)
	}
	t := &Texture{
		width:     0,
		ratioWith: 0.0,
		data:      data,
	}
	return t
}

func (t *Texture) Setup() {
	for _, d := range t.data {
		for _, k := range d {
			if k.color.A > 0 {
				t.cache = append(t.cache, k)
			}
		}
	}
}

type Textures struct {
	w           int
	h           int
	data        *image.RGBA
	compiled    []*Texture
	empty       color.RGBA
	size        uint
	compiledLen uint
	minDist     int
}

func NewTextures(w int, h int, minDist int) *Textures {
	if minDist < 1 {
		minDist = 1
	}
	t := &Textures{
		w:       w,
		h:       h,
		minDist: minDist,
		empty:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
	t.data = t.load()
	t.compiled, t.size = t.compile(t.data)

	for _, z := range t.compiled {
		z.Setup()
	}

	t.compiledLen = uint(len(t.compiled))
	return t
}

func (t *Textures) compile(data *image.RGBA) ([]*Texture, uint) {
	var out []*Texture
	if data == nil {
		return out, 0
	}
	maxX := data.Rect.Max.X
	maxY := data.Rect.Max.Y
	counter := 0
	compiled := NewTexture(maxX)

	for y := 0; y < maxY; y++ {
		if counter == maxX {
			if compiled != nil {
				compiled.ratioWith = float64(compiled.width) / (float64(maxX)) * 1.0
				out = append(out, compiled)
				//fmt.Println(len(out) -1, maxX, compiled.width, compiled.ratioWith)
				compiled = nil
			}
			counter = 0
		}
		w := 0
		for x := 0; x < maxX; x++ {
			if compiled == nil {
				compiled = NewTexture(maxX)
			}
			c := data.RGBAAt(x, y)
			if c.A > 0 {
				w++
			}
			compiled.data[x][counter] = NewColorData(t.minDist, c, x, counter)
		}
		if compiled != nil {
			if w > compiled.width {
				compiled.width = w
			}
		}
		counter++
	}
	if compiled != nil {
		out = append(out, compiled)
	}
	fmt.Println("total textures: ", len(out))
	return out, uint(maxX)
}

func (t *Textures) load() *image.RGBA {
	file, err := os.ReadFile("resources" + string(os.PathSeparator) + "textures.png")
	if err != nil {
		return nil
	}
	p, err := png.Decode(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	m := image.NewRGBA(p.Bounds())
	draw.Draw(m, m.Bounds(), p, image.Point{}, draw.Src)
	return m
}

func (t *Textures) Cache(id uint) []*ColorData {
	if id < t.compiledLen {
		return t.compiled[id].cache
	}
	return nil
}

func (t *Textures) RatioWidth(id uint) float64 {
	if id < t.compiledLen {
		return t.compiled[id].ratioWith
	}
	return 0.0
}

func (t *Textures) RGBAAt(id uint, tx uint, ty uint, distance float64, light float64) color.RGBA {
	if !t.Valid(id, tx, ty) {
		return t.empty
	}
	return t.compiled[id].data[tx][ty].Get(distance, light)
}

func (t *Textures) Size() uint {
	return t.size
}

func (t *Textures) Valid(id uint, tx uint, ty uint) bool {
	if id < t.compiledLen && tx < t.size && ty < t.size {
		return true
	}
	return false
}

func (t *Textures) Save() {
	size := int(t.Size())
	for imageId, compiled := range t.compiled {
		m := image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: size, Y: size}})
		for x, data := range compiled.data {
			line := x * 4
			for y, c := range data {
				cursor := y * (size * 4)
				idx := line + cursor
				m.Pix[idx] = c.color.R
				m.Pix[idx+1] = c.color.G
				m.Pix[idx+2] = c.color.B
				m.Pix[idx+3] = c.color.A
			}
		}
		f, err := os.Create("resources" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + "image" + strconv.Itoa(imageId) + ".png")
		if err != nil {
			log.Fatal(err)
		}
		if err := png.Encode(f, m); err != nil {
			f.Close()
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
