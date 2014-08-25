package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type jsonResult struct {
	TotalRows int       `json:"total_rows"`
	Offset    int       `json:"offset"`
	Rows      []jsonRow `json:"rows"`
}

type jsonRow struct {
	ID  string    `json:"id"`
	Key time.Time `json:"key"`
	Doc jsonDoc   `json:"doc"`
}

type jsonDoc struct {
	DistTags jsonDistTags           `json:"dist-tags"`
	Versions map[string]interface{} `json:"versions"`
}

type jsonDistTags struct {
	Latest string `json:"latest"`
}

type result struct {
	Name      string
	UpdatedAt time.Time
	Versions  []string
	Latest    string
}

func (jr jsonResult) toResults() ([]result, error) {
	results := []result{}
	for _, row := range jr.Rows {
		r := result{}
		r.Name = row.ID
		r.UpdatedAt = row.Key
		versions, err := row.Doc.toVersions()
		if err != nil {
			return nil, err
		}
		r.Versions = versions
		r.Latest = row.Doc.DistTags.Latest
		results = append(results, r)
	}
	return results, nil
}

func (jd jsonDoc) toVersions() ([]string, error) {
	versions := []string{}
	for version := range jd.Versions {
		versions = append(versions, version)
	}
	return versions, nil
}

func fetch(since time.Time) ([]result, error) {
	sinceJson, err := json.Marshal(since)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("include_docs", "true")
	q.Add("startkey", string(sinceJson))
	u := url.URL{
		Scheme:   "https",
		Host:     "skimdb.npmjs.com",
		Path:     "registry/_design/app/_view/updated",
		RawQuery: q.Encode(),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var r jsonResult
	err = decoder.Decode(&r)
	if err != nil {
		return nil, err
	}
	results, err := r.toResults()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func main() {
	for {
		results, err := fetch(time.Now())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v", results)
		time.Sleep(10 * time.Second)
	}
}
