package server

import (
	"net/http"
	"time"
)

func Listen() {

	
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello World"))
	})
	srv := http.Server{
		Addr: ":9091",
		Handler: mux, 
		DisableGeneralOptionsHandler: false,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Protocols: &http.Protocols{}, 
	}


	err := http.ListenAndServe(":9091", http.DefaultServeMux)
	if err != nil {
		panic(err)
	}
}