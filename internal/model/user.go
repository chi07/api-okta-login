package model

import "time"

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	OktaID    string    `json:"oktaID,omitempty" bson:"oktaID"`
	Name      string    `json:"name,omitempty" bson:"name"`
	Email     string    `json:"email,omitempty" bson:"email"`
	Status    string    `json:"status,omitempty" bson:"status"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updated_at"`
}

const StatusActivated = "activated"
