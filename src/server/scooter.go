package main

type Scooter struct {
	Id                   string `json:"id"`
	IsAvailable          bool   `json:"is_available"`
	TotalDistance        int64  `json:"total_distance"`
	CurrentReservationId string `json:"current_reservation_id"`
}

type ScooterReservation struct {
	ScooterId     string `json:"scooter_id"`
	ReservationId string `json:"reservation_id"`
}

type ScooterRelease struct {
	ScooterId     string `json:"scooter_id"`
	ReservationId string `json:"reservation_id"`
	Distance      int64  `json:"distance"`
}
