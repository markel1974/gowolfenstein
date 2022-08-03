package main

import (
	"github.com/markel1974/gowolfenstein/pixels"
	"time"
)

type Weapon struct {
	world        *World
	owner        string
	w            float64
	h            float64
	rate         int64
	mass         float64
	speed        float64
	lastInterval int64
	textureId    uint
}

func NewWeapon(world *World, textureId uint, w float64, h float64, mass float64, speed float64, rate int64) *Weapon {
	return &Weapon{
		world:        world,
		textureId:    textureId,
		mass:         mass,
		w:            w,
		h:            h,
		speed:        speed,
		rate:         (60 * 1000) / rate,
		lastInterval: 0,
	}
}

func (w *Weapon) SetOwner(owner string) {
	w.owner = owner
}

func (w *Weapon) Shoot(posX float64, posY float64, dir pixels.Vec) bool {
	epoch := time.Now().UnixNano() / int64(time.Millisecond)
	lastFire := epoch - w.lastInterval //TimeManager::GetTime() - w.lastInterval
	if lastFire < w.rate {
		return false
	}

	//TODO Weapon Animation
	//TODO Play Weapon Sound
	w.lastInterval = epoch
	bullet := NewBullet(w.world, w.owner, w.textureId, posX, posY, w.w, w.h, w.mass)
	//bullet.SetLight(3, 2.0)
	//bullet.SetLight(w.lightRadius, w.lightIntensity)
	bullet.Vx = dir.X * w.speed
	bullet.Vy = dir.Y * w.speed
	bullet.Friction = 0.99
	bullet.VxMin = 0.1
	bullet.VyMin = 0.1

	w.world.AddNode(bullet)

	return true
}
