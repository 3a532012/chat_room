package domain

import (
	"context"
	"time"
)

const (
	PEDDING = 0
	ACCEPT  = 1
)

type Friend struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	SenderID   string    `bson:"sender_id" json:"sender_id"`
	ReceiverID string    `bson:"receiver_id" json:"receiver_id"`
	Status     int       `bson:"status" json:"status"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

type FriendRepository interface {
	AddFriend(ctx context.Context, senderID string, receiverID string) error
	AcceptFriend(ctx context.Context, id string, receiverID string) error
	RemoveFriend(ctx context.Context, id string, userID string) error
	RequestFriendList(ctx context.Context, senderID string) error
	InviteFriendList(ctx context.Context, receiverID string) error
	FriendList(ctx context.Context, userID string) error
}

type FriendUsecase interface {
	AddFriend(ctx context.Context, senderID string, receiverTag string) error
	AcceptFriend(ctx context.Context, id string, receiverID string) error
	RemoveFriend(ctx context.Context, id string, userID string) error
	RequestFriendList(ctx context.Context, senderID string) error
	InviteFriendList(ctx context.Context, receiverID string) error
	FriendList(ctx context.Context, userID string) error
}
