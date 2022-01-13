package main

func main() {
	terrain := LoadTerrainFromFile("data/50x50.json")
	drone := InitDrone(Vector2{0, 0}, Vector2{49, 49}, 0)
	drone2 := InitDrone(Vector2{49, 49}, Vector2{0, 0}, 70)
	drone3 := InitDrone(Vector2{0, 0}, Vector2{49, 49}, 140)
	drones := []*Drone{drone, drone2, drone3}

	mg := ComputeGraph(terrain, drones)

	smuggler2 := Smuggler{Start: Vector2{0, 49}, End: Vector2{49, 0}}
	path2, _ := mg.ComputeSmugglerPath(&smuggler2)

	SaveResultsToFile(mg, &path2, "results/50x50.json")
}
