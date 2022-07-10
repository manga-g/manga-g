package prone

// import colly for scraping url
import (
    "fmt"
    "net/url"
    "time"

    "github.com/PuerkitoBio/goquery"
    _ "github.com/PuerkitoBio/goquery"
    "github.com/gocolly/colly"
)

type PageInfo struct {
    StatusCode int
    Links      map[string]int
}

type Scraper struct {
    Site     *url.URL
    SiteList []string
    Info     *PageInfo
    Selector string
    Attr     string
    Value    string
    Colly    *colly.Collector
}

//func init() {
//    // create a new scraper
//    s := Scraper{
//        Site:     &url.URL{Scheme: "https", Host: "www.google.com"},
//        Selector: "a[href]",
//        Attr:     "",
//        Value:    "",
//        Colly:    colly.NewCollector(colly.AllowedDomains("www.google.com")),
//    }
//    // scrape the site
//    s.Scrape()
//}

// func that returns a new scraper
func NewScraper(site string, selector string, attr string, value string) *Scraper {
    return &Scraper{
        Site:     &url.URL{Scheme: "https", Host: site},
        Selector: selector,
        Attr:     attr,
        Value:    value,
        Colly:    colly.NewCollector(colly.AllowedDomains(site)),
    }
}
func (scrape *Scraper) GetSite() string {
    return scrape.Site.String()
}

func (scrape *Scraper) GetSelector() string {
    return scrape.Selector
}
func (scrape *Scraper) AddSite(site string) {
    scrape.SiteList = append(scrape.SiteList, site)
}
func (scrape *Scraper) SetSite(site string) {
    //scrape.Site = &url.URL{Scheme: "https", Host: site}
    scrape.AddSite(site)
    scrape.Colly.AllowedDomains = scrape.SiteList
}

func (scrape *Scraper) GetSites() []string {
    return scrape.SiteList
}

func (scrape *Scraper) SetAttr(attr string) {
    scrape.Attr = attr
}

func (scrape *Scraper) SetSelector(selector string) {
    scrape.Selector = selector
}

func (scrape *Scraper) Scrape() {
    scrape.Colly.AllowedDomains = scrape.SiteList
    scrape.Colly.OnHTML("article", func(element *colly.HTMLElement) {
        metaTags := element.DOM.ParentsUntil("~").Find("meta")
        metaTags.Each(func(_ int, s *goquery.Selection) {})
    })

    scrape.Colly.OnHTML(scrape.Selector, func(element *colly.HTMLElement) {
        link := element.Attr(scrape.Attr)
        err := element.Request.Visit(element.Request.AbsoluteURL(link))
        if err != nil {
            return
        }
    })

    err := scrape.Colly.Limit(&colly.LimitRule{DomainGlob: "*", RandomDelay: 5 * time.Second})
    if err != nil {
        return
    }

    scrape.Colly.OnRequest(func(request *colly.Request) {
        // request.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
        fmt.Println("Visiting:", request.URL.String())
    })

    err2 := scrape.Colly.Visit(scrape.Site.String())
    if err2 != nil {
        fmt.Println("Error Visiting site:", err2)
        return
    }
}
