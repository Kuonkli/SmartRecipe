package handlers

import (
	"SmartRecipe/database"
	"SmartRecipe/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PostIngredientsHandler(c *gin.Context) {
	var req []*models.Ingredient
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, r := range req {
		err := database.Database.Ingredients.AddIngredient(r)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func PostRecipesHandler(c *gin.Context) {
	type FullRecipe struct {
		Title         string              `json:"title"`
		Description   string              `json:"description"`
		Instructions  string              `json:"instructions"`
		Kilocalories  float64             `json:"kilocalories"`
		Proteins      float64             `json:"proteins"`
		Fats          float64             `json:"fats"`
		Carbohydrates float64             `json:"carbohydrates"`
		Score         float64             `json:"score"`
		Ingredients   []models.Ingredient `json:"ingredients"`
	}

	var recipes []FullRecipe
	if err := c.ShouldBindJSON(&recipes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, req := range recipes {
		recipe := &models.Recipe{
			Title:         req.Title,
			Description:   req.Description,
			Instructions:  req.Instructions,
			Kilocalories:  req.Kilocalories,
			Proteins:      req.Proteins,
			Fats:          req.Fats,
			Carbohydrates: req.Carbohydrates,
			Score:         req.Score,
		}

		id, err := database.Database.Recipes.AddRecipe(recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, ingredient := range req.Ingredients {
			ingredient.RecipeId = id
			err = database.Database.Ingredients.AddIngredient(&ingredient)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}
	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func PostMealPlanHandler(c *gin.Context) {
	type UserData struct {
		Title              string `json:"title"`
		Gender             string `json:"gender"`
		Height             int    `json:"height"`
		Age                int    `json:"age"`
		Weight             int    `json:"weight"`
		PreferredCuisine   string `json:"preferred_cuisine"`
		ExcludedCategories string `json:"excluded_categories"`
		Goals              string `json:"goal"`
	}

	type PredictionResponse struct {
		Breakfast []int `json:"breakfast"`
		Lunch     []int `json:"lunch"`
		Dinner    []int `json:"dinner"`
	}

	var req UserData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error marshaling data"})
		return
	}

	resp, err := http.Post("http://localhost:5000/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error making request"})
		return
	}
	defer resp.Body.Close()

	var prediction PredictionResponse
	err = json.NewDecoder(resp.Body).Decode(&prediction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding response"})
		return
	}
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}
	mp := &models.MealPlan{
		UserId:      userId,
		Name:        req.Title,
		Description: "",
	}
	mealPlanId, err := database.Database.MealPlans.AddMealPlan(mp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var recipe models.PlanRecipes
	for _, id := range prediction.Breakfast {
		if id == 0 {
			break
		}
		recipe.PlanId = mealPlanId
		recipe.RecipeId = id
		recipe.MealTime = "breakfast"
		_, err := database.Database.PlanRecipes.AddPlanRecipes(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	for _, id := range prediction.Lunch {
		if id == 0 {
			break
		}
		recipe.PlanId = mealPlanId
		recipe.RecipeId = id
		recipe.MealTime = "lunch"
		_, err := database.Database.PlanRecipes.AddPlanRecipes(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	for _, id := range prediction.Dinner {
		if id == 0 {
			break
		}
		recipe.PlanId = mealPlanId
		recipe.RecipeId = id
		recipe.MealTime = "dinner"
		_, err := database.Database.PlanRecipes.AddPlanRecipes(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	endpoint := fmt.Sprintf("/c/meal_plans/%d", mealPlanId)
	c.Redirect(http.StatusFound, endpoint)
}

func PostRecipeImageHandler(c *gin.Context) {
	var req *models.RecipeImage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := database.Database.RecipeImages.AddRecipeImage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}
