package models

type BggRank struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

type AtlasResult struct {
	Id            string             `json:"id"`
	ThumbUrl      string             `json:"thumb_url"`
	ImageUrl      string             `json:"image_url"`
	BggId         int                `json:"bgg_id"`
	YearPublished int                `json:"year_published"`
	Ranks         map[string]BggRank `json:"bgg_rank"`
	Bayesaverage  float64            `json:"bgg_bayes"`
	Average       float64            `json:"bgg_average"`
}
