package user

import (
	"database-api/database"
	"database/sql"
	"fmt"
)

type userDeleted struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func DeleteUser(id int) (userDeleted, error) {
	db, err := database.ConectarBanco()
	if err != nil {
		return userDeleted{}, err
	}

	var deletedUser userDeleted
	err = db.QueryRow(`SELECT id, name, email FROM "users" WHERE id=$1`, id).Scan(
		&deletedUser.ID,
		&deletedUser.Name,
		&deletedUser.Email,
	)
	if err == sql.ErrNoRows {
		return userDeleted{}, fmt.Errorf("usuário %d não encontrado", id)
	}
	if err != nil {
		return userDeleted{}, err
	}

	_, err = db.Exec(`DELETE FROM "users" WHERE id=$1`, deletedUser.ID)
	if err != nil {
		return userDeleted{}, err
	}

	return deletedUser, nil
}
