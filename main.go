package main

import (
	"fmt"
	"math"
	"slices"
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

type Projection struct {
	xCenter float64
	yCenter float64
	scale   float64
	zNear   float64
}

func NewProjection(renderWidth, renderHeight int, scale, zNear float64) Projection {
	return Projection{
		xCenter: float64(renderWidth-1) / 2,
		yCenter: float64(renderHeight-1) / 2,
		scale:   scale,
		zNear:   zNear,
	}
}

func calculateInterpolationParameterZ(z1, z2, znear float64) float64 {
	return (znear - z1) / (z2 - z1)
}

func calculateInterpolationZ(z1, z2, t float64) float64 {
	return z1 + t*(z2-z1)
}

// linearInterpolation performs linear interpolation between two points (x1, y1) and (x2, y2) based on a parameter t (0 <= t <= 1).
func linearInterpolation(x1, y1, x2, y2, t float64) (float64, float64) {
	x := x1 + t*(x2-x1) // x = x1 + t * (x2 - x1)
	y := y1 + t*(y2-y1) // y = y1 + t * (y2 - y1)
	return x, y
}

func (p Projection) ClippingEdge(p1, p2 Vertex) (Vertex, Vertex, bool) {
	if p1.z < p.zNear && p2.z < p.zNear {
		return p1, p2, false
	}

	if p1.z < p.zNear && p2.z >= p.zNear {
		t := calculateInterpolationParameterZ(p1.z, p2.z, p.zNear)
		z := calculateInterpolationZ(p1.z, p2.z, t)
		x, y := linearInterpolation(p1.x, p1.y, p2.x, p2.y, t)
		return NewVertex(x, y, z), p2, true
	}

	if p1.z >= p.zNear && p2.z < p.zNear {
		t := calculateInterpolationParameterZ(p2.z, p1.z, p.zNear)
		z := calculateInterpolationZ(p2.z, p1.z, t)
		x, y := linearInterpolation(p2.x, p2.y, p1.x, p1.y, t)
		return p1, NewVertex(x, y, z), true
	}

	return p1, p2, true

}

// perspectiveProjection projects a 3D point (x, y, z) onto a 2D plane using perspective projection.
func (p Projection) perspectiveProjection(v Vertex) (float64, float64) {
	// perspective projection
	xPoint := p.xCenter + (v.x/v.z)*p.scale // pixelX=xCenter+x′*scale
	yPoint := p.yCenter - (v.y/v.z)*p.scale // pixelY=yCenter-y'*scale
	return xPoint, yPoint
}

func (p Projection) ProjectEdge(p1, p2 Vertex) (float64, float64, float64, float64) {
	// Project the rotated 3D points onto the 2D plane using perspective projection
	sx1, sy1 := p.perspectiveProjection(p1)
	sx2, sy2 := p.perspectiveProjection(p2)
	return sx1, sy1, sx2, sy2
}

type Vertex struct {
	x float64
	y float64
	z float64
}

func NewVertex(x, y, z float64) Vertex {
	return Vertex{x: x, y: y, z: z}
}

// temporarilyTranslate translates a point (x, y, z) to center the object3D/pivot at the origin for rotation.
func (v Vertex) temporarilyTranslate(xCenter, yCenter, zCenter float64) Vertex {
	return NewVertex(v.x-xCenter, v.y-yCenter, v.z-zCenter)
}

// translateItBack translates the coordinates back to their original position after rotation.
func (v Vertex) translateItBack(xCenter, yCenter, zCenter float64) Vertex {
	return NewVertex(v.x+xCenter, v.y+yCenter, v.z+zCenter)
}

func getCosSin(degrees float64) (float64, float64) {
	radians := degrees * (math.Pi / 180) // Convert degrees to radians
	return math.Cos(radians), math.Sin(radians)
}

type Transformer struct {
	Degrees      float64
	Rotation     func(v Vertex, degrees float64) Vertex
	AngularSpeed float64
}

func (t *Transformer) UpdateDegrees() {
	t.Degrees += t.AngularSpeed
}

// Rotation around X axis
func rotationX(v Vertex, degrees float64) Vertex {
	cos, sin := getCosSin(degrees)
	// y' = y*cos(theta) - z*sin(theta)
	// z' = y*sin(theta) + z*cos(theta)
	return NewVertex(v.x, v.y*cos-v.z*sin, v.y*sin+v.z*cos)
}

// Rotation around Y axis
func rotationY(v Vertex, degrees float64) Vertex {
	cos, sin := getCosSin(degrees)
	// x' = x*cos(theta) - z*sin(theta)
	// z' = x*sin(theta) + z*cos(theta)
	return NewVertex(v.x*cos-v.z*sin, v.y, v.x*sin+v.z*cos)
}

// Rotation around Z axis
func rotationZ(v Vertex, degrees float64) Vertex {
	cos, sin := getCosSin(degrees)
	// x' = x*cos(theta) - y*sin(theta)
	// y' = x*sin(theta) + y*cos(theta)
	return NewVertex(v.x*cos-v.y*sin, v.x*sin+v.y*cos, v.z)
}

func (v Vertex) TransformVertex(xCenter, yCenter, zCenter float64, transformers []Transformer) Vertex {
	vertex := v.temporarilyTranslate(xCenter, yCenter, zCenter)
	for _, t := range transformers {
		vertex = t.Rotation(vertex, t.Degrees)
	}
	vertex = vertex.translateItBack(xCenter, yCenter, zCenter)
	return vertex
}

type Edge struct {
	A int
	B int
}

func NewEdge(a, b int) Edge {
	return Edge{A: a, B: b}
}

// lineSteps calculates the number of steps needed to draw a line between two points (x1, y1) and (x2, y2) based on the greater distance in either the x or y direction.
// DDA algorithm is used to determine the number of steps required for drawing a line between two points in a 2D space. The number of steps is determined by the greater distance in either the x or y direction, ensuring that the line is drawn smoothly and accurately.
func lineSteps(x1, y1, x2, y2 float64) int {
	dx := math.Abs(x2 - x1)
	dy := math.Abs(y2 - y1)
	return int(math.Ceil(math.Max(dx, dy))) // Determine the number of steps based on the greater distance
}

type Object3D struct {
	vertices     []Vertex
	edges        []Edge
	transformers []Transformer
}

func NewObject3D(vertices []Vertex, edges []Edge, transformers []Transformer) *Object3D {
	return &Object3D{vertices: vertices, edges: edges, transformers: transformers}
}

func (o *Object3D) SetTransformers(transformers []Transformer) {
	o.transformers = transformers
}

func (o Object3D) BoundingBoxCenter() (float64, float64, float64) {
	verticesLen := len(o.vertices)
	xValues := make([]float64, verticesLen)
	yValues := make([]float64, verticesLen)
	zValues := make([]float64, verticesLen)

	for i, v := range o.vertices {
		xValues[i] = v.x
		yValues[i] = v.y
		zValues[i] = v.z
	}

	return (slices.Min(xValues) + slices.Max(xValues)) / 2, (slices.Min(yValues) + slices.Max(yValues)) / 2, (slices.Min(zValues) + slices.Max(zValues)) / 2
}

func (o Object3D) TransformAllVertices() []Vertex {
	xCenter, yCenter, zCenter := o.BoundingBoxCenter()
	transformedVertices := make([]Vertex, len(o.vertices))
	for i, v := range o.vertices {
		transformedVertices[i] = v.TransformVertex(xCenter, yCenter, zCenter, o.transformers)
	}
	return transformedVertices
}

func (o Object3D) Rasterize(render *Render, projection Projection) {
	transformedVertices := o.TransformAllVertices()
	for _, e := range o.edges {
		p1 := transformedVertices[e.A]
		p2 := transformedVertices[e.B]
		// Clipping edge
		p1, p2, show := projection.ClippingEdge(p1, p2)
		if !show {
			continue
		}
		// ProjectEdge
		sx1, sy1, sx2, sy2 := projection.ProjectEdge(p1, p2)
		// Draw the line corresponding to the edge
		render.DrawLine2D(sx1, sy1, sx2, sy2)
	}
}

func (o *Object3D) Update() {
	for i := range o.transformers {
		o.transformers[i].UpdateDegrees()
	}
}

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
