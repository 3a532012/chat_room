package friendUsecase

import (
	"context"
	"friend/domain"
	"log"
)

type friendUsecase struct {
	repository domain.FriendRepository
}

// NewMongodbFriendRepository ...
func NewFriendUsecase(repo domain.FriendRepository) domain.FriendUsecase {
	return &friendUsecase{
		repository: repo,
	}
}

func (friendUsecase *friendUsecase) FindByID(ctx context.Context, id string) (*domain.Friend, error) {
	friend, err := friendUsecase.repository.FindByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return friend, nil
}
func (friendUsecase *friendUsecase) FindByName(ctx context.Context, name string) (*domain.Friend, error) {
	friend, err := friendUsecase.repository.FindByName(ctx, name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return friend, nil
}
func (friendUsecase *friendUsecase) Store(ctx context.Context, friend *domain.Friend) (*domain.Friend, error) {

	newFriend, err := friendUsecase.repository.Store(ctx, friend)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newFriend, nil
}

func (friendUsecase *friendUsecase) All(ctx context.Context) ([]*domain.Friend, error) {
	friends, err := friendUsecase.repository.All(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return friends, nil
}

func (friendUsecase *friendUsecase) AddFriend(ctx context.Context, friendID string, tagID string) error {
	friend, err := friendUsecase.repository.FindByTagID(ctx, tagID)
	if err != nil {
		return err
	}
	if err := friendUsecase.repository.AddFriend(ctx, friendID, friend.ID); err != nil {
		return err
	}

	return nil
}

func (friendUsecase *friendUsecase) RemoveFriend(ctx context.Context, friendID string, tagID string) error {
	friend, err := friendUsecase.repository.FindByTagID(ctx, tagID)
	if err != nil {
		return err
	}
	if err := friendUsecase.repository.RemoveFriend(ctx, friendID, friend.ID); err != nil {
		return err
	}

	return nil
}

func (friendUsecase *friendUsecase) FindByTagID(ctx context.Context, tagID string) (*domain.Friend, error) {

	return friendUsecase.repository.FindByTagID(ctx, tagID)
}
func (friendUsecase *friendUsecase) UpdateTagID(ctx context.Context, friendID string, tagID string) (*domain.Friend, error) {

	return friendUsecase.repository.UpdateTagID(ctx, friendID, tagID)
}
