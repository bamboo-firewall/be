package domain

import "context"

type Profile struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type RequestUpdateProfile struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ProfileUsecase interface {
	GetProfileByID(c context.Context, userID string) (*Profile, error)
}
