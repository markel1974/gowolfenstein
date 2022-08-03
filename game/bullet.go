package main

import (
	"github.com/markel1974/gowolfenstein/game/matrix"
)

type Bullet struct {
	world          *World
	owner          string
	distance       float64
	textureId      uint
	lightRadius    int
	lightIntensity float64
	hit            bool
	counter        int
	*matrix.Entity
}

func NewBullet(world *World, owner string, textureId uint, x float64, y float64, w float64, h float64, mass float64) *Bullet {
	s := &Bullet{
		Entity:    matrix.NewEntity(x, y, w, h, mass),
		owner:     owner,
		world:     world,
		distance:  0,
		hit:       false,
		counter:   3,
		textureId: textureId,
	}

	s.SetLight(1, 1.5)
	return s
}

func (s *Bullet) GetType() string {
	return "bullet"
}

func (s *Bullet) Update() {
	if s.hit {
		if s.counter > 0 {
			//if s.counter == 1 {
			//
			s.counter--
		}
	}
}

func (s *Bullet) GetEntity() *matrix.Entity {
	return s.Entity
}

func (s *Bullet) GetTextureId() uint {
	return s.textureId
}

func (s *Bullet) GetDistance() float64 {
	return s.distance
}

func (s *Bullet) SetDistance(distance float64) {
	s.distance = distance
}

func (s *Bullet) GetLight() (int, float64) {
	return s.lightRadius, s.lightIntensity
}

func (s *Bullet) SetLight(radius int, intensity float64) {
	s.lightRadius = radius
	s.lightIntensity = intensity
	//fmt.Println(s.lightIntensity)
}

func (s *Bullet) Halt() {
	s.hit = true
	s.textureId = 126
}

func (s *Bullet) Collision(collider IWorldEntity) {
	if collider.GetEntity().Id == s.owner {
		s.Vx = 0
		s.Vy = 0
		return
	}
	if collider.GetType() == "wall" {
		s.Vx = 0
		s.Vy = 0
	} else {
		s.Vx *= 0.30
		s.Vy *= 0.30
	}
	s.hit = true
	s.textureId = 126
	s.SetLight(3, 2.5)
}

func (s *Bullet) Removable() bool {
	return s.counter == 0
}

func (s *Bullet) Owner() string {
	return s.owner
}
