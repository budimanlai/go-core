package entity

import "time"

type Otp struct {
	ID        int
	Handphone string
	TrxID     string
	PinCode   string
	Status    string
	CreatedAt time.Time
}
