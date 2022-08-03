package main

import "github.com/markel1974/gowolfenstein/game/matrix"

type Wall struct {
	world          *World
	distance       float64
	textureId      uint
	collision      int
	lightRadius    int
	lightIntensity float64
	*matrix.Entity
}

func NewWall(world *World, textureId uint, x float64, y float64, w float64, h float64, mass float64, collision int) *Wall {
	s := &Wall{
		Entity:    matrix.NewEntity(x, y, w, h, mass),
		world:     world,
		distance:  0,
		textureId: textureId,
		collision: collision,
	}
	return s
}

func (s *Wall) GetType() string {
	return "wall"
}

func (s *Wall) Update() {
}

func (s *Wall) GetEntity() *matrix.Entity {
	return s.Entity
}

func (s *Wall) GetTextureId() uint {
	return s.textureId
}

func (s *Wall) GetDistance() float64 {
	return s.distance
}

func (s *Wall) SetDistance(distance float64) {
	s.distance = distance
}

func (s *Wall) GetLight() (int, float64) {
	return s.lightRadius, s.lightIntensity
}

func (s *Wall) SetLight(radius int, intensity float64) {
	s.lightRadius = radius
	s.lightIntensity = intensity
}

func (s *Wall) Halt() {
}

func (s *Wall) Collision(collider IWorldEntity) {
}

func (s *Wall) Removable() bool {
	return false
}
