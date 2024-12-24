package scrapers

import (
	"SmartRecipe/database"
	"SmartRecipe/models"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

func extractFloatFromString(s string) (float64, error) {
	cleanedString := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == ',' {
			return r
		}
		return -1
	}, s)
	number, err := strconv.ParseFloat(strings.Replace(cleanedString, ",", ".", -1), 64)
	if err != nil {
		return 0, fmt.Errorf("error converting string '%s' to float64: %v", s, err)
	}
	return number, nil
}

func PerekrestokScraper() {
	// Создаем новый коллектор
	categoryCollector := colly.NewCollector(
		colly.CacheDir("./scrapers/colly_cache"),
		//colly.Async(true),
	)
	itemCollector := colly.NewCollector(
		colly.CacheDir("./scrapers/colly_cache"),
		//colly.Async(true),
	)
	infoCollector := colly.NewCollector(
		colly.CacheDir("./scrapers/colly_cache"),
		//colly.Async(true),
	)

	//categoryCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})
	//itemCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})
	//infoCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	log.Println("Начало парсинга категорий")
	// Парсинг категорий
	categoryCollector.OnHTML(".Box-sc-149qidf-0.jziFNa a", parsePerekrestokCategories)
	err := categoryCollector.Visit(basePerekrestokURL + "/cat/")
	if err != nil {
		log.Println("Ошибка посещения url каталога: ", err)
		return
	}
	categoryCollector.Wait()

	log.Println("Начало парсинга товаров")
	// Парсинг товаров в категориях
	itemCollector.OnHTML(".product-card__link", parsePerekrestokItems)
	for categoryURL := range categoriesURLs {
		err := itemCollector.Visit(categoryURL)
		if err != nil {
			log.Println("Посещения url категории: ", err)
			return
		}
	}
	itemCollector.Wait()
	log.Println("Начало парсинга информации о товарах")
	// Парсинг информации о товарах
	infoCollector.OnHTML(".sc-dlfnuX.gGbYZZ", parsePerekrestokInfo)
	for itemURL := range itemsURLs {
		err := infoCollector.Visit(itemURL)
		if err != nil {
			log.Println("Ошибка посещения url товара: ", err)
			return
		}
	}
	infoCollector.Wait()
}

func parsePerekrestokCategories(e *colly.HTMLElement) {
	log.Println("Категории парсятся...")
	mu.Lock()
	defer mu.Unlock()
	category := e.Attr("href")
	fullCategoryURL := basePerekrestokURL + category
	if _, exists := categoriesURLs[fullCategoryURL]; !exists {
		categoriesURLs[fullCategoryURL] = true
	}
}

func parsePerekrestokItems(e *colly.HTMLElement) {
	mu.Lock()
	defer mu.Unlock()
	item := e.Attr("href")
	fullItemURL := basePerekrestokURL + item
	if _, exists := itemsURLs[fullItemURL]; !exists {
		itemsURLs[fullItemURL] = true
	}
	log.Println("Товары парсятся...")
}

func parsePerekrestokInfo(e *colly.HTMLElement) {
	mu.Lock()
	defer mu.Unlock()

	if e.DOM.Find(".product-calories-item__value").Length() == 0 {
		return
	}

	name := e.ChildText(".sc-fubCzh.ibFUIH.product__title")
	if name == "" {
		return
	}
	category := e.ChildText(".sc-kstqJO.jena-DX.product-feature-link")
	subcategory := e.DOM.Find("a.sc-kstqJO.jena-DX.product-feature-link").Eq(4).Text()
	log.Println(category, subcategory)

	calories, err := extractFloatFromString(e.DOM.Find(".product-calories-item__value").Eq(0).Text())
	if err != nil {
		log.Println(err)
		return
	}
	fats, err := extractFloatFromString(e.DOM.Find(".product-calories-item__value").Eq(2).Text())
	if err != nil {
		log.Println(err)
	}
	protein, err := extractFloatFromString(e.DOM.Find(".product-calories-item__value").Eq(1).Text())
	if err != nil {
		log.Println(err)
	}
	carbohydrates, err := extractFloatFromString(e.DOM.Find(".product-calories-item__value").Eq(3).Text())
	if err != nil {
		log.Println(err)
	}
	price, err := extractFloatFromString(e.DOM.Find(".price-card-unit-value").Text())
	if err != nil {
		log.Println(err)
	}

	product := &models.Product{
		Title:         name,
		Category:      category,
		Subcategory:   subcategory,
		Price:         price,
		Kilocalories:  calories,
		Proteins:      protein,
		Fats:          fats,
		Carbohydrates: carbohydrates,
		ShopId:        1,
		URL:           e.Request.URL.String(),
	}
	log.Println("Информация о товарах парсится...", e.Request.URL.String())
	err = database.Database.Products.AddProduct(product)
	if err != nil {
		log.Println(err)
		return
	}
}
