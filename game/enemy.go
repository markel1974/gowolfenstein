package main

import (
	"fmt"
	"github.com/markel1974/gowolfenstein/game/matrix"
	"github.com/markel1974/gowolfenstein/pixels"
	"math/rand"
)

const (
	maxHP   = 15
	maxGore = 16
)

type EnemyState int

const (
	stateIdle   EnemyState = iota
	stateChase  EnemyState = iota
	stateAttack EnemyState = iota
	stateHurt   EnemyState = iota
	stateDeath  EnemyState = iota
)

type Enemy struct {
	world            *World
	hp               int
	state            EnemyState
	weapon           *Weapon
	attackRange      float64
	distanceToPlayer float64
	sightRange       float64

	sightTimeline  []uint
	damageTimeline []uint
	deadTimeline   []uint
	fireTimeline   []uint
	goreTimeline   []uint

	distance       float64
	textureId      uint
	lightRadius    int
	lightIntensity float64

	directions []float64

	latestX float64
	latestY float64

	gore int

	animator *Animator

	removable bool

	*matrix.Entity
}

func NewEnemy(world *World, posX float64, posY float64) *Enemy {
	//TODO DA COMPLETARE
	//TODO IN CASO DI ATTACCO LA DIREZIONE DEVE ESSERE VERSO IL GIOCATORE

	e := &Enemy{
		world:       world,
		Entity:      matrix.NewEntity(posX, posY, 0.5, 0.5, 50),
		hp:          maxHP,
		weapon:      nil,
		state:       stateIdle,
		attackRange: 3.0,
		sightRange:  16.0,
		distance:    0,
		textureId:   197,
		animator:    NewAnimator(),
	}

	e.Friction = 0.9

	minSpeed := 0.04 + rand.Float64()*(0.06-0.04)
	maxSpeed := 0.06 + rand.Float64()*(0.08-0.06)

	e.directions = append(e.directions, maxSpeed)
	e.directions = append(e.directions, -maxSpeed)
	e.directions = append(e.directions, minSpeed)
	e.directions = append(e.directions, -minSpeed)

	e.latestX = 0
	e.latestY = 0

	e.sightTimeline = append(e.sightTimeline, 188)
	e.sightTimeline = append(e.sightTimeline, 194)
	e.sightTimeline = append(e.sightTimeline, 195)
	e.sightTimeline = append(e.sightTimeline, 196)
	e.sightTimeline = append(e.sightTimeline, 197)

	e.fireTimeline = append(e.fireTimeline, 191)
	e.fireTimeline = append(e.fireTimeline, 192)
	e.fireTimeline = append(e.fireTimeline, 193)

	e.damageTimeline = append(e.damageTimeline, 190)
	e.damageTimeline = append(e.damageTimeline, 189)
	e.damageTimeline = append(e.damageTimeline, 184)

	e.deadTimeline = append(e.deadTimeline, 190)
	e.deadTimeline = append(e.deadTimeline, 189)
	e.deadTimeline = append(e.deadTimeline, 184)
	e.deadTimeline = append(e.deadTimeline, 185)
	e.deadTimeline = append(e.deadTimeline, 186)
	e.deadTimeline = append(e.deadTimeline, 187)

	//e.goreTimeline = append(e.goreTimeline,184)
	e.goreTimeline = append(e.goreTimeline, 185)
	e.goreTimeline = append(e.goreTimeline, 186)
	e.goreTimeline = append(e.goreTimeline, 187)

	return e
}

func (e *Enemy) GetType() string {
	return "enemy"
}

func (e *Enemy) GetEntity() *matrix.Entity {
	return e.Entity
}

func (e *Enemy) GetTextureId() uint {
	return e.textureId
}

func (e *Enemy) GetDistance() float64 {
	return e.distance
}

func (e *Enemy) SetDistance(distance float64) {
	e.distance = distance
}

func (e *Enemy) GetLight() (int, float64) {
	return e.lightRadius, e.lightIntensity
}

func (e *Enemy) SetLight(radius int, intensity float64) {
	e.lightRadius = radius
	e.lightIntensity = intensity
}

func (e *Enemy) Halt() {
}

func (e *Enemy) Collision(collider IWorldEntity) {
	switch collider.GetType() {
	case "bullet":
		bullet := collider.(*Bullet)
		if bullet.Owner() == e.Entity.Id {
			return
		}
		fmt.Println("Enemy hit from", bullet.Owner())
		e.hurt()
		e.damage(3)
	case "physical":
		physical := collider.(*Physical)
		fmt.Println("Enemy physical contact from player", physical.Owner(), physical.Kind())
		e.hurt()
		e.damage(1)
	}
	//return e.State == stateDeath
}

func (e *Enemy) Removable() bool {
	return e.removable
}

func (e *Enemy) checkSightLine() bool {
	return true
	//TODO
	/*
		Ray ray
		ray.m_Origin = GetPosition()
		ray.m_Direction = m_PlayerDirection
		ray.m_InvDirection = 1.0f / m_PlayerDirection
		return GameManager::Get()->GetLevel()->CheckEnemyRayCollision(ray)
	*/
}

func (e *Enemy) Update() {
	e.distanceToPlayer = e.Distance(e.world.player.GetEntity())
	e.faceCamera()

	e.textureId = e.animator.Play()

	switch e.state {
	case stateIdle:
		e.idle()
	case stateChase:
		e.chase()
	case stateAttack:
		e.attack()
	case stateHurt:
		e.hurt()
	case stateDeath:
		e.death()
	}
}

func (e *Enemy) SetWeapon(weapon *Weapon) {
	e.weapon = weapon
	e.weapon.SetOwner(e.Entity.Id)
}

func (e *Enemy) idle() {
	if e.distanceToPlayer < e.sightRange {
		if e.checkSightLine() {
			e.animator.Set(e.sightTimeline, 100, true)
			e.state = stateChase
		}
	}
}

func (e *Enemy) chase() {
	if e.distanceToPlayer <= e.attackRange {
		e.animator.SetStatic(191)
		e.state = stateAttack
		//e.latestX = -e.world.player.dir.X
		//e.latestY = -e.world.player.dir.Y
	} else if e.distanceToPlayer >= e.sightRange {
		e.animator.SetStatic(188)
		e.state = stateIdle
	}

	if e.latestX == 0 {
		e.latestX = e.getDirection()
	}
	if e.latestY == 0 {
		e.latestY = e.getDirection()
	}
	if e.Vx == 0 {
		e.Vx = e.latestX
	}
	if e.Vy == 0 {
		e.Vy = e.latestY
	}
	nextX := e.GetX() + e.Vx
	nextY := e.GetY() + e.Vy

	if e.world.HitWallEntity(nextX, e.GetY(), e.Vx, e.Vy, e) != nil {
		e.latestX = 0
		e.Vx = 0
	}
	if e.world.HitWallEntity(e.GetX(), nextY, e.Vx, e.Vy, e) != nil {
		e.latestY = 0
		e.Vx = 0
	}
	/*
		if e.world.HitWall(int(nextX), int(e.GetY())) {
			e.latestX = 0
			e.Vx = 0
		}
		if e.world.HitWall(int(e.GetX()), int(nextY)) {
			e.latestY = 0
			e.Vx = 0
		}
	*/
}

func (e *Enemy) attack() {
	if e.weapon == nil {
		return
	}
	player := e.world.player

	e.animator.SetStatic(192)

	if e.weapon.Shoot(e.Entity.GetX(), e.Entity.GetY(), pixels.Vec{X: -player.dir.X, Y: -player.dir.Y}) {
		e.animator.SetStatic(193)

		//TODO RIMUOVERE
		if e.distanceToPlayer > e.attackRange {
			e.animator.Set(e.sightTimeline, 100, true)
			e.state = stateChase
		}
		/*
			if e.checkSightline() {
				//e.Shoot()
			} else {
				e.State = stateChase
			}
		*/
	}
}

func (e *Enemy) damage(damage int) {
	e.hp -= damage
	if e.hp > 0 {
		e.animator.Set(e.damageTimeline, 100, false)
		e.state = stateHurt
		f := NewFurnishing(e.world, 174, e.GetX(), e.GetY(), 0, 0, 0)
		e.world.AddNode(f)
		return
	}
	if e.state != stateDeath {
		e.state = stateDeath
		e.animator.Set(e.deadTimeline, 100, false)
		//AudioManager::Get()->PlayEnemyDeath(m_Position)
		//GameManager::Get()->GetLevel()->SpawnAmmo(m_Position)
		return
	}
	if e.gore < maxGore {
		e.animator.Set(e.goreTimeline, 100, false)
		f := NewFurnishing(e.world, 174, e.GetX(), e.GetY(), 0, 0, 0)
		e.world.AddNode(f)
		e.gore++
		return
	}

	f := NewFurnishing(e.world, 170, e.GetX(), e.GetY(), 0, 0, 0)
	e.world.AddNode(f)

	e.removable = true

	//e.animator.SetStatic(170)
}

func (e *Enemy) hurt() {
	if e.state != stateDeath {
		e.state = stateAttack
	}
}

func (e *Enemy) death() {
	//TODO
}

func (e *Enemy) faceCamera() {
	/*
		float camera_angle = -atanf(m_PlayerDirection.z/m_PlayerDirection.x) + (90.0f*glm::pi<float>() / 180.0f);
		if m_PlayerDirection.x > 0 {
			camera_angle += glm::pi<float>()
		}
		SetRotation(glm::vec3(0, camera_angle, 0))
	*/
}

func (e *Enemy) getDirection() float64 {
	var min = 0
	var max = len(e.directions) - 1
	idx := rand.Intn(max-min) + min
	out := e.directions[idx]
	return out
}

/*
func (e * Enemy) checkSightLine() bool {
	Ray ray
	ray.m_Origin = GetPosition()
	ray.m_Direction = m_PlayerDirection
	ray.m_InvDirection = 1.0f / m_PlayerDirection

	return GameManager::Get()->GetLevel()->CheckEnemyRayCollision(ray)
}
*/
