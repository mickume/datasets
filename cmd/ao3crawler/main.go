package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	MIN_LINE_LENGTH = 20
	NEW_LINE_TOKEN  = "\n\n"
	API_MAX_DELAY   = 10
)

// aoc <output_dir> <input_file>
// Example: aoc datasets/ao3_harrypotter input.txt
func main() {

	if len(os.Args) < 3 {
		log.Fatal(fmt.Errorf("invalid arguments"))
	}

	output_dir := os.Args[1]
	input_file := os.Args[2]

	if output_dir == "." {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		output_dir = path
	}
	if !strings.HasPrefix(input_file, "/") {
		input_file = filepath.Join(output_dir, input_file)
	}

	dedupe(input_file)

	file, err := os.Open(input_file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read id's from the input file and retrieve the texts
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if len(id) > 0 && !strings.HasPrefix(id, "#") {
			if err := crawl(id, output_dir); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func crawl(id, path string) error {
	output := fmt.Sprintf("%s/raw/%s.txt", path, id)

	if _, err := os.Stat(output); errors.Is(err, os.ErrNotExist) {
		if err := fetch(id, output); err != nil {
			return err
		}
		randomPause(API_MAX_DELAY)
	}

	return nil
}

func fetch(id, output string) error {
	f, err := create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	c := colly.NewCollector()

	c.OnHTML("div.userstuff p", func(e *colly.HTMLElement) {
		f.WriteString(e.Text + NEW_LINE_TOKEN)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Retrieving ", r.URL)
	})

	url := fmt.Sprintf("https://archiveofourown.org/works/%s?view_full_work=true&view_adult=true", id)
	return c.Visit(url)
}

func dedupe(input string) error {

	dir := filepath.Dir(input)
	tempFile := filepath.Join(dir, "input.temp")
	if err := os.Rename(input, tempFile); err != nil {
		return err
	}

	of, err := os.Open(tempFile)
	if err != nil {
		return err
	}

	nf, err := create(input)
	if err != nil {
		return err
	}

	defer func() {
		nf.Close()
		of.Close()
		os.Remove(tempFile)
	}()

	// the dictionary used to keep unique ids
	ids := make(map[string]bool)

	// read id's from the input file and retrieve the texts
	scanner := bufio.NewScanner(of)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if len(id) > 0 && !strings.HasPrefix(id, "#") {
			if _, ok := ids[id]; !ok {
				if _, err := nf.WriteString(id + "\n"); err != nil {
					return err
				}
				ids[id] = true
			}
		}
	}

	return nil
}

// randomPause enforces a delay to avoid hitting API rate-limits
func randomPause(d int) {
	time.Sleep(time.Duration(rand.Intn(d)) * time.Second)
}

func create(path string) (*os.File, error) {

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(path)
}
