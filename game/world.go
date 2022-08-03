package main

import (
	"encoding/json"
	"fmt"
	"github.com/markel1974/gowolfenstein/game/matrix"
	"github.com/markel1974/gowolfenstein/pixels"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"
)

type WorldObject struct {
	ID        int     `json:"id"`
	X         int     `json:"x"`
	Y         float64 `json:"y"`
	Z         int     `json:"z"`
	Rotation  int     `json:"rotation,omitempty"`
	RW        float64 `json:"rw,omitempty"`
	Texture   int     `json:"texture,omitempty"`
	Collision int     `json:"collision,omitempty"`
	Type      int     `json:"type,omitempty"`
	Amount    int     `json:"amount,omitempty"`
	Cx        int     `json:"cx,omitempty"`
	Cy        int     `json:"cy,omitempty"`
	Cz        int     `json:"cz,omitempty"`
}

type WorldConfig struct {
	Data    [][]int        `json:"data"`
	Objects []*WorldObject `json:"objects"`
}

type World struct {
	d             [][]IWorldEntity
	screenWidth   int
	screenHeight  int
	textures      *Textures
	width         int
	height        int
	visibility    *Visibility
	wallsDistance []float64
	wallDistance  float64
	wallDefault   *Wall
	config        WorldConfig

	mainSurface *pixels.PictureRGBA

	player  *Player
	enemies []*Enemy
	lights  *Lights
	physics *Physics
	sprites *Sprites
}

func NewWorld(screenWidth int, screenHeight int, textures *Textures, mainSurface *pixels.PictureRGBA) *World {
	rand.Seed(time.Now().UTC().UnixNano())

	w := &World{
		textures:      textures,
		mainSurface:   mainSurface,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
		wallsDistance: make([]float64, screenWidth),
		wallDistance:  8.0,
	}
	_ = w.load()
	w.wallDefault = NewWall(w, 0, 0, 0, 1.0, 1.0, 1000, 1)
	w.width = len(w.d)
	w.height = len(w.d[0])
	w.visibility = NewVisibility(w.width, w.height)
	w.lights = NewLights(w.width, w.height)
	w.physics = NewPhysics(w)
	w.sprites = NewSprites(w, w.textures, w.lights, w.mainSurface, w.screenWidth, w.screenHeight)

	//TODO PLAYER TEST
	w.player = NewPlayer(w, 60.0, 34.0, 0.5, 0.5, 80.0)
	w.AddNode(w.player)

	//TODO ENEMY TEST
	for e := 0; e < 512; e++ {
		x := float64(rand.Intn(w.width))
		y := float64(rand.Intn(w.height))
		enemy := NewEnemy(w, x, y)
		enemyWeapon := NewWeapon(w, 124, w.textures.RatioWidth(156), w.textures.RatioWidth(156), 1, 1.0, 50)

		//enemyWeapon := NewWeapon(w, 123, w.textures.RatioWidth(120), w.textures.RatioWidth(120), 2, 0.4, 100)
		enemy.SetWeapon(enemyWeapon)
		w.enemies = append(w.enemies, enemy)
		w.AddNode(enemy)
	}

	//TODO LIGHT TEST
	for e := 0; e < 128; e++ {
		x := float64(rand.Intn(w.width))
		y := float64(rand.Intn(w.height))
		radius := 2 + rand.Intn(8-2)
		intensity := 1.5 + rand.Float64()*(2.5-1.5)
		light := NewLight(w, 150, x+0.5, y+0.5, radius, intensity)
		if fault := float64(rand.Intn(5)); fault >= 4 {
			light.SetFault(true)
		}
		w.AddNode(light)
	}

	for idx, obj := range w.config.Objects {
		textureId := uint(obj.Texture + 121)
		diameter := obj.RW
		if diameter <= 0 {
			diameter = w.textures.RatioWidth(textureId)
		}

		//TODO TEST RIMUOVERE!
		diameter = 0.5
		fmt.Println("diameter: ", idx, diameter)

		//offset := center
		//if obj.Y > 0 {
		//	offset = obj.Y
		//}
		cx := float64(obj.X) + (diameter / 2)
		cy := float64(obj.Z) + (diameter / 2)

		//mass :=  0.6 + rand.Float64() * (2 - 0.6)
		//TODO RIMUOVERE
		mass := 30 + rand.Float64()*(45-30)
		//mass = 5
		//mass := 0.7
		f := NewFurnishing(w, textureId, cx, cy, diameter, diameter, mass)

		//TODO FROM CONFIG
		//if textureId == 120 || textureId == 129 || textureId == 156 {
		//	f.SetLight(1, 0xc0)
		//}

		w.AddNode(f)
	}

	return w
}

func (w *World) load() error {
	data, err := ioutil.ReadFile("resources" + string(os.PathSeparator) + "levels" + string(os.PathSeparator) + "1.json")
	if err != nil {
		return nil
	}

	if err := json.Unmarshal(data, &w.config); err != nil {
		return err
	}
	w.d = make([][]IWorldEntity, len(w.config.Data))
	for x, line := range w.config.Data {
		w.d[x] = make([]IWorldEntity, len(line))
		for y, c := range line {
			if c > 0 {
				w.d[x][y] = NewWall(w, uint(c), float64(x), float64(y), 1.0, 1.0, 1000, 1)
			}
		}
	}

	return nil
}

func (w *World) GetWorld() [][]IWorldEntity {
	return w.d
}

func (w *World) GetWidth() int {
	return w.width
}

func (w *World) GetHeight() int {
	return w.height
}

func (w *World) GetWallTexture(x int, y int) uint {
	if !w.valid(x, y) {
		return 0
	}
	h := w.d[x][y]
	if h == nil {
		return 0
	}
	return h.GetTextureId()
}

/*
func (w * World) Set(x int, y int, object *WorldObject) {
	if object == nil {
		return
	}
	if !w.valid(x, y) {
		return
	}
	w.d[x][y] = object
}
*/

func (w *World) Unset(x int, y int) {
	if !w.valid(x, y) {
		return
	}
	t := w.d[x][y]
	if t == nil {
		return
	}

	//TODO DOOR DEFINITION!!!!!
	//TODO WALL ANIMATION
	if t.GetTextureId() >= 98 && t.GetTextureId() <= 109 {
		w.d[x][y] = nil
	}

	fmt.Println("CAN'T OPEN: ", t.GetTextureId())
}

func (w *World) AddVisibility(x int, y int, kind string) {
	w.visibility.Add(x, y, kind)
}

func (w *World) GetVisibility(x int, y int) bool {
	return w.visibility.Get(x, y)
}

func (w *World) UpdateAvailable() {
	w.visibility.Update()
}

func (w *World) GetWorldAvailable() (int, [][]*VisibilityData) {
	return w.visibility.GetId(), w.visibility.GetWorld()
}

func (w *World) SetWallDistance(x int, v float64) {
	w.wallsDistance[x] = v
	if x == w.screenWidth/2 {
		w.wallDistance = v
	}
}

func (w *World) GetWallDistance(x int) float64 {
	return w.wallsDistance[x]
}

func (w *World) GetMediumWallDistance() float64 {
	return w.wallDistance
}

func (w *World) HitWall(x int, y int) bool {
	if !w.valid(x, y) {
		return true
	}
	h := w.d[x][y]
	if h == nil {
		return false
	}
	return true
}

func (w *World) HitWallEntity(nextX float64, nextY float64, nextVx float64, nextVy float64, entity IWorldEntity) IWorldEntity {
	sp := entity.GetEntity()

	if math.Signbit(nextVx) {
		nextX -= sp.GetWidth()
	} else {
		nextX += sp.GetWidth()
	}
	if math.Signbit(nextVy) {
		nextY -= sp.GetWidth()
	} else {
		nextY += sp.GetWidth()
	}
	if !w.valid(int(nextX), int(nextY)) {
		return w.wallDefault
	}

	return w.d[int(nextX)][int(nextY)]
}

/*
func (w * World) HitWallEntity(x float64, y float64) IWorldEntity {
	if !w.valid(x, y) {
		return nil
	}
	h := w.d[x][y]
	if h == nil {
		return nil
	}
	return w.d[x][y]
}
*/

func (w *World) valid(x int, y int) bool {
	if x >= 0 && y >= 0 && x < w.width && y < w.height {
		return true
	}
	return false
}

func (w *World) UpdateWorld() {
	w.lights.Update()
	entities := w.physics.Update()
	w.sprites.Update(entities)
}

func (w *World) AddNode(entity IWorldEntity) {
	w.physics.AddNode(entity)
}

func (w *World) UpdateNode(entity IWorldEntity) {
	w.physics.UpdateNode(entity)
}

func (w *World) Nodes() map[matrix.IAABB]uint {
	return w.physics.Nodes()
}

func (w *World) Query(entity *matrix.Entity) []matrix.IAABB {
	return w.physics.Query(entity)
}

func (w *World) PhysicalCollision(src IWorldEntity, dir pixels.Vec, speed float64) {
	w.physics.DoPhysicalCollision(src, dir, speed)
}
