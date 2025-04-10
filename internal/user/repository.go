package user

import "fmt"

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByDocumentNumber(documentNumber string) (*User, error)
	GetUser(userID int) (*User, error)
	GetAllUsers() ([]User, error)
	SaveUser(user *User) error
	UpdateUser(user *User) error
}

type MemoryUserRepository struct {
	users []User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: []User{},
	}
}

func (r *MemoryUserRepository) GetUserByEmail(email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (r *MemoryUserRepository) GetUserByDocumentNumber(documentNumber string) (*User, error) {
	for _, user := range r.users {
		if user.DocumentNumber == documentNumber {
			return &user, nil
		}
	}
	return nil, nil
}

func (r *MemoryUserRepository) GetUser(userID int) (*User, error) {
	for _, user := range r.users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("usuário com ID %d não encontrado", userID)
}

func (r *MemoryUserRepository) GetAllUsers() ([]User, error) {
	return r.users, nil
}

func (r *MemoryUserRepository) SaveUser(user *User) error {
	r.users = append(r.users, *user)
	return nil
}

func (r *MemoryUserRepository) UpdateUser(user *User) error {
	for i, u := range r.users {
		if u.ID == user.ID {
			r.users[i] = *user
			return nil
		}
	}
	return fmt.Errorf("usuário com ID %d não encontrado", user.ID)
}
