package user

import (
	"sync"

	"github.com/google/uuid"
)

type NoUserFoundError struct{}

func (e NoUserFoundError) Error() string {
	return "no user found"
}

func NewStorage() *Storage {
	return &Storage{
		sync.Map{},
	}
}

type Storage struct {
	sync.Map
}

func (us *Storage) Get(id uuid.UUID) (*User, error) {
	if userI, ok := us.Load(id); ok {
		if user, ok := userI.(*User); ok {
			return user, nil
		}
	}

	return nil, NoUserFoundError{}
}

func (us *Storage) Create(user *User) {
	us.Store(user.ID, user)
}

func (us *Storage) Delete(id uuid.UUID) error {
	return us.Delete(id)
}
