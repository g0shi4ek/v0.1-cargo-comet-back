package domain

type CreateObservationRequest struct {
	CometID        *int    `json:"comet_id"`
	RightAscension float64 `json:"right_ascension" binding:"required"`
	Declination    float64 `json:"declination" binding:"required"`
	ObservedAt     string  `json:"observed_at" binding:"required"`
}

type UpdateObservationRequest struct {
	RightAscension float64 `json:"right_ascension" binding:"required"`
	Declination    float64 `json:"declination" binding:"required"`
	ObservedAt     string  `json:"observed_at" binding:"required"`
}

type CreateCometRequest struct {
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}

// чисто для бека
type UpdateCometRequest struct {
	Name                 string  `json:"name"`
	SemiMajorAxis        float64 `json:"semi_major_axis"`
	Eccentricity         float64 `json:"eccentricity"`
	RaanDeg              float64 `json:"raan_deg"`
	AscendingNodeLong    float64 `json:"ascending_node_long"`
	ArgumentOfPerihelion float64 `json:"argument_of_perihelion"`
	TrueAnomalyDeg       string  `json:"true_anomaly_deg"`
}
