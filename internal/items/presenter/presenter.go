package presenter

type ItemCreate struct {
	Title       string `json:"title" validate:"required" example:"item title"`
	Description string `json:"description" validate:"required" example:"item description"`
}

type ItemResponse struct {
	Id          uint   `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	OwnerId     uint   `json:"owner_id,omitempty"`
}

type ItemUpdate struct {
	Title       string `json:"title" example:"item title"`
	Description string `json:"description" example:"item description"`
}
