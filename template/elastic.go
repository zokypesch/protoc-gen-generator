package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplElastic = `package core

import (
	"context"
	"fmt"
	elastic "gopkg.in/olivere/elastic.v5"
	"log"
)

// ESCore for struct information
type ESCore struct {
	client    *elastic.Client
	indexName string
	mapping   string
	typeIndex string
}

// ESModule mdule of elastic search
type ESModule interface {
	CreateIndex(ctx context.Context) error
	AddDocument(ctx context.Context, ID string, body interface{}) error
	GetQueryTerm(key string, value interface{}) *elastic.TermQuery
	Query(ctx context.Context, termQuery *elastic.TermQuery, boolQuery *elastic.BoolQuery, offset int, size int, sortBy string, asc bool) ([]*elastic.SearchHit, int, error)
	DeleteDocumentByID(ctx context.Context, ID string) error
	UpdateDocument(ctx context.Context, ID string, script []*elastic.Script, upsert map[string]interface{}) (int, error)
	GetScriptLine(condition string, paramKey string, paramValue interface{}) *elastic.Script
	GenerateScriptLines() []*elastic.Script
	GenerateQueryTerms() []*elastic.TermQuery
	DeleteIndex(ctx context.Context) error
	GetClient() *elastic.Client
	GetByID(ctx context.Context, ID string) (*elastic.GetResult, error)
	GetBoolQuery(termQuery ...elastic.Query) *elastic.BoolQuery
	GetNewMultiMatchQuery(key string, multiValue []string, operator string) *elastic.MultiMatchQuery
	GetNewMatchQuery(key string, value interface{}) *elastic.MatchQuery
	QueryCount(boolQuery *elastic.BoolQuery) (int64, error)
	GetNewRangeQuery(key string, from interface{}, end interface{}) *elastic.RangeQuery
}

// NewEsCore func for create new ES CORE
func NewEsCore(clientAddr string, indexName string, mapping string, typeIndex string) ESModule {
	client, err := elastic.NewClient(elastic.SetURL(clientAddr))
	if err != nil {
		// Handle error
		panic(err)
	}

	return &ESCore{client, indexName, mapping, typeIndex}

}

// CreateIndex for creating a index
func (es *ESCore) CreateIndex(ctx context.Context) error {

	exists, err := es.client.IndexExists(es.indexName).Do(ctx)
	if err != nil {
		log.Println("erorr here")
		return err
	}

	if exists {
		return nil
	}

	createIndex, errCrIndex := es.client.CreateIndex(es.indexName).BodyString(es.mapping).Do(ctx)

	if errCrIndex != nil {
		log.Println("erorr when create index")
		return err
	}

	if !createIndex.Acknowledged {
		return fmt.Errorf("not acknowladge")
	}

	return nil
}

// GetQueryTerm for get query term
func (es *ESCore) GetQueryTerm(key string, value interface{}) *elastic.TermQuery {
	return elastic.NewTermQuery(key, value)
}

// AddDocument for adding document in elastic search
func (es *ESCore) AddDocument(ctx context.Context, ID string, body interface{}) error {
	_, err := es.client.Index().
		Index(es.indexName).
		Type(es.typeIndex).
		Id(ID).
		Refresh("wait_for").
		BodyJson(body).
		Do(ctx)

	_, errFlush := es.client.Flush().Index(es.indexName).Do(ctx)
	if errFlush != nil {
		log.Println("error flush")
	}

	return err
}

// GetBoolQuery for get boolquery
func (es *ESCore) GetBoolQuery(termQuery ...elastic.Query) *elastic.BoolQuery {
	boolQueryNew := elastic.NewBoolQuery()
	boolQueryNew.Must(termQuery...)
	// for _, v := range termQuery {
	// 	boolQueryNew.Must(v)
	// }

	return boolQueryNew
}

// GetNewRangeQuery for get new match query
func (es *ESCore) GetNewRangeQuery(key string, from interface{}, end interface{}) *elastic.RangeQuery {
	rangeQuery := elastic.NewRangeQuery(
		key,
	).From(from).To(end)

	return rangeQuery
}

// GetNewMultiMatchQuery for get new match query
func (es *ESCore) GetNewMultiMatchQuery(key string, multiValue []string, operator string) *elastic.MultiMatchQuery {
	multiQuery := elastic.NewMultiMatchQuery(
		key,
		multiValue...,
	).Type("phrase_prefix").Operator(operator)

	return multiQuery
}

// GetNewMatchQuery for get new match query
func (es *ESCore) GetNewMatchQuery(key string, value interface{}) *elastic.MatchQuery {
	matchQuery := elastic.NewMatchQuery(key, value)

	return matchQuery
}

// Query for query data
func (es *ESCore) Query(ctx context.Context, termQuery *elastic.TermQuery,
	boolQuery *elastic.BoolQuery, offset int, size int, sortBy string, asc bool) ([]*elastic.SearchHit, int, error) {

	query := es.client.Search().
		Index(es.indexName).
		Type(es.typeIndex).      // search in index
		Sort(sortBy, asc).       // sort by param field, ascending
		From(offset).Size(size). // take documents by param
		Pretty(true)             // pretty print request and response JSON

	if termQuery != nil {
		query.Query(termQuery)
	}

	if boolQuery != nil {
		query.Query(boolQuery)
	}

	searchResult, err := query.Do(ctx)

	if err != nil {
		return nil, 0, err
	}
	return searchResult.Hits.Hits, int(searchResult.TotalHits()), nil
}

// QueryCount for query count
func (es *ESCore) QueryCount(boolQuery *elastic.BoolQuery) (int64, error) {

	count, err := es.client.Count(es.indexName).
		Type(es.typeIndex).
		Query(boolQuery).
		// Pretty(true).
		Do(context.TODO())

	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetScriptLine for get scriptline
func (es *ESCore) GetScriptLine(condition string, paramKey string, paramValue interface{}) *elastic.Script {
	return elastic.NewScriptInline("ctx._source."+condition).Lang("painless").Param(paramKey, paramValue)
}

// UpdateDocument for updating a document
func (es *ESCore) UpdateDocument(ctx context.Context, ID string, script []*elastic.Script, upsert map[string]interface{}) (int, error) {
	query := es.client.Update().Index(es.indexName).Refresh("wait_for").Type(es.typeIndex).Id(ID)

	if script != nil {
		for _, v := range script {
			query.Script(v)
		}
	}
	update, err := query.Upsert(upsert).Do(ctx)
	if err != nil {
		return 0, err
	}

	_, errFlush := es.client.Flush().Index(es.indexName).Do(ctx)
	if errFlush != nil {
		log.Println("error flush")
	}

	return update.Version, nil

}

// DeleteDocumentByID for delete th elastic search
func (es *ESCore) DeleteDocumentByID(ctx context.Context, ID string) error {
	// Delete an index.
	_, err := es.client.Delete().
		Index(es.indexName).
		Type(es.typeIndex).
		Id(ID).
		Refresh("wait_for").
		Do(ctx)

	if err != nil {
		// Handle error
		return err
	}

	return nil
}

// GenerateScriptLines for generate scriptlines
func (es *ESCore) GenerateScriptLines() []*elastic.Script {
	var res []*elastic.Script

	return res
}

// GenerateQueryTerms for generate scriptlines
func (es *ESCore) GenerateQueryTerms() []*elastic.TermQuery {
	var res []*elastic.TermQuery

	return res
}

// DeleteIndex for deleting index
func (es *ESCore) DeleteIndex(ctx context.Context) error {
	_, err := es.client.DeleteIndex(es.indexName).Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

// GetByID for get ID document
func (es *ESCore) GetByID(ctx context.Context, ID string) (*elastic.GetResult, error) {
	get, err := es.client.Get().
		Index(es.indexName).
		Type(es.typeIndex).
		Id(ID).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if !get.Found {
		return nil, fmt.Errorf("data not found")
	}
	return get, nil
}

// GetClient for get client es
func (es *ESCore) GetClient() *elastic.Client {
	return es.client
}
`

var ListElastic = lib.List{
	FileType:     ".elastic.go",
	Template:     tmplElastic,
	Location:     "./core/",
	Lang:         "elastic",
	ReplaceQuote: true,
}
