package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gonum.org/v1/gonum/graph"
)

type JsonHeights struct {
	Heights [][]float64 `json:"heights"`
}

type JsonResults struct {
	MG   *SmugglerMovementGraph `json:smugglerGraph`
	Path [][]int                `json:smugglerPath`
}

func LoadTerrainFromFile(filename string) *Terrain {
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var heights JsonHeights

	json.Unmarshal(byteValue, &heights)

	terrain := InitTerrain(heights.Heights)

	return terrain
}

func SaveResultsToFile(mg *SmugglerMovementGraph, path *[]graph.Node, filename string) {
	resPath := make([][]int, len(*path))

	for i, node := range *path {
		_, x, y := mg.Terrain.SliceNodeToCoordByID(node.ID())
		resPath[i] = []int{x, y}
	}

	file, _ := json.MarshalIndent(JsonResults{MG: mg, Path: resPath}, "", " ")

	_ = ioutil.WriteFile(filename, file, 0644)
}
