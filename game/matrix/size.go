package matrix

type Size struct {
	w float64
	h float64
}

func NewSize(w float64, h float64) Size {
	return Size{w: w, h: h}
}

func (s *Size) GetWidth() float64 {
	return s.w
}

func (s *Size) GetHeight() float64 {
	return s.h
}
