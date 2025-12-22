package links

type CreateLink struct {
	URL         string `json:"url" validate:"required,url"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}
