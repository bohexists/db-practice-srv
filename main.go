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

	users, err := insertUser(db, "test", "test", "test")
	if err != nil {
		log.Fatal(err)
	}

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

func updateUser(db *sql.DB, id int, name string, password string, email string) error {
	_, err := db.Exec("UPDATE users SET name = $1, password = $2, email = $3 WHERE id = $4", name, password, email, id)
	if err != nil {
		return err
	}
	return nil
}

func deleteUser(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func insertUser(db *sql.DB, name string, password string, email string) (Users, error) {
	tx, err := db.Begin()
	if err != nil {
		return Users{}, err
	}

	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Fatal(errRollback)
			}
		} else {
			errCommit := tx.Commit()
			if errCommit != nil {
				log.Fatal(errCommit)
			}
		}
	}()

	_, err = tx.Exec("INSERT INTO users (name, password, email) VALUES ($1, $2, $3)", name, password, email)
	if err != nil {
		return Users{}, err
	}

	u, err := getUserByName(tx, name)
	if err != nil {
		return Users{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Users{}, err
	}

	return u, nil

}
