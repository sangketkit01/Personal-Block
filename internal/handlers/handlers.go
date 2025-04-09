package handlers

import (
	"github.com/sangketkit01/personal-block/internal/config"
	"github.com/sangketkit01/personal-block/internal/forms"
	"github.com/sangketkit01/personal-block/internal/helpers"
	"github.com/sangketkit01/personal-block/internal/models"
	"github.com/sangketkit01/personal-block/internal/render"
	"github.com/sangketkit01/personal-block/internal/repository"
	"net/http"
)

// Repo is the handlers repository
var Repo *Repository

// Repository holds app config and able handlers to access DatabaseRepo functions
type Repository struct {
	App *config.AppConfig
	DB  *repository.DBRepo
}

// NewRepository creates a new repository
func NewRepository(app *config.AppConfig, db *repository.DBRepo) *Repository {
	return &Repository{
		App: app,
		DB:  db,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

func (m *Repository) SignUpPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "signup.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) SignUpInsert(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.InfoLog.Println("cannot parse form")
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	username := form.Get("username")
	email := form.Get("email")
	phone := form.Get("phone")
	password := form.Get("password")
	confirmPassword := form.Get("confirm_password")

	form.MinLength("password", 8)
	form.IsEmail("email")

	if len(phone) != 10 {
		form.Error.Add("phone", "invalid phone number")

	}

	if password != confirmPassword {
		form.Error.Add("confirm_password", "password does not match")
	}

	if !form.Valid() {
		render.Template(w, r, "signup.page.tmpl", &models.TemplateData{
			Form: form,
		})

		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: password,
		Name: username,
	}

	err = m.DB.InsertUser(user)
	if err != nil {
		m.App.InfoLog.Println("cannot insert user")
		helpers.ServerError(w, err)
		return
	}

	m.App.InfoLog.Println("user created")
	m.App.Session.Put(r.Context(),"flash","Signup successfully")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}


func (m *Repository) LoginVerify(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	err = m.DB.LoginUser(username, password)
	if err != nil{
		m.App.Session.Put(r.Context(),"error","Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(),"flash","Login Successfully")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}