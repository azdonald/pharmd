package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
	"github.com/go-chi/chi/v5"
)

type serverImpl struct {
	authManager service.AuthServiceManager
}

func NewServer(authManager service.AuthServiceManager) ServerInterface {
	return &serverImpl{authManager: authManager}
}

func (s serverImpl) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jsonRequest := json.NewDecoder(r.Body)
	var loginRequest LoginRequest
	if err := jsonRequest.Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := s.authManager.Login(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	org, err := s.authManager.GetOrganisationByID(ctx, user.OrganisationID)
	if err != nil {
		log.Println("Error fetching organisation:", err)
		http.Error(w, "Failed to fetch organisation", http.StatusInternalServerError)
		return
	}

	accessToken, err := utils.CreateAccessToken(*user)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := utils.CreateRefreshToken(*user)
	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	expiresIn := 900 // 15 minutes
	response := AuthResponse{
		AccessToken:  &accessToken,
		ExpiresIn:    &expiresIn,
		RefreshToken: &refreshToken,
		User: &User{
			FirstName:         &user.FirstName,
			LastName:          &user.LastName,
			Id:                &user.ID,
			Email:             &user.Email,
			OrganisationName:  &org.Name,
		},
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s serverImpl) PostAuthLogout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s serverImpl) PostAuthRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jsonRequest := json.NewDecoder(r.Body)
	var refreshRequest RefreshRequest
	if err := jsonRequest.Decode(&refreshRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, err := utils.ExtractRefreshClaimFromToken(refreshRequest.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	user, err := s.authManager.GetUserByID(ctx, claims.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	org, err := s.authManager.GetOrganisationByID(ctx, user.OrganisationID)
	if err != nil {
		http.Error(w, "Failed to fetch organisation", http.StatusInternalServerError)
		return
	}

	accessToken, err := utils.CreateAccessToken(*user)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	newRefreshToken, err := utils.CreateRefreshToken(*user)
	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	expiresIn := 900
	response := AuthResponse{
		AccessToken:  &accessToken,
		ExpiresIn:    &expiresIn,
		RefreshToken: &newRefreshToken,
		User: &User{
			FirstName:         &user.FirstName,
			LastName:          &user.LastName,
			Id:                &user.ID,
			Email:             &user.Email,
			OrganisationName:  &org.Name,
		},
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s serverImpl) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jsonRequest := json.NewDecoder(r.Body)
	var registerRequest RegisterRequest
	if err := jsonRequest.Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	org := models.Organisation{
		Name: registerRequest.OrganisationName,
	}

	user := models.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Email:     string(registerRequest.Email),
		Password:  []byte(registerRequest.Password),
	}

	u, err := s.authManager.Register(ctx, org, user)
	if err != nil {
		log.Println("Registration error:", err)
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	orgData, err := s.authManager.GetOrganisationByID(ctx, u.OrganisationID)
	if err != nil {
		log.Println("Error fetching organisation:", err)
		http.Error(w, "Failed to fetch organisation", http.StatusInternalServerError)
		return
	}

	accessToken, err := utils.CreateAccessToken(*u)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := utils.CreateRefreshToken(*u)
	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	expiresIn := 900
	response := AuthResponse{
		AccessToken:  &accessToken,
		ExpiresIn:    &expiresIn,
		RefreshToken: &refreshToken,
		User: &User{
			FirstName:         &u.FirstName,
			LastName:          &u.LastName,
			Id:                &u.ID,
			Email:             &u.Email,
			OrganisationName:  &orgData.Name,
		},
	}

	utils.WriteResponse(ctx, w, response, http.StatusCreated)
}

func (s serverImpl) PutAuthChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	jsonRequest := json.NewDecoder(r.Body)
	var changePasswordRequest PutAuthChangePasswordJSONRequestBody
	if err := jsonRequest.Decode(&changePasswordRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.authManager.ChangePassword(ctx, userID, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword)
	if err != nil {
		http.Error(w, "Failed to change password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(ctx, w, nil, http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterAuthRoutes(r *chi.Mux) http.Handler {
	r.Post("/v1/register", wrapper.PostAuthRegister)
	r.Post("/v1/login", wrapper.PostAuthLogin)
	r.Post("/v1/logout", wrapper.PostAuthLogout)
	r.Post("/v1/refresh", wrapper.PostAuthRefresh)
	r.Put("/v1/change-password", wrapper.PutAuthChangePassword)
	return r
}
