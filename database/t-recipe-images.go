package database

import (
	"SmartRecipe/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type RecipeImagesTable struct {
	db *sql.DB
}

func newRecipeImagesTable(db *sql.DB) (*RecipeImagesTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS recipe_images (
        id SERIAL PRIMARY KEY,
        recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
        image_path VARCHAR(100)
    );`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating recipe_images table: %v", err)
	}
	log.Println("RecipeImages table created successfully")
	return &RecipeImagesTable{db}, nil
}

func (rit *RecipeImagesTable) AddRecipeImage(image *models.RecipeImage) (int, error) {
	var id int
	query := `INSERT INTO recipe_images (recipe_id, image_path) VALUES ($1, $2) RETURNING id;`
	err := rit.db.QueryRow(query, image.RecipeId, image.ImagePath).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error adding recipe image: %v", err)
	}
	return id, nil
}

func (rit *RecipeImagesTable) GetRecipeImageById(id int) ([]byte, error) {
	query := `SELECT image_path FROM recipe_images WHERE id = $1`
	row := rit.db.QueryRow(query, id)
	var image []byte
	if err := row.Scan(&image); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No recipe image found with id %d", id)
			return nil, fmt.Errorf("no recipe image found with id %d", id)
		}
		return nil, fmt.Errorf("error getting recipe image by id: %v", err)
	}
	return image, nil
}

func (rit *RecipeImagesTable) DeleteRecipeImage(id int) error {
	query := `DELETE FROM recipe_images WHERE id=$1`
	_, err := rit.db.Exec(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No recipe image found with id %d", id)
			return fmt.Errorf("no recipe image found with id %d", id)
		}
		return fmt.Errorf("error deleting recipe image: %v", err)
	}
	return nil
}

func (rit *RecipeImagesTable) EditRecipeImage(id int, newImage []byte) error {
	query := `UPDATE recipe_images SET image_path=$1 WHERE id=$2`
	_, err := rit.db.Exec(query, newImage, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No recipe image found with id %d", id)
			return fmt.Errorf("no recipe image found with id %d", id)
		}
		return fmt.Errorf("error editing recipe image: %v", err)
	}
	return nil
}
