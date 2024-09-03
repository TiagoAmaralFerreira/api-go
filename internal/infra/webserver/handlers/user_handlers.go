package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TiagoAmaralFerreira/api-go/internal/dto"
	"github.com/TiagoAmaralFerreira/api-go/internal/entity"
	"github.com/TiagoAmaralFerreira/api-go/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: userDB,
	}
}

// GetJWT godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body     dto.GetJWTInput  true  "user credentials"
// @Success      200  {object}  dto.GetJWTOutput
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /users/generate_token [post]

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt, ok := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	fmt.Print(jwt)
	if !ok {
		http.Error(w, "JWT auth not found in context", http.StatusInternalServerError)
		return
	}

	jwtExperiesIn, ok := r.Context().Value("jwtExperiesIn").(int)
	if !ok {
		http.Error(w, "JWT expiry time not found in context", http.StatusInternalServerError)
		return
	}

	var user dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !u.ValidatePassword(user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, tokenString, err := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExperiesIn)).Unix(),
	})
	if err != nil {
		http.Error(w, "Failed to create JWT token", http.StatusInternalServerError)
		return
	}

	accessToken := dto.GetJWTOutput{AcessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accessToken)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
