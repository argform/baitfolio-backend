package geo

import (
	"math"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= eps
}

func TestXToLon(t *testing.T) {
	tests := []struct {
		name string
		x float64
		n float64
		want float64
	}{
		{"left edge z0", 0, 1, -180},
		{"right edge z0", 1, 1, 180},
		{"left half z1", 0, 2, -180},
		{"center z1", 1, 2, 0},
		{"right half z1", 2, 2, 180},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xToLon(tt.x, tt.n)
			if !almostEqual(got, tt.want) {
				t.Fatalf("xToLon(%v, %v) = %v, want %v", tt.x, tt.n, got, tt.want)
			}
		})
	}
}

func TestYToLat(t *testing.T) {
	maxLat := 85.0511287798066

	tests := []struct {
		name string
		y float64
		n float64
		want float64
	}{
		{"top edge z0", 0, 1, maxLat},
		{"bottom edge z0", 1, 1, -maxLat},
		{"top edge z1", 0, 2, maxLat},
		{"equator z1", 1, 2, 0},
		{"bottom edge z1", 2, 2, -maxLat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := yToLat(tt.y, tt.n)
			if !almostEqual(got, tt.want) {
				t.Fatalf("yToLat(%v, %v) = %.12f, want %.12f", tt.y, tt.n, got, tt.want)
			}
		})
	}
}

func TestTileToBBox(t *testing.T) {
	maxLat := 85.0511287798066

	tests := []struct {
		name string
		tile Tile
		want BBox
	}{
		{
			name: "whole world z0",
			tile: Tile{X: 0, Y: 0, Z: 0},
			want: BBox{
				West:  -180,
				East:  180,
				North: maxLat,
				South: -maxLat,
			},
		},
		{
			name: "bottom-right quarter z1",
			tile: Tile{X: 1, Y: 1, Z: 1},
			want: BBox{
				West:  0,
				East:  180,
				North: 0,
				South: -maxLat,
			},
		},
		{
			name: "top-left quarter z1",
			tile: Tile{X: 0, Y: 0, Z: 1},
			want: BBox{
				West:  -180,
				East:  0,
				North: maxLat,
				South: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TileToBBox(tt.tile)

			if !almostEqual(got.West, tt.want.West) {
				t.Fatalf("West = %.12f, want %.12f", got.West, tt.want.West)
			}
			if !almostEqual(got.East, tt.want.East) {
				t.Fatalf("East = %.12f, want %.12f", got.East, tt.want.East)
			}
			if !almostEqual(got.North, tt.want.North) {
				t.Fatalf("North = %.12f, want %.12f", got.North, tt.want.North)
			}
			if !almostEqual(got.South, tt.want.South) {
				t.Fatalf("South = %.12f, want %.12f", got.South, tt.want.South)
			}
		})
	}
}
