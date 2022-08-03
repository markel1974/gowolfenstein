package main

import "github.com/markel1974/gowolfenstein/game/matrix"

type Furnishing struct {
	world          *World
	distance       float64
	textureId      uint
	lightRadius    int
	lightIntensity float64
	*matrix.Entity
}

func NewFurnishing(world *World, textureId uint, x float64, y float64, w float64, h float64, mass float64) *Furnishing {
	s := &Furnishing{
		Entity:    matrix.NewEntity(x, y, w, h, mass),
		world:     world,
		distance:  0,
		textureId: textureId,
	}
	return s
}

func (s *Furnishing) GetType() string {
	return "furnishing"
}

func (s *Furnishing) Update() {
}

func (s *Furnishing) GetEntity() *matrix.Entity {
	return s.Entity
}

func (s *Furnishing) GetTextureId() uint {
	return s.textureId
}

func (s *Furnishing) GetDistance() float64 {
	return s.distance
}

func (s *Furnishing) SetDistance(distance float64) {
	s.distance = distance
}

func (s *Furnishing) GetLight() (int, float64) {
	return s.lightRadius, s.lightIntensity
}

func (s *Furnishing) SetLight(radius int, intensity float64) {
	s.lightRadius = radius
	s.lightIntensity = intensity
}

func (s *Furnishing) Halt() {
}

func (s *Furnishing) Collision(collider IWorldEntity) {
}

func (s *Furnishing) Removable() bool {
	return false
}
