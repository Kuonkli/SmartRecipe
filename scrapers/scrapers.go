package scrapers

import "sync"

var (
	categoriesURLs     = make(map[string]bool)
	itemsURLs          = make(map[string]bool)
	basePerekrestokURL = "https://www.perekrestok.ru"
	mu                 sync.Mutex
)

func StartScraping() {
	PerekrestokScraper()
}

func UpdatePrices() {

}
