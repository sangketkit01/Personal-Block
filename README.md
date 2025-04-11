
# Personal-Block 📝

**Personal-Block** is a personal blog web application where users can write posts, view details, leave comments, and like comments. The app includes user authentication, profile management, and post interaction features with a clean and simple interface.

---

## 🌟 Features

- 🏠 **Home Page**
  - Textarea to write new posts
  - Display all blocks (posts) from all users

- 📚 **My Block**
  - Display only the blocks written by the current user

- 💬 **Read Block Page**
  - Show block details
  - Display all comments and their like counts
  - Allow users to like/unlike comments

- 👤 **Profile Page**
  - Edit user profile
  - Change password

- 🔐 **Authentication**
  - Login / Signup / Logout
  - Session management with SCS
  - CSRF protection with nosurf

---

## 🧱 Tech Stack

| Layer       | Technology           |
|-------------|----------------------|
| Backend     | Go (Golang)          |
| Frontend    | Go Templates (`.tmpl`) |
| Routing     | [Chi](https://github.com/go-chi/chi)         |
| Handlers     | [net/http](https://pkg.go.dev/net/http)         |
| Sessions    | [scs](https://github.com/alexedwards/scs)    |
| CSRF Token  | [nosurf](https://github.com/justinas/nosurf) |
| Database    | PostgreSQL           |
| Migrations  | [Soda (Buffalo)](https://gobuffalo.io/en/docs/db/soda) |
| Architecture| Hexagonal (Ports & Adapters) |
| Caching     | Implemented          |

---

## 🔁 Route Overview

```go
mux.Get("/login", handlers.Repo.LoginPage)
mux.Post("/login/verify", handlers.Repo.LoginVerify)
mux.Get("/signup", handlers.Repo.SignUpPage)
mux.Post("/signup/insert", handlers.Repo.SignUpInsert)

mux.Group(func (r chi.Router) {
    r.Use(Auth)
    r.Get("/", handlers.Repo.Home)
    r.Get("/logout", handlers.Repo.Logout)

    r.Get("/myblock", handlers.Repo.MyBlock)
    r.Get("/profile", handlers.Repo.ProfilePage)

    r.Post("/update-profile", handlers.Repo.UpdateProfile)
    r.Post("/update-password", handlers.Repo.UpdatePassword)

    r.Post("/new-post", handlers.Repo.NewPost)

    r.Post("/insert-like/{id}/{user_id}", handlers.Repo.InsertLike)
    r.Post("/remove-like/{id}/{user_id}", handlers.Repo.RemoveLike)
    r.Post("/insert-comment-like/{id}/{user_id}", handlers.Repo.InsertCommentLike)
    r.Post("/remove-comment-like/{id}/{user_id}", handlers.Repo.RemoveCommentLike)

    r.Get("/read-block/{id}", handlers.Repo.ReadBlock)
    r.Post("/insert-comment/{block_id}", handlers.Repo.InsertComment)
})
```

---

## 🧰 Helper Functions

- Server-side error handling
- Client-side error handling
- Not-found response handler
- Utility functions to improve code clarity

---

## 🚀 Installation Guide

### 1. Clone the project

```bash
git clone https://github.com/sangketkit01/Personal-Block.git
cd Personal-Block
```

### 2. Create `database.yml`

Create a file `database.yml` and `.env` in the project root :

```yaml
development:
  dialect: postgres
  database: postgres
  user: postgres
  password: 
  host: 127.0.0.1
  pool: 5

test:
  url: {{envOr "TEST_DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/myapp_test"}}

production:
  url: {{envOr "DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/myapp_production"}}
```

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=
```

### 3. Install Soda and run the project

```bash
go install github.com/gobuffalo/pop/soda@latest
soda migrate

cd cmd
go run .
```

---

## 📂 Project Structure

```
Personal-Block/
│
├── cmd/
│   ├── main.go
│   ├── middleware.go
│   └── route.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── driver/
│   │   └── driver.go
│   ├── forms/
│   │   └── forms.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── helpers/
│   │   └── helpers.go
│   ├── models/
│   │   ├── models.go
│   │   └── templatedata.go
│   ├── render/
│   │   └── render.go
│   └── repository/
│       └── db.go
│
├── migrations/
├── static/
├── templates/
├── .env
├── .gitignore
├── database.yml
├── go.mod
├── go.sum
└── README.md


---

## 📄 License

This project is licensed under the MIT License.  
Feel free to use, modify, and share it.

---

## 💡 Author

Developed by [@sangketkit01](https://github.com/sangketkit01)
