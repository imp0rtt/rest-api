package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"rest-api/internal/apperror"
	"rest-api/internal/handlers"
	"rest-api/pkg/logging"
	"strings"
)

var _ handlers.Handler = &handler{}

const (
	usersURL  = "/users"
	userURL   = "/users/:uuid"
	mainURL   = "/main"
	regURL    = "/reg"
	logURL    = "/log"
	panelURL  = "/panel"
	getaccURL = "/acc"
)

type handler struct {
	logger  *logging.Logger
	service *Service
}

func NewHandler(logger *logging.Logger, service *Service) *handler {
	return &handler{logger: logger, service: service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, mainURL, apperror.Middleware(h.MainPage))
	router.HandlerFunc(http.MethodGet, regURL, apperror.Middleware(h.Registration))
	router.HandlerFunc(http.MethodGet, logURL, apperror.Middleware(h.LogIn))
	router.HandlerFunc(http.MethodGet, panelURL, apperror.Middleware(h.Panel))
	router.HandlerFunc(http.MethodPost, getaccURL, apperror.Middleware(h.GetAcc))
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *handler) MainPage(w http.ResponseWriter, r *http.Request) error {
	err := h.service.MainPage(w)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) Registration(w http.ResponseWriter, r *http.Request) error {
	err := h.service.RegPage(w)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) LogIn(w http.ResponseWriter, r *http.Request) error {
	err := h.service.LogPage(w)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) Panel(w http.ResponseWriter, r *http.Request) error {
	err := h.service.PanelPage(w)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) GetAcc(w http.ResponseWriter, r *http.Request) error {
	err := h.service.GetAcc(r)
	if err != nil {
		return err
	}

	//TODO redirect to panel page
	return nil
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	err := h.service.GetList(w)
	if err != nil {
		return err
	}
	return nil
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	err := h.service.CreateUser(w, r)
	if err != nil {
		return fmt.Errorf("this is API error")
	}
	return nil
}
func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	id := strings.Split(r.URL.Path, "/")[2]

	h.logger.Info(id)
	err := h.service.GetUserByUUID(w, r, id)
	if err != nil {
		return err
	}
	return nil
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id := strings.Split(r.URL.Path, "/")[2]
	h.logger.Debug(id)
	err := h.service.UpdateUser(w, r, id)
	if err != nil {
		return err
	}
	return nil
}
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is PartiallyUpdateUser"))
	return nil

}
func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	id := strings.Split(r.URL.Path, "/")[2]
	h.logger.Debug(id)

	err := h.service.DeleteUser(w, r, id)
	if err != nil {
		return err
	}
	return nil
}
