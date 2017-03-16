package betalist

type BetalistType struct {
	Path             string
	Company_name     string
	Description_html string
	Commitment       string
	City             string
	Location         string
	Country          string
	Title            string
	Source_Id        string
	Remote           bool
	Created_at_i     uint64
	Tags             []string `json:"_tags"`
}
