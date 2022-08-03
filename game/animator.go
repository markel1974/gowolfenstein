package main

import "time"

type Animator struct {
	idx          int
	animation    []uint
	rate         int64
	lastInterval int64
	repeat       bool
	count        int
	defaultId    uint
}

func NewAnimator() *Animator {
	return &Animator{
		idx:       0,
		animation: nil,
		rate:      0,
		repeat:    true,
		count:     0,
	}
}

func (a *Animator) SetStatic(id uint) {
	a.defaultId = id
	a.animation = nil
}

func (a *Animator) Set(animation []uint, rate int, repeat bool) {
	a.animation = animation
	a.count = len(animation) - 1
	a.rate = int64(rate)
	a.repeat = repeat
	a.lastInterval = 0
	a.idx = 0
}

func (a *Animator) Play() uint {
	if a.animation == nil {
		return a.defaultId
	}
	if !a.repeat && a.idx == a.count {
		return a.animation[a.count]
	}
	epoch := time.Now().UnixNano() / int64(time.Millisecond)
	last := epoch - a.lastInterval
	if last < a.rate {
		return a.animation[a.idx]
	}
	a.lastInterval = epoch
	a.idx++
	if a.idx > a.count {
		a.idx = 0
	}
	return a.animation[a.idx]
}
