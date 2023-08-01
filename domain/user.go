package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required|email"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type DeleteUserRequest struct {
	ID string `json:"id" binding:"required"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Role     string             `bson:"role" json:"role"`
	Password string             `bson:"password" json:"-"`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	Update(c context.Context, user *User) error
	DeleteById(c context.Context, id string) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
	GetTotal(c context.Context) (int64, error)
}

type UserUsecase interface {
	Create(c context.Context, user *User) error
	Update(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	DeleteById(c context.Context, id string) error
	GetUserByEmail(c context.Context, email string) (User, error)
	GetUserByID(c context.Context, id string) (User, error)
}
