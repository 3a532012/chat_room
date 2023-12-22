package domain

import "context"

type User struct {
	ID       string `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string `bson:"name" json:"name"`
	Password string `bson:"password" json:"password"`
	Token    string `bson:"token,omitempty" json:"token,omitempty"`
	IsOnline bool   `bson:"isOnline,omitempty" json:"isOnline,omitempty"`
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	Store(ctx context.Context, user *User) (*User, error)
	All(ctx context.Context) ([]*User, error)
}

type UserUsecase interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	Store(ctx context.Context, user *User) (*User, error)
	All(ctx context.Context) ([]*User, error)
}
