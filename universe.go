package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Ovni interface {
	getRigidBody() *RigidBody
	draw(*imdraw.IMDraw)
}

type Asteroid struct {
	rigidBody *RigidBody
}

func (a Asteroid) getRigidBody() *RigidBody {
	return a.rigidBody
}

func (a Asteroid) draw(imag *imdraw.IMDraw) {
	imag.Color = colornames.Brown
	a.rigidBody.draw(imag)
}

func updateUniverse(dt float64, ovnis []Ovni) []Ovni {
	for i, o := range ovnis {
		r := o.getRigidBody()
		r.physics(dt)
		if r.body.Max.Y < 0 {
			ovnis[i] = ovnis[len(ovnis)-1]
			ovnis = ovnis[:len(ovnis)-1]
		}
	}

	// TODO: less "random" random generation
	if rand.Intn(100) == 1 {
		ovnis = append(ovnis, Asteroid{
			rigidBody: NewRigidBodyBySize(
				rand.Float64()*(width/2),
				height+10,
				50.0, 50.0,
				pixel.V(0, -100),
			),
		})
	}

	return ovnis
}
