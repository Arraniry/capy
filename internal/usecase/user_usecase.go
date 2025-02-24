package usecase

type UserUsecase struct {
	repo UserRepository
}

type UserRepository interface {
	GetAll() ([]User, error)
	GetByID(id uint) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (u *UserUsecase) GetAll() ([]User, error) {
	return u.repo.GetAll()
}

func (u *UserUsecase) GetByID(id uint) (*User, error) {
	return u.repo.GetByID(id)
}

func (u *UserUsecase) Create(user *User) error {
	// TODO: Add validation
	return u.repo.Create(user)
}

func (u *UserUsecase) Update(user *User) error {
	// TODO: Add validation
	return u.repo.Update(user)
}

func (u *UserUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
