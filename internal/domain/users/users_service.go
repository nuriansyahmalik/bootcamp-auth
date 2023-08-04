package users

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type UsersService interface {
	Create(requestFormat UserRequestFormat, id uuid.UUID) (user Users, err error)
	Login(requestFormat LoginRequestFormat) (user Users, err error)
	ResolveByID(id uuid.UUID) (user Users, err error)
	Update(id uuid.UUID, requestFormat UserRequestFormat) (user Users, err error)
}

type UsersServiceImpl struct {
	UsersRepository UsersRepository
	Config          *configs.Config
}

func ProvideUsersServiceImpl(usersRepository UsersRepository, config *configs.Config) *UsersServiceImpl {
	return &UsersServiceImpl{
		UsersRepository: usersRepository,
		Config:          config,
	}
}

func (u *UsersServiceImpl) Create(requestFormat UserRequestFormat, id uuid.UUID) (user Users, err error) {
	user, err = user.UsersRequestFormat(requestFormat, id)
	if err != nil {
		return
	}
	if err != nil {
		return user, failure.BadRequest(err)
	}
	err = u.UsersRepository.Create(user)
	return
}

func (u *UsersServiceImpl) Login(requestFormat LoginRequestFormat) (user Users, err error) {
	user, err = user.LoginRequestFormat(requestFormat)
	if err != nil {
		return
	}
	user, err = u.UsersRepository.ResolveByUsername(user.Username)
	if err != nil {
		return user, failure.BadRequest(err)
	}
	return
}

func (u *UsersServiceImpl) ResolveByID(id uuid.UUID) (user Users, err error) {
	user, err = u.UsersRepository.ResolveByID(id)
	if err != nil {
		return user, failure.BadRequest(err)
	}
	return
}

func (u *UsersServiceImpl) Update(id uuid.UUID, requestFormat UserRequestFormat) (user Users, err error) {
	user, err = u.UsersRepository.ResolveByID(id)
	if err != nil {
		return
	}
	err = user.Update(requestFormat)
	if err != nil {
		return
	}
	err = u.UsersRepository.Update(user)
	return

}
