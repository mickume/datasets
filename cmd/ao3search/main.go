package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	MAX_PAGES     = 25
	API_MAX_DELAY = 10

	QS = "/works?commit=Sort+and+Filter&page=%d&work_search[complete]=T&work_search[crossover]=&work_search[date_from]=&work_search[date_to]=&work_search[excluded_tag_names]=&work_search[language_id]=%s&work_search[other_tag_names]=&work_search[query]=&work_search[sort_column]=revised_at&work_search[words_from]=%d&work_search[words_to]="
)

// aos <output_file> <search string> <max_page_to_query>
// Example: aos input.txt "Harry Potter" 1
func main() {
	page := -1
	i := 1

	if len(os.Args) < 3 {
		log.Fatal(fmt.Errorf("invalid arguments"))
	}

	output_file := os.Args[1]
	tags := os.Args[2]

	if len(os.Args) == 4 {
		p, _ := strconv.Atoi(os.Args[3])
		page = p
	} else {
		page = MAX_PAGES
	}

	ids := 0
	for i <= page {
		fmt.Printf("Retrieving IDs on page %d\n", i)

		n, err := search(queryString(tags, "en", 5000, i), output_file)
		if err != nil {
			log.Fatal(err)
		}
		i++
		ids += n

		if n < 20 {
			break
		}
		randomPause(API_MAX_DELAY * 2) // ultra slow, this is not important !
	}

	fmt.Printf("Found %d stories\n", ids)
}

func search(query string, output_file string) (int, error) {
	listCollector := colly.NewCollector()
	listCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*archiveofourown.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	storyCollector := listCollector.Clone()

	if err := createPath(output_file); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(output_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	storiesFound := 0

	// Find and visit all links
	listCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "/works/") {
			parts := strings.Split(link, "/")
			if len(parts) == 3 {
				id := parts[2]
				if !strings.ContainsAny(id, "?#") && id != "search" {
					f.WriteString(fmt.Sprintf("%s\n", id))
					storiesFound++
				}
			}
		}
	})

	storyCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Story Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// start the crawling
	f.WriteString(fmt.Sprintf("# %s\n", query))
	err = listCollector.Visit(query)

	return storiesFound, err
}

// randomPause enforces a delay to avoid hitting API rate-limits
func randomPause(d int) {
	time.Sleep(time.Duration(rand.Intn(d)) * time.Second)
}

func createPath(path string) error {

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func queryString(tags, language string, wordsFrom, page int) string {
	return "https://archiveofourown.org/tags/" + strings.ReplaceAll(tags, " ", "%20") + fmt.Sprintf(QS, page, language, wordsFrom)
}
