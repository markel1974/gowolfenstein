package main

import "github.com/markel1974/gowolfenstein/game/matrix"

type IWorldEntity interface {
	GetType() string
	GetEntity() *matrix.Entity
	GetAABB() *matrix.AABB
	GetDistance() float64
	SetDistance(float64)
	GetTextureId() uint
	GetLight() (int, float64)
	SetLight(radius int, intensity float64)
	Halt()
	Collision(collider IWorldEntity)
	Update()
	Removable() bool
}
