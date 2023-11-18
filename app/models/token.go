package models

import "github.com/golang-jwt/jwt"

type Token struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	UserID string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty" bson:"user_name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Active bool `json:"active,omitempty" bson:"active,omitempty"`
	*jwt.StandardClaims
}