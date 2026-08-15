package main

// linearInterpolation performs linear interpolation between two points (x1, y1) and (x2, y2) based on a parameter t (0 <= t <= 1).
func linearInterpolation(x1, y1, x2, y2, t float64) (float64, float64) {
	x := x1 + t*(x2-x1) // x = x1 + t * (x2 - x1)
	y := y1 + t*(y2-y1) // y = y1 + t * (y2 - y1)
	return x, y
}
