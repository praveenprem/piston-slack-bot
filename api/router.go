package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type (
	Tls struct {
		Enable      bool   `json:"enable"`
		CertPath    string `json:"certPath"`
		CertKeyPath string `json:"certKeyPath"`
		Proto       string `json:"proto"`
	}
	Server struct {
		Ip      string `json:"ip"`
		Port    int    `json:"port"`
		Timeout int    `json:"timeout"`
		Tls     Tls    `json:"tls"`
	}
)

var (
	apiVersion = "v0"
)

func (s *Server) Start() http.Server {
	handler := appRecovery(middleware(newRouter()))
	return http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.Ip, s.Port),
		Handler:           handler,
		ReadHeaderTimeout: time.Duration(s.Timeout) * time.Second,
		WriteTimeout:      time.Duration(s.Timeout) * time.Second,
	}
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)
	router.NotFoundHandler = http.HandlerFunc(notFound)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method...).
			Path(fmt.Sprintf("/%s/%s", apiVersion, route.Path)).
			Name(route.Name).
			Handler(handler).
			Headers(route.Headers...)
	}
	router.Use(mux.CORSMethodMiddleware(router))
	return router
}

func appRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Catastrophic failure detected")
				_, _ = fmt.Fprintln(os.Stderr, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
				//TODO: Add Error response to show on Slack
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "auth") {
			w.Header().Add("Content-Type", "application/json")
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s Method Not Allowed", r.Method, r.RequestURI)
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "%s not allowd", r.Method)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s Not Found", r.Method, r.RequestURI)

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Page not found")

}
