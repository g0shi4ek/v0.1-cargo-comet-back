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
	Comet          *Comet    `json:"comet,omitempty" gorm:"foreignKey:CometID"`
	IsHorizontal   bool      `json:"is_horizontal"`
}

type Comet struct {
	ID                   int        `json:"id" gorm:"primaryKey"`
	UserID               int        `json:"user_id"`
	Name                 string     `json:"name"`
	PhotoURL             string     `json:"photo_url"`
	SemiMajorAxis        float64    `json:"semi_major_axis"`
	Eccentricity         float64    `json:"eccentricity"`
	RaanDeg              float64    `json:"raan_deg"`
	AscendingNodeLong    float64    `json:"ascending_node_long"`
	ArgumentOfPerihelion float64    `json:"argument_of_perihelion"`
	OrbitActual          bool       `json:"orbit_actual"`
	TrueAnomalyDeg       float64    `json:"true_anomaly_deg"`
	MinApproachDate      *time.Time `json:"min_approach_date"`
	MinApproachDistance  *float64   `json:"min_approach_distance"`
	CloseActual          bool       `json:"close_actual"`
	CalculatedAt         time.Time  `json:"calculated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type CalculationRequest struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	UserID       int    `json:"user_id"`
	CometID      int    `json:"comet_id"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

type OrbitalElements struct {
	SemiMajorAxis        float64
	Eccentricity         float64
	RaanDeg              float64
	AscendingNodeLong    float64
	ArgumentOfPerihelion float64
	TrueAnomalyDeg       float64
}

type CloseApproach struct {
	Date     time.Time
	Distance float64 // в а.е.
}

type TrajectoryPoint struct {
	Time time.Time `json:"time"`
	X    float64   `json:"x"` // Гелиоцентрическая координата X (а.е.)
	Y    float64   `json:"y"` // Гелиоцентрическая координата Y (а.е.)
	Z    float64   `json:"z"` // Гелиоцентрическая координата Z (а.е.)
}

// Trajectory содержит траектории кометы и Земли
type Trajectory struct {
	CometTrajectory []TrajectoryPoint `json:"comet_trajectory"`
	EarthTrajectory []TrajectoryPoint `json:"earth_trajectory"`
}
