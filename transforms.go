package main

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
