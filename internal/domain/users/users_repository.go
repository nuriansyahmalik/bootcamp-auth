package users

import (
	"database/sql"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	usersQueries = struct {
		selectUsers string
		insertUsers string
		updateUsers string
	}{
		selectUsers: `
			SELECT
				u.id,
				u.username,
				u.name,
				u.password,
				u.role
			FROM user u`,
		insertUsers: `
			INSERT INTO user(
				id,
			    username,
			    name,
			    password,
			    role
			) VALUES (
				:id,
			    :username,
			    :name,
			    :password,
			    :role)`,
		updateUsers: `
			UPDATE user
			SET
			    username=:username,
				name=:name
			WHERE id = :id`,
	}
)

type UsersRepository interface {
	Create(user Users) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (user Users, err error)
	ResolveByUsername(username string) (user Users, err error)
	Update(users Users) (err error)
}
type UsersRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideUsersRepositoryMySQL(db *infras.MySQLConn) *UsersRepositoryMySQL {
	return &UsersRepositoryMySQL{DB: db}
}
func (u *UsersRepositoryMySQL) Create(user Users) (err error) {
	stmt, err := u.DB.Write.PrepareNamed(usersQueries.insertUsers)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (u *UsersRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = u.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM user WHERE user.id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (u *UsersRepositoryMySQL) ResolveByID(id uuid.UUID) (user Users, err error) {
	err = u.DB.Read.Get(
		&user,
		usersQueries.selectUsers+" WHERE u.id = ?",
		id.String())
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("user")
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (u *UsersRepositoryMySQL) ResolveByUsername(username string) (user Users, err error) {
	err = u.DB.Read.Get(&user, usersQueries.selectUsers+" WHERE u.username = ?", username)
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("user")
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *UsersRepositoryMySQL) Update(users Users) (err error) {
	exists, err := r.ExistsByID(users.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if !exists {
		err = failure.NotFound("user")
		logger.ErrorWithStack(err)
		return
	}

	// transactionally update the Foo
	// strategy:
	// 1. delete all the Foo's items
	// 2. create a new set of Foo's items
	// 3. update the Foo
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, users); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

// txUpdate updates a Foo transactionally, given the *sqlx.Tx param.
func (r *UsersRepositoryMySQL) txUpdate(tx *sqlx.Tx, users Users) (err error) {
	stmt, err := tx.PrepareNamed(usersQueries.updateUsers)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(users)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
