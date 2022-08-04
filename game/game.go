package main

import (
	"flag"
	"github.com/markel1974/gowolfenstein/pixels"
	"time"
)

type Game struct {
	fullscreen   bool
	showMap      bool
	screenWidth  int
	screenHeight int
	scale        float64

	mainSurface *pixels.PictureRGBA
	world       *World
	textures    *Textures
	walls       *Walls

	miniMap    *MiniMap
	mainMatrix pixels.Matrix
	mainSprite *pixels.Sprite
	mapMatrix  pixels.Matrix
	mapSprite  *pixels.Sprite
	mapAngle   float64
	weapons    []*Weapon
}

func NewGame() *Game {
	return &Game{
		fullscreen:   false,
		showMap:      true,
		screenWidth:  1980 / 4,
		screenHeight: 1080 / 4,
		scale:        2.0,
	}
}

func (g *Game) Setup(c pixels.Vec) {
	g.mainSurface = pixels.NewPictureRGBA(pixels.R(float64(0), float64(0), float64(g.screenWidth), float64(g.screenHeight)))
	g.textures = NewTextures(g.screenWidth, g.screenHeight, 6)
	g.world = NewWorld(g.screenWidth, g.screenHeight, g.textures, g.mainSurface)

	g.walls = NewWalls(g.world, g.textures, g.world.lights, g.mainSurface, g.screenWidth, g.screenHeight)

	g.miniMap = NewMiniMap(g.world)

	g.mainSprite = pixels.NewSprite()
	g.mainSprite.SetCached(pixels.CacheModeUpdate)
	g.mainSprite.Set(g.mainSurface, g.mainSurface.Bounds())
	g.mainMatrix = pixels.IM.Moved(c).Scaled(c, g.scale)

	mapSurface := g.miniMap.GetPicture()
	g.mapSprite = pixels.NewSprite()
	g.mapSprite.SetCached(pixels.CacheModeUpdate)
	g.mapSprite.Set(mapSurface, mapSurface.Bounds())
	g.mapAngle = 0.0

	weapon1 := NewWeapon(g.world, 124, g.textures.RatioWidth(120), g.textures.RatioWidth(120), 1, 1.0, 120)
	weapon2 := NewWeapon(g.world, 124, g.textures.RatioWidth(129), g.textures.RatioWidth(129), 1, 1.2, 120)
	weapon3 := NewWeapon(g.world, 124, g.textures.RatioWidth(156), g.textures.RatioWidth(156), 1, 1.5, 1000)

	g.weapons = append(g.weapons, weapon1)
	g.weapons = append(g.weapons, weapon2)
	g.weapons = append(g.weapons, weapon3)
}

func (g *Game) Run() {
	cfg := pixels.WindowConfig{
		Bounds:      pixels.R(0, 0, float64(g.screenWidth)*g.scale, float64(g.screenHeight)*g.scale),
		VSync:       true,
		Undecorated: false,
		Smooth:      false,
	}

	if g.fullscreen {
		cfg.Monitor = pixels.PrimaryMonitor()
	}

	win, err := pixels.NewGLWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := win.Bounds().Center()

	g.Setup(c)

	var last time.Time
	const move = 4.5
	const moveSlow = 1.2
	const framerate = 30

	var currentTimer float64
	var lastTimer float64
	var speed float64
	mouseEnabled := true

	//win.SetTitle("wolf (fps: " + strconv.Itoa(counter) + ")")

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		doubleSpeed := false

		if mouseEnabled && win.MouseInsideWindow() {
			mousePos := win.MousePosition()
			mousePrevPos := win.MousePreviousPosition()
			if mousePos.X != mousePrevPos.X || mousePos.Y != mousePrevPos.Y {
				mouseX := mousePos.X - mousePrevPos.X
				mouseY := mousePos.Y - mousePrevPos.Y
				g.world.player.mouseMove(mouseX, mouseY)
			}
		}

		for v := range win.KeysPressed() {
			switch v {
			case pixels.KeyEscape:
				return
			case pixels.KeyM:
				g.showMap = true
			case pixels.KeyN:
				g.showMap = false
			case pixels.KeyUp, pixels.KeyW:
				g.world.player.moveForward(move * dt * speed)
			case pixels.KeyDown, pixels.KeyS:
				g.world.player.moveBackwards(move * dt * speed)
			case pixels.KeyRight:
				g.world.player.turnRight(moveSlow * dt * speed)
			case pixels.KeyLeft:
				g.world.player.turnLeft(moveSlow * dt * speed)
			case pixels.KeyA:
				g.world.player.moveLeft(move * dt * speed)
			case pixels.KeyD:
				g.world.player.moveRight(move * dt * speed)
			case pixels.KeyLeftShift, pixels.KeyRightShift:
				doubleSpeed = true
			case pixels.KeySpace:
				g.world.player.Fire()
			}
		}

		if doubleSpeed {
			speed = 0.8
		} else {
			speed = 0.4
		}

		if win.JustPressed(pixels.KeyJ) {
			g.world.player.SetWeapon(g.weapons[0])
			g.world.player.Fire()
		}
		if win.JustPressed(pixels.KeyK) {
			g.world.player.SetWeapon(g.weapons[1])
			g.world.player.Fire()
		}
		if win.JustPressed(pixels.KeyL) {
			g.world.player.SetWeapon(g.weapons[2])
			g.world.player.Fire()
		}
		if win.JustPressed(pixels.MouseButton1) {
			g.world.player.SetWeapon(g.weapons[2])
			g.world.player.Fire()
		}
		if win.JustPressed(pixels.KeyP) {
			g.world.player.Kick()
		}
		if win.JustPressed(pixels.Key1) {
			g.world.player.SetLight(8, 4)
		}
		if win.JustPressed(pixels.Key2) {
			g.world.player.SetLight(0, 0)
		}
		if win.JustPressed(pixels.KeyTab) {
			g.world.player.Open()
		}
		if win.JustPressed(pixels.KeyM) {
			mouseEnabled = false
		}
		//if win.JustPressed(pixelgl.KeyS) { g.Setup(c) }

		currentTimer = pixels.GLGetTime()
		if currentTimer-lastTimer >= 1.0/framerate {
			lastTimer = currentTimer

			//win.Clear(color.Black)
			g.world.UpdateAvailable()
			g.walls.Update()
			//g.lights.Update()
			g.world.UpdateWorld()
			//g.sprites.Update(entities)

			g.mainSprite.Draw(win, g.mainMatrix)
			if g.showMap {
				g.miniMap.Update(g.world.player.GetX(), g.world.player.GetY())
				mapAngle := 1.57 + g.world.player.dir.Angle() //3.14 / 2 => 1.57
				if mapAngle != g.mapAngle {
					g.mapAngle = mapAngle
					mc := pixels.Vec{X: g.miniMap.x, Y: g.miniMap.y}
					g.mapMatrix = pixels.IM.Moved(mc).Rotated(mc, mapAngle).ScaledXY(pixels.ZV, pixels.MakeVec(-g.scale, g.scale))
				}
				g.mapSprite.Draw(win, g.mapMatrix)
			}
		}
		win.Update()
	}
}

func main() {
	g := NewGame()

	flag.BoolVar(&g.fullscreen, "f", g.fullscreen, "fullscreen")
	flag.IntVar(&g.screenWidth, "w", g.screenWidth, "width")
	flag.IntVar(&g.screenHeight, "h", g.screenHeight, "height")
	flag.Float64Var(&g.scale, "s", g.scale, "scale")
	flag.Parse()

	//g.Setup()

	pixels.GLRun(g.Run)
}
