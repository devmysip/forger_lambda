package models

type UserActivity struct {
	Date     string `json:"date"`
	Day      string `json:"day,omitempty"`
	Activity []struct {
		ChapterNo string `json:"chapter_no"`
		VerseNo   string `json:"verse_no"`
	} `json:"activity"`
}
