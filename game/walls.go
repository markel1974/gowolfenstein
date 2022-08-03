package main

import (
	"github.com/markel1974/gowolfenstein/pixels"
	"math"
)

type Walls struct {
	world        *World
	screenWidth  int
	screenHeight int
	textures     *Textures
	camera       []float64

	dirX         float64
	dirY         float64
	planeX       float64
	planeY       float64
	x            int
	cameraX      float64
	posX         float64
	posY         float64
	worldX       int
	worldY       int
	rayDirX      float64
	rayDirY      float64
	wallDistance float64
	wallX        float64
	side         bool
	drawStart    int
	drawEnd      int
	lineHeight   int

	currentWorldX    int
	currentWorldY    int
	currentLight     float64
	currentTextureId uint
	floorWorldX      int
	floorWorldY      int
	floorLight       float64
	ceilLight        float64
	floorTextureId   uint
	ceilTextureId    uint

	lights *Lights
	pd     *pixels.PictureRGBA
}

func NewWalls(world *World, textures *Textures, lights *Lights, pd *pixels.PictureRGBA, screenWidth int, screenHeight int) *Walls {
	e := &Walls{
		world:        world,
		textures:     textures,
		lights:       lights,
		pd:           pd,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
	for x := 0; x < e.screenWidth; x++ {
		e.camera = append(e.camera, 2*float64(x)/float64(e.screenWidth)-1)
	}
	return e
}

func (f *Walls) Reset() {
	f.dirX = 0
	f.dirY = 0
	f.planeX = 0
	f.planeY = 0
	f.x = 0
	f.cameraX = 0
	f.posX = 0
	f.posY = 0
	f.worldX = 0
	f.worldY = 0
	f.rayDirX = 0
	f.rayDirY = 0
	f.drawStart = 0
	f.drawEnd = 0
	f.lineHeight = 0
	f.wallDistance = 0
	f.wallX = 0
	f.side = false
	f.currentWorldX = -1
	f.currentWorldY = -1
	f.currentLight = 1
	f.currentTextureId = 0
	f.floorWorldX = -1
	f.floorWorldY = -1
	f.floorLight = 1
	f.ceilLight = 1
	f.floorTextureId = 0
	f.ceilTextureId = 0
}

func (f *Walls) Update() {
	f.Reset()
	player := f.world.player

	f.posX = player.GetX()
	f.posY = player.GetY()
	f.dirX = player.dir.X
	f.dirY = player.dir.Y
	f.planeX = player.plane.X
	f.planeY = player.plane.Y

	for f.x, f.cameraX = range f.camera {
		f.rayCasting()
		f.drawWall()
		f.drawFloorCeil()
	}
}

func (f *Walls) rayCasting() {
	f.worldX = int(f.posX)
	f.worldY = int(f.posY)
	f.rayDirX = f.dirX + f.planeX*f.cameraX
	f.rayDirY = f.dirY + f.planeY*f.cameraX
	deltaDistX := math.Sqrt(1.0 + (f.rayDirY*f.rayDirY)/(f.rayDirX*f.rayDirX))
	deltaDistY := math.Sqrt(1.0 + (f.rayDirX*f.rayDirX)/(f.rayDirY*f.rayDirY))
	stepX := 0
	stepY := 0
	sideDistX := 0.0
	sideDistY := 0.0

	if f.rayDirX < 0 {
		stepX = -1
		sideDistX = (f.posX - float64(f.worldX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(f.worldX) + 1.0 - f.posX) * deltaDistX
	}
	if f.rayDirY < 0 {
		stepY = -1
		sideDistY = (f.posY - float64(f.worldY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(f.worldY) + 1.0 - f.posY) * deltaDistY
	}

	f.world.AddVisibility(f.worldX, f.worldY, "wall")

	for rayHit := false; !rayHit; {
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			f.worldX += stepX
			f.side = false
		} else {
			sideDistY += deltaDistY
			f.worldY += stepY
			f.side = true
		}
		if rayHit = f.world.HitWall(f.worldX, f.worldY); !rayHit {
			f.world.AddVisibility(f.worldX, f.worldY, "wall")
		}
	}
	if f.side {
		f.wallDistance = (float64(f.worldY) - f.posY + (1-float64(stepY))/2) / f.rayDirY
		f.wallX = f.posX + f.wallDistance*f.rayDirX
	} else {
		f.wallDistance = (float64(f.worldX) - f.posX + (1-float64(stepX))/2) / f.rayDirX
		f.wallX = f.posY + f.wallDistance*f.rayDirY
	}
	f.wallX -= math.Floor(f.wallX)
	f.lineHeight = int(float64(f.screenHeight) / f.wallDistance)
	if f.lineHeight < 1 {
		f.lineHeight = 1
	}
	baseHeight := f.lineHeight/2 + (f.screenHeight)/2
	f.drawStart = -baseHeight
	if f.drawStart < 0 {
		f.drawStart = 0
	}
	f.drawEnd = baseHeight
	if f.drawEnd >= f.screenHeight {
		f.drawEnd = f.screenHeight - 1
	}
	f.world.SetWallDistance(f.x, f.wallDistance)
}

func (f *Walls) drawWall() {
	texX := int(f.wallX * float64(f.textures.Size()))
	if !f.side && f.rayDirX > 0 {
		texX = int(f.textures.Size()) - texX - 1
	}
	if f.side && f.rayDirY < 0 {
		texX = int(f.textures.Size()) - texX - 1
	}
	if f.currentWorldX != f.worldX || f.currentWorldY != f.worldY {
		f.currentWorldX = f.worldX
		f.currentWorldY = f.worldY
		f.currentTextureId = f.world.GetWallTexture(f.worldX, f.worldY)
		f.currentLight = f.lights.Get(f.worldX, f.worldY)
	}
	h1 := f.screenHeight * 128
	h2 := f.lineHeight * 128
	distance := f.wallDistance
	if f.side {
		distance *= 2
	}

	/*
		z := f.drawStart
		if f.currentTextureId == 98 || f.currentTextureId == 99 {
			if _test_ == 0 {
				_test_ = z
				_test_ *= 1000
			} else {
				_test_ ++
				z = _test_ / 1000
			}
		}
	*/

	for y := f.drawStart; y < f.drawEnd+1; y++ {
		d := ((y * 256) - h1) + h2
		texY := (((d * int(f.textures.Size())) / f.lineHeight) / 256) % int(f.textures.Size())
		c := f.textures.RGBAAt(f.currentTextureId, uint(texX), uint(texY), distance, f.currentLight)
		f.pd.SetColor(f.x, y, c)
	}
}

func (f *Walls) drawFloorCeil() {
	var floorWallX float64
	var floorWallY float64

	if !f.side && f.rayDirX > 0 {
		floorWallX = float64(f.worldX)
		floorWallY = float64(f.worldY) + f.wallX
	} else if !f.side && f.rayDirX < 0 {
		floorWallX = float64(f.worldX) + 1.0
		floorWallY = float64(f.worldY) + f.wallX
	} else if f.side && f.rayDirY > 0 {
		floorWallX = float64(f.worldX) + f.wallX
		floorWallY = float64(f.worldY)
	} else {
		floorWallX = float64(f.worldX) + f.wallX
		floorWallY = float64(f.worldY) + 1.0
	}

	for y := f.drawEnd + 1; y < f.screenHeight; y++ {
		currentDist := float64(f.screenHeight) / (2.0*float64(y) - float64(f.screenHeight))
		//distWall := perpWallDist
		//distPlayer := 0.0
		//weight := (currentDist - distPlayer) / (distWall - distPlayer)
		weight := currentDist / f.wallDistance
		currentFloorX := weight*floorWallX + (1.0-weight)*f.posX
		currentFloorY := weight*floorWallY + (1.0-weight)*f.posY
		sizeTx := currentFloorX * float64(f.textures.Size())
		sizeTy := currentFloorY * float64(f.textures.Size())
		baseTx := uint(sizeTx) % f.textures.Size()
		baseTy := uint(sizeTy) % f.textures.Size()
		tx := int(sizeTx) / int(f.textures.Size())
		ty := int(sizeTy) / int(f.textures.Size())

		if f.floorWorldX != tx || f.floorWorldY != ty {
			f.floorWorldX = tx
			f.floorWorldY = ty
			//floorTextureId = f.world.GetFloorTexture(floorWorldX, floorWorldY)
			//ceilTextureId = f.world.GetCeilTexture(floorWorldX, floorWorldY)
			f.floorTextureId = uint(f.floorWorldX)
			f.ceilTextureId = uint(f.floorWorldY)
			f.floorLight = f.lights.Get(f.floorWorldX, f.floorWorldY)
			f.ceilLight = f.lights.Get(f.floorWorldX, f.floorWorldY)
		}

		ft1 := f.textures.RGBAAt(f.floorTextureId, baseTx, baseTy, currentDist, f.floorLight)
		f.pd.SetColor(f.x, y, ft1)

		ct1 := f.textures.RGBAAt(f.ceilTextureId, baseTx, baseTy, currentDist, f.ceilLight)
		f.pd.SetColor(f.x, f.screenHeight-y, ct1)
	}
}
