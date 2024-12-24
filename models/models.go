package models

type Ingredient struct {
	Id          int    `json:"id"`
	RecipeId    int    `json:"recipe_id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	Quantity    string `json:"quantity"`
}

type MealPlan struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PersonalRecipe struct {
	Id         int  `json:"id"`
	UserId     int  `json:"user_id"`
	RecipeId   int  `json:"recipe_id"`
	Score      int  `json:"score"`
	IsFavorite bool `json:"is_favorite"`
	IsTried    bool `json:"is_tried"`
}

type Product struct {
	Id            int     `json:"id"`
	Title         string  `json:"title"`
	Category      string  `json:"category"`
	Subcategory   string  `json:"subcategory"`
	Price         float64 `json:"price"`
	Kilocalories  float64 `json:"kilocalories"`
	Proteins      float64 `json:"proteins"`
	Fats          float64 `json:"fats"`
	Carbohydrates float64 `json:"carbohydrates"`
	ShopId        int     `json:"shop_id"`
	URL           string  `json:"url"`
}

type Recipe struct {
	Id            int            `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Instructions  string         `json:"instructions"`
	Kilocalories  float64        `json:"kilocalories"`
	Proteins      float64        `json:"proteins"`
	Fats          float64        `json:"fats"`
	Carbohydrates float64        `json:"carbohydrates"`
	Score         float64        `json:"score"`
	Images        []*RecipeImage `json:"images"`
}

type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Shop struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	RoleId   int    `json:"role"`
}

type RecipeImage struct {
	Id        int    `json:"id"`
	RecipeId  int    `json:"recipe_id"`
	ImagePath string `json:"image_path"`
}

type PlanRecipes struct {
	Id       int    `json:"id"`
	PlanId   int    `json:"plan_id"`
	RecipeId int    `json:"recipe_id"`
	MealTime string `json:"meal_time"`
}
