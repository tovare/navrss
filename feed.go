// Package feed for connecting to the jobs API
package feed

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Google Cloud Function free tier:
// gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http --env-vars-file .env.yaml
// https://us-central1-nav-stillinger.cloudfunctions.net/RSSFeed

// RSSFeed is an HTTP Cloud Function.
func RSSFeed(w http.ResponseWriter, r *http.Request) {

	municipal := r.FormValue("municipal")
	county := r.FormValue("county")
	size := r.FormValue("size")
	view := r.FormValue("view")

	// query Ads and populate jobs.
	var jobs Jobs
	{
		// Build query for job adverts.
		var bearer string
		if bearer = os.Getenv("ARBEIDAPI"); bearer == "" {
			// Below is the public token, replace with private token for continued service.
			bearer = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJwdWJsaWMudG9rZW4udjFAbmF2Lm5vIiwiYXVkIjoiZmVlZC1hcGktdjEiLCJpc3MiOiJuYXYubm8iLCJpYXQiOjE1NTc0NzM0MjJ9.jNGlLUF9HxoHo5JrQNMkweLj_91bgk97ZebLdfx3_UQ"
		}
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
		var body []byte
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
	switch view {

	case "html-headlines":
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("User-Agent", "TovAreRSS/1.0 <mail@tovare.com>")
		fmt.Fprint(w, jobs.renderHtmHeadlines())

	default:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Header().Set("User-Agent", "TovAreRSS/1.0 <mail@tovare.com>")
		fmt.Fprint(w, jobs.renderRSS())
	}
}

// Default rendering RSS
func (jobs Jobs) renderRSS() string {
	// Prepare result
	res := RSS{
		Version: "2.0",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Channel: Channel{
			Title:         "Based on data from Arbeidsplassen",
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
	return xml.Header + string(output)
}

// HTML headlines for last positions.
func (jobs Jobs) renderHtmHeadlines() string {

	t := `
<!DOCTYPE html>
<html>
 <meta charset="UTF-8">
 <head>Siste stillinger</head>
 <body>
   <table>
	 <thead>
	    <tr><th>stilling</th></tr>
	 </thead>
	 <tbody>
	   {{range .Items}}
		  <tr>
		  </tr>
	   {{end}}
	 </tbody>
   </table>  
 </body>
</html>
`
	//tem, _ := template.New("headlines").Parse(t)

	return t
}
