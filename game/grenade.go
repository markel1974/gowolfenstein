package main

import (
	"github.com/markel1974/gowolfenstein/game/matrix"
)

type Grenade struct {
	world          *World
	owner          string
	distance       float64
	textureId      uint
	lightRadius    int
	lightIntensity float64
	destroy        int
	*matrix.Entity
}

func NewGrenade(world *World, owner string, textureId uint, x float64, y float64, w float64, h float64, mass float64) *Grenade {
	s := &Grenade{
		Entity:    matrix.NewEntity(x, y, w, h, mass),
		owner:     owner,
		world:     world,
		destroy:   3,
		distance:  0,
		textureId: textureId,
	}
	return s
}

func (s *Grenade) GetType() string {
	return "grenade"
}

func (s *Grenade) Update() {
}

func (s *Grenade) GetEntity() *matrix.Entity {
	return s.Entity
}

func (s *Grenade) GetTextureId() uint {
	return s.textureId
}

func (s *Grenade) GetDistance() float64 {
	return s.distance
}

func (s *Grenade) SetDistance(distance float64) {
	s.distance = distance
}

func (s *Grenade) GetLight() (int, float64) {
	return s.lightRadius, s.lightIntensity
}

func (s *Grenade) SetLight(radius int, intensity float64) {
	s.lightRadius = radius
	s.lightIntensity = intensity
}

func (s *Grenade) Halt() bool {
	return true
}

func (s *Grenade) Collision(collider IWorldEntity) bool {
	if collider.GetType() == "wall" {
		return true
	}
	return s.doDestroy()
}

func (s *Grenade) Owner() string {
	return s.owner
}

func (s *Grenade) doDestroy() bool {
	if s.destroy > 0 {
		s.destroy--
		if s.destroy == 0 {
			return true
		}
	}
	return false
}
