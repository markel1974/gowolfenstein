package main

type lightData struct {
	id        int
	intensity float64
	size      int
}

func newLightData() *lightData {
	return &lightData{}
}

type Lights struct {
	id   int
	w    int
	h    int
	data [][]*lightData
}

func NewLights(w int, h int) *Lights {
	data := make([][]*lightData, w)
	for x := 0; x < w; x++ {
		data[x] = make([]*lightData, h)
		for y := 0; y < h; y++ {
			data[x][y] = newLightData()
		}
	}
	return &Lights{
		id:   0,
		w:    w,
		h:    h,
		data: data,
	}
}

func (l *Lights) Update() {
	l.id++
}

func (l *Lights) Set(centerX int, centerY int, diameter int, intensity float64) {
	//TODO TEST
	//diameter = 3
	//intensity = 2.5

	if diameter < 1 {
		l.set(centerX, centerY, intensity)
		return
	}
	if diameter%2 == 0 {
		diameter++
	}
	radius := diameter / 2

	//fmt.Println("-------------------------")
	//fmt.Println("light", lightX, lightY)
	//fmt.Println("center", centerX, centerY)

	for lx := centerX - radius; lx <= centerX+radius; lx++ {
		if lx < 0 {
			continue
		}
		//if lx == centerX {
		//	fmt.Println("---------- CENTER --------------", lx)
		//} else {
		//	fmt.Println("--------------------------------", lx)
		//}
		for ly := centerY - radius; ly <= centerY+radius; ly++ {
			if ly < 0 {
				continue
			}
			distance := (centerX-lx)*(centerX-lx) + (centerY-ly)*(centerY-ly)
			targetI := intensity - (float64(distance) * 0.3)
			if targetI < 1 {
				targetI = 1
			}
			l.set(lx, ly, targetI)
			//fmt.Println("distance (", lx, ly, ")", distance, targetI)
		}
	}
}

func (l *Lights) Get(x int, y int) float64 {
	if l.valid(x, y) {
		light := l.data[x][y]
		if light.id == l.id {
			return light.intensity
		}
	}
	return 1.0
}

func (l *Lights) set(x int, y int, intensity float64) {
	if !l.valid(x, y) {
		return
	}
	light := l.data[x][y]
	if light.id == l.id {
		if intensity > light.intensity {
			light.intensity = intensity
		}
	} else {
		light.id = l.id
		light.intensity = intensity
	}
}

func (l *Lights) valid(x int, y int) bool {
	if x >= 0 && x < l.w && y >= 0 && y < l.h {
		return true
	}
	return false
}
