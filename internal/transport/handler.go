package http

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dionaditya/go-production-ready-api/internal/comment"
	"github.com/dionaditya/go-production-ready-api/internal/user"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	Router *mux.Router
	DB     *gorm.DB
}

type CommentControlelr struct {
	Service *comment.Service
}

type UserController struct {
	Service *user.Service
}

type Response struct {
	Message string
	Error   string
}

type ResponseAPI struct {
	Result interface{}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/register", "/api/health", "/api/login", "/api/refresh-token"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                                               //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(Response{Message: "Missing auth token"}); err != nil {
				panic(err)
			}
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(Response{Message: "Invalid/Malformed auth token"}); err != nil {
				panic(err)
			}
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_CODE")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(Response{Message: "Malinformed authentification token"}); err != nil {
				panic(err)
			}
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(Response{Message: "token invalid"}); err != nil {
				panic(err)
			}
			return
		}

		log.Info(claims["email"])
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func LogginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).Info("handled request")
		next.ServeHTTP(w, r)
	})
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) SetupRoutes() {
	log.Info("Setting up routes")
	commentService := comment.NewService(h.DB)
	commentController := &CommentControlelr{
		Service: commentService,
	}
	userService := user.NewService(h.DB)
	userController := &UserController{
		Service: userService,
	}
	h.Router = mux.NewRouter()
	h.Router.Use(LogginMiddleware)
	h.Router.Use(AuthMiddleware)
	h.Router.HandleFunc("/api/refresh-token", userController.RefreshToken).Methods("POST")
	h.Router.HandleFunc("/api/register", userController.Register).Methods("POST")
	h.Router.HandleFunc("/api/login", userController.Login).Methods("POST")
	h.Router.HandleFunc("/api/comment", commentController.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment", commentController.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", commentController.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", commentController.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", commentController.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charse=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "I am alive"}); err != nil {
			panic(err)
		}

	})
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json; charse=UTF-8")

	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		panic(err)
	}
	return
}
