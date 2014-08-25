package nodeify

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type fetchResult struct {
	Rows []fetchResultRow `json:"rows"`
}

type fetchResultRow struct {
	ID  string            `json:"id"`
	Key time.Time         `json:"key"`
	Doc fetchResultRowDoc `json:"doc"`
}

type fetchResultRowDoc struct {
	DistTags fetchResultRowDocTags `json:"dist-tags"`
}

type fetchResultRowDocTags struct {
	Latest string `json:"latest"`
}

type Fetcher struct {
	u url.URL
}

func NewFetcher(rawurl string) (*Fetcher, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Fetcher{*u}, nil
}

func (f Fetcher) Fetch(since time.Time) ([]Module, error) {
	sinceBytes, err := since.MarshalJSON()
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("include_docs", "true")
	q.Add("startkey", string(sinceBytes))
	f.u.RawQuery = q.Encode()
	log.Print(f.u.String())
	resp, err := http.Get(f.u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return unmarshal(data)
}

func unmarshal(data []byte) ([]Module, error) {
	var modules []Module
	var fr fetchResult
	err := json.Unmarshal(data, &fr)
	if err != nil {
		return nil, err
	}
	for _, row := range fr.Rows {
		var module Module
		module.Name = row.ID
		module.UpdatedAt = row.Key
		module.LatestVersion = row.Doc.DistTags.Latest
		modules = append(modules, module)
	}
	return modules, nil
}
