package user

import (
	"e-course/domain"
	"e-course/pkg/resp"
	"e-course/pkg/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

// FindByEmail implements domain.UserUsecase.
func (u *userUsecase) FindByEmail(email string) (*domain.User, *resp.ErrorResp) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, &resp.ErrorResp{
				Code: 404,
				Err:  nil,
			}
		}
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Create(data domain.UserCreateRequestBody) (*domain.User, *resp.ErrorResp) {
	existedUser, err := u.userRepo.FindByEmail(data.Email)
	if utils.IsErrorNot404(err) {
		return nil, err
	}

	if existedUser != nil {
		return nil, &resp.ErrorResp{
			Code: 409,
			Err:  errors.New("email is already registered"),
		}
	}

	hashedPwd, errHashedPwd := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if errHashedPwd != nil {
		return nil, &resp.ErrorResp{
			Code: 500,
			Err:  errHashedPwd,
		}
	}

	user := domain.User{
		Name:         data.Email,
		Email:        data.Email,
		Password:     string(hashedPwd),
		CodeVerified: utils.RandString(32),
	}

	dataUser, err := u.userRepo.Create(user)
	if err != nil {
		return nil, &resp.ErrorResp{
			Code: 500,
			Err:  errHashedPwd,
		}
	}

	return dataUser, nil
}

func NewUserUsacase(userUC domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: userUC}
}