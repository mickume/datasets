package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	endOfLine = "\n"
)

var (
	stopWords = []string{
		"notes:",
		"summary:",
		"chapter text",
		"disclaimer:",
		"disclaimers:",
		"main article:",
		"https://",
		"http://",
		"****",
		"....",
		". . .",
		"—--",
		"author note",
		"warnings:",
		"quick notice:",
	}

	quote1 byte = '\''
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal(fmt.Errorf("invalid arguments"))
	}

	min_length := 15
	base_dir := os.Args[1]
	namespace := os.Args[2]
	input_file := os.Args[3]

	if base_dir == "." {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		base_dir = path
	}
	if !strings.HasPrefix(input_file, "/") {
		input_file = filepath.Join(base_dir, namespace, input_file)
	}

	file, err := os.Open(input_file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	num := 0
	var l int64

	// read id's from the input file, clean & move the files
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if len(id) > 0 && !strings.HasPrefix(id, "#") {
			source := fmt.Sprintf("%s/.cache/%s.txt", base_dir, id)
			target := fmt.Sprintf("%s/%s/data/%s.txt", base_dir, namespace, id)

			n, err := cleanAndRewrite(source, target, min_length)
			if err != nil {
				log.Fatal(err)
			}

			l = l + int64(n)
			num++
		}
	}

	b, u := formatBytes(l)
	fmt.Printf("Cleaned %d files. Total length=%d%s.\n", num, b, u)
}

func cleanAndRewrite(source, target string, length int) (int, error) {
	n := 0  // chars written
	lc := 0 // line count

	reader, err := os.Open(source)
	if err != nil {
		return 0, err
	}

	dst, err := create(target)
	if err != nil {
		return 0, err
	}
	writer := bufio.NewWriter(dst)

	defer func() {
		reader.Close()
		writer.Flush()
		dst.Close()
	}()

	scanner := bufio.NewScanner(reader)
	sc := 0

	for scanner.Scan() {
		line, l, skipped := cleanString(scanner.Text(), length)

		if !skipped {
			writer.WriteString(fmt.Sprintf("%s%s", line, endOfLine))
			sc = 0

			n = n + l
		} else {
			if lc == 0 {
				writer.WriteString(fmt.Sprintf("%s%s", line, endOfLine))
			}
			//if sc == 1 {
			//	writer.WriteString(paragraphToken)
			//}
			sc++
		}
		lc++
	}

	return n, nil
}

func cleanString(s string, length int) (string, int, bool) {
	line := strings.Trim(s, " ")

	if len(line) == 0 {
		return "", 0, true
	}

	checks := strings.ToLower(line)

	for _, word := range stopWords {
		if strings.Contains(checks, word) {
			return "", 0, true
		}
	}

	//if len(line) == 0 {
	//	return "", 0, true
	//}

	if line[0] == quote1 {
		line = "\"" + line[1:]
	}

	line = strings.ReplaceAll(line, "***", "")
	line = strings.ReplaceAll(line, "__", "")
	line = strings.ReplaceAll(line, "~*~", "")
	line = strings.ReplaceAll(line, "''", "\" ")
	line = strings.ReplaceAll(line, "‘", "\"")
	line = strings.ReplaceAll(line, "’ ", "\"")
	line = strings.ReplaceAll(line, "“", "\"")
	line = strings.ReplaceAll(line, "”", "\"")
	line = strings.ReplaceAll(line, "' ", "\" ")
	line = strings.ReplaceAll(line, " '", " \"")
	line = strings.ReplaceAll(line, ".'", ".\"")
	line = strings.ReplaceAll(line, "…", "...")
	line = strings.ReplaceAll(line, "\"\"", "\"")

	return line, len(line), len(line) < length
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

func formatBytes(l int64) (int64, string) {
	unit := "b"

	if l < 10*1024 {
		unit = "b"
	} else if l >= 10*1024 && l < 1000*1024 {
		l = l / 1024
		unit = "Kb"
	} else {
		l = l / (1024 * 1024)
		unit = "Mb"
	}

	return l, unit
}
