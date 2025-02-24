package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	usecase UserUsecase
}

type UserUsecase interface {
	GetAll() ([]User, error)
	GetByID(id uint) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

func NewUserHandler(usecase UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", h.GetAll).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/users", h.Create).Methods("POST")
	r.HandleFunc("/users/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/users/{id}", h.Delete).Methods("DELETE")
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.usecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}
