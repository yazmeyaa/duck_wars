package game

type Coordinate struct {
	X uint16
	Y uint16
}

type Cell struct {
	BuildingType string
	Health       uint32
}

type Field struct {
	Cells map[Coordinate]Cell
}
