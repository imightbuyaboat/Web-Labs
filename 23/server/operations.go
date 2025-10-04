package main

import (
	"fmt"
	"math"
)

func Add(x, y float64) float64 {
	return x + y
}

func Sub(x, y float64) float64 {
	return x - y
}

func Mult(x, y float64) float64 {
	return x * y
}

func Div(x, y float64) (float64, error) {
	if y == 0.0 {
		return .0, fmt.Errorf("div by 0")
	}
	return x / y, nil
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return .0, fmt.Errorf("x must be >= 0")
	}
	return math.Sqrt(x), nil
}

func Percent(x, percent float64) (float64, error) {
	if percent < 0 || percent > 100 {
		return .0, fmt.Errorf("percent must be in [0, 100]")
	}
	return x * percent / 100, nil
}

func Round(x float64, y int64) float64 {
	pow := math.Pow(10, float64(y))
	return math.Round(x*pow) / pow
}

func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}
