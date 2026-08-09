package index

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

func BuildMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()
	indexMapping.DefaultMapping = documentMapping

	titleFieldMapping := bleve.NewTextFieldMapping()
	titleFieldMapping.Analyzer = "en"
	titleFieldMapping.Store = true
	titleFieldMapping.Index = true
	documentMapping.AddFieldMappingsAt("title", titleFieldMapping)

	bodyFieldMapping := bleve.NewTextFieldMapping()
	bodyFieldMapping.Analyzer = "en"
	bodyFieldMapping.Store = false
	bodyFieldMapping.Index = true
	documentMapping.AddFieldMappingsAt("body", bodyFieldMapping)

	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Analyzer = "keyword"
	idFieldMapping.Store = true
	idFieldMapping.Index = true
	documentMapping.AddFieldMappingsAt("id", idFieldMapping)

	return indexMapping
}
