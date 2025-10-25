package domain

import (
	"time"
)


type Observation struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	UserID         int       `json:"user_id"`
	CometID        *int      `json:"comet_id"`
	RightAscension float64   `json:"right_ascension"`
	Declination    float64   `json:"declination"`
	ObservedAt     time.Time `json:"observed_at"`
	PhotoURL       string    `json:"photo_url"`
	Comet          *Comet    `json:"comet,omitempty" gorm:"foreignKey:CometID"`
}

type Comet struct {
	ID                    int        `json:"id" gorm:"primaryKey"`
	UserID               int        `json:"user_id"`
	Name                 string     `json:"name"`
	SemiMajorAxis        float64    `json:"semi_major_axis"`
	Eccentricity         float64    `json:"eccentricity"`
	Inclination          float64    `json:"inclination"`
	AscendingNodeLong    float64    `json:"ascending_node_long"`
	ArgumentOfPerihelion float64    `json:"argument_of_perihelion"`
	TimeOfPerihelion     time.Time  `json:"time_of_perihelion"`
	MinApproachDate      *time.Time `json:"min_approach_date"`
	MinApproachDistance  *float64   `json:"min_approach_distance"`
	CalculatedAt         time.Time  `json:"calculated_at"`
}

type CalculationRequest struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	UserID       int       `json:"user_id"`
	CometID      int       `json:"comet_id"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
}

type OrbitalElements struct {
	SemiMajorAxis        float64
	Eccentricity         float64
	Inclination          float64
	AscendingNodeLong    float64
	ArgumentOfPerihelion float64
	TimeOfPerihelion     time.Time
}

type CloseApproach struct {
	Date     time.Time
	Distance float64 // в а.е.
}
