package excel

import (
	"math"
	"testing"
)

func TestCalcularSomaQuadrados(t *testing.T) {
	tests := []struct {
		name    string
		valores []float64
		want    float64
	}{
		{"Empty slice", []float64{}, 0},
		{"Single value", []float64{3}, 9},
		{"Multiple values", []float64{1, 2, 3}, 14},
		{"Zeros", []float64{0, 0}, 0},
		{"Negative values", []float64{-2, -4}, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcularSomaQuadrados(tt.valores); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("calcularSomaQuadrados() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcularMedia(t *testing.T) {
	tests := []struct {
		name    string
		valores []float64
		want    float64
	}{
		{"Empty slice", []float64{}, 0},
		{"Single value", []float64{5}, 5},
		{"Multiple positive values", []float64{2, 4, 6}, 4},
		{"Mix of positive and negative", []float64{-2, 0, 2}, 0},
		{"Fractional values", []float64{1.5, 2.5}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcularMedia(tt.valores); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("calcularMedia() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcularDesvioPadrao(t *testing.T) {
	tests := []struct {
		name    string
		valores []float64
		want    float64
	}{
		{"Empty slice", []float64{}, 0},
		{"Single value", []float64{5}, 0},
		{"Identical values", []float64{5, 5, 5}, 0},
		{"Different values", []float64{1, 2, 3}, 1},
		{"Larger uniform range", []float64{2, 4, 4, 4, 5, 5, 7, 9}, 2.138089935299395},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcularDesvioPadrao(tt.valores)
			if math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("calcularDesvioPadrao() = %v, want %v", got, tt.want)
			}
		})
	}
}
