package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "strconv"
	"github.com/gorilla/mux"
	"os"
	"time"
	jwt "github.com/golang-jwt/jwt/v4"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func newAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPhandleFunc(s.handleUser))
	router.HandleFunc("/user/{id}",withJWTAuth(makeHTTPhandleFunc(s.handleGetUserById)))

	log.Println("JSON API server running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUser(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateUser(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}
func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
    	id,err := getID(r)
		if err != nil {
			return err
		}
		user , err := s.store.getUserById(id)
		
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, user)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteUser(w,r)
	}

	return fmt.Errorf("method not allowed %s" , r.Method)

}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
    users , err := s.store.getUsers()
	
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}
func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
    createUserRequest := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
	   return err
	}

    user := newUser(createUserRequest.Email,createUserRequest.UserName)

	if err := s.store.createUser(user); err != nil {
		return nil
	}

	tokenString, err := createJWT(user)
	if err != nil {
		return err
	}

	fmt.Println("JWT token : ", tokenString)
	return WriteJSON(w,http.StatusOK,user)
}
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id,err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.deleteUser(id); err != nil {
		return err
	}
	return WriteJSON(w , http.StatusOK , map[string]int{"deleted":id})
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}


func withJWTAuth(hendlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter , r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		tokenString := r.Header.Get("x-jwt-token")
		_,err := validateJWT(tokenString)
		if err != nil {
			WriteJSON(w,http.StatusForbidden,ApiError{Error: "invalid token"})
			return 
		}
		hendlerFunc(w,r)
	}
}

func createJWT (user *User) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"Issuer":  "go-api",
		"userId": user.ID,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([] byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token,error){
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString , func(token *jwt.Token) (interface{} , error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); ! ok {
			return nil , fmt.Errorf("Unexpected signinng method : %v" , token.Header["alg"])
		}
		return []byte(secret) , nil
	})
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPhandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle the error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int , error){
	idStr := mux.Vars(r)["id"]
    id,err := strconv.Atoi(idStr)
	if err != nil {
		return id , fmt.Errorf("ivalid id given %s" , idStr)
	}
	return id , nil
}