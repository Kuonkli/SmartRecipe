package database

import (
	"SmartRecipe/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type RecipesTable struct {
	db *sql.DB
}

func newRecipesTable(db *sql.DB) (*RecipesTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS recipes (
        id SERIAL PRIMARY KEY,
        title VARCHAR(100),
        description TEXT,
        instructions TEXT,
        kilocalories DOUBLE PRECISION,
        proteins DOUBLE PRECISION,
        fats DOUBLE PRECISION,
        carbohydrates DOUBLE PRECISION,
        score DOUBLE PRECISION
    );`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating recipes table: %v", err)
	}
	log.Println("Recipes table created successfully")
	return &RecipesTable{db}, nil
}

func (rt *RecipesTable) GetRecipes(limit int, filter string) ([]*models.Recipe, error) {
	query := `
    SELECT id, title, description, instructions, kilocalories, proteins, fats, carbohydrates, score
    FROM recipes`
	if filter == "recent" {
		query += `
		ORDER BY id DESC
		LIMIT $1`
	} else if filter == "high-caloric" {
		query += `
		ORDER BY kilocalories DESC
		LIMIT $1`
	} else if filter == "low-caloric" {
		query += `
		ORDER BY kilocalories
		LIMIT $1`
	} else if filter == "title" {
		query += `
		ORDER BY title
		LIMIT $1`
	} else if filter == "score" {
		query += `
		ORDER BY score DESC
		LIMIT $1`
	} else {
		query += `
		ORDER BY id
    	LIMIT $1`
	}
	rows, err := rt.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting recipes: %v", err)
	}
	defer rows.Close()
	var recipes []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description, &recipe.Instructions, &recipe.Kilocalories, &recipe.Proteins, &recipe.Fats, &recipe.Carbohydrates, &recipe.Score); err != nil {
			return nil, fmt.Errorf("error scanning recipe: %v", err)
		}
		imagesQuery := `SELECT id, recipe_id, image_path FROM recipe_images WHERE recipe_id = $1`
		images, err := rt.db.Query(imagesQuery, recipe.Id)
		if err != nil {
			return nil, fmt.Errorf("error getting recipe images: %v", err)
		}
		var Images []*models.RecipeImage
		for images.Next() {
			var image models.RecipeImage
			if err := images.Scan(&image.Id, &image.RecipeId, &image.ImagePath); err != nil {
				return nil, fmt.Errorf("error scanning recipe image: %v", err)
			}
			Images = append(Images, &image)
		}
		recipe.Images = Images
		recipes = append(recipes, &recipe)
	}
	return recipes, nil
}

func (rt *RecipesTable) GetRecipesByTitle(limit int, filter, search string) ([]*models.Recipe, error) {
	query := `
    SELECT id, title, description, instructions, kilocalories, proteins, fats, carbohydrates, score
    FROM recipes
    WHERE (title ILIKE '%' || $1 || '%' OR SIMILARITY(title, $1) > 0.2)`
	if filter == "recent" {
		query += `
		ORDER BY id DESC
		LIMIT $2`
	} else if filter == "high-caloric" {
		query += `
		ORDER BY kilocalories DESC
		LIMIT $2`
	} else if filter == "low-caloric" {
		query += `
		ORDER BY kilocalories
		LIMIT $2`
	} else if filter == "title" {
		query += `
		ORDER BY title
		LIMIT $2`
	} else if filter == "score" {
		query += `
		ORDER BY score DESC
		LIMIT $2`
	} else {
		query += `
		ORDER BY SIMILARITY(title, $1) DESC
    	LIMIT $2`
	}
	rows, err := rt.db.Query(query, search, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting recipes: %v", err)
	}
	defer rows.Close()
	var recipes []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description, &recipe.Instructions, &recipe.Kilocalories, &recipe.Proteins, &recipe.Fats, &recipe.Carbohydrates, &recipe.Score); err != nil {
			return nil, fmt.Errorf("error scanning recipe: %v", err)
		}
		imagesQuery := `SELECT id, recipe_id, image_path FROM recipe_images WHERE recipe_id = $1`
		images, err := rt.db.Query(imagesQuery, recipe.Id)
		if err != nil {
			return nil, fmt.Errorf("error getting recipe images: %v", err)
		}
		var Images []*models.RecipeImage
		for images.Next() {
			var image models.RecipeImage
			if err := images.Scan(&image.Id, &image.RecipeId, &image.ImagePath); err != nil {
				return nil, fmt.Errorf("error scanning recipe image: %v", err)
			}
			Images = append(Images, &image)
		}
		recipe.Images = Images
		recipes = append(recipes, &recipe)
	}
	return recipes, nil
}

func (rt *RecipesTable) GetRecipeById(id int) (*models.Recipe, error) {
	query := `
    SELECT title, description, instructions, kilocalories, proteins, fats, carbohydrates, score 
    FROM recipes 
    WHERE id = $1`
	row := rt.db.QueryRow(query, id)
	var title, description, instructions string
	var kilocalories, proteins, fats, carbohydrates, score float64
	if err := row.Scan(&title, &description, &instructions, &kilocalories, &proteins, &fats, &carbohydrates, &score); err != nil {
		return nil, err
	}
	imagesQuery := `SELECT id, recipe_id, image_path FROM recipe_images WHERE recipe_id = $1`
	images, err := rt.db.Query(imagesQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error getting recipe images: %v", err)
	}
	var Images []*models.RecipeImage
	for images.Next() {
		var image models.RecipeImage
		if err := images.Scan(&image.Id, &image.RecipeId, &image.ImagePath); err != nil {
			return nil, fmt.Errorf("error scanning recipe image: %v", err)
		}
		Images = append(Images, &image)
	}
	r := &models.Recipe{
		Id:            id,
		Title:         title,
		Description:   description,
		Instructions:  instructions,
		Kilocalories:  kilocalories,
		Proteins:      proteins,
		Fats:          fats,
		Carbohydrates: carbohydrates,
		Score:         score,
		Images:        Images,
	}

	return r, nil
}

func (rt *RecipesTable) AddRecipe(r *models.Recipe) (int, error) {
	var id int
	query := `INSERT INTO recipes (title, description, instructions, kilocalories, proteins, fats, carbohydrates, score) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`
	err := rt.db.QueryRow(query, r.Title, r.Description, r.Instructions, r.Kilocalories, r.Proteins, r.Fats, r.Carbohydrates, r.Score).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error adding recipe: %v", err)
	}
	return id, nil
}

func (rt *RecipesTable) DeleteRecipe(id int) error {
	query := `DELETE FROM recipes WHERE id=$1 `
	_, err := rt.db.Exec(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No recipe found with id %d", id)
			return fmt.Errorf("no recipe found with id %d", id)
		}
		return fmt.Errorf("error deleting recipe: %v", err)
	}
	return nil
}

func (rt *RecipesTable) EditRecipe(r *models.Recipe) error {
	query := `UPDATE recipes SET title=$1, description=$2, instructions=$3, kilocalories=$4, proteins=$5, fats=$6, carbohydrates=$7, score=$8 WHERE id=$9`
	_, err := rt.db.Exec(query, r.Title, r.Description, r.Instructions, r.Kilocalories, r.Proteins, r.Fats, r.Carbohydrates, r.Score, r.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No recipe found with id %d", r.Id)
			return fmt.Errorf("no recipe found with id %d", r.Id)
		}
		return fmt.Errorf("error editing recipe: %v", err)
	}
	return nil
}

func (rt *RecipesTable) UpdateScore(id int) error {
	query := `
        SELECT score 
        FROM personal_recipes 
        WHERE recipe_id = $1`
	rows, err := rt.db.Query(query, id)
	if err != nil {
		return fmt.Errorf("error updating recipe score: %v", err)
	}
	defer rows.Close()

	var sum, count int
	for rows.Next() {
		var score int
		if err := rows.Scan(&score); err != nil {
			return fmt.Errorf("error updating recipe score: %v", err)
		}
		sum += score
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error updating recipe rows: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("no scores found for recipe with id %d", id)
	}

	finalScore := float64(sum) / float64(count)
	query = `UPDATE recipes SET score=$1 WHERE id=$2`
	_, err = rt.db.Exec(query, finalScore, id)
	if err != nil {
		return fmt.Errorf("error updating recipe score: %v", err)
	}

	return nil
}
