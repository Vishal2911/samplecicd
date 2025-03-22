package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// User represents the user model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory store
type Store struct {
	users []User
	mutex sync.Mutex
}

var store = Store{
	users: make([]User, 0),
}

func main() {
	router := mux.NewRouter()

	// CRUD Endpoints
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", router)
}

// Create User
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	user.ID = len(store.users) + 1
	store.users = append(store.users, user)

	json.NewEncoder(w).Encode(user)
}

// Get All Users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	store.mutex.Lock()
	defer store.mutex.Unlock()

	json.NewEncoder(w).Encode(store.users)
}

// Get User by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	for _, user := range store.users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// Update User
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	for i, user := range store.users {
		if user.ID == id {
			updatedUser.ID = id
			store.users[i] = updatedUser
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// Delete User
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	for i, user := range store.users {
		if user.ID == id {
			store.users = append(store.users[:i], store.users[i+1:]...)
			w.Write([]byte(`{"message": "User deleted successfully"}`))
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}