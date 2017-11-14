package snowcapapi

import (
	"time"
)

type Event struct {
	UUID        string    `json:"uuid"`
	CreatedOn   time.Time `json:"createdOn"`
	TimeStamp   time.Time `json:"timeStamp"`
	Type        string    `json:"type"`
	Depth       float64   `json:"depth"`
	DumpDepth   float64   `json:"dumpDepth"`
	DepthMin    float64   `json:"depthMin"`
	DepthMax    float64   `json:"depthMax"`
	Temp        float64   `json:"temp"`
	TempMin     float64   `json:"tempMin"`
	TempMax     float64   `json:"tempMax"`
	Pressure    float64   `json:"pressure"`
	Humidity    float64   `json:"humidity"`
	Description string    `json:"description"`
}
