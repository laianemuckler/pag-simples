package user

import "fmt"

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserUsecase {
    return &UserService{
        repo: repo,
    }
}

func (s *UserService) GetUser(userID int) (*User, error) {
    return s.repo.GetUser(userID)
}

func (s *UserService) GetAllUsers() ([]User, error) {
    return s.repo.GetAllUsers()
}

func (s *UserService) ValidateUniqueUser(cpf, email string) error {
	userByDoc, _ := s.repo.GetUserByDocumentNumber(cpf)
	if userByDoc != nil {
			return fmt.Errorf("CPF já cadastrado: %s", cpf)
	}

	userByEmail, _ := s.repo.GetUserByEmail(email)
	if userByEmail != nil {
			return fmt.Errorf("e-mail já cadastrado: %s", email)
	}

	return nil
}

func (s *UserService) SaveUser(user *User) error {
	if err := s.ValidateUniqueUser(user.DocumentNumber, user.Email); err != nil {
			return err
	}

	return s.repo.SaveUser(user)
}
