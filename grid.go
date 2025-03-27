package main

import (
	"image"
	"math"

	"github.com/gopxl/pixel"
	"gonum.org/v1/gonum/mat"
)

type Grid struct {
	chunkSie int
	grid     *mat.Dense
}

func NewGrid(sizeX, sizeY, chunkSize int) *Grid {
	if (sizeX/chunkSize)*chunkSize != sizeX || (sizeY/chunkSize)*chunkSize != sizeY {
		panic("not divisible")
	}

	return &Grid{
		chunkSie: chunkSize,
		grid:     mat.NewDense(sizeX/chunkSize+1, sizeY/chunkSize+1, nil),
	}
}

func (g *Grid) Decay(mult float64) {
	g.grid.Scale(mult, g.grid)
}

func (g *Grid) AddPheremone(at pixel.Vec, amount float64) {
	gpx := int(math.Round(at.X / float64(g.chunkSie)))
	gpy := int(math.Round(at.Y / float64(g.chunkSie)))
	if gpx >= 0 && gpy >= 0 && gpx < g.grid.RawMatrix().Cols && gpy < g.grid.RawMatrix().Rows {
		g.grid.Set(gpx, gpy, g.grid.At(gpx, gpy)+amount)
	}
}

func (g *Grid) Sense(at pixel.Vec) float64 {
	x := int(math.Round(at.X / float64(g.chunkSie)))
	y := int(math.Round(at.Y / float64(g.chunkSie)))
	if x < 0 || x >= g.grid.RawMatrix().Cols {
		return 0
	}
	if y < 0 || y >= g.grid.RawMatrix().Rows {
		return 0
	}
	return g.grid.At(x, y)
}

func (g *Grid) Sprite() (*pixel.Sprite, pixel.Matrix) {
	img := image.NewRGBA(image.Rect(0, 0, g.grid.RawMatrix().Rows, g.grid.RawMatrix().Cols))
	for x := range g.grid.RawMatrix().Rows {
		for y := range g.grid.RawMatrix().Cols {
			v := g.grid.At(x, y)
			if v > 1 {
				v = 1
			}
			img.Set(x, g.grid.RawMatrix().Cols-1-y, pixel.RGB(v, v, v))
		}
	}
	pic := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pic, pic.Bounds()), pixel.IM.Scaled(pixel.ZV, float64(g.chunkSie)).Moved(pixel.V(300, 300))
}
