package matrix

type Point struct {
	x float64
	y float64
}

func NewPointFloat(x float64, y float64) Point {
	return Point{
		x: x,
		y: y,
	}
}

func (p *Point) AddTo(x float64, y float64) {
	p.x += x
	p.y += y
}

func (p *Point) AddToX(x float64) {
	p.x += x
}

func (p *Point) AddToY(y float64) {
	p.y += y
}

func (p *Point) MoveTo(x float64, y float64) {
	p.x = x
	p.y = y
}

func (p *Point) MoveToX(x float64) {
	p.x = x
}

func (p *Point) MoveToY(y float64) {
	p.y = y
}

func (p *Point) GetX() float64 {
	return p.x
}

func (p *Point) GetY() float64 {
	return p.y
}
