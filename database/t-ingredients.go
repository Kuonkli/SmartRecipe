package database

import (
	"SmartRecipe/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type IngredientsTable struct {
	db *sql.DB
}

func newIngredientsTable(db *sql.DB) (*IngredientsTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS ingredients (
        id SERIAL PRIMARY KEY,
        recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
        title VARCHAR(100),
        product_category VARCHAR(100),
        product_subcategory VARCHAR(100),
        quantity VARCHAR(100)
    );`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating ingredients table: %v", err)
	}
	log.Println("Ingredients table created successfully")
	return &IngredientsTable{db}, nil
}

func (it *IngredientsTable) AddIngredient(ingredient *models.Ingredient) error {
	// Проверка существования рецепта в таблице recipes
	var exists int
	checkRecipeQuery := `SELECT 1 FROM recipes WHERE id = $1`
	err := it.db.QueryRow(checkRecipeQuery, ingredient.RecipeId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("recipe with id %d does not exist", ingredient.RecipeId)
		}
		return fmt.Errorf("error checking if recipe exists: %v", err)
	}

	// Проверка существования категории и подкатегории в таблице продуктов
	var productExists bool
	checkProductCategoryQuery := `SELECT EXISTS(SELECT 1 FROM products WHERE category = $1 AND subcategory = $2)`
	err = it.db.QueryRow(checkProductCategoryQuery, ingredient.Category, ingredient.SubCategory).Scan(&productExists)
	if err != nil {
		return fmt.Errorf("error checking if product category exists: %v", err)
	}

	if !productExists {
		ingredient.Category = "unknown"
		ingredient.SubCategory = "unknown"
	}

	query := `INSERT INTO ingredients (recipe_id, title, product_category, product_subcategory, quantity) VALUES ($1, $2, $3, $4, $5)`
	_, err = it.db.Exec(query, ingredient.RecipeId, ingredient.Title, ingredient.Category, ingredient.SubCategory, ingredient.Quantity)
	if err != nil {
		return fmt.Errorf("error adding ingredient: %v", err)
	}
	return nil
}

func (it *IngredientsTable) GetIngredientById(id int) (*models.Ingredient, error) {
	if id == 0 {
		return nil, fmt.Errorf("error getting ingredient by id, id is empty")
	}
	query := `
    SELECT recipe_id, title, product_category, product_subcategory, quantity 
    FROM ingredients 
    WHERE id = $1`
	row := it.db.QueryRow(query, id)
	var recipeId int
	var title, productCategory, productSubcategory, quantity string
	if err := row.Scan(&recipeId, &title, &productCategory, &productSubcategory, &quantity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No ingredient found with id %d", id)
			return nil, fmt.Errorf("no ingredient found with id %d", id)
		}
		return nil, fmt.Errorf("error getting ingredient by id: %v", err)
	}

	ingredient := &models.Ingredient{
		Id:          id,
		RecipeId:    recipeId,
		Title:       title,
		Category:    productCategory,
		SubCategory: productSubcategory,
		Quantity:    quantity,
	}

	return ingredient, nil
}

func (it *IngredientsTable) GetIngredientsByRecipeId(recipeId int) ([]*models.Ingredient, error) {
	query := `
    SELECT id, recipe_id, title, product_category, product_subcategory, quantity 
    FROM ingredients 
    WHERE recipe_id = $1`
	rows, err := it.db.Query(query, recipeId)
	if err != nil {
		return nil, fmt.Errorf("error getting ingredients by recipe id: %v", err)
	}
	defer rows.Close()
	var ingredients []*models.Ingredient
	for rows.Next() {
		var ingredient models.Ingredient
		if err := rows.Scan(&ingredient.Id, &ingredient.RecipeId, &ingredient.Title, &ingredient.Category, &ingredient.SubCategory, &ingredient.Quantity); err != nil {
			return nil, fmt.Errorf("error scanning ingredient: %v", err)
		}
		ingredients = append(ingredients, &ingredient)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows result: %v", err)
	}
	return ingredients, nil
}

func (it *IngredientsTable) EditIngredient(ingredient *models.Ingredient) error {
	// Проверка существования рецепта в таблице recipes
	var exists int
	checkRecipeQuery := `SELECT 1 FROM recipes WHERE id = $1`
	err := it.db.QueryRow(checkRecipeQuery, ingredient.RecipeId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("recipe with id %d does not exist", ingredient.RecipeId)
		}
		return fmt.Errorf("error checking if recipe exists: %v", err)
	}

	// Проверка существования категории и подкатегории в таблице продуктов
	var productExists bool
	checkProductCategoryQuery := `SELECT EXISTS(SELECT 1 FROM products WHERE category = $1 AND subcategory = $2)`
	err = it.db.QueryRow(checkProductCategoryQuery, ingredient.Category, ingredient.SubCategory).Scan(&productExists)
	if err != nil {
		return fmt.Errorf("error checking if product category exists: %v", err)
	}

	if !productExists {
		ingredient.Category = "unknown"
		ingredient.SubCategory = "unknown"
	}

	query := `UPDATE ingredients SET recipe_id=$1, title=$2, product_category=$3, product_subcategory=$4, quantity=$5 WHERE id=$6`
	_, err = it.db.Exec(query, ingredient.RecipeId, ingredient.Title, ingredient.Category, ingredient.SubCategory, ingredient.Quantity, ingredient.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No ingredient found with id %d", ingredient.Id)
			return fmt.Errorf("no ingredient found with id %d", ingredient.Id)
		}
		return fmt.Errorf("error editing ingredient: %v", err)
	}
	return nil
}

func (it *IngredientsTable) DeleteIngredient(id int) error {
	query := `DELETE FROM ingredients WHERE id=$1`
	_, err := it.db.Exec(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No product found with id %d", id)
			return fmt.Errorf("no product found with id %d", id)
		}
		return fmt.Errorf("error deleting product: %v", err)
	}
	return nil
}
