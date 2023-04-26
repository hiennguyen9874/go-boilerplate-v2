package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/auth"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/middleware"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users/presenter"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/jwt"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/responses"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/utils"
)

type authHandler struct {
	cfg     *config.Config
	usersUC users.UserUseCaseI
	logger  logger.Logger
}

func CreateAuthHandler(uc users.UserUseCaseI, cfg *config.Config, logger logger.Logger) auth.Handlers {
	return &authHandler{cfg: cfg, usersUC: uc, logger: logger}
}

// SignIn godoc
// @Summary Sign In
// @Description Sign in, get access token for future requests.
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "email"
// @Param password formData string true "password"
// @Success 200 {object} presenter.Token
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/login [post]
func (h *authHandler) SignIn() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(presenter.UserSignIn)

		r.ParseMultipartForm(0) //nolint:errcheck
		user.Email = r.FormValue("username")
		user.Password = r.FormValue("password")

		err := utils.ValidateStruct(r.Context(), user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignIn(
			r.Context(),
			user.Email,
			user.Password,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "bearer",
		})
	}
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Get new access token from refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} presenter.Token
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/refresh [get]
func (h *authHandler) RefreshToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		refreshToken := middleware.TokenFromHeader(r)

		accessToken, refreshToken, err := h.usersUC.Refresh(ctx, refreshToken)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "bearer",
		})
	}
}

// GetPublicKey godoc
// @Summary Get public key
// @Description Get rsa public key to decode token.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} presenter.PublicKey
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/publickey [get]
func (h *authHandler) GetPublicKey() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		publicKeyAccessToken, err := jwt.DecodeBase64(h.cfg.Jwt.AccessTokenPublicKey)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		publicKeyRefreshToken, err := jwt.DecodeBase64(h.cfg.Jwt.RefreshTokenPublicKey)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.PublicKey{
			PublicKeyAccessToken:  string(publicKeyAccessToken[:]),
			PublicKeyRefreshToken: string(publicKeyRefreshToken[:]),
		})
	}
}

// Logout godoc
// @Summary Logout
// @Description Logout, remove current refresh token in db.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/logout [get]
func (h *authHandler) Logout() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		refreshToken := middleware.TokenFromHeader(r)

		err := h.usersUC.Logout(ctx, refreshToken)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}
	}
}

// LogoutAllToken godoc
// @Summary Logout all session
// @Description Logout all session of this user.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/logoutall [get]
func (h *authHandler) LogoutAllToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		refreshToken := middleware.TokenFromHeader(r)

		id, err := h.usersUC.ParseIdFromRefreshToken(ctx, refreshToken)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = h.usersUC.LogoutAll(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}
	}
}

// VerifyEmail godoc
// @Summary Verify user
// @Description Verify user using code from email.
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "offset" Format(code)
// @Success 200 {object} string
// @Failure 400	{object} responses.ErrorResponse
// @Router /auth/verifyemail [get]
func (h *authHandler) VerifyEmail() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		verificationCode := q.Get("code")

		err := h.usersUC.Verify(ctx, verificationCode)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse("Email verified successfully"))
	}
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Forgot password, code will send to email.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgotPassword body presenter.ForgotPassword true "Forgot Password"
// @Success 200 {object} string
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/forgotpassword [post]
func (h *authHandler) ForgotPassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		forgotPassword := new(presenter.ForgotPassword)

		err := json.NewDecoder(r.Body).Decode(&forgotPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), forgotPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		err = h.usersUC.ForgotPassword(ctx, forgotPassword.Email)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r,
			responses.CreateSuccessResponse("You will receive a reset email if user with that email exist"))
	}
}

// ResetPassword godoc
// @Summary Reset Password
// @Description Reset Password, using code from email.
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "code" Format(code)
// @Param resetPassword body presenter.ResetPassword true "Reset Password"
// @Success 200 {object} string
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/resetpassword [patch]
func (h *authHandler) ResetPassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		resetToken := q.Get("code")

		resetPassword := new(presenter.ResetPassword)

		err := json.NewDecoder(r.Body).Decode(&resetPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), resetPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		err = h.usersUC.ResetPassword(
			ctx,
			resetToken,
			resetPassword.NewPassword,
			resetPassword.ConfirmPassword,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r,
			responses.CreateSuccessResponse("Password data updated successfully, please re-login"))
	}
}
