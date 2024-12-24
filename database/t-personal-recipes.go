package database

import (
	"SmartRecipe/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type PersonalRecipesTable struct {
	db *sql.DB
}

func newPersonalRecipesTable(db *sql.DB) (*PersonalRecipesTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS personal_recipes (
        id SERIAL PRIMARY KEY,
        user_id INT REFERENCES users(id) ON DELETE CASCADE,
        recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
        score INT,
        is_favorite BOOLEAN,
        is_tried BOOLEAN
    );`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating personal_recipes table: %v", err)
	}
	log.Println("PersonalRecipes table created successfully")
	return &PersonalRecipesTable{db}, nil
}

func (prt *PersonalRecipesTable) GetPersonalRecipesByUserId(userId int, limit int) ([]*models.PersonalRecipe, error) {
	query := `
        SELECT id, user_id, recipe_id, score, is_favorite, is_tried 
        FROM personal_recipes 
        WHERE user_id = $1
        LIMIT $2`
	rows, err := prt.db.Query(query, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting user personal recipes: %v", err)
	}
	defer rows.Close()

	var recipes []*models.PersonalRecipe
	for rows.Next() {
		var pr models.PersonalRecipe
		if err := rows.Scan(&pr.Id, &pr.UserId, &pr.RecipeId, &pr.Score, &pr.IsFavorite, &pr.IsTried); err != nil {
			return nil, fmt.Errorf("error scanning user personal recipe: %v", err)
		}
		recipes = append(recipes, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows result: %v", err)
	}

	return recipes, nil
}

func (prt *PersonalRecipesTable) GetPersonalRecipesByRecipeId(recipeId int) ([]models.PersonalRecipe, error) {
	query := `
        SELECT id, user_id, recipe_id, score, is_favorite, is_tried 
        FROM personal_recipes 
        WHERE recipe_id = $1`
	rows, err := prt.db.Query(query, recipeId)
	if err != nil {
		return nil, fmt.Errorf("error getting user personal recipes: %v", err)
	}
	defer rows.Close()

	var recipes []models.PersonalRecipe
	for rows.Next() {
		var pr models.PersonalRecipe
		if err := rows.Scan(&pr.Id, &pr.UserId, &pr.RecipeId, &pr.Score, &pr.IsFavorite, &pr.IsTried); err != nil {
			return nil, fmt.Errorf("error scanning user personal recipe: %v", err)
		}
		recipes = append(recipes, pr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows result: %v", err)
	}

	return recipes, nil
}

func (prt *PersonalRecipesTable) GetTriedRecipes(userId int, limit int) ([]models.Recipe, error) {
	query := `
        SELECT r.id, r.title, r.description, r.instructions, r.kilocalories, r.proteins, r.fats, r.carbohydrates, r.score
        FROM personal_recipes pr
        JOIN recipes r ON pr.recipe_id = r.id
        WHERE pr.user_id = $1 AND pr.is_tried = true
        LIMIT $2`
	rows, err := prt.db.Query(query, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting tried recipes: %v", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description, &recipe.Instructions, &recipe.Kilocalories, &recipe.Proteins, &recipe.Fats, &recipe.Carbohydrates, &recipe.Score); err != nil {
			return nil, fmt.Errorf("error scanning tried recipe: %v", err)
		}
		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows result: %v", err)
	}

	return recipes, nil
}

func (prt *PersonalRecipesTable) GetRatedRecipes(userId int, limit int) ([]models.Recipe, error) {
	query := `
        SELECT r.id, r.title, r.description, r.instructions, r.kilocalories, r.proteins, r.fats, r.carbohydrates, r.score
        FROM personal_recipes pr
        JOIN recipes r ON pr.recipe_id = r.id
        WHERE pr.user_id = $1 AND pr.score IS NOT NULL AND pr.score BETWEEN 1 AND 5 
        LIMIT $2`
	rows, err := prt.db.Query(query, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting rated recipes: %v", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description, &recipe.Instructions, &recipe.Kilocalories, &recipe.Proteins, &recipe.Fats, &recipe.Carbohydrates, &recipe.Score); err != nil {
			return nil, fmt.Errorf("error scanning rated recipe: %v", err)
		}
		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows result: %v", err)
	}

	return recipes, nil
}

func (prt *PersonalRecipesTable) GetFavoriteRecipes(userId int, limit int) ([]models.Recipe, error) {
	query := `
        SELECT r.id, r.title, r.description, r.instructions, r.kilocalories, r.proteins, r.fats, r.carbohydrates, r.score
        FROM personal_recipes pr
        JOIN recipes r ON pr.recipe_id = r.id
        WHERE pr.user_id = $1 AND pr.is_favorite = true
        LIMIT $2`
	rows, err := prt.db.Query(query, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting favorite recipes: %v", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description, &recipe.Instructions, &recipe.Kilocalories, &recipe.Proteins, &recipe.Fats, &recipe.Carbohydrates, &recipe.Score); err != nil {
			return nil, fmt.Errorf("error scanning favorite recipe: %v", err)
		}
		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows result: %v", err)
	}

	return recipes, nil
}

func (prt *PersonalRecipesTable) AddPersonalRecipe(pr *models.PersonalRecipe) error {
	// Проверка существования рецепта в таблице recipes
	var exists int
	checkRecipeQuery := `SELECT 1 FROM recipes WHERE id = $1`
	err := prt.db.QueryRow(checkRecipeQuery, pr.RecipeId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("recipe with id %d does not exist", pr.RecipeId)
		}
		return fmt.Errorf("error checking if recipe exists: %v", err)
	}

	// Проверка существования пользователя в таблице users
	checkUserQuery := `SELECT 1 FROM users WHERE id = $1`
	err = prt.db.QueryRow(checkUserQuery, pr.UserId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with id %d does not exist", pr.UserId)
		}
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	// Вставка записи в таблицу personal_recipes
	query := `INSERT INTO personal_recipes (user_id, recipe_id, score, is_favorite, is_tried) VALUES ($1, $2, $3, $4, $5);`
	_, err = prt.db.Exec(query, pr.UserId, pr.RecipeId, pr.Score, pr.IsFavorite, pr.IsTried)
	if err != nil {
		return fmt.Errorf("error adding personal recipe: %v", err)
	}
	return nil
}

func (prt *PersonalRecipesTable) GetPersonalRecipeById(id int) (*models.PersonalRecipe, error) {
	if id == 0 {
		return nil, fmt.Errorf("error getting personal recipe by id, id is empty")
	}
	query := `
    SELECT user_id, recipe_id, score, is_favorite, is_tried 
    FROM personal_recipes 
    WHERE id = $1`
	row := prt.db.QueryRow(query, id)
	var userId, recipeId, score int
	var isFavorite, isTried bool
	if err := row.Scan(&userId, &recipeId, &score, &isFavorite, &isTried); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No personal recipe found with id %d", id)
			return nil, fmt.Errorf("no personal recipe found with id %d", id)
		}
		return nil, fmt.Errorf("error getting personal recipe by id: %v", err)
	}

	pr := &models.PersonalRecipe{
		Id:         id,
		UserId:     userId,
		RecipeId:   recipeId,
		Score:      score,
		IsFavorite: isFavorite,
		IsTried:    isTried,
	}

	return pr, nil
}

func (prt *PersonalRecipesTable) Exists(userId, recipeId int) (int, error) {
	query := `SELECT id FROM personal_recipes WHERE user_id = $1 AND recipe_id = $2`
	row := prt.db.QueryRow(query, userId, recipeId)

	var id int
	err := row.Scan(&id)
	if err != nil { // В случае если не найдено возваращаем -1
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}
		return -1, fmt.Errorf("error checking if personal recipe exists: %v", err)
	}
	return id, nil // Если найдено возвращаем айдишник
}

func (prt *PersonalRecipesTable) DeletePersonalRecipe(id int) error {
	query := `DELETE FROM personal_recipes WHERE id=$1`
	_, err := prt.db.Exec(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No personal recipe found with id %d", id)
			return fmt.Errorf("no personal recipe found with id %d", id)
		}
		return fmt.Errorf("error deleting personal recipe: %v", err)
	}
	return nil
}

func (prt *PersonalRecipesTable) EditPersonalRecipe(pr *models.PersonalRecipe) error {
	var exists int
	checkRecipeQuery := `SELECT 1 FROM recipes WHERE id = $1`
	err := prt.db.QueryRow(checkRecipeQuery, pr.RecipeId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("recipe with id %d does not exist", pr.RecipeId)
		}
		return fmt.Errorf("error checking if recipe exists: %v", err)
	}

	// Проверка существования пользователя в таблице users
	checkUserQuery := `SELECT 1 FROM users WHERE id = $1`
	err = prt.db.QueryRow(checkUserQuery, pr.UserId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with id %d does not exist", pr.UserId)
		}
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	query := `UPDATE personal_recipes SET user_id=$1, recipe_id=$2, score=$3, is_favorite=$4, is_tried=$5 WHERE id=$6`
	_, err = prt.db.Exec(query, pr.UserId, pr.RecipeId, pr.Score, pr.IsFavorite, pr.IsTried, pr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No personal recipe found with id %d", pr.Id)
			return fmt.Errorf("no personal recipe found with id %d", pr.Id)
		}
		return fmt.Errorf("error editing personal recipe: %v", err)
	}
	return nil
}
