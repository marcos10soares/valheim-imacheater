package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
	colly "github.com/gocolly/colly/v2"
)

type Item struct {
	Name         string
	Type         string
	InternalID   string
	Usage        string
	URL          string
	ImageURL     string
	Weight       string
	Stack        string
	Teleportable string
}

const fName = "../www/items_list.json" // to replace existing
// const fName = "./items_list.json"

func main() {
	// flags for operation mode
	scrapeMode := flag.Bool("scrape", false, "mode to scrape wiki, bool, -scrape to use")
	downloadMode := flag.Bool("download", false, "mode to download images, bool, -download to use")
	flag.Parse()

	if *scrapeMode == false && *downloadMode == false {
		fmt.Println("This program needs at least one argument. Use -h or --help for help.")
		os.Exit(1)
	}

	if *scrapeMode {
		fmt.Println("starting scraper...")
		scrape()
	}
	if *downloadMode {
		fmt.Println("downloading images...")
		downloadImages()
	}

}

func scrape() {
	var items []Item
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	c := colly.NewCollector(
		// Visit only domains: valheim.fandom.com, www.valheim.fandom.com
		colly.AllowedDomains("valheim.fandom.com", "www.valheim.fandom.com"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./valheim_fandom_cache"),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		// filter links by starting with "/wiki/" and that do not contain ":" or "?"
		if strings.HasPrefix(url, "/wiki/") && (!strings.Contains(url, ":")) && (!strings.Contains(url, "?")) {
			e.Request.Visit(url)
		}
	})

	// print "visiting" on envery request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// check if this url contains an item and extract details of the item
	// On every element with mw-parser-output class call callback
	c.OnHTML(`.mw-parser-output`, func(e *colly.HTMLElement) {
		attrs := e.ChildAttrs("thead tr .pi-data-label", "data-source")

		// check if is item
		is_item := false
		for _, attr := range attrs {
			if attr == "weight" {
				is_item = true
			}
		}

		if is_item {
			item := Item{
				Name: e.ChildText("aside h2.pi-title"),
				// Type:         item_type,
				// InternalID:   item_id,
				// Usage:        item_usage,
				URL: e.Request.URL.String(),
				// ImageURL:     item_img_url,
				// Weight:       item_weight,
				// Stack:        item_stack,
				// Teleportable: item_teleportable,
			}
			item.ImageURL = e.ChildAttr(".pi-image img", "src")

			keys := e.ChildAttrs(".pi-item.pi-data", "data-source")
			values := e.ChildTexts(".pi-item.pi-data .pi-data-value.pi-font")

			// find value of attributes
			for i, key := range keys {
				if key == "id" {
					item.InternalID = values[i]
				}
				if key == "type" {
					item.Type = values[i]
				}
				if key == "usage" {
					item.Usage = values[i]
				}
				if key == "teleport" {
					item.Teleportable = values[i]
				}
			}

			attrs_values := e.ChildTexts("tbody tr .pi-data-value")

			// find value of attributes
			for i, attr := range attrs {
				if attr == "weight" {
					item.Weight = attrs_values[i]
				}
				if attr == "stack" {
					item.Stack = attrs_values[i]
				}
			}

			items = append(items, item)
		}
	})

	c.Visit("https://valheim.fandom.com/wiki/Items_List")

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(items)
}

func downloadImages() {
	// Open our jsonFile
	jsonFile, err := os.Open(fName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + jsonFile.Name())
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our items array
	var items []Item

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'items' which we defined above
	json.Unmarshal(byteValue, &items)

	// create and start new bar
	bar := pb.StartNew(len(items))

	for _, item := range items {
		// fmt.Println(item.InternalID, item.ImageURL)
		err := DownloadFile("../www/img/items/"+item.InternalID+".png", item.ImageURL)
		if err != nil {
			panic(err)
		}
		// fmt.Println("Downloaded: " + item.InternalID)
		bar.Increment()
	}
	fmt.Println("Download complete.")
	bar.Finish()
}

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
