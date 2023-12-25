package domain

import "context"

type User struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Password string   `bson:"password" json:"password"`
	TagID    string   `bson:"tag-id,omitempty" json:"tag-id,omitempty"`
	Friends  []string `bson:"friends" json:"friends"`
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByTagID(ctx context.Context, tagID string) (*User, error)
	UpdateTagID(ctx context.Context, userID string, tagID string) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	Store(ctx context.Context, user *User) (*User, error)
	AddFriend(ctx context.Context, userID string, tagID string) error
	RemoveFriend(ctx context.Context, userID string, tagID string) error
	TagIDExists(ctx context.Context, tagID string) (bool, error)
	All(ctx context.Context) ([]*User, error)
}

type UserUsecase interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByTagID(ctx context.Context, tagID string) (*User, error)
	UpdateTagID(ctx context.Context, userID string, tagID string) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	Store(ctx context.Context, user *User) (*User, error)
	AddFriend(ctx context.Context, userID string, tagID string) error
	RemoveFriend(ctx context.Context, userID string, tagID string) error
	All(ctx context.Context) ([]*User, error)
}
