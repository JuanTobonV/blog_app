package store

import (
	"database/sql"
	"net/http"

	"github.com/JuanTobonV/blog_app/internal/model"
)

type IUserStore interface { //Definimos el contrato que tendremos de las operaciones de la DB
	getAll() ([]*model.User, error)
	getById(id int) (*model.User, error)
	updateById(id int) (*model.User, error)
	create(user *model.User) (*model.User, error)
	delete(id int) (string, error) 
	
}

type userStore struct {
	db *sql.DB
}

func New(db *sql.DB) IUserStore {
	return &userStore {
		db: db,
	}
}

func (s *userStore) create(user *model.User) (*model.User, error) {
	q := `INSER INTO users (username, password`
	
	
	return nil, nil
}



func (s *userStore) getAll() ([]*model.User, error) {
	q := `SELECT id, username, blogs FROM users`

	rows, err := s.db.Query(q)

	if err != nil {
		return nil,err
	}

	defer rows.Close()

	var users []*model.User

	for rows.Next() { //Next prepares the next result row for reading with the [Rows.Scan] method. It returns true on success, or false if there is no next result row or an error happened while preparing it. [Rows.Err] should be consulted to distinguish between the two cases.
		u := &model.User{} // instanciamos el usuario

		if err := rows.Scan(&u.Id, &u.Username, &u.Blogs); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil

}


func (s * userStore) getById(id int) (*model.User, error){ return nil, nil}

func (s * userStore) updateById(id int) (*model.User, error){ return nil, nil}
func (s * userStore) delete(id int) (string, error) { return "", nil}


