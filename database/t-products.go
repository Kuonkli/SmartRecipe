package database

import (
	"SmartRecipe/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type ProductsTable struct {
	db *sql.DB
}

func newProductsTable(db *sql.DB) (*ProductsTable, error) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        title VARCHAR(100),
        category VARCHAR(100),
        subcategory VARCHAR(100),
        price DOUBLE PRECISION,
        kilocalories DOUBLE PRECISION,
        proteins DOUBLE PRECISION,
        fats DOUBLE PRECISION,
        carbohydrates DOUBLE PRECISION,
        shop_id INT REFERENCES shops(id) ON DELETE CASCADE,
        url VARCHAR(250)
    );`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating products table: %v", err)
	}
	log.Println("Products table created successfully")
	return &ProductsTable{db}, nil
}

func (pt *ProductsTable) AddProduct(p *models.Product) error {
	// Проверка существования магазина в таблице shops
	var exists int
	checkShopQuery := `SELECT 1 FROM shops WHERE id = $1`
	err := pt.db.QueryRow(checkShopQuery, p.ShopId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("shop with id %d does not exist", p.ShopId)
		}
		return fmt.Errorf("error checking if shop exists: %v", err)
	}

	// Вставка записи в таблицу products
	query := `INSERT INTO products (title, category, subcategory, price, kilocalories, fats, proteins, carbohydrates, shop_id, url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`
	_, err = pt.db.Exec(query, p.Title, p.Category, p.Subcategory, p.Price, p.Kilocalories, p.Fats, p.Proteins, p.Carbohydrates, p.ShopId, p.URL)
	if err != nil {
		return fmt.Errorf("error adding product: %v", err)
	}
	return nil
}

func (pt *ProductsTable) GetProductById(id int) (*models.Product, error) {
	if id == 0 {
		return nil, fmt.Errorf("error getting product by id, id is empty")
	}
	query := `
    SELECT title, category, subcategory, price, kilocalories, fats, proteins, carbohydrates, shop_id, url 
    FROM products 
    WHERE id = $1`
	row := pt.db.QueryRow(query, id)
	var title, category, subcategory, url string
	var price, kilocalories, proteins, fats, carbohydrates float64
	var shopId int
	if err := row.Scan(&title, &category, &subcategory, &price, &kilocalories, &proteins, &fats, &carbohydrates, &shopId, &url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no product found with id %d", id)
		}
		return nil, fmt.Errorf("error getting product by id: %v", err)
	}

	p := &models.Product{
		Id:            id,
		Title:         title,
		Category:      category,
		Subcategory:   subcategory,
		Price:         price,
		Kilocalories:  kilocalories,
		Proteins:      proteins,
		Fats:          fats,
		Carbohydrates: carbohydrates,
		ShopId:        shopId,
		URL:           url,
	}

	return p, nil
}

func (pt *ProductsTable) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE id=$1`
	_, err := pt.db.Exec(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No product found with id %d", id)
			return fmt.Errorf("no product found with id %d", id)
		}
		return fmt.Errorf("error deleting product: %v", err)
	}
	return nil
}

func (pt *ProductsTable) GetIngredientPriceFromShop(shopId int, ingredient *models.Ingredient) (float64, error) {
	query := `
    SELECT price
    FROM products
    WHERE shop_id = $1
    AND (title ILIKE '%' || $2 || '%' OR SIMILARITY(title, $2) > 0.009)
    AND (category ILIKE '%' || $3 || '%' OR SIMILARITY(category, $3) > 0.5)
    AND (subcategory ILIKE '%' || $4 || '%' OR SIMILARITY(subcategory, $4) > 0.5)
    ORDER BY
        CASE
            WHEN title ILIKE '%' || $2 || '%' THEN 3
            WHEN category ILIKE '%' || $3 || '%' THEN 2
            WHEN subcategory ILIKE '%' || $4 || '%' THEN 1
            ELSE 0
        END DESC,
        SIMILARITY(title, $2) DESC,
        SIMILARITY(category, $3) DESC,
        SIMILARITY(subcategory, $4) DESC
    LIMIT 8`

	rows, err := pt.db.Query(query, shopId, ingredient.Title, ingredient.Category, ingredient.SubCategory)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var sum float64
	var count int
	for rows.Next() {
		var price float64
		if err := rows.Scan(&price); err != nil {
			return 0, err
		}
		sum += price
		count++
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}
	// Check for division by zero
	if count == 0 {
		return 0, nil
	}
	return sum / float64(count), nil
}

func (pt *ProductsTable) EditProduct(p *models.Product) error {
	// Проверка существования магазина в таблице shops
	var exists int
	checkShopQuery := `SELECT 1 FROM shops WHERE id = $1`
	err := pt.db.QueryRow(checkShopQuery, p.ShopId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("shop with id %d does not exist", p.ShopId)
		}
		return fmt.Errorf("error checking if shop exists: %v", err)
	}

	// Обновление записи в таблице products
	query := `UPDATE products SET title=$1, category=$2, subcategory=$3, price=$4, kilocalories=$5, proteins=$6, fats=$7, carbohydrates=$8, shop_id=$9, url=$10 WHERE id=$11`
	_, err = pt.db.Exec(query, p.Title, p.Category, p.Subcategory, p.Price, p.Kilocalories, p.Proteins, p.Fats, p.Carbohydrates, p.ShopId, p.URL, p.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No product found with id %d", p.Id)
			return fmt.Errorf("no product found with id %d", p.Id)
		}
		return fmt.Errorf("error editing product: %v", err)
	}
	return nil
}
