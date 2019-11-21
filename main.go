package main

func main() {
	terrain := LoadTerrainFromFile("data/50x50.json")
	drone := InitDrone(Vector2{0, 0}, Vector2{0, 2}, 0)
	drones := []*Drone{drone}

	mg := ComputeGraph(terrain, drones)

	smuggler2 := Smuggler{Start: Vector2{40, 40}, End: Vector2{40, 41}}
	path2, _ := mg.ComputeSmugglerPath(&smuggler2)

	SaveResultsToFile(mg, &path2, "results/50x50.json")

}
