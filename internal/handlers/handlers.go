package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sangketkit01/personal-block/internal/config"
	"github.com/sangketkit01/personal-block/internal/forms"
	"github.com/sangketkit01/personal-block/internal/helpers"
	"github.com/sangketkit01/personal-block/internal/models"
	"github.com/sangketkit01/personal-block/internal/render"
	"github.com/sangketkit01/personal-block/internal/repository"
	"golang.org/x/crypto/bcrypt"
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


	user, err := m.DB.LoginUser(username, password)
	if err != nil{
		m.App.InfoLog.Println(err)
		m.App.Session.Put(r.Context(),"error","Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(),"flash","Logged in Successfully")
	m.App.Session.Put(r.Context(), "user" , user )
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request){
	blocks, err := m.DB.GetAllBlocks(r)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	user := m.App.Session.Get(r.Context(),"user").(models.User)

	data := make(map[string]interface{})
	data["blocks"] = blocks
	data["user"] = user

	render.Template(w,r,"home.page.tmpl",&models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewPost(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	user := m.App.Session.Get(r.Context(),"user").(models.User)
	content := r.Form.Get("block-content")	

	err = m.DB.InsertNewBlock(user.ID,content)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(),"flash","Create post successfully")
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

func (m *Repository) MyBlock(w http.ResponseWriter, r *http.Request){
	user := m.App.Session.Get(r.Context(),"user").(models.User)

	blocks, err := m.DB.GetBlockByUserID(user.ID)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	data := make(map[string]interface{})
	data["blocks"] = blocks

	render.Template(w,r,"myblock.page.tmpl",&models.TemplateData{
		Data: data,
	})
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request){
	_ = m.App.Session.Destroy(r.Context())

	m.App.Session.Put(r.Context(),"flash","Logged out succuessfully")
	http.Redirect(w,r,"/login",http.StatusTemporaryRedirect)
}

func (m *Repository) ProfilePage(w http.ResponseWriter, r *http.Request){
	user := m.App.Session.Get(r.Context(),"user").(models.User)

	data := make(map[string]interface{})
	data["user"] = user

	render.Template(w,r,"profile.page.tmpl",&models.TemplateData{
		Data: data,
	})
}

func (m *Repository) UpdateProfile(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	user := m.App.Session.Get(r.Context(),"user").(models.User)
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")

	form := forms.New(r.PostForm)
	if len(phone) != 10 || !form.IsNumeric("phone") {
		m.App.InfoLog.Println("Invalid phone")
		m.App.Session.Put(r.Context(),"error","Invalid phone")
		http.Redirect(w,r,"/profile",http.StatusSeeOther)
		return
	}

	user.Name = name
	user.Email = email
	user.Phone = phone
	
	err = m.DB.UpdateProfile(user)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	m.App.Session.Put(r.Context(),"user",user)
	m.App.Session.Put(r.Context(),"flash","Updated profile successfully")
	http.Redirect(w,r,"/profile",http.StatusSeeOther)
}

func (m *Repository) UpdatePassword(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	user := m.App.Session.Get(r.Context(),"user").(models.User)

	oldPassword := r.Form.Get("old-password")
	newPassword := r.Form.Get("new-password")
	confirmPassword := r.Form.Get("confirm-password")

	if len(newPassword) < 8 {
		m.App.Session.Put(r.Context(),"error","Password must have at least 8 characters")
		http.Redirect(w,r,"/profile",http.StatusSeeOther)
		return
	}

	if confirmPassword != newPassword {
		m.App.Session.Put(r.Context(),"error","Password does not match")
		http.Redirect(w,r,"/profile",http.StatusSeeOther)
		return
	}

	err = m.DB.UpdateUserPassword(user.ID,oldPassword,newPassword)
	if err == bcrypt.ErrMismatchedHashAndPassword{
		m.App.Session.Put(r.Context(),"error","Invalid Password")
		http.Redirect(w,r,"/profile",http.StatusSeeOther)
		return
	}else if err != nil{
		helpers.ServerError(w, err)
		return
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword),bcrypt.DefaultCost)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	user.Password = string(newHashedPassword)

	m.App.Session.Put(r.Context(),"user",user)
	m.App.Session.Put(r.Context(),"flash","Updated password successfully")
	http.Redirect(w,r,"/profile",http.StatusSeeOther)
}

func (m *Repository) InsertLike(w http.ResponseWriter, r *http.Request){
	blockID, err := strconv.Atoi(chi.URLParam(r,"id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r,"user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.InsertLikeByPostIDUserID(blockID,userID)
	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	
	http.Redirect(w,r,r.Referer(),http.StatusSeeOther)
}

func (m *Repository) RemoveLike(w http.ResponseWriter, r *http.Request){
	blockID, err := strconv.Atoi(chi.URLParam(r,"id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r,"user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.RemoveLikeByPostIDUserID(blockID,userID)
	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	
	http.Redirect(w,r,r.Referer(),http.StatusSeeOther)
}