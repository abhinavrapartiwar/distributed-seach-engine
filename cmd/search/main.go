package main

import (
	"distributed-search-engine/internal/config"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/blevesearch/bleve/v2"
)

func search(index bleve.Index, text string) {
	query := bleve.NewMatchQuery(text)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"title"}
	start := time.Now()
	searchResult, err := index.Search(search)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf("Error while getting search result %v", err)
		return
	}
	log.Printf(
		"[search] query completed: query=%q hits=%d elapsed=%s",
		text,
		len(searchResult.Hits),
		elapsed,
	)
	for _, hit := range searchResult.Hits {
		fmt.Printf("ID is : %s\n", hit.ID)
		fmt.Printf("Score is : %f\n", hit.Score)
		fmt.Printf("Title is : %s\n", hit.Fields["title"])
	}
}

func main() {
	indexPath := flag.String(
		"index",
		config.DefaultIndexPath,
		"path to the Bleve index",
	)

	queryText := flag.String(
		"query",
		"Country",
		"search query",
	)
	flag.Parse()
	log.Printf(
		"[search] opening index: path=%q",
		*indexPath,
	)
	index, err := bleve.Open(*indexPath)

	if err != nil {
		log.Fatalf(
			"[search] failed to open index: path=%q err=%v",
			*indexPath,
			err,
		)
	}
	defer index.Close()
	log.Printf("Opened new index at %s", *indexPath)
	search(index, *queryText)
}
