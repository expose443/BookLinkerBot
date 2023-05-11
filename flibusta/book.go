package flibusta

import (
	"fmt"

	"github.com/gocolly/colly"
)

var (
	BookStoreLink  = "https://flibusta.site"
	BookSearchLink = "https://flibusta.site/booksearch?ask=%s&chb=on" // chb - checkbox checked only book (search only book)
	BookInfo       = "https://flibusta.site%s"
	Format         = "/b/" // only books link
	BookDownload   = "https://flibusta.site%s/fb2"
)

func GetBookLinks(bookName string) string {
	c := colly.NewCollector(
		colly.AllowedDomains("https://flibusta.site", "flibusta.site"),
	)
	storage := "Book Name | Link to Downlaod\n"
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if len(link) > 3 && link[:3] == Format {
			storage += e.Text + " | " + fmt.Sprintf("http://flibusta.site%s/fb2\n", link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL.String())
	})
	c.Visit(fmt.Sprintf(BookSearchLink, bookName))
	return storage
}
