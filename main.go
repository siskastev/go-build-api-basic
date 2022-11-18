package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Pet struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Owner *Owner `json:"owner"`
}

type Owner struct {
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
}

//fake DB
var pets []Pet

//check name and kind of pet is empty
func (p *Pet) isEmpty() bool {
	return p.Kind == "" && p.Name == ""
}

func main() {
	fmt.Println("API-BASIC")
	//seeding
	pets = append(pets,
		Pet{"1", "Cat Persia", "Momo", 2, &Owner{"Siska Stevani", "Malang"}},
		Pet{"2", "Dog Labrador", "Haki", 2, &Owner{"Aris Antonius", "Sampit"}})

	//routing
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/pets", getAllPets).Methods("GET")
	r.HandleFunc("/pet/{id}", getPet).Methods("GET")
	r.HandleFunc("/pet", create).Methods("POST")
	r.HandleFunc("/pet/{id}", update).Methods("PUT")
	r.HandleFunc("/pet/{id}", delete).Methods("DELETE")

	//listen to a port
	log.Fatal(http.ListenAndServe(":8000", r))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Learn Build API"))
	w.Write([]byte("Following Courses LCO by Hitesth Choundary"))
}

func getAllPets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pets)

}

func getPet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params from request
	params := mux.Vars(r)

	//find id then return the response
	for _, pet := range pets {
		if pet.ID == params["id"] {
			json.NewEncoder(w).Encode(pet)
			return
		}
	}
	json.NewEncoder(w).Encode("Data not Found")
}

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//check body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	var pet Pet
	_ = json.NewDecoder(r.Body).Decode(&pet)
	if pet.isEmpty() {
		json.NewEncoder(w).Encode("No Data")
		return
	}

	//check when name and kind is duplicate
	for _, v := range pets {
		if v.Kind == pet.Kind && v.Name == pet.Name {
			json.NewEncoder(w).Encode("Already Exists")
			return
		}
	}

	// generate unique ID
	rand.Seed(time.Now().UnixNano())
	pet.ID = strconv.Itoa(rand.Intn(100))
	// append pet into pets
	pets = append(pets, pet)
	json.NewEncoder(w).Encode(pet)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, pet := range pets {
		if pet.ID == params["id"] {
			pets = append(pets[:index], pets[index+1:]...)
			var pet Pet
			if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
				return
			}
			pet.ID = params["id"]
			pets = append(pets, pet)
			json.NewEncoder(w).Encode(pets[index])
			return
		}
	}
	json.NewEncoder(w).Encode("Data not Found")
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, pet := range pets {
		if pet.ID == params["id"] {
			pets = append(pets[:index], pets[index+1:]...)
			json.NewEncoder(w).Encode("Success Delete Data")
			return
		}
	}
	json.NewEncoder(w).Encode("Data not Found")
}
