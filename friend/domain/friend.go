package domain

import "context"

type Friend struct {
	ID       string `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string `bson:"name" json:"name"`
	Password string `bson:"password" json:"password"`
	Token    string `bson:"token,omitempty" json:"token,omitempty"`
	IsOnline bool   `bson:"isOnline,omitempty" json:"isOnline,omitempty"`
}

type FriendRepository interface {
	FindByID(ctx context.Context, id string) (*Friend, error)
	FindByName(ctx context.Context, name string) (*Friend, error)
	Store(ctx context.Context, friend *Friend) (*Friend, error)
	All(ctx context.Context) ([]*Friend, error)
}

type FriendUsecase interface {
	FindByID(ctx context.Context, id string) (*Friend, error)
	FindByName(ctx context.Context, name string) (*Friend, error)
	Store(ctx context.Context, friend *Friend) (*Friend, error)
	All(ctx context.Context) ([]*Friend, error)
}
