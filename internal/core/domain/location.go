package domain

type LocationUpdate struct {
	Lat       float64 `json:"lat" binding:"required" firestore:"lat"`
	Lng       float64 `json:"lng" binding:"required" firestore:"lng"`
	Timestamp string  `json:"timestamp" firestore:"timestamp"`
	Speed     float64 `json:"speed" firestore:"speed"`
}
