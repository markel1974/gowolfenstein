package main

import (
	"fmt"
	"github.com/markel1974/gowolfenstein/game/matrix"
	"github.com/markel1974/gowolfenstein/pixels"
	"math"
)

const pV = 0.01

type Player struct {
	world          *World
	dir            pixels.Vec
	plane          pixels.Vec
	distance       float64
	lightRadius    int
	lightIntensity float64
	weapon         *Weapon

	hitCount int
	*matrix.Entity
}

func NewPlayer(world *World, x float64, y float64, w float64, h float64, mass float64) *Player {
	s := &Player{
		world:    world,
		dir:      pixels.MakeVec(-1.0, 0.0),
		plane:    pixels.MakeVec(0.0, 0.66),
		distance: 0,
		Entity:   matrix.NewEntity(x, y, w, h, mass),
		hitCount: 0,
	}
	return s
}

func (as *Player) GetType() string {
	return "player"
}

func (as *Player) Update() {
	if as.hitCount > 0 {
		as.hitCount--
		as.SetLight(11, 10)
	} else {
		as.SetLight(0, 0)
	}
}

func (as *Player) GetEntity() *matrix.Entity {
	return as.Entity
}

func (as *Player) GetTextureId() uint {
	return 0
}

func (as *Player) GetDistance() float64 {
	return as.distance
}

func (as *Player) SetDistance(distance float64) {
	as.distance = distance
}

func (as *Player) GetLight() (int, float64) {
	return as.lightRadius, as.lightIntensity
}

func (as *Player) SetLight(radius int, intensity float64) {
	as.lightRadius = radius
	as.lightIntensity = intensity
}

func (as *Player) Open() {
	x := as.dir.X
	y := as.dir.Y
	targetX := as.GetX()
	targetY := as.GetY()
	found := false
	if as.world.HitWall(int(as.GetX()+x), int(as.GetY())) {
		targetX = as.GetX() + x
		found = true
	}
	if as.world.HitWall(int(as.GetX()), int(as.GetY()+y)) {
		targetY = as.GetY() + y
		found = true
	}
	if found {
		as.world.Unset(int(targetX), int(targetY))
	}
}

func (as *Player) Halt() {
}

func (as *Player) Collision(collider IWorldEntity) {
	fmt.Println("player hit from", collider.GetType())
	if collider.GetType() != "bullet" {
		return
	}
	bullet := collider.(*Bullet)
	if bullet.Owner() == as.Entity.Id {
		return
	}

	as.hitCount = 10
	fmt.Println("player hit from", bullet.Owner())

	//as.SetLight(3, 4)
}

func (as *Player) Removable() bool {
	return false
}

func (as *Player) update(moved bool) {
	if moved {
		as.Friction = 0.9
		as.world.UpdateNode(as)
	} else {
		as.Friction = 0.5
	}
}

func (as *Player) mouseMove(mouseX float64, mouseY float64) {
	mouseX *= 0.01
	mouseY *= 0.01
	const factor = 1
	if mouseX > factor {
		mouseX = factor
	} else if mouseX < -factor {
		mouseX = -factor
	}
	if mouseY > factor {
		mouseY = factor
	} else if mouseY < -factor {
		mouseY = -factor
	}

	as.turnRight(mouseX)
}

func (as *Player) moveForward(s float64) {
	moved := false
	x := as.dir.X * s
	y := as.dir.Y * s

	/*
		if as.Collider != nil {
			if as.Collider.Intersect(as.GetX() + x, as.GetY() + y, as.GetWidth(), as.GetWidth()) {
				as.Vx = pV
				as.Vx = pV
				as.Friction = 0.0
				as.world.UpdateNode(as)
				return
			}
		}
	*/

	if as.world.HitWallEntity(as.GetX()+x, as.GetY(), as.dir.X, as.dir.Y, as) == nil {
		as.AddToX(x)
		as.Vx = pV
		moved = true
	}
	if as.world.HitWallEntity(as.GetX(), as.GetY()+y, as.dir.X, as.dir.Y, as) == nil {
		as.AddToY(y)
		as.Vy = pV
		moved = true
	}
	as.update(moved)

	/*
		moved := false
		if as.world.GetMediumWallDistance() >= 0.3 {
			x := as.dir.X * s
			y := as.dir.Y * s
			if !as.world.HitWall(int(as.GetX()+x), int(as.GetY())) {
				as.AddToX(x)
				as.Vx = pV
				moved = true
			}
			if !as.world.HitWall(int(as.GetX()), int(as.GetY()+y)) {
				as.AddToY(y)
				as.Vy = pV
				moved = true
			}
		}
		as.update(moved)
	*/
}

func (as *Player) moveLeft(s float64) {
	x := as.plane.X * s
	y := as.plane.Y * s
	moved := false
	if as.world.HitWallEntity(as.GetX()-x, as.GetY(), -as.plane.X, -as.plane.Y, as) == nil {
		as.AddToX(-x)
		as.Vx = -pV
		moved = true
	}
	if as.world.HitWallEntity(as.GetX(), as.GetY()-y, -as.plane.X, -as.plane.Y, as) == nil {
		as.AddToY(-y)
		as.Vy = -pV
		moved = true
	}
	as.update(moved)

	/*
		x := as.plane.X * s
		y := as.plane.Y * s
		moved := false
		if as.world.HitWallEntity(as.GetX() -x, as.GetY() -y, -as.plane.X, -as.plane.Y, as) == nil {
			as.AddToX(-x)
			as.AddToY(-y)
			as.Vx = -pV
			as.Vy = -pV
			moved = true
		}
		as.update(moved)
	*/

	/*
		x := as.plane.X * s
		y := as.plane.Y * s
		moved := false
		if !as.world.HitWall(int(as.GetX() - x), int(as.GetY())) {
			as.AddToX(-x)
			as.Vx = -pV
			moved = true
		}
		if !as.world.HitWall(int(as.GetX()), int(as.GetY() - y)) {
			as.AddToY(-y)
			as.Vy = -pV
			moved = true
		}
		as.update(moved)
	*/
}

func (as *Player) moveBackwards(s float64) {
	moved := true
	x := as.dir.X * s
	y := as.dir.Y * s
	if as.world.HitWallEntity(as.GetX()-x, as.GetY(), -as.dir.X, -as.dir.Y, as) == nil {
		as.AddToX(-x)
		as.Vx = -pV
		moved = true
	}
	if as.world.HitWallEntity(as.GetX(), as.GetY()-y, -as.dir.X, -as.dir.Y, as) == nil {
		as.AddToY(-y)
		as.Vy = -pV
		moved = true
	}
	as.update(moved)
	/*
		x := as.dir.X * s
		y := as.dir.Y * s
		moved := false
		if !as.world.HitWall(int(as.GetX() - x), int(as.GetY())) {
			as.AddToX(-x)
			as.Vx = -pV
			moved = true
		}
		if !as.world.HitWall(int(as.GetX()), int(as.GetY() - y)) {
			as.AddToY(-y)
			as.Vy = -pV
			moved = true
		}
		as.update(moved)
	*/
}

func (as *Player) moveRight(s float64) {
	x := as.plane.X * s
	y := as.plane.Y * s
	moved := false
	if as.world.HitWallEntity(as.GetX()+x, as.GetY(), as.plane.X, as.plane.Y, as) == nil {
		as.AddToX(x)
		as.Vx = pV
		moved = true
	}
	if as.world.HitWallEntity(as.GetX(), as.GetY()+y, as.plane.X, as.plane.Y, as) == nil {
		as.AddToY(y)
		as.Vy = pV
		moved = true
	}
	as.update(moved)

	/*
		x := as.plane.X * s
		y := as.plane.Y * s
		moved := false
		if as.world.HitWallEntity(as.GetX() +x, as.GetY() +y, as.plane.X, as.plane.Y, as) == nil {
			as.AddToX(x)
			as.AddToY(y)
			as.Vx = pV
			as.Vy = pV
			moved = true
		}
		as.update(moved)
	*/

	/*
		x := as.plane.X * s
		y := as.plane.Y * s
		moved := false
		if !as.world.HitWall(int(as.GetX() + x), int(as.GetY())) {
			as.AddToX(x)
			as.Vx = pV
			moved = true
		}
		if !as.world.HitWall(int(as.GetX()), int(as.GetY() + y)) {
			as.AddToY(y)
			as.Vy = pV
			moved = true
		}
		as.update(moved)
	*/
}

func (as *Player) turnRight(s float64) {
	cosS := math.Cos(-s)
	sinS := math.Sin(-s)

	oldDirX := as.dir.X
	as.dir.X = as.dir.X*cosS - as.dir.Y*sinS
	as.dir.Y = oldDirX*sinS + as.dir.Y*cosS
	oldPlaneX := as.plane.X
	as.plane.X = as.plane.X*cosS - as.plane.Y*sinS
	as.plane.Y = oldPlaneX*sinS + as.plane.Y*cosS
	as.Vy = pV

	as.update(true)
}

func (as *Player) turnLeft(s float64) {
	cosS := math.Cos(s)
	sinS := math.Sin(s)

	oldDirX := as.dir.X
	as.dir.X = as.dir.X*cosS - as.dir.Y*sinS
	as.dir.Y = oldDirX*sinS + as.dir.Y*cosS
	oldPlaneX := as.plane.X
	as.plane.X = as.plane.X*cosS - as.plane.Y*sinS
	as.plane.Y = oldPlaneX*sinS + as.plane.Y*cosS
	as.Vy = pV
	as.update(true)
}

func (as *Player) SetWeapon(weapon *Weapon) {
	as.weapon = weapon
	as.weapon.SetOwner(as.Entity.Id)
}

func (as *Player) Fire() {
	if as.weapon == nil {
		return
	}
	as.weapon.Shoot(as.GetX(), as.GetY(), as.dir)
}

func (as *Player) Kick() {
	mass := 1.0
	speed := 1.0
	f := NewPhysical(as.world, as.Entity.Id, "kick", as.GetX(), as.GetY(), 0.5, 0.5, mass)
	as.world.PhysicalCollision(f, as.dir, speed)
}
