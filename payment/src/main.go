package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
}

type PaymentRequest struct {
	ID       int     `json:"customer_id,omitempty"`
	Currency string  `json: "currency,omitempty"`
	Value    float32 `json: "value,omitempty"`
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":9000", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.makePayment).Methods("POST")
	a.Router.HandleFunc("/health", a.healthCheck).Methods("GET")

}

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"alive": "true"})
}

func (a *App) makePayment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
		return
	}
	println(string(body))
	var d PaymentRequest
	err = json.Unmarshal(body, &d)
	if err != nil {
		println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
		return
	}

	if d.Currency == "" || d.ID < 1 || d.Value < 0 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid Request. Missing required fields"})
		return
	}
	response := rand.Intn(2)
	if response == 0 {
		respondWithJSON(w, http.StatusOK, map[string]string{"result": "true"})
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "false"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	a := App{}
	a.Initialize()
	a.Run(":9000")
}
