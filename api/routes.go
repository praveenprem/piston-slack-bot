package api

import (
	"github.com/praveenprem/testbed-slack-bot/api/controllers"
	"net/http"
)

type (
	Route struct {
		Name        string
		Headers     []string
		Method      []string
		Path        string
		Secure      bool
		HandlerFunc http.HandlerFunc
	}
	Routes []Route
)

var routes = Routes{
	Route{
		Name:        "Help",
		Headers:     []string{"Content-type", "application/x-www-form-urlencoded"},
		Method:      []string{http.MethodPost, http.MethodOptions},
		Path:        "help",
		Secure:      false,
		HandlerFunc: controllers.Help,
	},
	Route{
		Name:        "Languages",
		Headers:     []string{"Content-type", "application/x-www-form-urlencoded"},
		Method:      []string{http.MethodPost, http.MethodOptions},
		Path:        "languages",
		Secure:      false,
		HandlerFunc: controllers.Languages,
	},
	Route{
		Name:        "Execute",
		Headers:     []string{"Content-type", "application/x-www-form-urlencoded"},
		Method:      []string{http.MethodPost, http.MethodOptions},
		Path:        "execute",
		Secure:      false,
		HandlerFunc: controllers.Execute,
	},
	Route{
		Name:        "Authorise",
		//Headers:     []string{"Accept", "*"},
		Method:      []string{http.MethodGet, http.MethodOptions},
		Path:        "auth",
		Secure:      false,
		HandlerFunc: controllers.Authorise,
	},
}
