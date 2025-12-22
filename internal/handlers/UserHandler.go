package handlers

import (
	"encoding/json"
	"go-microservice/internal/models"
	"go-microservice/internal/services"
	"go-microservice/internal/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service     *services.UserService
	auditLogger *utils.AuditLogger
	notifier    *utils.Notifier
}

func NewUserHandler(s *services.UserService, auditLogger *utils.AuditLogger,
	notifier *utils.Notifier) *UserHandler {
	return &UserHandler{
		service:     s,
		auditLogger: auditLogger,
		notifier:    notifier,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Сохранение пользователя (присваивается ID)
	savedUser := h.service.Create(user)

	h.auditLogger.Log(utils.AuditEvent{
		Action: "CREATE",
		Entity: "user",
		ID:     savedUser.ID,
		Time:   time.Now(),
	})

	h.notifier.Send(utils.Notification{
		Type: "USER_CREATED",
		Data: savedUser.ID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedUser)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.service.GetAll())
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := h.service.GetByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	updateUser, err := h.service.Update(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(updateUser)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
