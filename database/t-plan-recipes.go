package database

import (
	"SmartRecipe/models"
	"database/sql"
	"fmt"
	"log"
)

type PlanRecipesTable struct {
	db *sql.DB
}

func newPlanRecipesTable(db *sql.DB) (*PlanRecipesTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS plan_recipes (
        id SERIAL PRIMARY KEY,
		plan_id INT REFERENCES meal_plans(id) ON DELETE CASCADE,
		recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
		meal_time VARCHAR(50)  
    );`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating plan_recipes table: %v", err)
	}
	log.Println("PlanRecipes table created successfully")
	return &PlanRecipesTable{db}, nil
}

func (prt *PlanRecipesTable) GetPlanRecipes(planId int) ([]*models.PlanRecipes, error) {
	query := `SELECT id, plan_id, recipe_id, meal_time FROM plan_recipes WHERE plan_id = $1`
	rows, err := prt.db.Query(query, planId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var planRecipes []*models.PlanRecipes
	for rows.Next() {
		var planRecipe models.PlanRecipes
		if err := rows.Scan(&planRecipe.Id, &planRecipe.PlanId, &planRecipe.RecipeId, &planRecipe.MealTime); err != nil {
			return nil, fmt.Errorf("error getting plan recipes: %v", err)
		}
		planRecipes = append(planRecipes, &planRecipe)
	}
	if len(planRecipes) == 0 {
		return nil, sql.ErrNoRows
	}
	return planRecipes, nil
}

func (prt *PlanRecipesTable) AddPlanRecipes(pr *models.PlanRecipes) (int, error) {
	var id int
	query := `INSERT INTO plan_recipes (plan_id, recipe_id, meal_time) VALUES ($1, $2, $3) RETURNING id`
	err := prt.db.QueryRow(query, pr.PlanId, pr.RecipeId, pr.MealTime).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error adding plan recipes: %v", err)
	}
	return id, nil
}
