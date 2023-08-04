package users

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Name     string    `db:"name"`
	Password string    `db:"password"`
	Role     string    `db:"role"`
}

type (
	UserRequestFormat struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	LoginRequestFormat struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)
type UserResponseFormat struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
	Name     string    `json:"name,omitempty"`
	Password string    `json:"password,omitempty"`
	Role     string    `json:"role,omitempty"`
	Token
}
type Token struct {
	Token string `json:"token"`
}

func (u *Users) Update(req UserRequestFormat) (err error) {
	u.Username = req.Username
	u.Name = req.Name
	return
}

func (u Users) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.ToResponseFormat())
}

func (u *Users) UsersRequestFormat(req UserRequestFormat, id uuid.UUID) (user Users, err error) {
	id, _ = uuid.NewV4()
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return Users{}, err
	}
	user = Users{
		ID:       id,
		Username: req.Username,
		Name:     req.Name,
		Password: hashedPassword,
		Role:     req.Role,
	}
	//NOTE: GA USAH
	users := make([]Users, 0)
	users = append(users, user)
	return
}

func (u *Users) ToResponseFormat() UserResponseFormat {
	token, err := jwt.GenerateJWT(u.ID, u.Username, u.Role)
	if err != nil {
		log.Error().Msg("error generate token")
	}
	return UserResponseFormat{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
		Token:    Token{Token: token},
	}
}

func (u *Users) LoginRequestFormat(req LoginRequestFormat) (user Users, err error) {
	user = Users{
		Username: req.Username,
		Password: req.Password,
	}
	users := make([]Users, 0)
	users = append(users, user)
	return
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
