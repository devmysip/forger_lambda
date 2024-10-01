package models

type Read struct {
	Chapter  int   `json:"chapter"`
	Verses   []int `json:"verses"`
	Progress int   `json:"progress"`
}

type User struct {
	Email          string     `json:"email"`
	DisplayName    *string    `json:"display_name,omitempty"`
	ProfileURL     *string    `json:"profile_url,omitempty"`
	FCMToken       *string    `json:"fcm_token,omitempty"`
	ClientEndpoint *string    `json:"client_endpoint,omitempty"`
	LastRead       *string    `json:"last_read"`
	Reads          []Read     `json:"reads"`
	UpdatedAt      string     `json:"updated_at,omitempty"`
	CreatedAt      string     `json:"created_at,omitempty"`
	Update         *AppUpdate `json:"app_update,omitempty"`
}

type UpdateRead struct {
	ChapterNo int `json:"chapter_no"`
	VerseNo   int `json:"verse_no"`
}

type AppUpdate struct {
	BuildNo     int    `json:"build_no"`
	ForceUpdate int    `json:"force_update"`
	SoftUpdate  int    `json:"soft_update"`
	Title       string `json:"title"`
	Message     string `json:"message"`
}

var GitaChapters = []int{
	47,
	72,
	43,
	42,
	29,
	47,
	30,
	28,
	34,
	42,
	55,
	20,
	35,
	27,
	20,
	24,
	28,
	78,
}
