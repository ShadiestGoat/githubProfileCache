package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type respRepo struct {
	Name string `json:"name"`
	FullName string `json:"full_name"`
	IsFork bool `json:"fork"`
	Size int `json:"size"`
	Lang string `json:"language"`
}

type langType string

const (
	T_BACKEND  langType = "backend"
	T_FRONTEND langType = "frontend"
)

var languageTypeInfo = map[string][]langType{
	"HTML": {T_FRONTEND},
	"CSS": {T_FRONTEND},
	"EJS": {T_FRONTEND},
	"Go": {T_BACKEND},
	"JavaScript": {T_FRONTEND, T_BACKEND},
	"TypeScript": {T_FRONTEND, T_BACKEND},
	"Lua": {T_BACKEND},
	"Python": {T_BACKEND},
	"React (JS)": {T_FRONTEND},
	"React (TS)": {T_FRONTEND},
	"Svelte": {T_FRONTEND},
	"Shell": {T_BACKEND},
}

type typeInfo struct {
	total int
	m map[string]int
}

func main() {
	godotenv.Load()
	
	go StartCacheLoop()
	
	r := chi.NewRouter()

	r.Get(`/`, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Authentication`) != os.Getenv(`PASSWORD`) {
			w.WriteHeader(401)
			w.Write([]byte(`{"error": "Not Authorized"}`))
			return
		}
		c.RLock()
		p, _ := json.Marshal(c.m)
		c.RUnlock()
		w.Write(p)
	})

	http.ListenAndServe(":" + os.Getenv("PORT"), r)
}
