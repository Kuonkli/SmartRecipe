package server

import (
	"SmartRecipe/api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type APIServer struct {
	Addr string
}

func (s *APIServer) Run() error {
	app := gin.Default()

	// Настройка CORS middleware
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "X-Refresh-Token", "Content-Type"},
		ExposeHeaders:    []string{"Authorization", "X-Refresh-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	app.Use(cors.New(config))

	app.Static("/image-storage", "./image-storage")

	app.POST("/signup", handlers.SignUp)
	app.POST("/login", handlers.Login)
	app.POST("/refresh", handlers.RefreshToken)
	api := app.Group("/api", handlers.AuthMiddleware())
	{
		api.GET("/smart_recipe", handlers.GetMainPageHandler)
		recipes := api.Group("/recipes")
		{
			recipes.GET("/list", handlers.GetRecipesHandler)
			recipes.GET("/:id", handlers.GetRecipePageHandler)
			recipes.GET("/:id/shops/nearest", handlers.GetNearestShopHandler)
			recipes.GET("/:id/shops/cheapest", handlers.GetCheapestShopHandler)
			recipes.PUT(":id/edit/user_recipes", handlers.PutPersonalRecipe)
		}
		api.GET("/about", handlers.GetAboutUsHandler)
		user := api.Group("/user")
		{
			user.GET("/profile", handlers.GetUserProfileHandler)
			user.GET("/tried_recipes", handlers.GetTriedRecipesHandler)
			user.GET("/favorite_recipes", handlers.GetFavoriteRecipesHandler)
			user.GET("/rated_recipes", handlers.GetRatedRecipesHandler)
			user.POST("/logout", handlers.Logout)
		}
		mealPlans := api.Group("/meal_plans")
		{
			mealPlans.GET("", handlers.GetMealPlansHandler)
			mealPlans.GET("/:id", handlers.GetPlanRecipesHandler)
			mealPlans.POST("/generate", handlers.PostMealPlanHandler)
			mealPlans.DELETE("/delete", handlers.DeleteMealPlanHandler)
		}
		admin := api.Group("/admin", handlers.AdminMiddleware())
		{
			usersManagement := admin.Group("/users")
			{
				usersManagement.GET("/")
				usersManagement.DELETE("/delete/:id")
				usersManagement.PUT("/update_role/:id")
			}
			recipeManagement := admin.Group("/recipes")
			{
				recipeManagement.POST("/add/recipes", handlers.PostRecipesHandler)
				recipeManagement.POST("/add/ingredients", handlers.PostIngredientsHandler)
				recipeManagement.POST("/add/images", handlers.PostRecipeImageHandler)
				recipeManagement.DELETE("/delete/image", handlers.DeleteRecipeImageHandler)
				recipeManagement.DELETE("/delete/recipe", handlers.DeleteRecipeHandler)
				recipeManagement.DELETE("/delete/ingredient", handlers.DeleteIngredientHandler)
			}
		}
	}

	err := app.Run(s.Addr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
