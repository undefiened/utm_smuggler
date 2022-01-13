package main

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

type Smuggler struct {
	Start Vector2
	End   Vector2
}

var PossibleSmugglerDirections = [9][2]int{
	{0, 0},

	{1, 0},
	{-1, 0},

	{0, 1},
	{0, -1},

	{1, 1},
	{1, -1},
	{-1, 1},
	{-1, -1},
}

type SmugglerMovementGraph struct {
	Terrain          *Terrain
	Drones           []*Drone
	VisibilitySlices []*TerrainVisibility
	Graph            *simple.DirectedGraph
}

func (s *Smuggler) GetStartCoords() (int, int) {
	return int(math.Floor(s.Start.X)), int(math.Floor(s.Start.Y))
}

func (s *Smuggler) GetEndCoords() (int, int) {
	return int(math.Floor(s.End.X)), int(math.Floor(s.End.Y))
}

func (mg *SmugglerMovementGraph) ComputeSmugglerPath(smuggler *Smuggler) ([]graph.Node, float64) {
	fromX, fromY := smuggler.GetStartCoords()
	shortestPaths := path.DijkstraFrom(mg.Graph.Node(mg.Terrain.SliceCoordToNode(0, fromX, fromY)), mg.Graph)

	toX, toY := smuggler.GetEndCoords()

	for slice := 0; slice < len(mg.VisibilitySlices); slice++ {
		path, weight := shortestPaths.To(mg.Terrain.SliceCoordToNode(slice, toX, toY))
		if path != nil {
			return path, weight
		}
	}
	// path, _ := shortestPaths.To(self.g.Node(self.v).ID())
	return nil, math.Inf(1)
}

func ComputeGraph(terrain *Terrain, drones []*Drone) *SmugglerMovementGraph {
	visibilitySlices := ComputeVisibilitySlices(terrain, drones)
	g := simple.NewDirectedGraph()

	for slice := 0; slice < len(visibilitySlices); slice++ {
		for x := 0; x < terrain.Width; x++ {
			for y := 0; y < terrain.Height; y++ {
				g.AddNode(simple.Node(terrain.SliceCoordToNode(slice, x, y)))
			}
		}
	}

	for slice := 0; slice < len(visibilitySlices)-1; slice++ {
		visibilitySlice := visibilitySlices[slice]

		for x := 0; x < terrain.Width; x++ {
			for y := 0; y < terrain.Height; y++ {
				isVisible := visibilitySlice.Visible[y][x]
				if !isVisible {
					for _, pair := range PossibleSmugglerDirections {
						new_x := x + pair[1]
						new_y := y + pair[0]
						if new_y < terrain.Height && new_y >= 0 && new_x < terrain.Width && new_x >= 0 {
							isAnotherVisible := visibilitySlices[slice+1].Visible[new_y][new_x]
							if !isAnotherVisible {
								e := g.NewEdge(g.Node(terrain.SliceCoordToNode(slice, x, y)), g.Node(terrain.SliceCoordToNode(slice+1, new_x, new_y)))
								g.SetEdge(e)
							}
						}
					}

				}

			}
		}
	}

	smugglerData := SmugglerMovementGraph{
		Terrain:          terrain,
		VisibilitySlices: visibilitySlices,
		Drones:           drones,
		Graph:            g,
	}

	return &smugglerData
}

func ComputeVisibilitySlices(terrain *Terrain, drones []*Drone) []*TerrainVisibility {
	minTimeF, maxTimeF := FindMinMaxTimes(drones)
	_, maxTime := int(math.Ceil(minTimeF)), int(math.Floor(maxTimeF))
	numberOfSlices := maxTime

	visibilitySlices := make([]*TerrainVisibility, numberOfSlices+1)

	for time := 0; time <= maxTime; time++ {
		visibilitySlices[time] = terrain.ComputeVisibilityAtTime(float64(time), drones)
	}

	return visibilitySlices
}

func FindMinMaxTimes(drones []*Drone) (float64, float64) {
	min := drones[0].StartTime
	max := drones[0].EndTime

	for _, drone := range drones {
		if drone.StartTime < min {
			min = drone.StartTime
		}

		if drone.EndTime > max {
			max = drone.EndTime
		}
	}

	return min, max
}
