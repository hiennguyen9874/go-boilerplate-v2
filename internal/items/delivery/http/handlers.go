package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/items"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/items/presenter"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/middleware"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/responses"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/utils"
)

type itemHandler struct {
	cfg     *config.Config
	itemsUC items.ItemUseCaseI
	logger  logger.Logger
}

func CreateItemHandler(uc items.ItemUseCaseI, cfg *config.Config, logger logger.Logger) items.Handlers {
	return &itemHandler{cfg: cfg, itemsUC: uc, logger: logger}
}

// Create godoc
// @Summary Create Item
// @Description Create new item.
// @Tags items
// @Accept json
// @Produce json
// @Param item body presenter.ItemCreate true "Add item"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /item [post]
func (h *itemHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		item := new(presenter.ItemCreate)

		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(ctx, item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		newItem, err := h.itemsUC.CreateWithOwner(
			ctx,
			user.Id,
			mapModel(item),
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		itemResponse := *mapModelResponse(newItem)
		render.Respond(w, r, responses.CreateSuccessResponse(itemResponse))
	}
}

// Get godoc
// @Summary Read item
// @Description Get item by ID.
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item Id"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /item/{id} [get]
func (h *itemHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		item, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && item.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(item)))
	}
}

// GetMulti godoc
// @Summary Read Items
// @Description Retrieve items.
// @Tags items
// @Accept json
// @Produce json
// @Param limit query int false "limit" Format(limit)
// @Param offset query int false "offset" Format(offset)
// @Success 200 {object} responses.SuccessResponse[[]presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /item [get]
func (h *itemHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		var items []*models.Item
		if user.IsSuperUser {
			items, err = h.itemsUC.GetMulti(ctx, limit, offset)
		} else {
			items, err = h.itemsUC.GetMultiByOwnerId(ctx, user.Id, limit, offset)
		}
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(items)))
	}
}

// Delete godoc
// @Summary Delete item
// @Description Delete an item by ID.
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item Id"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /item/{id} [delete]
func (h *itemHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		item, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && item.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		err = h.itemsUC.DeleteWithoutGet(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(item)))
	}
}

// Update godoc
// @Summary Update item
// @Description Update an item by ID.
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item Id"
// @Param item body presenter.ItemUpdate true "Update item"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /item/{id} [put]
func (h *itemHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		item := new(presenter.ItemUpdate)

		err = json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		dbItem, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && dbItem.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		values := make(map[string]interface{})
		if item.Title != "" {
			values["title"] = item.Title
		}
		if item.Description != "" {
			values["description"] = item.Description
		}

		updatedItem, err := h.itemsUC.Update(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedItem)))
	}
}

func mapModel(exp *presenter.ItemCreate) *models.Item {
	return &models.Item{
		Title:       exp.Title,
		Description: exp.Description,
	}
}

func mapModelResponse(exp *models.Item) *presenter.ItemResponse {
	return &presenter.ItemResponse{
		Id:          exp.Id,
		Title:       exp.Title,
		Description: exp.Description,
		OwnerId:     exp.OwnerId,
	}
}

func mapModelsResponse(exp []*models.Item) []*presenter.ItemResponse {
	out := make([]*presenter.ItemResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
