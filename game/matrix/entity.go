package matrix

import (
	"math"
)

const (
	friction = 0.9
)

type Entity struct {
	Rect
	Id       string
	Mass     float64
	VxMin    float64
	Vx       float64
	VyMin    float64
	Vy       float64
	Friction float64
	G        float64
	GForce   float64
	impulse  float64
	//breaker  *Entity
	Collider *Entity
}

func calcG(e *Entity) float64 {
	if e.GForce == 0.0 {
		return 0.0
	}
	return (math.Abs(e.Vx) + math.Abs(e.Vy)) * e.GForce
}

func CalcDistance(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	d := (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)
	sd := math.Sqrt(d)
	if sd <= 0 {
		sd = 0.01
	}
	return sd
}

func NewEntity(x float64, y float64, w float64, h float64, mass float64) *Entity {
	a := &Entity{
		Id:   NextUUId(),
		Rect: NewRect(x, y, w, h, 1.0),
		Mass: mass,
		Vx:   0.0,
		Vy:   0.0,

		Friction: friction,
		G:        0.0,
		GForce:   0.0,
		VxMin:    0.001,
		VyMin:    0.001,
		impulse:  0.001,
		//breaker: nil,
	}
	return a
}

func (e *Entity) Invalidate() {
	//e.clearBreaker()
	e.clearCollider()
}

func (e *Entity) HasCollision(obj2 *Entity) bool {
	return e.rectIntersect(e.point.x, e.point.y, e.size.w, e.size.h, obj2.point.x, obj2.point.y, obj2.size.w, obj2.size.h)
}

func (e *Entity) rectIntersect(x1 float64, y1 float64, w1 float64, h1 float64, x2 float64, y2 float64, w2 float64, h2 float64) bool {
	if x2 > w1+x1 || x1 > w2+x2 || y2 > h1+y1 || y1 > h2+y2 {
		return false
	}
	return true
}

func (e *Entity) Distance(collider *Entity) float64 {
	distance := CalcDistance(e.center.x, e.center.y, collider.center.x, collider.center.y)
	return distance
}

func (e *Entity) SetupCollision(collider *Entity) {
	e.Collider = collider
	collider.Collider = e

	distance := e.Distance(collider)
	vecCollision := Point{x: collider.center.x - e.center.x, y: collider.center.y - e.center.y}
	vecCollisionNorm := Point{x: vecCollision.x / distance, y: vecCollision.y / distance}

	speed := 0.01
	//if e.breaker != collider {
	//	e.breaker = collider
	//	collider.breaker = e
	e.Friction = friction
	e.GForce = 0.0
	relVx := math.Abs(e.Vx - collider.Vx)
	relVy := math.Abs(e.Vy - collider.Vy)
	speed = relVx*math.Abs(vecCollisionNorm.x) + relVy*math.Abs(vecCollisionNorm.y)
	//}

	impulse := 2 * speed / (e.Mass + collider.Mass)

	e.Vx -= impulse * collider.Mass * vecCollisionNorm.x
	e.Vy -= impulse * collider.Mass * vecCollisionNorm.y
	collider.Vx += impulse * e.Mass * vecCollisionNorm.x
	collider.Vy += impulse * e.Mass * vecCollisionNorm.y
}

func (e *Entity) SetupInelasticCollision(collider *Entity) {
	e.Collider = collider
	//e.breaker = collider
	if collider != nil {
		collider.Collider = e
		//	collider.breaker.Collider = e
	}
	e.Friction = 0.7
	e.Vx = -e.Vx
	e.Vy = -e.Vy
	e.G = 0.0
}

func (e *Entity) isMoving() bool {
	if e.Vx == 0 && e.Vy == 0 {
		return false
	}
	return true
}

func (e *Entity) hit(collider *Entity) bool {
	if collider == nil {
		return false
	}
	distance := CalcDistance(e.center.x, e.center.y, collider.center.x, collider.center.y)
	if distance > collider.GetWidth() {
		return false
	}
	return true
}

func (e *Entity) clearCollider() {
	if e.Collider != nil {
		if e.Collider.Collider == e {
			e.Collider.Collider = nil
		}
		e.Collider = nil
	}
}

/*
func (e * Entity) clearBreaker() {
	if e.breaker != nil {
		if e.breaker.breaker != e {
			e.breaker.breaker = nil
		}
		e.breaker = nil
	}
}
*/

func (e *Entity) Compute() bool {
	if e.Collider != nil {
		if distance := e.Distance(e.Collider); distance >= e.GetWidth()/2+e.Collider.GetWidth()/2 {
			e.clearCollider()
		}
	}
	//if e.breaker != nil {
	//	if !e.hit(e.breaker) {
	//		e.clearBreaker()
	//	}
	//}
	if !e.isMoving() {
		e.G = 0.0
		return false
	}
	e.Vx *= e.Friction
	e.Vy *= e.Friction
	if math.Abs(e.Vx) < e.VxMin {
		e.Vx = 0.0
	}
	if math.Abs(e.Vy) < e.VyMin {
		e.Vy = 0.0
	}
	if !e.isMoving() {
		e.G = 0.0
		return false
	}
	e.G = calcG(e)
	return true
}

func (e *Entity) MoveTest() (float64, float64) {
	x := e.point.x + e.Vx
	y := e.point.y + e.Vy
	return x, y
}

func (e *Entity) Move() {
	e.AddTo(e.Vx, e.Vy)
}

/*
func (e *Entity) Hit() bool {
	if e.lastCollider == nil {
		return false
	}
	distance := calcDistance(e.center.x, e.center.y, e.lastCollider.center.x, e.lastCollider.center.y)
	if distance > e.lastCollider.GetWidth() {
		return false
	}
	return true
}
*/

/*
func (e *Entity) HitRect(x float64, y float64, w float64, h float64) bool {
	if e.cb == nil {
		return false
	}
	r := NewRect(x, y, w, h, 1.0)
	distance := calcDistance(r.center.x, r.center.y, e.lastCollider.center.x, e.lastCollider.center.y)
	if distance > e.cb.GetWidth() {
		return false
	}
	return true
}
*/
