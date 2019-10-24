// Package feed ... Process for harvesting jobs into firestore for deeper analysis.
package feed

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// UpdateAndStore newe data from positions.
// gcloud functions deploy UpdateAndStore --memory=128 --runtime go111 --trigger-topic oppdater
func UpdateAndStore(ctx context.Context, m PubSubMessage) error {

	// Initialize the firebase client.
	//ctx := context.Background()
	client, err := firestore.NewClient(ctx, "nav-stillinger")
	if err != nil {
		log.Fatal(err)
	}

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
		q.Add("size", "50") // Smarter query.
		q.Add("page", "0")
		req.URL.RawQuery = q.Encode()

		// Query service.
		client := http.DefaultClient
		resp, err := client.Do(req)
		// "success is not final, failure is not fatal: it is the
		//   courage to continue that counts"
		var body []byte
		if err != nil {
			log.Print(err)
		} else {
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Print(err)
			}
		}
		err = json.Unmarshal(body, &jobs)
		if err != nil {
			log.Print(err)
			jobs = Jobs{}
		}
	}

	// TODO: Expiration.

	// Add positions to database.
	for _, j := range jobs.Content {
		_, err := client.Collection("test").Doc(j.UUID).Set(ctx, j)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}
	}
	return nil

}
