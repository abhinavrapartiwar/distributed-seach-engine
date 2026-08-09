package main

import (
	"distributed-search-engine/internal/config"
	"distributed-search-engine/internal/document"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// const (
// 	defaultInputFile  = `C:\Coding\Core\Go Lang\distributed-seach-engine\data\raw\enwiki-20260701-pages-articles-multistream1.xml-p1p41242`
// 	defaultOutputFile = `C:\Coding\Core\Go Lang\distributed-seach-engine\data\processed\output.jsonl`
// 	maxBodyLen        = 2000
// )

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

	if len(text) > config.DefaultMaxBodyLength {
		text = text[:config.DefaultMaxBodyLength]
	}
	return strings.TrimSpace(text)
}

func main() {
	var (
		inputFile  = flag.String("input", config.DefaultInputPath, "Wikipedia XML dump path")
		outputFile = flag.String("output", config.DefaultOutputPath, "JSONL output path")
		maxDocs    = flag.Int("max-docs", 100000, "Maximum documents to write; 0 means unlimited")
	)
	flag.Parse()

	openFile, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf(
			"[parser] failed to open input file: path=%q err=%v",
			*inputFile,
			err,
		)
	}
	defer openFile.Close()

	out, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf(
			"[parser] failed to create output file: path=%q err=%v",
			*outputFile,
			err,
		)
	}
	defer out.Close()

	decoder := xml.NewDecoder(openFile)

	var (
		pagesRead        int
		documentsWritten int
		pagesSkipped     int
	)
	log.Printf("Starting XML Parser...")
	start := time.Now()
	encoder := json.NewEncoder(out)
	for {
		tok, err := decoder.Token()

		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf(
				"[parser] failed to read XML token: err=%v",
				err,
			)
		}
		startElement, ok := tok.(xml.StartElement)

		if !ok || startElement.Name.Local != "page" {
			continue
		}

		var p document.Page

		if err := decoder.DecodeElement(&p, &startElement); err != nil {
			log.Printf(
				"[parser] skipped page: failed to decode page: err=%v",
				err,
			)
			pagesSkipped++
			continue
		}
		pagesRead++

		if pagesRead%10000 == 0 {
			log.Printf("Read %d pages from XML...", pagesRead)
		}

		cleanedBody := cleanWikiText(p.Revision.Text)
		if len([]rune(cleanedBody)) < config.DefaultMinBodyLength {
			log.Printf(
				"[parser] skipped page: reason=body_below_minimum title=%q body_length=%d minimum=%d",
				p.Title,
				len([]rune(cleanedBody)),
				config.DefaultMinBodyLength,
			)
			pagesSkipped++
			continue
		}

		doc := document.Document{
			ID:    fmt.Sprintf("%d", pagesRead),
			Title: p.Title,
			Body:  cleanedBody,
		}

		if err := encoder.Encode(doc); err != nil {
			log.Printf(
				"[parser] failed to encode document: doc_id=%q err=%v",
				doc.ID,
				err,
			)
			continue
		}
		documentsWritten++

		if *maxDocs > 0 && documentsWritten >= *maxDocs {
			log.Printf(
				"[parser] document limit reached: documents_written=%d limit=%d",
				documentsWritten,
				*maxDocs,
			)
			break
		}
	}
	elapsed := time.Since(start)
	log.Printf(
		"[parser] parsing completed: pages_read=%d documents_written=%d pages_skipped=%d elapsed=%s",
		pagesRead,
		documentsWritten,
		pagesSkipped,
		elapsed,
	)
}
