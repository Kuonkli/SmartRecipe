package database

import (
	"SmartRecipe/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type RolesTable struct {
	db *sql.DB
}

func newRolesTable(db *sql.DB) (*RolesTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS roles (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) UNIQUE
    );
    INSERT INTO roles (name) VALUES ('Client') ON CONFLICT DO NOTHING;
	INSERT INTO roles (name) VALUES ('Administrator') ON CONFLICT DO NOTHING;
    `
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating roles table: %v", err)
	}
	log.Println("Roles table created successfully with default roles")
	return &RolesTable{db}, nil
}

func (rt *RolesTable) GetRoleById(id int) (*models.Role, error) {
	if id == 0 {
		return nil, fmt.Errorf("error getting role by id, id is empty")
	}

	query := `SELECT name FROM roles WHERE id = $1`
	row := rt.db.QueryRow(query, id)
	var name string
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No role found with id %d", id)
			return nil, fmt.Errorf("no role found with id %d", id)
		}
		return nil, fmt.Errorf("error getting role by id: %v", err)
	}

	r := &models.Role{
		Id:   id,
		Name: name,
	}

	return r, nil
}
