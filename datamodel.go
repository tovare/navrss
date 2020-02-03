// Package feed Shared data structures.
package feed

import (
	"encoding/xml"
	"time"
)

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

// Job represents a job advert.
type Job struct {
	UUID          string    `json:"uuid" firestore:"name,omitempty"`
	Published     time.Time `json:"published" firestore:"published,omitempty"`
	Expires       time.Time `json:"expires" firestore:"expires,omitempty"`
	WorkLocations []struct {
		Country    string `json:"country" firestore:"country,omitempty"`
		City       string `json:"city" firestore:"city,omitempty"`
		PostalCode string `json:"postalCode" firestore:"postalcode,omitempty"`
		County     string `json:"county" firestore:"country,omitempty"`
		Municipal  string `json:"municipal" firestore:"municipal,omitempty"`
	} `json:"workLocations"`
	Title                string   `json:"title" firestore:"title,omitempty"`
	Description          string   `json:"description" firestore:"description,omitempty"` // html5
	Source               string   `json:"source" firestore:"source,omitempty"`
	ApplicationDue       string   `json:"applicationDue" firestore:"applicationdue,omitempty"`
	Occupations          []string `json:"occupations" firestore:"occupations,omitempty"`
	OccupationCategories []struct {
		Level1 string `json:"level1" firestore:"level1"`
		Level2 string `json:"level2" firestore:"level2"`
	} `json:"occupationCategories" firestore:"occupationcategories,omitempty"`
	Link           string `json:"link" firestore:"link,omitempty"`
	Engagementtype string `json:"engagementtype" firestore:"engagementtype,omitempty"`
	Extent         string `json:"extent" firestore:"extent,omitempty"`
	Starttime      string `json:"starttime" firestore:"starttime,omitempty"`
	Positioncount  string `json:"positioncount" firestore:"positioncount,omitempty"`
	Sector         string `json:"sector" firestore:"sector,omitempty"`
	Employer       struct {
		Name        string `json:"name" firestore:"name,omitempty"`
		Orgnr       string `json:"orgnr" firestore:"orgnr,omitempty"`
		Description string `json:"description" firestore:"description,omitempty"`
	} `json:"employer" firestore:"employer,omitempty"`
}

// Jobs is a series of jobs.
type Jobs struct {
	Content       []Job  `json:"content" firestore:"content,omitempty"`
	TotalElements int    `json:"totalElements" totalElements:"name,omitempty"`
	PageNumber    int    `json:"pageNumber" pageNumbere:"name,omitempty"`
	PageSize      int    `json:"pageSize" pageSize:"name,omitempty"`
	TotalPages    int    `json:"totalPages" totalPages:"name,omitempty"`
	First         bool   `json:"first" firestore:"first,omitempty"`
	Last          bool   `json:"last" firestore:"last,omitempty"`
	Sort          string `json:"sort" firestore:"sort,omitempty"`
}
