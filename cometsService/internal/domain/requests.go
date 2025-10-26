package domain

type CreateObservationRequest struct {
	CometID        *int    `json:"comet_id"`
	RightAscension float64 `json:"right_ascension" binding:"required"`
	Declination    float64 `json:"declination" binding:"required"`
	ObservedAt     string  `json:"observed_at" binding:"required"`
	IsHorizontal   bool    `json:"is_horizontal"`
}

type UpdateObservationRequest struct {
	RightAscension float64 `json:"right_ascension" binding:"required"`
	Declination    float64 `json:"declination" binding:"required"`
	ObservedAt     string  `json:"observed_at" binding:"required"`
}

type CreateCometRequest struct {
	Name     string `form:"name" binding:"required"`
	PhotoURL string `json:"photo_url"`
}

// чисто для бека
type UpdateCometRequest struct {
	Name                 string  `json:"name"`
	PhotoURL             string  `json:"photo_url"`
	SemiMajorAxis        float64 `json:"semi_major_axis"`
	Eccentricity         float64 `json:"eccentricity"`
	RaanDeg              float64 `json:"raan_deg"`
	AscendingNodeLong    float64 `json:"ascending_node_long"`
	ArgumentOfPerihelion float64 `json:"argument_of_perihelion"`
	TrueAnomalyDeg       string  `json:"true_anomaly_deg"`
}

type GetTrajectoryRequest struct {
	StartTime string `form:"start_time" binding:"required"` // "2006-01-02T15:04:05Z"
	EndTime   string `form:"end_time" binding:"required"`   // "2006-01-02T15:04:05Z"
	NumPoints int    `form:"num_points" binding:"required,min=10,max=1000"`
}
