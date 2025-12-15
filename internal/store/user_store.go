package store

import (
	"database/sql"

	"github.com/JuanTobonV/blog_app/internal/model"
)

type IUserStore interface { //Definimos el contrato que tendremos de las operaciones de la DB
	GetAll() ([]*model.User, error)
	GetById(id int) (*model.User, error)
	UpdateById(user *model.User) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Delete(id int) (string, error) 
	GetByUsername(username string) (*model.User, error)
	
}

type userStore struct {
	db *sql.DB
}

func New(db *sql.DB) IUserStore {
	return &userStore {
		db: db,
	}
}

func (s *userStore) Create(user *model.User) (*model.User, error) {
    // SQL query with placeholders ($1, $2)
    q := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
    
    // Execute query and scan returned ID into user.Id
    err := s.db.QueryRow(q, user.Username, user.Password).Scan(&user.Id)
    
    if err != nil {
        return nil, err  // Return error if something fails
    }
    
    return user, nil  // Return user with new ID
}



func (s *userStore) GetAll() ([]*model.User, error) {
    // Select only columns that exist in users table
    q := `SELECT id, username FROM users`
    
    rows, err := s.db.Query(q)  // Execute query, get multiple rows
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // Always close rows when done
    
    var users []*model.User  // Create empty slice of user pointers
    
    // Loop through each row
    for rows.Next() {
        u := &model.User{}  // Create new user for each row
        
        // Scan columns into user fields
        if err := rows.Scan(&u.Id, &u.Username); err != nil {
            return nil, err
        }
        
        users = append(users, u)  // Add user to slice
    }
    
    return users, nil
}


func (s *userStore) GetById(id int) (*model.User, error){
	 q := `SELECT id, username, password FROM users WHERE id = $1`

	 user := &model.User{}

	 err := s.db.QueryRow(q, id).Scan(&user.Id, &user.Username, &user.Password)

	 if err == sql.ErrNoRows {
		return nil, nil
	 }

	 if err != nil {
		return nil, err
	 }

	 return user, nil
}

func (s *userStore) GetByUsername(username string) (*model.User, error) {
	q := `SELECT id, username, password FROM users WHERE username = $1`

	user := &model.User{}

	err := s.db.QueryRow(q, username).Scan(&user.Id, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userStore) UpdateById(user *model.User) (*model.User, error){ 
	q := `UPDATE users SET username = $1, password = $2 WHERE id = $3`

	_, err := s.db.Exec(q, user.Username, user.Password, user.Id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userStore) Delete(id int) (string, error) { 

	q := `DELETE FROM users WHERE id = $1`

	_, err := s.db.Exec(q, id)

	return "User deleted sucessfully!", err

}


