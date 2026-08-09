package main

import (
	"distributed-search-engine/internal/config"
	indexBuilder "distributed-search-engine/internal/index"
	"flag"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
)

func main() {
	inputPath := flag.String(
		"input",
		config.DefaultOutputPath,
		"path to the processed JSONL file",
	)

	indexPath := flag.String(
		"index",
		config.DefaultIndexPath,
		"path to the Bleve index",
	)
	resetIndex := flag.Bool(
		"reset",
		false,
		"remove the existing index before indexing",
	)

	flag.Parse()
	if *resetIndex {
		if err := os.RemoveAll(*indexPath); err != nil {
			log.Fatalf(
				"[indexer] failed to remove existing index: path=%q err=%v",
				*indexPath,
				err,
			)
		}
	}

	indexMapping := indexBuilder.BuildMapping()
	index, err := bleve.New(*indexPath, indexMapping)

	if err != nil {
		log.Fatalf(
			"[indexer] failed to create index: path=%q err=%v",
			*indexPath,
			err,
		)
	}
	defer index.Close()
	log.Printf(
		"[indexer] index created: path=%q",
		*indexPath,
	)
	if err := indexBuilder.IndexDocument(index, *inputPath); err != nil {
		log.Fatalf(
			"[indexer] failed to index document: err=%v",
			err,
		)
	}

}
