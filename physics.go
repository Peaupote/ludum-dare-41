package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const (
	moveSpeed = 50
)

// RigidBody represents any physical beeing
type RigidBody struct {
	// physics
	body     pixel.Rect
	velocity pixel.Vec
}

// NewRigidBodyBySize create a new RigidBody
func NewRigidBodyBySize(x, y, w, h float64, v pixel.Vec) *RigidBody {
	return &RigidBody{
		body:     pixel.R(x, y, x+w, y+h),
		velocity: v,
	}
}

func (r *RigidBody) physics(dt float64) {
	r.body = r.body.Moved(r.velocity.Scaled(dt))
}

func (r *RigidBody) draw(t *imdraw.IMDraw) {
	t.Push(r.body.Max)
	t.Push(r.body.Min)
	t.Rectangle(0)
}

func (r *RigidBody) hit(rect pixel.Rect) bool {
	// TODO: seems not to work
	return r.body.Norm().Intersect(rect.Norm()).Area() != 0
}
