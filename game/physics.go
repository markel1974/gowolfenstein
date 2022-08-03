package main

import (
	"github.com/markel1974/gowolfenstein/game/matrix"
	"github.com/markel1974/gowolfenstein/pixels"
	"sort"
)

type Physics struct {
	world      *World
	tree       *matrix.AABBTree
	lights     *Lights
	collisions []IWorldEntity
}

func NewPhysics(world *World) *Physics {
	return &Physics{
		world:      world,
		lights:     world.lights,
		collisions: make([]IWorldEntity, 1024),
		tree:       matrix.NewAABBTree(uint(world.width * world.height)),
	}
}

func (w *Physics) Update() []IWorldEntity {
	nodes := w.tree.Nodes()
	entities := make([]IWorldEntity, 0, len(nodes))

	if len(nodes) > len(w.collisions) {
		w.collisions = make([]IWorldEntity, len(nodes))
	}

	idx := 0
	for z := range nodes {
		target := z.(IWorldEntity)
		if !target.GetEntity().Compute() {
			target.Halt()
			continue
		}
		if w.doApplyMove(target) {
			w.collisions[idx] = target
			idx++
		}
	}

	if idx > 0 {
		cb := make(map[string]bool)
		for x := 0; x < idx; x++ {
			sp := w.collisions[x]
			if w.doSetupElasticCollision(cb, sp) {
				w.doApplyMove(sp)
			}
		}
	}

	player := w.world.player
	for z := range nodes {
		target := z.(IWorldEntity)
		sp := target.GetEntity()
		lightRadius, lightIntensity := target.GetLight()
		if lightRadius > 0 {
			w.lights.Set(int(sp.GetX()), int(sp.GetY()), lightRadius, lightIntensity)
		}
		dist := (player.GetX()-sp.GetX())*(player.GetX()-sp.GetX()) + (player.GetY()-sp.GetY())*(player.GetY()-sp.GetY())
		target.SetDistance(dist)
		target.Update()
		if target.Removable() {
			w.RemoveNode(target)

			continue
		}
		if sp != player.Entity {
			if w.world.GetVisibility(int(sp.GetX()), int(sp.GetY())) {
				w.world.AddVisibility(int(sp.GetX()), int(sp.GetY()), target.GetType())
				if target.GetDistance() < spriteMinDistance {
					entities = append(entities, target)
				}
			}
		}
	}

	if len(entities) > 0 {
		sort.SliceStable(entities, func(i, j int) bool { return entities[i].GetDistance() > entities[j].GetDistance() })
	}

	return entities
}

func (w *Physics) doApplyMove(target IWorldEntity) bool {
	ret := false
	entity := target.GetEntity()
	x, y := entity.MoveTest()
	if e := w.world.HitWallEntity(x, y, entity.Vx, entity.Vy, target); e != nil {
		//if  e := w.world.HitWallEntity2(x, y); e != nil{
		target.Collision(e)
		entity.SetupInelasticCollision(nil)
	} else {
		ret = true
	}
	entity.Move()
	w.tree.UpdateObject(target)
	return ret
}

func (w *Physics) doSetupElasticCollision(cb map[string]bool, collider IWorldEntity) bool {
	q := w.tree.QueryOverlaps(collider)
	ret := false

	if len(q) > 0 {
		colliderEntity := collider.GetEntity()

		for _, z := range q {
			target := z.(IWorldEntity)
			targetEntity := target.GetEntity()

			collisionId := targetEntity.Id + "|" + colliderEntity.Id
			collisionReverseId := colliderEntity.Id + "|" + targetEntity.Id
			if _, ok := cb[collisionId]; ok {
				continue
			}
			cb[collisionId] = true
			cb[collisionReverseId] = true

			if colliderEntity == targetEntity {
				continue
			}
			if targetEntity.GetWidth() == 0 && targetEntity.GetWidth() == 0 {
				continue
			}

			collider.Collision(target)
			target.Collision(collider)
			targetEntity.SetupCollision(colliderEntity)

			ret = true
		}
	}
	return ret
}

func (w *Physics) DoPhysicalCollision(src IWorldEntity, dir pixels.Vec, speed float64) {
	entity := src.GetEntity()
	position := matrix.NewEntity(entity.GetX()+dir.X, entity.GetY()+dir.Y, entity.GetWidth(), entity.GetHeight(), entity.Mass)
	q := w.tree.QueryOverlaps(position)
	if len(q) > 0 {
		for _, z := range q {
			target := z.(IWorldEntity)
			var targetEntity = target.GetEntity()
			if targetEntity == entity {
				continue
			}
			if targetEntity.GetWidth() == 0 && targetEntity.GetWidth() == 0 {
				continue
			}
			target.Collision(src)

			entity.Vx = dir.X * speed
			entity.Vy = dir.Y * speed
			entity.Friction = 0.99

			targetEntity.SetupCollision(entity)
			targetEntity.GForce = 100
			targetEntity.G = targetEntity.GForce
		}
	}
}

func (w *Physics) AddNode(entity IWorldEntity) {
	w.tree.InsertObject(entity)
}

func (w *Physics) UpdateNode(entity IWorldEntity) {
	w.tree.UpdateObject(entity)
}

func (w *Physics) RemoveNode(entity IWorldEntity) {
	entity.GetEntity().Invalidate()
	w.tree.RemoveObject(entity)
}

func (w *Physics) Nodes() map[matrix.IAABB]uint {
	return w.tree.Nodes()
}

func (w *Physics) Query(entity *matrix.Entity) []matrix.IAABB {
	q := w.tree.QueryOverlaps(entity)
	return q
}
