package userUsecase

import (
	"context"
	"log"
	"user/domain"
)

type userUsecase struct {
	repository domain.UserRepository
}

// NewMongodbUserRepository ...
func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		repository: repo,
	}
}

func (userUsecase *userUsecase) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := userUsecase.repository.FindByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}
func (userUsecase *userUsecase) FindByName(ctx context.Context, name string) (*domain.User, error) {
	user, err := userUsecase.repository.FindByName(ctx, name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}
func (userUsecase *userUsecase) Store(ctx context.Context, user *domain.User) (*domain.User, error) {

	newUser, err := userUsecase.repository.Store(ctx, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newUser, nil
}

func (userUsecase *userUsecase) All(ctx context.Context) ([]*domain.User, error) {
	users, err := userUsecase.repository.All(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (userUsecase *userUsecase) AddFriend(ctx context.Context, userID string, tagID string) error {
	user, err := userUsecase.repository.FindByTagID(ctx, tagID)
	if err != nil {
		return err
	}
	if err := userUsecase.repository.AddFriend(ctx, userID, user.ID); err != nil {
		return err
	}

	return nil
}

func (userUsecase *userUsecase) RemoveFriend(ctx context.Context, userID string, tagID string) error {
	user, err := userUsecase.repository.FindByTagID(ctx, tagID)
	if err != nil {
		return err
	}
	if err := userUsecase.repository.RemoveFriend(ctx, userID, user.ID); err != nil {
		return err
	}

	return nil
}

func (userUsecase *userUsecase) FindByTagID(ctx context.Context, tagID string) (*domain.User, error) {

	return userUsecase.repository.FindByTagID(ctx, tagID)
}
func (userUsecase *userUsecase) UpdateTagID(ctx context.Context, userID string, tagID string) (*domain.User, error) {

	return userUsecase.repository.UpdateTagID(ctx, userID, tagID)
}
