package database

import (
	"SmartRecipe/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type ShopsTable struct {
	db *sql.DB
}

func newShopsTable(db *sql.DB) (*ShopsTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS shops (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) UNIQUE NOT NULL,
        score DOUBLE PRECISION
    );`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating shops table: %v", err)
	}
	log.Println("Shops table created successfully")
	return &ShopsTable{db}, nil
}

func (st *ShopsTable) AddShop(s *models.Shop) error {
	query := `INSERT INTO shops (name, score) VALUES ($1, $2);`
	_, err := st.db.Exec(query, s.Name, s.Score)
	if err != nil {
		return fmt.Errorf("error adding shop: %v", err)
	}
	return nil
}

func (st *ShopsTable) GetShopById(id int) (*models.Shop, error) {
	if id == 0 {
		return nil, fmt.Errorf("error getting shop by id, id is empty")
	}
	query := `
    SELECT name, score 
    FROM shops 
    WHERE id = $1`
	row := st.db.QueryRow(query, id)
	var name string
	var score float64
	if err := row.Scan(&name, &score); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No shop found with id %d", id)
			return nil, fmt.Errorf("no shop found with id %d", id)
		}
		return nil, fmt.Errorf("error getting shop by id: %v", err)
	}

	s := &models.Shop{
		Id:    id,
		Name:  name,
		Score: score,
	}

	return s, nil
}

func (st *ShopsTable) Exists(name string) (int, error) {
	query := `
    SELECT id
    FROM shops
    WHERE name ILIKE '%' || $1 || '%' OR SIMILARITY(name, $1) > 0.4
    ORDER BY
        CASE
            WHEN name ILIKE '%' || $1 || '%' THEN 2
            WHEN SIMILARITY(name, $1) > 0.4 THEN 1
            ELSE 0
        END DESC,
        SIMILARITY(name, $1) DESC
    LIMIT 1`

	var id int
	err := st.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}

func (st *ShopsTable) DeleteShop(id int) error {
	query := `DELETE FROM shops WHERE id=$1`
	_, err := st.db.Exec(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No shop found with id %d", id)
			return fmt.Errorf("no shop found with id %d", id)
		}
		return fmt.Errorf("error deleting shop: %v", err)
	}
	return nil
}

func (st *ShopsTable) EditShop(s *models.Shop) error {
	query := `UPDATE shops SET name=$1, score=$2 WHERE id=$3`
	_, err := st.db.Exec(query, s.Name, s.Score, s.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No shop found with id %d", s.Id)
			return fmt.Errorf("no shop found with id %d", s.Id)
		}
		return fmt.Errorf("error editing shop: %v", err)
	}
	return nil
}
