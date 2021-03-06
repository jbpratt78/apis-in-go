package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/alice"
)

type city struct {
	Name string
	Area uint64
}

// Middleware to check content type as JSON
func filterContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Currently in the check content type middleware")
		// filtering requests by MIME type
		if r.Header.Get("Content-type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - Unsupported Media Type. Please send JSON"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// Middleware to add server timestamp for response cookie
func setServerTimeCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		// setting cookie to each and every response
		cookie := http.Cookie{Name: "Server-Time (UTC)", Value: strconv.FormatInt(time.Now().Unix(), 10)}
		http.SetCookie(w, &cookie)
		log.Println("Currently in the set server time middleware")
	})
}

func mainLogic(w http.ResponseWriter, r *http.Request) {
	// check if method is POST
	if r.Method == "POST" {
		var temp city
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&temp)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		// response creation logic goes here, for not its just print
		fmt.Printf("Got %s city with area of %d sq miles\n", temp.Name, temp.Area)
		// report everything went fine
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 - Created"))
	} else {
		// say method not allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}
}

func main() {
	mainLogicHandler := http.HandlerFunc(mainLogic)
	//http.Handle("/city", filterContentType(setServerTimeCookie(mainLogicHandler)))
	chain := alice.New(filterContentType, setServerTimeCookie).Then(mainLogicHandler)
	http.Handle("/city", chain)
	http.ListenAndServe(":8000", nil)
}

// curl -H "Content-Type: application/json" -X POST http://localhost:8000/city -d '{"name":"Boston", "area":89}'
