package user

import (
	"database-api/database"
	"database-api/utils"
	"fmt"
)

type createUserModel struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(name, email, password string) (createUserModel, error) {

	db, err := database.ConectarBanco()
	if err != nil {
		return createUserModel{}, err
	}

	var existId int
	err = db.QueryRow(`SELECT id FROM "users" WHERE email=$1`, email).Scan(&existId)
	if err == nil {
		return createUserModel{}, fmt.Errorf("Email já está sendo usado.")
	}

	var created_user createUserModel

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return createUserModel{}, err
	}

	err = db.QueryRow(
		`INSERT INTO "users"(name,email,password) VALUES ($1, $2, $3) RETURNING name, email, password`,
		name, email, hashedPassword,
	).Scan(&created_user.Name, &created_user.Email, &created_user.Password)

	if err != nil {
		return createUserModel{}, err
	}

	return created_user, err
}
