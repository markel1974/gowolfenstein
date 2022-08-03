package main

type VisibilityData struct {
	id   int
	kind string
}

func NewVisibilityData() *VisibilityData {
	return &VisibilityData{}
}

type Visibility struct {
	id   int
	w    int
	h    int
	data [][]*VisibilityData
}

func NewVisibility(w int, h int) *Visibility {
	data := make([][]*VisibilityData, w)
	for x := 0; x < w; x++ {
		data[x] = make([]*VisibilityData, h)
		for y := 0; y < h; y++ {
			data[x][y] = NewVisibilityData()
		}
	}
	return &Visibility{
		id:   0,
		w:    w,
		h:    h,
		data: data,
	}
}

func (w *Visibility) Update() {
	w.id++
}

func (w *Visibility) GetId() int {
	return w.id
}

func (w *Visibility) GetWorld() [][]*VisibilityData {
	return w.data
}

func (w *Visibility) Add(x int, y int, kind string) {
	w.set(x, y, kind)
	/*
		w.set(x - 1, y - 1)
		w.set(x - 1, y)
		w.set(x, y - 1)
		w.set(x + 1, y + 1)
		w.set(x + 1, y)
		w.set(x, y + 1)
	*/
}

func (w *Visibility) set(x int, y int, kind string) {
	if !w.valid(x, y) {
		return
	}
	d := w.data[x][y]
	d.id = w.id
	d.kind = kind
}

func (w *Visibility) Get(x int, y int) bool {
	if w.valid(x, y) {
		d := w.data[x][y]
		return d.id == w.id
	}
	return false
}

func (w *Visibility) valid(x int, y int) bool {
	if x >= 0 && x < w.w && y >= 0 && y < w.h {
		return true
	}
	return false
}
