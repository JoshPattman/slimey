package main

import (
	"math"

	"github.com/gopxl/pixel"
)

type Particle struct {
	Pos           pixel.Vec
	Dir           pixel.Vec
	RotationSpeed float64
	Speed         float64
}

func (p *Particle) UpdateDir(grid *Grid, deltaTime float64) {
	sensorLeftPosition := p.Pos.Add(p.Dir.Scaled(20).Rotated(-math.Pi / 5))
	sensorRightPosition := p.Pos.Add(p.Dir.Scaled(20).Rotated(math.Pi / 5))

	sensorLeftReading := grid.Sense(sensorLeftPosition)
	sensorRightReading := grid.Sense(sensorRightPosition)

	if sensorLeftReading > sensorRightReading {
		p.Dir = p.Dir.Rotated(-p.RotationSpeed * deltaTime)
	} else if sensorLeftReading < sensorRightReading {
		p.Dir = p.Dir.Rotated(p.RotationSpeed * deltaTime)
	}
}

func (p *Particle) UpdatePos(deltaTime float64, maxX, maxY float64) {
	p.Pos = p.Pos.Add(p.Dir.Scaled(p.Speed * deltaTime))
	if p.Pos.X < 0 && p.Dir.X < 0 {
		p.Dir.X *= -1
	}
	if p.Pos.Y < 0 && p.Dir.Y < 0 {
		p.Dir.Y *= -1
	}
	if p.Pos.X > maxX && p.Dir.X > 0 {
		p.Dir.X *= -1
	}
	if p.Pos.Y > maxY && p.Dir.Y > 0 {
		p.Dir.Y *= -1
	}
}
