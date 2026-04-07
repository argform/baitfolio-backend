package geo

import "math"

type Tile struct {
	X int
	Y int
	Z int
}

type BBox struct {
	West float64
	East float64
	North float64
	South float64
}

func xToLon(x, n float64) float64 {
	return x/n*360.0 - 180.0
}

func yToLat(y, n float64) float64 {
	return math.Atan(math.Sinh(math.Pi*(1.0-2.0*y/n))) * 180.0 / math.Pi
}

func TileToBBox(t Tile) BBox {
	n := math.Exp2(float64(t.Z))
	x := float64(t.X)
	y := float64(t.Y)

	return BBox{
		West: xToLon(x, n),
		East: xToLon(x+1, n),
		North: yToLat(y, n),
		South: yToLat(y+1, n),
	}
}