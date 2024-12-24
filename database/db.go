package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

var Database *DatabaseRepository

type DatabaseRepository struct {
	Users           *UsersTable
	Recipes         *RecipesTable
	PersonalRecipes *PersonalRecipesTable
	MealPlans       *MealPlansTable
	PlanRecipes     *PlanRecipesTable
	Ingredients     *IngredientsTable
	Products        *ProductsTable
	Shops           *ShopsTable
	Roles           *RolesTable
	RecipeImages    *RecipeImagesTable
}

func initConfig() {
	viper.SetConfigName("dbconfig")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func NewDatabaseInstance() {
	initConfig()
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error opening database: ", err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
		return
	}
	log.Println("Successfully connected to database")
	Database = &DatabaseRepository{}
	Database.connectTables(db)
}

func (st *DatabaseRepository) connectTables(db *sql.DB) {
	var err error
	if st.Roles, err = newRolesTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.Users, err = newUsersTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.Recipes, err = newRecipesTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.PersonalRecipes, err = newPersonalRecipesTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.MealPlans, err = newMealPlansTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.RecipeImages, err = newRecipeImagesTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.PlanRecipes, err = newPlanRecipesTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.Shops, err = newShopsTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.Products, err = newProductsTable(db); err != nil {
		log.Fatal(err)
		return
	}
	if st.Ingredients, err = newIngredientsTable(db); err != nil {
		log.Fatal(err)
		return
	}
}
