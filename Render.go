package main

import (
	"fmt"
	"math"
	"slices"
)

type Point2D struct {
	PixelX int
	PixelY int
}

func NewPoint2D(px, py int) Point2D {
	return Point2D{PixelX: px, PixelY: py}
}

// Render represents a 2D grid of strings that can be drawn to the console.
type Render struct {
	width       int
	height      int
	frameBuffer [][]string
}

func NewRender(width, height int) *Render {
	// Create a 2D slice of strings to represent the render buffer
	frameBuffer := make([][]string, height)
	// Fill the buffer with " " to represent empty space
	for i := range frameBuffer {
		frameBuffer[i] = slices.Repeat([]string{WHITECHAR}, width)
	}
	return &Render{
		width:       width,
		height:      height,
		frameBuffer: frameBuffer,
	}
}

func (r *Render) SetPoint(point2D Point2D) {
	r.frameBuffer[point2D.PixelY][point2D.PixelX] = POINTCHAR
}

func (r *Render) ClearFrameBuffer() {
	for i := range r.frameBuffer {
		r.frameBuffer[i] = slices.Repeat([]string{WHITECHAR}, r.width)
	}
}

func (r *Render) DrawLine2D(sx1, sy1, sx2, sy2 float64) {
	steps := lineSteps(sx1, sy1, sx2, sy2)

	if steps == 0 {
		pixelSX := int(math.Round(sx1))
		pixelSY := int(math.Round(sy1))

		// If the projected endpoints coincide, paint a single pixel
		if pixelSX < r.width && pixelSY < r.height && pixelSX >= 0 && pixelSY >= 0 {
			r.SetPoint(NewPoint2D(pixelSX, pixelSY))
		}
		return
	}

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)                   // Calculate the interpolation parameter t (0 <= t <= 1)
		x, y := linearInterpolation(sx1, sy1, sx2, sy2, t) // Get the interpolated point

		pixelX := int(math.Round(x))
		pixelY := int(math.Round(y))

		if pixelX < r.width && pixelY < r.height && pixelX >= 0 && pixelY >= 0 { // Ensure the point is within bounds
			r.SetPoint(NewPoint2D(pixelX, pixelY))
		}
	}
}

func (r *Render) Draw() {
	for _, row := range r.frameBuffer {
		for _, character := range row {
			fmt.Print(character, " ")
		}
		fmt.Println()
	}
}

// lineSteps calculates the number of steps needed to draw a line between two points (x1, y1) and (x2, y2) based on the greater distance in either the x or y direction.
// DDA algorithm is used to determine the number of steps required for drawing a line between two points in a 2D space. The number of steps is determined by the greater distance in either the x or y direction, ensuring that the line is drawn smoothly and accurately.
func lineSteps(x1, y1, x2, y2 float64) int {
	dx := math.Abs(x2 - x1)
	dy := math.Abs(y2 - y1)
	return int(math.Ceil(math.Max(dx, dy))) // Determine the number of steps based on the greater distance
}
