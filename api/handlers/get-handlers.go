package handlers

import (
	"SmartRecipe/database"
	"SmartRecipe/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func GetIndexHandler(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/c/smart_recipe")
}

func GetAuthPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{"message": "success"})
}

func GetMainPageHandler(c *gin.Context) {
	recipes, err := database.Database.Recipes.GetRecipes(3, "recent")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "success",
		"recently_added": recipes,
	})
}

func GetRecipesHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be greater than 0"})
		return
	}
	filter := c.Query("filter")
	search, exists := c.GetQuery("search")
	var recipes []*models.Recipe
	if !exists {
		recipes, err = database.Database.Recipes.GetRecipes(limit*10, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		recipes, err = database.Database.Recipes.GetRecipesByTitle(limit*10, filter, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(recipes) == 0 {
			c.JSON(http.StatusNoContent, gin.H{
				"message": "no recipes found",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"recipes": recipes,
	})
}

func GetRecipePageHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Id not found"})
		c.Abort()
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	recipeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	recipe, err := database.Database.Recipes.GetRecipeById(recipeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ingredients, err := database.Database.Ingredients.GetIngredientsByRecipeId(recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := database.Database.PersonalRecipes.Exists(userId, recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var score int
	var isFavorite, isTried bool
	if id == -1 {
		score = 0
		isFavorite = false
		isTried = false
	} else {
		pr, err := database.Database.PersonalRecipes.GetPersonalRecipeById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		score = pr.Score
		isFavorite = pr.IsFavorite
		isTried = pr.IsTried
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"recipe":      recipe,
		"ingredients": ingredients,
		"user_score":  score,
		"is_favorite": isFavorite,
		"is_tried":    isTried,
	})
}

const (
	geocoderAPIKey      = "8667f992-2d31-4a0a-8917-b9e93f3e0e52"
	suggestAPIKey       = "f53ad791-d7fd-461b-8658-72d38ee32627"
	yandexSuggestURL    = "https://suggest-maps.yandex.ru/v1/suggest"
	yandexGeocodeURL    = "https://geocode-maps.yandex.ru/1.x/"
	suggestResultsLimit = 10
)

type coordinates struct {
	Lat, Lon float64
}

type store struct {
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Distance float64 `json:"distance"`
}

type suggestResponse struct {
	Results []struct {
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Distance struct {
			Text  string  `json:"text"`
			Value float64 `json:"value"`
		} `json:"distance"`
		Address struct {
			FormattedAddress string `json:"formatted_address"`
			Component        []struct {
				Name string   `json:"name"`
				Kind []string `json:"kind"`
			} `json:"component"`
		} `json:"address"`
	} `json:"results"`
}

type GeoObjectCollection struct {
	FeatureMember []struct {
		GeoObject struct {
			Point struct {
				Pos string `json:"pos"`
			} `json:"Point"`
		} `json:"GeoObject"`
	} `json:"featureMember"`
}

func GetNearestShopHandler(c *gin.Context) {
	var clientAddress struct {
		Address string `json:"address"`
	}
	// Parse JSON request
	if err := c.ShouldBindJSON(&clientAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	clientAddress.Address = strings.Replace(clientAddress.Address, " ", "+", -1)

	// Get coordinates of client address
	clientCoords, err := getCoordinates(clientAddress.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Get suggestions for stores
	stores, err := getSuggestions(clientCoords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error getting suggestions: %s", err.Error())})
		return
	}
	sort.Slice(stores, func(i, j int) bool {
		return stores[i].Distance < stores[j].Distance
	})

	var nearestShop *models.Shop
	var nearestShopDistance float64
	var nearestShopAddress string
	for _, s := range stores {
		id, err := database.Database.Shops.Exists(s.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if id != -1 {
			nearestShop, err = database.Database.Shops.GetShopById(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			nearestShopDistance = s.Distance
			nearestShopAddress = s.Address
			break
		}
	}
	if nearestShop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No shop found"})
		return
	}

	recipeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if recipeId < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipe Id must be greater than 0"})
		return
	}
	recipe, err := database.Database.Recipes.GetRecipeById(recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ingredients, err := database.Database.Ingredients.GetIngredientsByRecipeId(recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type p struct {
		Title string  `json:"title"`
		Price float64 `json:"approximate_price"`
	}
	var res []p
	var sum float64
	for _, i := range ingredients {
		price, err := database.Database.Products.GetIngredientPriceFromShop(nearestShop.Id, i)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		res = append(res, p{i.Title, price})
		sum += price
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "success",
		"shop_id":    nearestShop.Id,
		"shop_name":  nearestShop.Name,
		"distance":   nearestShopDistance,
		"address":    nearestShopAddress,
		"recipe":     recipe.Title,
		"products":   res,
		"cart_price": sum,
	})
}

func getCoordinates(address string) (coordinates, error) {
	var coords coordinates

	// Create URL for geocode API
	geocodeURL := fmt.Sprintf("%s?apikey=%s&geocode=%s&format=json", yandexGeocodeURL, geocoderAPIKey, url.QueryEscape(address))

	// Make HTTP request
	resp, err := http.Get(geocodeURL)
	if err != nil {
		return coords, err
	}
	defer resp.Body.Close()

	// Decode JSON response
	var result struct {
		Response struct {
			GeoObjectCollection GeoObjectCollection `json:"GeoObjectCollection"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return coords, err
	}

	// Extract coordinates from response
	pos := result.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos
	var lon, lat float64
	fmt.Sscanf(pos, "%f %f", &lon, &lat)
	coords = coordinates{Lat: lat, Lon: lon}
	return coords, nil
}

func getSuggestions(coords coordinates) ([]store, error) {
	var stores []store

	// Create URL for suggest API
	suggestURL := fmt.Sprintf("%s?apikey=%s&text=Supermarket&spn=0.05,0.05&lang=ru_RU&results=%d&strict_bounds=1&print_address=1&ll=%f,%f&ull=%f,%f",
		yandexSuggestURL, suggestAPIKey, suggestResultsLimit, coords.Lon, coords.Lat, coords.Lon, coords.Lat)

	// Make HTTP request
	resp, err := http.Get(suggestURL)
	if err != nil {
		return stores, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return stores, err
	}

	// Decode JSON response
	var result suggestResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return stores, err
	}

	// Extract stores information from response
	for _, item := range result.Results {
		store := store{
			Name:     item.Title.Text,
			Address:  item.Address.FormattedAddress,
			Distance: item.Distance.Value,
		}
		stores = append(stores, store)
	}

	return stores, nil
}

func GetCheapestShopHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "cheapest shop"})
}

func GetAboutUsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{"message": "about us"})
}

func GetUserProfileHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := database.Database.Users.GetUserById(userId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"message": "success",
		"user":    user,
	})
}

func GetTriedRecipesHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	q := c.DefaultQuery("limit", "1")
	limit, err := strconv.Atoi(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be greater than 0"})
		return
	}
	triedRecipes, err := database.Database.PersonalRecipes.GetTriedRecipes(userId, limit*10)

	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"tried_recipes": triedRecipes,
	})
}

func GetFavoriteRecipesHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	q := c.DefaultQuery("limit", "1")
	limit, err := strconv.Atoi(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if limit < 1 {
		limit = 1
	}
	favoriteRecipes, err := database.Database.PersonalRecipes.GetFavoriteRecipes(userId, limit*10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":          "success",
		"favorite_recipes": favoriteRecipes,
	})
}

func GetRatedRecipesHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	q := c.DefaultQuery("limit", "1")
	limit, err := strconv.Atoi(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if limit < 1 {
		limit = 1
	}
	ratedRecipes, err := database.Database.PersonalRecipes.GetRatedRecipes(userId, limit*10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"rated_recipes": ratedRecipes,
	})
}

func GetMealPlansHandler(c *gin.Context) {
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
	q := c.DefaultQuery("limit", "1")
	limit, err := strconv.Atoi(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if limit < 1 {
		limit = 1
	}
	mealPlans, err := database.Database.MealPlans.GetUserMealPlans(userId, limit*10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "success",
		"meal_plans": mealPlans,
	})
}

func GetPlanRecipesHandler(c *gin.Context) {
	var response struct {
		Breakfast struct {
			Recipes []*models.Recipe `json:"recipes"`
		} `json:"breakfast"`
		Lunch struct {
			Recipes []*models.Recipe `json:"recipes"`
		} `json:"lunch"`
		Dinner struct {
			Recipes []*models.Recipe `json:"recipes"`
		} `json:"dinner"`
	}
	planId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid plan ID"})
		return
	}
	planRecipes, err := database.Database.PlanRecipes.GetPlanRecipes(planId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "meal plan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, pr := range planRecipes {
		recipe, err := database.Database.Recipes.GetRecipeById(pr.RecipeId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		switch pr.MealTime {
		case "breakfast":
			response.Breakfast.Recipes = append(response.Breakfast.Recipes, recipe)
		case "lunch":
			response.Lunch.Recipes = append(response.Lunch.Recipes, recipe)
		case "dinner":
			response.Dinner.Recipes = append(response.Dinner.Recipes, recipe)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown meal time"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"meal_plan": response,
	})
}
