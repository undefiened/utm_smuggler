package main

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/graph/simple"
)

const DronesSpeed = 1    // px/tick
const PixelInMeters = 50 // meters per pixel
const DronesHeight = 1   // relative

type Terrain struct {
	Width   int
	Height  int
	Heights [][]float64
}

type TerrainVisibility struct {
	Visible [][]bool
}

type Vector2 struct {
	X, Y float64
}

type Drone struct {
	Origin      Vector2
	Destination Vector2
	StartTime   float64
	EndTime     float64
	Velocity    Vector2
}

func InitTerrain(heights [][]float64) *Terrain {
	t := Terrain{}
	t.Heights = heights
	t.Height = len(heights)
	t.Width = len(heights[0])

	return &t
}

func InitDrone(origin Vector2, destination Vector2, startTime float64) *Drone {
	d := Drone{Origin: origin, Destination: destination, StartTime: startTime}
	d.EndTime = d.GetEndTime()
	v := Vector2{}
	v.X = destination.X - origin.X
	v.Y = destination.Y - origin.Y
	v.Normalize()
	d.Velocity = v
	return &d
}

func (x *Vector2) DistanceTo(y *Vector2) float64 {
	d := (x.X-y.X)*(x.X-y.X) + (x.Y-y.Y)*(x.Y-y.Y)
	return math.Sqrt(float64(d))
}

func (x *Vector2) Normalize() {
	length := math.Sqrt(float64(x.X*x.X + x.Y*x.Y))
	x.X = x.X / length
	x.Y = x.Y / length
}

func (d *Drone) GetEndTime() float64 {
	distance := d.Destination.DistanceTo(&d.Origin)
	return distance/DronesSpeed + d.StartTime
}

func (d *Drone) ExistsAtTime(time float64) bool {
	if time >= d.StartTime && time <= d.EndTime {
		return true
	} else {
		return false
	}
}

func (d *Drone) GetPositionAtTime(time float64) (*Vector2, error) {
	if time > d.EndTime || time < d.StartTime {
		return nil, errors.New("Drone does not exist")
	}

	p := Vector2{
		X: d.Origin.X + d.Velocity.X*DronesSpeed*(time-d.StartTime),
		Y: d.Origin.Y + d.Velocity.Y*DronesSpeed*(time-d.StartTime)}

	return &p, nil
}

func (d *Drone) GetVisibilityLine(t *Terrain, time float64, p Vector2) ([]Vector2, []bool) {
	currentPosition, _ := d.GetPositionAtTime(time)
	points := Bresenham(*currentPosition, p)
	visibility := make([]bool, len(points))

	for i, point := range points {
		if t.GetHeight(point) > (DronesHeight+t.GetHeight(*currentPosition)-t.GetHeight(p))*PointToPointDistance(point, p)/PointToPointDistance(p, *currentPosition)+t.GetHeight(p) {
			return nil, visibility
		} else {
			visibility[i] = true //TODO: fix it
		}
	}

	return points, visibility
}

func PointToPointDistance(x, y Vector2) float64 {
	return math.Sqrt(math.Pow(x.X-y.X, 2) + math.Pow(x.Y-y.Y, 2))
}

func Bresenham(p, q Vector2) []Vector2 {
	points := make([]Vector2, 0)
	x1 := int(math.Floor(p.X))
	y1 := int(math.Floor(p.Y))
	x2 := int(math.Floor(q.X))
	y2 := int(math.Floor(q.Y))

	ystep := 0
	xstep := 0
	error := 0
	errorprev := 0

	y := y1
	x := x1
	ddy := 0
	ddx := 0
	dx := x2 - x1
	dy := y2 - y1

	points = append(points, p)

	if dy < 0 {
		ystep = -1
		dy = -dy
	} else {
		ystep = 1
	}

	if dx < 0 {
		xstep = -1
		dx = -dx
	} else {
		xstep = 1
	}

	ddy = 2 * dy
	ddx = 2 * dx

	if ddx >= ddy {
		error = dx
		errorprev = error

		for i := 0; i < dx; i++ {
			x = x + xstep
			error = error + ddy
			if error > ddx {
				y = y + ystep
				error = error - ddx

				if error+errorprev < ddx {
					points = append(points, Vector2{float64(x), float64(y - ystep)})
				} else if error+errorprev > ddx {
					points = append(points, Vector2{float64(x - xstep), float64(y)})
				} else {
					points = append(points, Vector2{float64(x), float64(y - ystep)})
					points = append(points, Vector2{float64(x - xstep), float64(y)})
				}
			}

			points = append(points, Vector2{float64(x), float64(y)})

			errorprev = error
		}
	} else {
		error = dy
		errorprev = error

		for i := 0; i < dy; i++ {
			y = y + ystep
			error = error + ddx

			if error > ddy {
				x = x + xstep
				error = error - ddy

				if error+errorprev < ddy {
					points = append(points, Vector2{float64(x - xstep), float64(y)})
				} else if error+errorprev > ddy {
					points = append(points, Vector2{float64(x), float64(y - ystep)})
				} else {
					points = append(points, Vector2{float64(x - xstep), float64(y)})
					points = append(points, Vector2{float64(x), float64(y - ystep)})
				}
			}

			points = append(points, Vector2{float64(x), float64(y)})
			errorprev = error
		}
	}

	return points
}

func (t *Terrain) ComputeVisibilityAtTime(time float64, drones []*Drone) *TerrainVisibility {
	v := make([][]bool, t.Height)

	for i := range v {
		v[i] = make([]bool, t.Width)
	}

	for x := 0; x < t.Width; x++ {
		for y := 0; y < t.Height; y++ {
			v[y][x] = false
		}
	}

	for x := 0; x < t.Width; x++ {
		for y := 0; y < t.Height; y++ {
			if !v[y][x] {
				for _, drone := range drones {
					if drone.ExistsAtTime(time) {
						points, visibility := drone.GetVisibilityLine(t, time, Vector2{float64(x), float64(y)})

						for i := 0; i < len(points); i++ {
							if visibility[i] {
								point := points[i]
								v[int(point.Y)][int(point.X)] = true
							}
						}
					}
				}
			}
		}
	}

	return &TerrainVisibility{Visible: v}
}

func (t *Terrain) GetHeight(p Vector2) float64 {
	return t.Heights[int(math.Floor(p.Y))][int(math.Floor(p.X))]
}

func (t *Terrain) SliceNodeToCoord(node *simple.Node) (int, int, int) {
	return t.SliceNodeToCoordByID(node.ID())
}

func (t *Terrain) SliceNodeToCoordByID(ID int64) (int, int, int) {
	slice := int(math.Floor(float64(ID) / float64(t.Width*t.Height)))
	localID := int(ID) - slice*t.Width*t.Height
	y := int(math.Floor(float64(localID) / float64(t.Width)))
	x := int(localID) - y*t.Width
	return slice, x, y
}

func (t *Terrain) SliceCoordToNode(slice, x, y int) int64 {
	r := int64(slice*t.Width*t.Height + y*t.Width + x)
	return r
}
