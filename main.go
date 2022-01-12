package main

import (
	"fmt"
	GT "groupietracker/src"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//Set timer
	start := time.Now()
	fmt.Println("We are launching the Groupie-tracker project\n...Please wait the loading...")
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	log.SetOutput(file)
	//Parse JSON filter
	go  GT.Parse()
	err = GT.Parse()

	if err != nil {
		log.Println(err.Error())
		return
	}
	//Path of Routes
	//Main page handler
	http.HandleFunc("/", GT.IndexPage)
	http.HandleFunc("/artist/", GT.ArtistPage)
	//Search page handler
	http.HandleFunc("/search", GT.SearchBar)
	//Filter page handler
	http.HandleFunc("/filter/", GT.FilterHandle)
	//Favicon page handler
	http.HandleFunc("/favicon.ico/", GT.FaviconHandler)
	//For static files Directory settings
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//Print timer
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Time elapsed for launching project :", elapsed)
	//LocalHost
	fmt.Println("Server is listening on port 8080...\nHttp Status :", http.StatusOK)
	//Open favorite browser
	GT.Openbrowser("http://localhost:8080/")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println(err)
		return
	}
}
