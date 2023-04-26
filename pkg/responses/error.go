package responses

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"
)

// render.Renderer Render() interface method
func (e *Response[D]) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Error.Status)
	return nil
}

func CreateErrorResponse(err error) render.Renderer {
	parsedErr := httpErrors.ParseErrors(err)

	return &Response[*string]{
		Data: nil,
		Error: &httpErrors.ErrResponse{
			Err:        parsedErr.GetErr(),
			Status:     parsedErr.GetStatus(),
			StatusText: parsedErr.GetStatusText(),
			Msg:        parsedErr.GetMsg(),
		},
		IsSuccess: false,
	}
}
