package user

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"
	"rest-api/internal/apperror"
	"rest-api/pkg/logging"
)

type Service struct {
	storage  Storage
	logger   *logging.Logger
	apperror apperror.AppError
}

func NewService(storage Storage, logger *logging.Logger) *Service {
	return &Service{storage: storage, logger: logger}
}

func (s *Service) MainPage(w http.ResponseWriter) error {

	tmpl, err := template.ParseFiles("public/html/main/index.html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}
	s.logger.Info("main html page was returned")

	return nil
}

func (s *Service) RegPage(w http.ResponseWriter) error {

	tmpl, err := template.ParseFiles("public/html/auth/reg.html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}
	s.logger.Info("reg html page was returned")

	return nil
}

func (s *Service) PanelPage(w http.ResponseWriter) error {

	tmpl, err := template.ParseFiles("public/html/panel/panel.html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}
	s.logger.Info("panel html page was returned")

	return nil
}

func (s *Service) GetAcc(r *http.Request) error {
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	user, err := s.storage.FindByEmail(context.Background(), email)
	if err != nil {
		return err
	}
	pass := user.PasswordHash
	if pass == password {
		return nil
	} else {
		return err
	}
}

func (s *Service) LogPage(w http.ResponseWriter) error {

	tmpl, err := template.ParseFiles("public/html/auth/log.html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}
	s.logger.Info("log html page was returned")

	return nil
}

func (s *Service) GetList(w http.ResponseWriter) error {
	users, err := s.storage.FindAll(context.Background())
	if err != nil {
		return err
	}

	main := filepath.Join("public", "html", "panel", "usersDynamicPage.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(main)
	if err != nil {
		return err
	}
	//исполняем именованный шаблон "users", передавая туда массив со списком пользователей
	err = tmpl.ExecuteTemplate(w, "users", users)
	if err != nil {
		return err
	}
	s.logger.Info("html page was returned")

	return nil
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) error {
	email := r.URL.Query().Get("email")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	user := User{
		ID:           "",
		Email:        email,
		Username:     username,
		PasswordHash: password,
	}
	_, err := s.storage.Create(context.Background(), user)
	if err != nil {
		return err
	}
	s.logger.Debug("new user was created")
	return nil
}
func (s *Service) GetUserByUUID(w http.ResponseWriter, r *http.Request, id string) error {
	user, err := s.storage.FindOne(context.Background(), id)
	if err != nil {
		return err
	}

	//w.Write([]byte(user.ID))
	//w.Write([]byte(user.Email))
	//w.Write([]byte(user.Username))
	//w.Write([]byte(user.PasswordHash))

	main := filepath.Join("public", "html", "panel", "userByUUID.html")

	tmpl, err := template.ParseFiles(main)
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "user", user)
	if err != nil {
		return err
	}
	s.logger.Info("html page was returned")

	return nil
}

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request, id string) error {
	user, err := s.storage.FindOne(context.Background(), id)
	if err != nil {
		return err
	}
	email := r.URL.Query().Get("email")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	user.Email = email
	user.Username = username
	user.PasswordHash = password

	err = s.storage.Update(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) PartiallyUpdateUser() error {
	//TODO Partially update user
	return nil
}
func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request, id string) error {
	err := s.storage.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}
