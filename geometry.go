package main

import (
	"math"
	"slices"
)

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
