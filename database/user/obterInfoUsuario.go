package user

import (
	"database-api/database"
	"database/sql"
	"fmt"
)

type foundUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func FindUserById(id int) (foundUser, error) {
	db, err := database.ConectarBanco()
	if err != nil {
		return foundUser{}, err
	}

	var user foundUser
	err = db.QueryRow(`SELECT id, name, email FROM "users" WHERE id=$1`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		return foundUser{}, fmt.Errorf("usuário %d não encontrado", id)
	}
	if err != nil {
		return foundUser{}, err
	}

	return user, nil
}
