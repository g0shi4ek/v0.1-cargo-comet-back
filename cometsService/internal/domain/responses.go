package domain

import "time"

type ObservationResponse struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	CometID        *int      `json:"comet_id"`
	RightAscension float64   `json:"right_ascension"`
	Declination    float64   `json:"declination"`
	ObservedAt     time.Time `json:"observed_at"`
	PhotoURL       string    `json:"photo_url"`
	CreatedAt      time.Time `json:"created_at"`
}

type CometCreatedResponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type CometResponse struct {
	ID                   int        `json:"id"`
	UserID               int        `json:"user_id"`
	Name                 string     `json:"name"`
	SemiMajorAxis        *float64   `json:"semi_major_axis"`
	Eccentricity         *float64   `json:"eccentricity"`
	Inclination          *float64   `json:"inclination"`
	AscendingNodeLong    *float64   `json:"ascending_node_long"`
	ArgumentOfPerihelion *float64   `json:"argument_of_perihelion"`
	TimeOfPerihelion     *time.Time `json:"time_of_perihelion"`
	MinApproachDate      *time.Time `json:"min_approach_date"`
	MinApproachDistance  *float64   `json:"min_approach_distance"`
	CalculatedAt         time.Time  `json:"calculated_at"`
}

type CometOrbitResponse struct {
	ID                   int        `json:"id"`
	SemiMajorAxis        *float64   `json:"semi_major_axis"`
	Eccentricity         *float64   `json:"eccentricity"`
	Inclination          *float64   `json:"inclination"`
	AscendingNodeLong    *float64   `json:"ascending_node_long"`
	ArgumentOfPerihelion *float64   `json:"argument_of_perihelion"`
	TimeOfPerihelion     *time.Time `json:"time_of_perihelion"`
}

type CometDistanceResponse struct {
	ID                   int        `json:"id"`
	MinApproachDate      *time.Time `json:"min_approach_date"`
	MinApproachDistance  *float64   `json:"min_approach_distance"`
	CalculatedAt         time.Time  `json:"calculated_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
