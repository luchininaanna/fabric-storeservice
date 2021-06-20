package data

import "time"

type FabricData struct {
	ID        string
	Name      string
	Amount    float32
	Cost      float32
	CreatedAt time.Time
	UpdatedAt *time.Time
}
