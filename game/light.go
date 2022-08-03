package main

import (
	"github.com/markel1974/gowolfenstein/game/matrix"
	"math/rand"
	"time"
)

type Light struct {
	world          *World
	distance       float64
	textureId      uint
	lightRadius    int
	lightIntensity float64
	faultNext      int64
	fault          bool
	*matrix.Entity
}

func NewLight(world *World, textureId uint, x float64, y float64, radius int, intensity float64) *Light {
	s := &Light{
		Entity:    matrix.NewEntity(x, y, 0, 0, 0),
		world:     world,
		distance:  0,
		textureId: textureId,
		fault:     false,
		faultNext: 0,
	}
	s.SetLight(radius, intensity)
	return s
}

func (s *Light) SetFault(fault bool) {
	s.fault = fault
}

func (s *Light) GetType() string {
	return "light"
}

func (s *Light) Update() {
	if s.fault {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		if now > s.faultNext {
			if s.lightRadius == 0 {
				s.lightRadius = 3
				s.lightIntensity = 4
			} else {
				s.lightRadius = 0
				s.lightIntensity = 0
			}
			offset := rand.Intn(1000-100) + 100
			s.faultNext = now + int64(offset)
		}
	}
}

func (s *Light) GetEntity() *matrix.Entity {
	return s.Entity
}

func (s *Light) GetTextureId() uint {
	return s.textureId
}

func (s *Light) GetDistance() float64 {
	return s.distance
}

func (s *Light) SetDistance(distance float64) {
	s.distance = distance
}

func (s *Light) GetLight() (int, float64) {
	return s.lightRadius, s.lightIntensity
}

func (s *Light) SetLight(radius int, intensity float64) {
	s.lightRadius = radius
	s.lightIntensity = intensity
}

func (s *Light) Halt() {
}

func (s *Light) Collision(collider IWorldEntity) {
}

func (s *Light) Removable() bool {
	return false
}
