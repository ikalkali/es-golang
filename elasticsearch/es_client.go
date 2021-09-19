package elasticsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

var (
	Client esClientInterface = &esClient{}
)

type esClientInterface interface {
	setClient(*elastic.Client)
	Index(string, string, interface{}) (*elastic.IndexResponse, error)
	Get(string, string, string) (*elastic.GetResult, error)
	GetAll(string, int64, int64) (*elastic.SearchResult, error)
	Search(string, elastic.Query) (*elastic.SearchResult, error)
}

type esClient struct {
	client *elastic.Client
}

func Init() {
	log := logrus.New()
	fmt.Println("Initializing elasticsearch client...")
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log),
		elastic.SetInfoLog(log),
		elastic.SetSniff(false),
	)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	Client.setClient(client)
	logrus.Info("Elasticsearch client initialized")
}

func (c *esClient) setClient(client *elastic.Client) {
	c.client = client
}

func (c *esClient) Get(index string, docType string, id string) (*elastic.GetResult, error) {
	ctx := context.Background()
	result, err := c.client.Get().Index(index).Type(docType).Id(id).Do(ctx)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return result, err
}

func (c *esClient) Index(index string, docType string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err := c.client.Index().Index(index).Type(docType).BodyJson(doc).Do(ctx)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return result, err
}

func (c *esClient) Search(index string, query elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	fmt.Println(query)
	result, err := c.client.Search(index).
		Query(query).
		RestTotalHitsAsInt(true).
		Do(ctx)
	if err != nil {
		logrus.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) GetAll(index string, limit int64, offset int64) (*elastic.SearchResult, error) {
	ctx := context.Background()
	q := elastic.NewMatchAllQuery()
	fmt.Println("INDEX DI CLIENT", index)
	var (
		result *elastic.SearchResult
		err    error
	)
	if limit == 0 && offset == 0 {
		result, err = c.client.Search(index).Query(q).RestTotalHitsAsInt(true).Do(ctx)
	} else {
		result, err = c.client.Search(index).From(int(offset)).Size(int(limit)).Query(q).RestTotalHitsAsInt(true).Do(ctx)
	}
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return result, err
}
