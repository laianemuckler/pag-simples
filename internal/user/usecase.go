package user

type UserUsecase interface {
    GetUser(userID int) (*User, error)
    GetAllUsers() ([]User, error)
		ValidateUniqueUser(cpf, email string) error
		SaveUser(user *User) error
}