package main

import (
	"github.com/markel1974/gowolfenstein/pixels"
	"math"
)

const (
	spriteMinDistance = 200
	spriteMaxScale    = 16
)

type Sprites struct {
	world            *World
	w                int
	h                int
	screenHalfWidth  int
	screenHalfHeight int
	textures         *Textures
	lights           *Lights
	pd               *pixels.PictureRGBA
}

func NewSprites(world *World, textures *Textures, lights *Lights, pd *pixels.PictureRGBA, screenWidth int, screenHeight int) *Sprites {
	s := &Sprites{
		world:            world,
		textures:         textures,
		lights:           lights,
		pd:               pd,
		w:                screenWidth,
		h:                screenHeight,
		screenHalfWidth:  screenWidth / 2,
		screenHalfHeight: screenHeight / 2,
	}
	return s
}

func (s *Sprites) Update(entities []IWorldEntity) {
	player := s.world.player
	invDet := 1.0 / (player.plane.X*player.dir.Y - player.dir.X*player.plane.Y) //required for correct matrix multiplication
	for _, entity := range entities {
		s.drawSprite(player, entity, invDet)
	}
}

func (s *Sprites) drawSprite(player *Player, entity IWorldEntity, invDet float64) {
	sp := entity.GetEntity()
	//translate sprite position to relative to camera
	spriteX := sp.GetX() - player.GetX()
	spriteY := sp.GetY() - player.GetY()

	//transform sprite with the inverse camera matrix
	// [ planeX   dirX ] -1                                       [ dirY      -dirX ]
	// [               ]       =  1/(planeX*dirY-dirX*planeY) *   [                 ]
	// [ planeY   dirY ]                                          [ -planeY  planeX ]

	transformX := invDet * (player.dir.Y*spriteX - player.dir.X*spriteY)
	transformY := invDet * (-player.plane.Y*spriteX + player.plane.X*spriteY) //this is actually the depth inside the screen, that what Z is in 3D
	if transformY <= 0 {
		return
	}

	spriteScreenX := int(float64(s.screenHalfWidth) * (1 + transformX/transformY))

	spriteHeight := int(float64(s.h) / transformY) //using 'transformY' instead of the real distance prevents fisheye
	spriteWidth := spriteHeight
	spriteHalfWidth := spriteWidth / 2
	spriteHalfHeight := spriteHeight / 2

	drawStartX := -spriteHalfWidth + spriteScreenX
	if drawStartX >= s.w || drawStartX <= -s.w {
		return
	}

	drawEndX := spriteHalfWidth + spriteScreenX
	drawEndXCount := drawEndX - drawStartX
	if drawEndXCount <= 0 {
		return
	}

	drawStartY := -spriteHalfHeight + s.screenHalfHeight
	if drawStartY <= -s.h || drawStartY >= s.h {
		return
	}

	drawEndY := spriteHalfHeight + s.screenHalfHeight
	drawEndYCount := drawEndY - drawStartY
	if drawEndYCount <= 0 {
		return
	}

	spriteLight := s.lights.Get(int(sp.GetX()), int(sp.GetY()))

	//s.pd.SetRGBA(drawStartX, drawStartY, color.RGBA{G: 255, A: 255})
	//s.pd.SetRGBA(drawEndX, drawEndY, color.RGBA{G: 255, A: 255})

	scaleX := float64(s.textures.Size()) / float64(drawEndXCount)
	scaleY := float64(s.textures.Size()) / float64(drawEndYCount)

	if scaleX <= 0 || scaleY <= 0 {
		return
	}
	pixelsX := int(math.Ceil(1 / scaleX))
	pixelsY := int(math.Ceil(1 / scaleY))
	if pixelsX <= 0 || pixelsY <= 0 {
		return
	}

	if pixelsX > 1 || pixelsY > 1 {
		//fmt.Println("TOTAL", pixelsX, pixelsY)
		if pixelsX > spriteMaxScale {
			pixelsX = spriteMaxScale
		}
		if pixelsY > spriteMaxScale {
			pixelsY = spriteMaxScale
		}

		for _, cd := range s.textures.Cache(entity.GetTextureId()) {
			xp := int(math.Floor(float64(cd.x) / scaleX))
			yp := int(math.Floor(float64(cd.y) / scaleY))
			c := cd.Get(entity.GetDistance(), spriteLight)
			for x := xp; x < xp+pixelsX; x++ {
				targetX := x + drawStartX
				if targetX < 0 || targetX >= s.w {
					continue
				}
				if transformY >= s.world.GetWallDistance(targetX) {
					continue
				}
				for y := yp; y < yp+pixelsY; y++ {
					targetY := y + drawStartY
					if targetY < 0 || targetY >= s.h {
						continue
					}
					s.pd.SetColor(targetX, targetY, c)
				}
			}
		}
	} else {
		for _, cd := range s.textures.Cache(entity.GetTextureId()) {
			xp := int(math.Floor(float64(cd.x) / scaleX))
			yp := int(math.Floor(float64(cd.y) / scaleY))
			targetX := xp + drawStartX
			if targetX < 0 || targetX >= s.w {
				continue
			}
			if transformY >= s.world.GetWallDistance(targetX) {
				continue
			}
			targetY := yp + drawStartY
			if targetY < 0 || targetY >= s.h {
				continue
			}
			s.pd.SetColor(targetX, targetY, cd.Get(entity.GetDistance(), spriteLight))
		}
	}
}

/*
func (s * Sprites) drawSprite(player * player, entity IWorldEntity, invDet float64) {
	sp := entity.GetEntity()
	//translate sprite position to relative to camera
	spriteX := sp.GetX() - player.GetX()
	spriteY := sp.GetY() - player.GetY()

	//transform sprite with the inverse camera matrix
	// [ planeX   dirX ] -1                                       [ dirY      -dirX ]
	// [               ]       =  1/(planeX*dirY-dirX*planeY) *   [                 ]
	// [ planeY   dirY ]                                          [ -planeY  planeX ]

	transformX := invDet * (player.dir.Y * spriteX - player.dir.X * spriteY)
	transformY := invDet * (-player.plane.Y * spriteX + player.plane.X * spriteY) //this is actually the depth inside the screen, that what Z is in 3D
	if transformY <= 0 {
		return
	}

	spriteScreenX := int(float64(s.screenHalfWidth) * (1 + transformX / transformY))

	spriteHeight := int(float64(s.h) / transformY) //using 'transformY' instead of the real distance prevents fisheye
	spriteWidth := spriteHeight
	spriteHalfWidth := spriteWidth / 2

	drawStartX := -spriteHalfWidth + spriteScreenX
	if drawStartX < 0 { drawStartX = 0 }
	drawEndX := spriteHalfWidth + spriteScreenX
	if drawEndX > s.w { drawEndX = s.w }

	stripeTotal := drawEndX - drawStartX
	if stripeTotal <= 0 {
		return
	}

	spriteLight := s.lights.Get(int(sp.GetX()), int(sp.GetY()))
	spriteHalfHeight := spriteHeight / 2
	drawStartY := -spriteHalfHeight + s.screenHalfHeight
	if drawStartY < 0 { drawStartY = 0 }
	drawEndY := spriteHalfHeight + s.screenHalfHeight
	if drawEndY > s.h { drawEndY = s.h}


	dispatch := s.createDispatcher(stripeTotal, drawStartX, drawEndX)

	fmt.Println("dispatch", dispatch)

	var wg sync.WaitGroup
	wg.Add(stripeTotal)

	for stripe := drawStartX; stripe < drawEndX; stripe++ {
		go func(stripe int) {
			if transformY < s.world.GetWallDistance(stripe) {
				if !s.cache.hasInterval(stripe, drawStartY, drawEndY-1) {
					drawCount := 0
					texX := uint((256 * (stripe - (-spriteWidth / 2 + spriteScreenX)) * int(s.textureSize) / spriteWidth) / 256)
					for spriteScreenY := drawStartY; spriteScreenY < drawEndY; spriteScreenY++ {
						if !s.cache.hasPixel(stripe, spriteScreenY) {
							//256 and 128 factors to avoid floats
							d := (spriteScreenY)*256 - s.h*128 + spriteHeight*128
							texY := uint(((d * int(s.textureSize)) / spriteHeight) / 256)
							c := s.textures.RGBAAt(entity.GetTextureId(), texX, texY % s.textureSize, entity.GetDistance(), spriteLight)
							if c.A > 0 {
								s.pd.SetRGBA(stripe, spriteScreenY - int(sp.G), c)
								s.cache.setPixel(stripe, spriteScreenY)
								drawCount++
							}
						}
					}
					s.cache.setInterval(stripe, drawStartY, drawCount)
				}
			}
			wg.Done()
		}(stripe)
	}
	wg.Wait()
}
*/

/*
func (s * Sprites) createDispatcher(stripeTotal int, start int, end int) [][]int{
	cpus := runtime.GOMAXPROCS(0)
	var dispatch [][] int
	var dispatchData []int

	dispatchCount := stripeTotal / cpus
	for stripe := start; stripe < end; stripe++ {
		if len(dispatchData) >= dispatchCount {
			dispatch = append(dispatch, dispatchData)
			dispatchData = nil
		}
		dispatchData = append(dispatchData, stripe)
	}
	if len(dispatchData) > 0 {
		dispatch = append(dispatch, dispatchData)
	}
	return dispatch
}
*/

/*
func (s * Sprites) drawSprite(player * player, entity IWorldEntity, invDet float64) {
	sp := entity.GetEntity()
	//translate sprite position to relative to camera
	spriteX := sp.GetX() - player.GetX()
	spriteY := sp.GetY() - player.GetY()

	//transform sprite with the inverse camera matrix
	// [ planeX   dirX ] -1                                       [ dirY      -dirX ]
	// [               ]       =  1/(planeX*dirY-dirX*planeY) *   [                 ]
	// [ planeY   dirY ]                                          [ -planeY  planeX ]

	transformX := invDet * (player.dir.Y * spriteX - player.dir.X * spriteY)
	transformY := invDet * (-player.plane.Y * spriteX + player.plane.X * spriteY) //this is actually the depth inside the screen, that what Z is in 3D
	if transformY <= 0 {
		return
	}

	spriteScreenX := int(float64(s.screenHalfWidth) * (1 + transformX / transformY))

	spriteHeight := int(float64(s.h) / transformY) //using 'transformY' instead of the real distance prevents fisheye
	spriteWidth := spriteHeight
	spriteHalfWidth := spriteWidth / 2

	drawStartX := -spriteHalfWidth + spriteScreenX
	if drawStartX < 0 { drawStartX = 0 }
	drawEndX := spriteHalfWidth + spriteScreenX
	if drawEndX > s.w { drawEndX = s.w }

	stripeTotal := drawEndX - drawStartX
	if stripeTotal <= 0 {
		return
	}

	spriteLight := s.lights.Get(int(sp.GetX()), int(sp.GetY()))
	spriteHalfHeight := spriteHeight / 2
	drawStartY := -spriteHalfHeight + s.screenHalfHeight
	if drawStartY < 0 { drawStartY = 0 }
	drawEndY := spriteHalfHeight + s.screenHalfHeight
	if drawEndY > s.h { drawEndY = s.h}
	//TODO OPTIMIZE!!!!
	dispatch := s.createDispatcher(stripeTotal, drawStartX, drawEndX)

	var wg sync.WaitGroup
	wg.Add(len(dispatch))

	for _, d := range dispatch {
		go func(stripes []int) {
			for _, stripe := range stripes {
				if transformY < s.world.GetWallDistance(stripe) {
					if !s.cache.hasInterval(stripe, drawStartY, drawEndY-1) {
						drawCount := 0
						texX := uint((256 * (stripe - (-spriteWidth/2 + spriteScreenX)) * int(s.textureSize) / spriteWidth) / 256)
						for spriteScreenY := drawStartY; spriteScreenY < drawEndY; spriteScreenY++ {
							if !s.cache.hasPixel(stripe, spriteScreenY) {
								//256 and 128 factors to avoid floats
								d := (spriteScreenY)*256 - s.h*128 + spriteHeight*128
								texY := uint(((d * int(s.textureSize)) / spriteHeight) / 256)
								c := s.textures.RGBAAt(entity.GetTextureId(), texX, texY%s.textureSize, entity.GetDistance(), spriteLight)
								if c.A > 0 {
									s.pd.SetRGBA(stripe, spriteScreenY-int(sp.G), c)
									s.cache.setPixel(stripe, spriteScreenY)
									drawCount++
								}
							}
						}
						s.cache.setInterval(stripe, drawStartY, drawCount)
					}
				}
			}

			wg.Done()
		}(d)
	}
	wg.Wait()
}
*/
