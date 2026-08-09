package index

import (
	"distributed-search-engine/internal/document"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
)

func IndexDocument(index bleve.Index, inputPath string) error {
	jsonlFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf(
			"open JSONL input %q: %w",
			inputPath,
			err,
		)
	}
	defer jsonlFile.Close()
	decoder := json.NewDecoder(jsonlFile)
	start := time.Now()
	docsIndexed := 0

	for {
		var doc document.Document

		if err := decoder.Decode(&doc); err != nil {
			log.Printf("[indexer] skipped document while decoding: reason= - %v", err)
			if err == io.EOF {
				break
			}
			continue
		}
		if err := index.Index(doc.ID, doc); err != nil {
			return fmt.Errorf(
				"index document %q: %w",
				doc.ID,
				err,
			)
		}
		docsIndexed++
		if docsIndexed%5000 == 0 {
			log.Printf(
				"[indexer] indexing progress: documents_indexed=%d",
				docsIndexed,
			)
		}
	}
	elapsed := time.Since(start)
	log.Printf(
		"[indexer] indexing completed: documents_indexed=%d elapsed=%s",
		docsIndexed,
		elapsed,
	)
	return nil
}
