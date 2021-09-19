package books

type Books struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publisher   string `json:"publisher"`
	ReleaseDate string `json:"release_date"`
}
