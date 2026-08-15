package main

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
