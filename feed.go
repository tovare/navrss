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
		"uuid": "ddf8fe54-853c-4a4a-9969-07342e71e4ba",
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
		"description": "<p>Org. nr: 976677218 Stillingsident: 4144720310 Presentasjon av stillingen:</p>\n<p>Ved Dalgård skole og ressurssenter er det ledig en fast 100 prosent stilling som lærer med tegnspråkkompetanse.</p>\n<p>Dalgård skole og ressurssenter består av et stort mangfold med 640 elever og 150 ansatte. De aller fleste elevene tilhører nærskolen på 1.-7. trinn. Ressurssenteret gir opplæring til elever med store og sammensatte behov på 1.-10. trinn fra hele Trondheim kommune. SFO er en naturlig del av vår virksomhet. Som medarbeider på Dalgård skole og ressurssenter er du en del av et rikt kompetansemiljø.</p>\n<p><a href=\"https://www.trondheim.kommune.no/org/oppvekst/skoler/dalgard-skole/\" rel=\"nofollow\">Les mer om oss Dalgård skole og ressurssenter på våre hjemmesider</a>.</p>\n<h3>Arbeidsoppgaver</h3>\n<ul><li>Lærer for alle elever, også med tegnspråk som førstespråk</li><li>Deltakelse i et hjem-skole samarbeid og samarbeid med ulike instanser etter behov</li><li>Deltakelse i teamarbeid med både pedagoger og miljøpersonalet</li><li>Deltakelse i skolens pedagogiske utviklingsarbeid</li><li>Ta i et tak og finne gode løsninger i skolehverdagen</li></ul>\n<h3>Kvalifikasjoner</h3>\n<ul><li>Grunnskolelærerutdanning 1.-7.</li><li>Formell kompetanse i tegnspråk</li><li>God kompetanse i relasjonsbygging ift alle typer elever er nødvendig</li><li>Ønskelig med erfaring fra undervisning der tegnspråk brukes i klasserommet</li><li>Ønskelig med erfaring fra undervisning og kommunikasjon på tegnspråk</li><li>Skolen har svømmehall. Det er nødvendig at den som tilsettes kan være med elever i basseng.</li></ul>\n<h3>Utdanningsretning</h3>\n<ul><li>Pedagogikk</li></ul>\n<h3>Utdanningsnivå</h3>\n<ul><li>Høyskole / Universitet</li></ul>\n<h3>Personlige egenskaper</h3>\n<ul><li>Skolen søker etter lærere som synliggjør sitt engasjement for skoleutvikling</li><li>Vi søker lærere som ønsker å være i utvikling både faglig og sosialt, og som ser tydelig klasseledelse og arbeid med å skape gode relasjoner i gruppa som en sentral forutsetning for elevers læring og utvikling</li><li>Samarbeidsevne, fleksibilitet og godt humør er personlige egenskaper som verdsettes høyt</li><li>Skolen har elever med ulike hjelpebehov. De som tilsettes må kunne yte praktisk hjelp i hverdagen, og må også være forberedt på å bistå elever med utagerende atferd. Søkere med relevant tilleggskompetanse/ personlige egenskaper/ interesser i forhold til Dalgård skoles elevmangfold ønskes</li><li>Ettersom noen av våre elever har behov for bistand i garderobesituasjoner og med daglig hygiene, ønsker vi ansatte av begge kjønn</li></ul>\n<h3>Language</h3>\n<ul><li>Norsk</li></ul>\n<h3>Vi tilbyr</h3>\n<ul><li>Et godt arbeidsmiljø i en levende kunnskapsvirksomhet, der det å jobbe tverrfaglig og i fellesskap er sentralt</li><li>En arbeidsplass der et utvidet kulturbegrep og det gode humør betyr mye for vår skoles identitet</li><li>Lønn etter sentrale avtaler og Trondheim kommunes lokale avtaler</li><li>Gode pensjons- og forsikringsordninger</li></ul>\n<h3>Kontaktinformasjon</h3>\n<p>Elin Gjønnes, Avdelingsleder, 90753795<br />\nGro Kristiansen, Rektor, 72545002 / 91642234</p>\n<h3>Arbeidssted</h3>\n<p>Anders Wigens veg 3<br />\n7024 Trondheim</p>\n<p>Søk på stillingen:<a href=\"https://candidate.webcruiter.com/cv?advertid&#61;4144720310&amp;language&#61;nb&amp;link_source_id&#61;17\" rel=\"nofollow\">Klikk her</a></p>\n<h3>Nøkkelinformasjon:</h3>\n<p>Arbeidsgiver:Trondheim kommune</p>\n<p>Referansenr.:4144720310<br />\nStillingsprosent: 100%<br />\nStartdato: 01.12.2019<br />\nSøknadsfrist: 10.11.2019</p>\n",
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
		"link": "https://arbeidsplassen.nav.no/stillinger/stilling/ddf8fe54-853c-4a4a-9969-07342e71e4ba",
		"employer": {
		  "name": "Dalgård skole og ressurssenter, Trondheim kommune",
		  "orgnr": "976677218",
		  "description": "<p><strong>Åpen -- kompetent -- modig</strong></p>\n<p>Trondheim kommune er landets tredje største kommune, og en av Norges største arbeidsplasser nord for Dovre. Hver dag arbeider om lag 13 000 medarbeidere med å gi byens innbyggere gode helse- og velferdstjenester, utviklende oppvekst- og utdanningsmuligheter, et godt bymiljø og en rekke kultur- og fritidstilbud. Trondheim kommune er en framtidsrettet og moderne arbeidsgiver som samarbeider nært med byens forsknings-, universitets- og høgskolemiljø.</p>\n<p>Les mer om <a href=\"https://www.trondheim.kommune.no/\" rel=\"nofollow\">Trondheim kommune</a></p>\n<p>Trondheim kommune ønsker mangfold og oppfordrer kvalifiserte kandidater til å søke uansett alder, kjønn, etnisitet og funksjonsnedsettelser. Etter søknadsfristens utløp blir det utarbeidet offentlig søkerliste.<br />\nGyldig politiattest må framlegges i tilsettinger der det er påkrevd.</p>\n",
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
