package server

import (
	"SmartRecipe/api/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

type APIServer struct {
	Addr string
}

func (s *APIServer) Run() error {
	app := gin.Default()
	//app.LoadHTMLGlob("frontend/*.html")
	app.GET("/", handlers.GetIndexHandler)

	openGroup := app.Group("/o")
	{
		openGroup.POST("/signup", handlers.SignUp)
		openGroup.POST("/login", handlers.Login)
		openGroup.GET("/refresh", handlers.RefreshToken)
		openGroup.GET("/auth", handlers.GetAuthPageHandler) // Авторизация

	}
	closeGroup := app.Group("/c", handlers.AuthMiddleware())
	{
		closeGroup.GET("/smart_recipe", handlers.GetMainPageHandler)
		//closeGroup.Static("/images", "./frontend/images")
		//closeGroup.Static("/js", "./frontend/js")
		//closeGroup.Static("/css", "./frontend/css")
		recipes := closeGroup.Group("/recipes")
		{
			recipes.GET("/list", handlers.GetRecipesHandler)
			recipes.GET("/:id", handlers.GetRecipePageHandler)
			recipes.GET("/:id/shops/nearest", handlers.GetNearestShopHandler)
			recipes.GET("/:id/shops/cheapest", handlers.GetCheapestShopHandler)
			recipes.PUT(":id/edit/user_recipes", handlers.PutPersonalRecipe)
		}
		closeGroup.GET("/about", handlers.GetAboutUsHandler)
		user := closeGroup.Group("/user")
		{
			user.GET("", handlers.GetUserProfileHandler)
			user.GET("/tried_recipes", handlers.GetTriedRecipesHandler)
			user.GET("/favorite_recipes", handlers.GetFavoriteRecipesHandler)
			user.GET("/rated_recipes", handlers.GetRatedRecipesHandler)
			user.POST("/logout", handlers.Logout)
		}
		mealPlans := closeGroup.Group("/meal_plans")
		{
			mealPlans.GET("", handlers.GetMealPlansHandler)
			mealPlans.GET("/:id", handlers.GetPlanRecipesHandler)
			mealPlans.POST("/generate", handlers.PostMealPlanHandler)
			mealPlans.DELETE("/delete", handlers.DeleteMealPlanHandler)
		}
		admin := closeGroup.Group("/admin", handlers.AdminMiddleware())
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
