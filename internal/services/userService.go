package services

import (
	"errors"
	"go-microservice/internal/models"
	"sync"

	"github.com/google/uuid"
)

type UserService struct {
	mu    sync.RWMutex
	users map[string]models.User
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]models.User),
	}
}

func (service *UserService) Create(user models.User) models.User {
	service.mu.Lock()
	defer service.mu.Unlock()

	user.ID = uuid.NewString() // Был хороший способ в проекте структуры нужно глянуть
	service.users[user.ID] = user

	return user
}

func (service *UserService) GetAll() []models.User {
	service.mu.RLock()
	defer service.mu.RUnlock()

	result := make([]models.User, 0, len(service.users))

	for _, u := range service.users {
		result = append(result, u)
	}

	return result
}

func (service *UserService) GetByID(id string) (models.User, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	u, ok := service.users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return u, nil
}

func (service *UserService) Update(id string, user models.User) (models.User, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	old, ok := service.users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}

	old.Name = user.Name
	old.Email = user.Email

	service.users[id] = old

	return old, nil
}

func (service *UserService) Delete(id string) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if _, ok := service.users[id]; !ok {
		return errors.New("user not found")
	}

	delete(service.users, id)

	return nil
}
