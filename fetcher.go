package nodeify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Fetcher interface {
	Fetch(since time.Time) ([]Module, error)
}

type HTTPFetcher struct {
	URL *url.URL
}

func (f HTTPFetcher) Fetch(since time.Time) ([]Module, error) {
	sb, err := since.MarshalJSON()
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("include_docs", "true")
	q.Add("startkey", string(sb))
	f.URL.RawQuery = q.Encode()
	resp, err := http.Get(f.URL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("HTTPFetcher: bad status code: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return unmarshal(data)
}

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
