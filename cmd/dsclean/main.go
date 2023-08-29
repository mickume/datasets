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
	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("invalid arguments"))
	}

	output := os.Args[1]
	length := 15

	if err := CleanupFiles(output, length); err != nil {
		log.Fatal(err)
	}
}

func CleanupFiles(output string, length int) error {

	// scan the dir for files to cleanup

	files, err := os.ReadDir(fmt.Sprintf("%s/raw", output))
	if err != nil {
		return err
	}

	num := 0
	var l int64
	for _, f := range files {
		fname := f.Name()
		if strings.HasSuffix(fname, ".txt") {
			id := strings.Split(fname, ".")[0]
			source := fmt.Sprintf("%s/raw/%s", output, fname)
			target := fmt.Sprintf("%s/data/%s.txt", output, id)

			n, err := CleanAndRewrite(source, target, length)
			if err != nil {
				return err
			}

			l = l + int64(n)
			num++
		}
	}

	b, u := FormatBytes(l)
	fmt.Printf("Cleaned %d files. Total length=%d%s.\n", num, b, u)

	return nil
}

func CleanAndRewrite(source, target string, length int) (int, error) {
	n := 0  // chars written
	lc := 0 // line count

	reader, err := os.Open(source)
	if err != nil {
		return 0, err
	}

	dst, err := Create(target)
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
		line, l, skipped := CleanString(scanner.Text(), length)

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

func CleanString(s string, length int) (string, int, bool) {
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

func Create(path string) (*os.File, error) {

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(path)
}

func FormatBytes(l int64) (int64, string) {
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
