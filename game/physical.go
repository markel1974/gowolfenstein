package main

import (
	"github.com/markel1974/gowolfenstein/game/matrix"
)

type Physical struct {
	world    *World
	owner    string
	kind     string
	distance float64
	*matrix.Entity
}

func NewPhysical(world *World, owner string, kind string, x float64, y float64, w float64, h float64, mass float64) *Physical {
	s := &Physical{
		Entity:   matrix.NewEntity(x, y, w, h, mass),
		owner:    owner,
		kind:     kind,
		world:    world,
		distance: 0,
	}
	return s
}

func (s *Physical) GetType() string {
	return "physical"
}

func (s *Physical) Update() {
}

func (s *Physical) GetEntity() *matrix.Entity {
	return s.Entity
}

func (s *Physical) GetTextureId() uint {
	return 0
}

func (s *Physical) GetDistance() float64 {
	return s.distance
}

func (s *Physical) SetDistance(distance float64) {
	s.distance = distance
}

func (s *Physical) GetLight() (int, float64) {
	return 0, 0
}

func (s *Physical) SetLight(_ int, _ float64) {
}

func (s *Physical) Halt() {
}

func (s *Physical) Collision(_ IWorldEntity) {
}

func (s *Physical) Owner() string {
	return s.owner
}

func (s *Physical) Kind() string {
	return s.kind
}

func (s *Physical) Removable() bool {
	return false
}
