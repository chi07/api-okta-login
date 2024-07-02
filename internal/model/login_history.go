package model

import "time"

type LoginHistory struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	OktaID    string    `json:"oktaID" bson:"oktaID"`
	Email     string    `json:"email" bson:"email"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

const StatusSuccess = "success"
