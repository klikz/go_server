package model

type Component struct {
	Available float64 `json:"available"`
	ID        int

	Code          string
	Name          string
	Checkpoint    string
	Checkpoint_id int
	Unit          string
	Specs         string
	Photo         string
	Time          string
	Type          string
	Type_id       int
	Weight        float64
}
