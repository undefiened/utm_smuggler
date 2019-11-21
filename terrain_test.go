package main

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/simple"
)

func TestBresenham(t *testing.T) {
	points := Bresenham(Vector2{0, 0}, Vector2{5, 5})
	assert.Equal(t, 16, len(points))
	assert.Equal(t, Vector2{0, 0}, points[0])
	assert.Equal(t, Vector2{5, 5}, points[15])

	points2 := Bresenham(Vector2{0, 0}, Vector2{1, 3})
	assert.Equal(t, 6, len(points2))
	assert.Equal(t, Vector2{0, 0}, points2[0])
	assert.Equal(t, Vector2{1, 3}, points2[5])
}

func CreateTestTerrain1() *Terrain {
	a := [][]float64{
		{0, 0, 1000, 0, 0},
		{0, 0, 1000, 0, 0},
		{0, 0, 1000, 0, 0},
		{1000, 1000, 1000, 0, 0},
		{0, 0, 1000, 0, 0},
	}

	return InitTerrain(a)
}

func TestCoordConversion(t *testing.T) {
	ter := CreateTestTerrain1()

	n1 := simple.Node(0)
	s, x, y := ter.SliceNodeToCoord(&n1)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)
	assert.Equal(t, 0, s)

	n2 := ter.SliceCoordToNode(0, 4, 2)
	assert.Equal(t, int64(14), n2)
	n3 := simple.Node(n2)
	s2, x2, y2 := ter.SliceNodeToCoord(&n3)
	assert.Equal(t, 4, x2)
	assert.Equal(t, 2, y2)
	assert.Equal(t, 0, s2)

	n4 := ter.SliceCoordToNode(0, 4, 4)
	assert.Equal(t, int64(24), n4)
	n5 := simple.Node(n4)
	s3, x3, y3 := ter.SliceNodeToCoord(&n5)
	assert.Equal(t, 4, x3)
	assert.Equal(t, 4, y3)
	assert.Equal(t, 0, s3)

	n6 := ter.SliceCoordToNode(1, 0, 0)
	assert.Equal(t, int64(25), n6)
	n7 := simple.Node(n6)
	s4, x4, y4 := ter.SliceNodeToCoord(&n7)
	assert.Equal(t, 0, x4)
	assert.Equal(t, 0, y4)
	assert.Equal(t, 1, s4)
}

func TestVisibility(t *testing.T) {
	terrain := CreateTestTerrain1()
	drone := InitDrone(Vector2{0, 0}, Vector2{1, 1}, 0)
	drones := []*Drone{drone}

	vis := terrain.ComputeVisibilityAtTime(0, drones)

	fmt.Print(vis.Visible)

	assert.Equal(t, 5, len(vis.Visible))
	assert.Equal(t, 5, len(vis.Visible[0]))
	assert.Equal(t, true, vis.Visible[0][0])
	assert.Equal(t, true, vis.Visible[0][1])
	assert.Equal(t, true, vis.Visible[1][0])
	assert.Equal(t, true, vis.Visible[1][1])
	assert.Equal(t, true, vis.Visible[2][0])
	assert.Equal(t, true, vis.Visible[2][1])

	assert.Equal(t, true, vis.Visible[0][2])

	falseValues := [][]int{
		{0, 3},
		{0, 4},
		{1, 3},
		{1, 4},
		{2, 3},
		{2, 4},
		{3, 3},
		{3, 4},
		{4, 3},
		{4, 4},

		{4, 0},
		{4, 1},
		{4, 2},
	}

	for _, pair := range falseValues {
		assert.Equal(t, false, vis.Visible[pair[0]][pair[1]])
	}
}

func TestAlgorithm1(t *testing.T) {
	terrain := CreateTestTerrain1()
	drone := InitDrone(Vector2{0, 0}, Vector2{0, 2}, 0)
	drones := []*Drone{drone}

	mg := ComputeGraph(terrain, drones)

	smuggler1 := Smuggler{Start: Vector2{0, 0}, End: Vector2{1, 1}}
	path1, weight1 := mg.ComputeSmugglerPath(&smuggler1)

	assert.Nil(t, path1)
	assert.Equal(t, math.Inf(1), weight1)

	smuggler2 := Smuggler{Start: Vector2{4, 0}, End: Vector2{4, 1}}
	path2, weight2 := mg.ComputeSmugglerPath(&smuggler2)

	assert.Equal(t, 2, len(path2))
	assert.Equal(t, float64(1), weight2)
}

func TestSaving(t *testing.T) {
	terrain := CreateTestTerrain1()
	drone := InitDrone(Vector2{0, 0}, Vector2{0, 2}, 0)
	drones := []*Drone{drone}

	mg := ComputeGraph(terrain, drones)

	smuggler2 := Smuggler{Start: Vector2{4, 0}, End: Vector2{4, 1}}
	path2, weight2 := mg.ComputeSmugglerPath(&smuggler2)

	assert.Equal(t, 2, len(path2))
	assert.Equal(t, float64(1), weight2)

	SaveResultsToFile(mg, &path2, "results/test_saving.json")

	//TODO: add checking of the save file
}

func TestTerrainLoad(t *testing.T) {
	terrain := LoadTerrainFromFile("data/test_terrain.json")

	assert.Equal(t, 5, terrain.Width)
	assert.Equal(t, 5, terrain.Height)

	assert.Equal(t, float64(0), terrain.GetHeight(Vector2{0, 0}))
	assert.Equal(t, float64(1000), terrain.GetHeight(Vector2{2, 0}))
	assert.Equal(t, float64(0), terrain.GetHeight(Vector2{4, 4}))
}
