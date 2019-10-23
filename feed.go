// Package feed for connecting to the jobs API
package feed

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Google Cloud Function free tier:
// gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http
// https://us-central1-nav-stillinger.cloudfunctions.net/RSSFeed

//***********************************************
// MODELS FOR RSS, Jobs
//***********************************************

// RSS is for serializing RSS v2.0 feeds. Tool: https://www.onlinetool.io/xmltogo/
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel is part of the rss document.
type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

// Item is part of the rss document.
type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
	//Author      string `xml:"author"`
	Description string `xml:"description"`
	Content     Cdata  `xml:"content:encoded"`
	//Encoded     string `xml:"encoded"`
}

// Cdata creates a CDATA
type Cdata struct {
	Value string `xml:",cdata"`
}

// Jobs represents a list of job adverts.
type Jobs struct {
	Content []struct {
		UUID          string    `json:"uuid"`
		Published     time.Time `json:"published"`
		Expires       time.Time `json:"expires"`
		WorkLocations []struct {
			Country    string `json:"country"`
			City       string `json:"city"`
			PostalCode string `json:"postalCode"`
			County     string `json:"county"`
			Municipal  string `json:"municipal"`
		} `json:"workLocations"`
		Title          string   `json:"title"`
		Description    string   `json:"description"` // html5
		Source         string   `json:"source"`
		ApplicationDue string   `json:"applicationDue"`
		Occupations    []string `json:"occupations"`
		Link           string   `json:"link"`
		Employer       struct {
			Name        string `json:"name"`
			Orgnr       string `json:"orgnr"`
			Description string `json:"description"`
		} `json:"employer"`
	} `json:"content"`
	TotalElements int    `json:"totalElements"`
	PageNumber    int    `json:"pageNumber"`
	PageSize      int    `json:"pageSize"`
	TotalPages    int    `json:"totalPages"`
	First         bool   `json:"first"`
	Last          bool   `json:"last"`
	Sort          string `json:"sort"`
}

//***********************************************
//  DATA FOR SILENT FAILURE.
//***********************************************

var errorData = []byte(`
{
	"content": [
	  {
		"uuid": "",
		"published": "2019-10-22T22:00:00Z",
		"expires": "2019-11-09T23:00:00Z",
		"updated": "2019-10-22T22:00:00.088822Z",
		"workLocations": [
		  {
			"country": "NORGE",
			"address": "Anders Wigens veg 3",
			"city": "TRONDHEIM",
			"postalCode": "7024",
			"county": "TRØNDELAG",
			"municipal": "TRONDHEIM"
		  }
		],
		"title": "Lærer med tegnspråk",
		"description": "\n",
		"sourceurl": null,
		"source": "XML_STILLING",
		"applicationDue": "2019-11-10T00:00",
		"occupationCategories": [
		  {
			"level1": "Utdanning",
			"level2": "Grunnskole"
		  },
		  {
			"level1": "Utdanning",
			"level2": "Instruktører  og pedagoger"
		  }
		],
		"jobtitle": null,
		"link": "https://arbeidsplassen.nav.no/stillinger/",
		"employer": {
		  "name": "Trondheim kommune",
		  "orgnr": "",
		  "description": "",
		  "homepage": null
		},
		"engagementtype": "Fast",
		"extent": "Heltid",
		"starttime": "2019-12-01T00:00",
		"positioncount": "1",
		"sector": "Offentlig"
	  }
	],
	"totalElements": 5000,
	"pageNumber": 0,
	"pageSize": 1,
	"totalPages": 5000,
	"first": true,
	"last": false,
	"sort": "published:desc"
  }
`)

// RSSFeed is an HTTP Cloud Function.
func RSSFeed(w http.ResponseWriter, r *http.Request) {

	municipal := r.FormValue("municipal")
	county := r.FormValue("county")
	size := r.FormValue("size")

	// query Ads and populate jobs.
	var jobs Jobs
	{
		// Build query for job adverts.
		// Below is the public token, replace with private token for continued service.
		bearer := "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJwdWJsaWMudG9rZW4udjFAbmF2Lm5vIiwiYXVkIjoiZmVlZC1hcGktdjEiLCJpc3MiOiJuYXYubm8iLCJpYXQiOjE1NTc0NzM0MjJ9.jNGlLUF9HxoHo5JrQNMkweLj_91bgk97ZebLdfx3_UQ"
		url := "https://arbeidsplassen.nav.no/public-feed/api/v1/ads"
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", bearer)

		q := req.URL.Query()
		if municipal != "" {
			q.Add("municipal", municipal)
		}
		if municipal != "" {
			q.Add("county", county)
		}
		if size != "" {
			q.Add("size", size)
		} else {
			q.Add("size", "50")
		}
		q.Add("page", "0")
		req.URL.RawQuery = q.Encode()

		// Query service.
		client := http.DefaultClient
		resp, err := client.Do(req)
		// "success is not final, failure is not fatal: it is the
		//   courage to continue that counts"
		body := errorData
		if err != nil {
			// TODO: Fallback
			log.Print(err)
		} else {
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				// TODO: Fallback
				log.Print(err)
			}
		}
		err = json.Unmarshal(body, &jobs)
		if err != nil {
			log.Print(err)
			jobs = Jobs{}
		}
	}

	// Prepare result
	res := RSS{
		Version: "2.0",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Channel: Channel{
			Title:         "Arbeidsplassen",
			Link:          "https://arbeidsplassen.no",
			Language:      "no",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         make([]Item, 0),
		},
	}
	for _, v := range jobs.Content {
		i := Item{
			Title:   v.Title,
			Link:    v.Link,
			GUID:    v.Link,
			PubDate: v.Published.Format(time.RFC1123Z),
			//Author:  nil,
			Content: Cdata{v.Description},
		}
		res.Channel.Items = append(res.Channel.Items, i)
	}

	output, err := xml.MarshalIndent(res, "  ", "    ")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("User-Agent", "TovAreRSS/1.0 <mail@tovare.com>")
	fmt.Fprint(w, xml.Header+string(output))
}
