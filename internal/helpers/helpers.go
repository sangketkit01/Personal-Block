package helpers

import (
	"fmt"
	"github.com/sangketkit01/personal-block/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w,
		fmt.Sprintf("%s\n%s\n%s", http.StatusText(http.StatusInternalServerError), err.Error(), debug.Stack()),
		http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, err error) {
	app.InfoLog.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user")
	return exists
}
