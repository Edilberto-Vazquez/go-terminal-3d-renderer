package main

import (
	"fmt"
	"time"
)

const (
	GREENMATRIX = "\033[38;5;46m"
	MAGENTANEON = "\033[38;5;201m"
	WHITE       = "\033[1;97m"
	RESET       = "\033[0m"
)
const WHITECHAR = GREENMATRIX + "-" + RESET
const POINTCHAR = WHITE + "*" + RESET

func main() {
	object3D := NewObject3D(
		[]Vertex{
			NewVertex(-4, 4, 10),
			NewVertex(4, 4, 10),
			NewVertex(4, -4, 10),
			NewVertex(-4, -4, 10),

			NewVertex(-4, 4, 20),
			NewVertex(4, 4, 20),
			NewVertex(4, -4, 20),
			NewVertex(-4, -4, 20),
		},
		[]Edge{
			NewEdge(0, 1),
			NewEdge(1, 2),
			NewEdge(2, 3),
			NewEdge(3, 0),

			NewEdge(4, 5),
			NewEdge(5, 6),
			NewEdge(6, 7),
			NewEdge(7, 4),

			NewEdge(0, 4),
			NewEdge(1, 5),
			NewEdge(2, 6),
			NewEdge(3, 7),
		},
		[]Transformer{
			{Degrees: 30, Rotation: rotationX},
			{Degrees: 70, Rotation: rotationY, AngularSpeed: 10},
			{Degrees: 15, Rotation: rotationZ},
		},
	)
	render := NewRender(23, 23)
	projection := NewProjection(render.width, render.height, 20, 5)

	for i := 0; i <= 360; i++ {
		object3D.Rasterize(render, projection)
		// A partir del segundo frame, regresar el cursor
		// 20 filas hacia arriba.
		if i > 0 {
			fmt.Printf("\033[%dA", render.height)
		}

		render.Draw()
		object3D.Update()
		time.Sleep(100 * time.Millisecond) // Pause for 100 milliseconds before the next frame
		render.ClearFrameBuffer()
	}
}
