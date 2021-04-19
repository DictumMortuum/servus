package models

type AtlasResult struct {
	Id            string `json:"id"`
	ThumbUrl      string `json:"thumb_url"`
	ImageUrl      string `json:"image_url"`
	BggId         int    `json:"bgg_id"`
	YearPublished int    `json:"year_published"`
	Rank          string `json:"bgg_rank"`
}
