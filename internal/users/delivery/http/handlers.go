package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/middleware"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users/presenter"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/responses"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/utils"
)

type userHandler struct {
	cfg     *config.Config
	usersUC users.UserUseCaseI
	logger  logger.Logger
}

func CreateUserHandler(uc users.UserUseCaseI, cfg *config.Config, logger logger.Logger) users.Handlers {
	return &userHandler{cfg: cfg, usersUC: uc, logger: logger}
}

// Create godoc
// @Summary Create User
// @Description Create new user.
// @Tags users
// @Accept json
// @Produce json
// @Param user body presenter.UserCreate true "Add user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user [post]
func (h *userHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(presenter.UserCreate)

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		newUser, err := h.usersUC.CreateUser(
			r.Context(),
			mapModel(user),
			user.ConfirmPassword,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		userResponse := *mapModelResponse(newUser)

		render.Respond(w, r, responses.CreateSuccessResponse(userResponse))
	}
}

// Get godoc
// @Summary Read user
// @Description Get user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id} [get]
func (h *userHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := h.usersUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// GetMulti godoc
// @Summary Read Users
// @Description Retrieve users.
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "limit" Format(limit)
// @Param offset query int false "offset" Format(offset)
// @Success 200 {object} responses.SuccessResponse[[]presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user [get]
func (h *userHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		users, err := h.usersUC.GetMulti(r.Context(), limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(users)))
	}
}

// Delete godoc
// @Summary Delete user
// @Description Delete an user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id} [delete]
func (h *userHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := h.usersUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// Update godoc
// @Summary Update user
// @Description Update an user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Param user body presenter.UserUpdate true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id} [put]
func (h *userHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user := new(presenter.UserUpdate)

		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		values := make(map[string]interface{})
		if user.Name != "" {
			values["name"] = user.Name
		}

		updatedUser, err := h.usersUC.Update(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdatePassword godoc
// @Summary Update password user
// @Description Update password user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Param user body presenter.UserUpdatePassword true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id}/updatepass [patch]
func (h *userHandler) UpdatePassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user := new(presenter.UserUpdatePassword)

		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		updatedUser, err := h.usersUC.UpdatePassword(
			r.Context(),
			id,
			user.OldPassword,
			user.NewPassword,
			user.ConfirmPassword,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// Me godoc
// @Summary Read user me
// @Description Get user me.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/me [get]
func (h *userHandler) Me() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// UpdateMe godoc
// @Summary Update user me
// @Description Update user me.
// @Tags users
// @Accept json
// @Produce json
// @Param user body presenter.UserUpdate true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/me [put]
func (h *userHandler) UpdateMe() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		userUpdate := new(presenter.UserUpdate)

		err = json.NewDecoder(r.Body).Decode(&userUpdate)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), userUpdate)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		values := make(map[string]interface{})
		if userUpdate.Name != "" {
			values["name"] = userUpdate.Name
		}

		updatedUser, err := h.usersUC.Update(r.Context(), user.Id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdatePasswordMe godoc
// @Summary Update password user me
// @Description Update password user me.
// @Tags users
// @Accept json
// @Produce json
// @Param user body presenter.UserUpdatePassword true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/me/updatepass [patch]
func (h *userHandler) UpdatePasswordMe() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		userUpdate := new(presenter.UserUpdatePassword)

		err = json.NewDecoder(r.Body).Decode(&userUpdate)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), userUpdate)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		updatedUser, err := h.usersUC.UpdatePassword(
			r.Context(),
			user.Id,
			userUpdate.OldPassword,
			userUpdate.NewPassword,
			userUpdate.ConfirmPassword,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// LogoutAllAdmin godoc
// @Summary Logout all of user
// @Description Logout all session of user with id.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id}/logoutall [get]
func (h *userHandler) LogoutAllAdmin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		err = h.usersUC.LogoutAll(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}
	}
}

func mapModel(exp *presenter.UserCreate) *models.User {
	return &models.User{
		Name:        exp.Name,
		Email:       exp.Email,
		Password:    exp.Password,
		IsActive:    true,
		IsSuperUser: false,
	}
}

func mapModelResponse(exp *models.User) *presenter.UserResponse {
	return &presenter.UserResponse{
		Id:          exp.Id,
		Name:        exp.Name,
		Email:       exp.Email,
		CreatedAt:   exp.CreatedAt,
		UpdatedAt:   exp.UpdatedAt,
		IsActive:    exp.IsActive,
		IsSuperUser: exp.IsSuperUser,
		Verified:    exp.Verified,
	}
}

func mapModelsResponse(exp []*models.User) []*presenter.UserResponse {
	out := make([]*presenter.UserResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
