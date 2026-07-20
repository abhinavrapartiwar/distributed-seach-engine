package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type page struct {
	Title    string   `xml:"title"`
	Revision revision `xml:"revision"`
}

type revision struct {
	Text string `xml:"text"`
}

type document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

const (
	maxBodyLen = 2000
	minBodyLen = 50
)

func cleanWikiText(raw string) string {
	r := strings.NewReplacer(
		"'''", "",
		"''", "",
		"[[", "",
		"]]", "",
		"{{", "",
		"}}", "",
		"[", "",
		"]", "",
		"\n", " ",
		"\t", " ",
	)
	text := r.Replace(raw)

	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	if len(text) > maxBodyLen {
		text = text[:maxBodyLen]
	}
	return strings.TrimSpace(text)
}

func main() {
	inputFile := `C:\Coding\Core\Go Lang\distributed-seach-engine\data\raw\enwiki-20260701-pages-articles-multistream1.xml-p1p41242`
	maxDocs := 100000
	outputFile := `C:\Coding\Core\Go Lang\distributed-seach-engine\data\processed\output.jsonl`

	openFile, err := os.Open(inputFile)
	out, e := os.Create(outputFile)

	if e != nil {
		log.Printf("err while opening the output file")
	}
	defer out.Close()

	if err != nil {
		log.Printf("err while opening the input file")
	}

	defer openFile.Close()

	decoder := xml.NewDecoder(openFile)

	var pagesRead int
	log.Printf("Starting XML Parser...")
	start := time.Now()
	encoder := json.NewEncoder(out)
	for {
		tok, err := decoder.Token()

		if err != nil {
			break
		}
		startElement, ok := tok.(xml.StartElement)

		if !ok || startElement.Name.Local != "page" {
			continue
		}

		var p page

		if err := decoder.DecodeElement(&p, &startElement); err != nil {
			log.Printf("Decode Error, skipping it.. - %v", err)
			continue
		}
		pagesRead++

		if pagesRead%10000 == 0 {
			log.Printf("Read %d pages from XML...", pagesRead)
		}

		if maxDocs > 0 && pagesRead >= maxDocs {
			log.Printf("Hit max Docs capapcity, stopping reader at page (%d)", pagesRead)
			break
		}

		cleanedBody := cleanWikiText(p.Revision.Text)
		doc := document{
			ID:    fmt.Sprintf("%d", pagesRead),
			Title: p.Title,
			Body:  cleanedBody,
		}

		if err := encoder.Encode(doc); err != nil {
			log.Printf("Encode Error, skipping it.. - %v", err)
			continue
		}

	}
	elapsed := time.Since(start)
	log.Printf("Finished reading %d pages in %v", pagesRead, elapsed)
}
