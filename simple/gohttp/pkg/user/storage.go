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

func (st *Storage) Get(id uuid.UUID) (*User, error) {
	if userI, ok := st.Load(id); ok {
		if user, ok := userI.(*User); ok {
			return user, nil
		}
	}

	return nil, NoUserFoundError{}
}

func (st *Storage) Create(user *User) {
	st.Store(user.ID, user)
}

func (st *Storage) Delete(id uuid.UUID) error {
	return st.Delete(id)
}
