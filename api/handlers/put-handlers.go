package handlers

import (
	"SmartRecipe/database"
	"SmartRecipe/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PutPersonalRecipe(c *gin.Context) {
	var req struct {
		Score      int  `json:"score"`
		IsFavorite bool `json:"is_favorite"`
		IsTried    bool `json:"is_tried"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Id not found"})
		c.Abort()
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user Id format"})
		return
	}
	recipeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	personalRecipeId, err := database.Database.PersonalRecipes.Exists(userId, recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if personalRecipeId != -1 {
		existingRecipe, err := database.Database.PersonalRecipes.GetPersonalRecipeById(personalRecipeId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		existingRecipe.IsTried = req.IsTried
		existingRecipe.IsFavorite = req.IsFavorite
		existingRecipe.Score = req.Score
		if err := database.Database.PersonalRecipes.EditPersonalRecipe(existingRecipe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !existingRecipe.IsTried && !existingRecipe.IsFavorite && existingRecipe.Score == 0 {
			err := database.Database.PersonalRecipes.DeletePersonalRecipe(existingRecipe.Id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if existingRecipe.Score != 0 {
			err := database.Database.Recipes.UpdateScore(recipeId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Recipe updated"})
	} else {
		newRecipe := &models.PersonalRecipe{
			UserId:     userId,
			RecipeId:   recipeId,
			Score:      req.Score,
			IsTried:    req.IsTried,
			IsFavorite: req.IsFavorite,
		}
		if err := database.Database.PersonalRecipes.AddPersonalRecipe(newRecipe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if newRecipe.Score != 0 {
			err := database.Database.Recipes.UpdateScore(recipeId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusCreated, gin.H{"message": "New recipe added"})
	}
}
