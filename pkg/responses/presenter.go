package responses

import "github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"

type Response[D any] struct {
	Data      D                       `json:"data"`
	Error     *httpErrors.ErrResponse `json:"error,omitempty"`
	IsSuccess bool                    `json:"is_success"`
}

// Just for swag
type SuccessResponse[D any] struct {
	Data      D    `json:"data"`
	IsSuccess bool `json:"is_success" example:"true"`
}

type ErrorResponse struct {
	Error     *httpErrors.ErrResponse `json:"error"`
	IsSuccess bool                    `json:"is_success" example:"false"`
}
