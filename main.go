package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Users struct {
	ID           int
	Name         string
	Password     string
	Email        string
	RegisteredAT time.Time
}

func main() {

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=your_password dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	users, err := getUsers(db)

	fmt.Println(users)

}

func getUsers(db *sql.DB) ([]Users, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	users := make([]Users, 0)
	for rows.Next() {
		u := Users{}
		err := rows.Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.RegisteredAT)
		if err != nil {
			return nil, err
		}
		users = append(users, u)

	}
	return users, nil
}

func getUserByID(db *sql.DB, id int) (Users, error) {
	var u Users
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", 1).
		Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.RegisteredAT)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No rows found")
		}
		return Users{}, err
	}
	return u, err
}

func getUserByName(db *sql.DB, name string) (Users, error) {
	var u Users
	err := db.QueryRow("SELECT * FROM users WHERE name = $1", name).
		Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.RegisteredAT)
	if err != nil {
		return Users{}, err
	}
	return u, err
}

func createUser(db *sql.DB, name string, password string, email string) error {
	_, err := db.Exec("INSERT INTO users (name, password, email) VALUES ($1, $2, $3)", name, password, email)
	if err != nil {
		return err
	}
	return nil
}
