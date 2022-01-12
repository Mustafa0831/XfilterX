package groupietracker

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

//IndexPage is handling and parsing base page
func IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/index.html")
	switch {
	case r.URL.Path != "/":
		ErrorHandler(w, "404: Not Found", http.StatusNotFound)
	case r.Method != http.MethodGet:
		ErrorHandler(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
	case err != nil:
		ErrorHandler(w, "500: Internal Server Error", http.StatusInternalServerError)
	default:
		tmpl.Execute(w, Data)
	}
}

//ArtistPage is handling and parsing Artist's each page
func ArtistPage(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Path[len("/artist/"):]
	if artistID == "" {
		ErrorHandler(w, "404: Not Found", http.StatusNotFound)
		return
	}
	id, err1 := strconv.Atoi(artistID)
	tpl, err2 := template.ParseFiles("static/templates/artist.html")

	switch {
	case r.Method != http.MethodGet:
		ErrorHandler(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
	case err1 != nil || err2 != nil:
		ErrorHandler(w, "500: Internal Server Error", http.StatusInternalServerError)
	case id <= 0 || id > len(Data):
		ErrorHandler(w, "404: Not Found", http.StatusNotFound)
	default:
		artist := Data[id-1]
		tpl.Execute(w, artist)
	}
}

//ErrorHandler is parsing errors to special template. In case of missing template,error is parsed through http.Error() function.
func ErrorHandler(w http.ResponseWriter, status string, errorcase int) {
	tmpl, err := template.ParseFiles("static/templates/error.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "500: Internal server error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(errorcase)
		tmpl.Execute(w, status)
	}
}

//SearchBar is looking for by members,creationdate, firtstalbum and relationdata
func SearchBar(w http.ResponseWriter, r *http.Request) {
	var searchResult []Artist
	searchOption := r.FormValue("options")
	searchText := r.FormValue("textFind")

	//Suggestion
	if strings.Contains(searchText, " -> ") {
		searchOption = searchText[strings.Index(searchText, " -> ")+len(" -> "):]
		searchText = searchText[:strings.Index(searchText, " -> ")]
	}
	//Fill Data
	for _, copy := range Data {
		switch searchOption {
		case "Artist":
			if strings.Contains((copy.Name), (searchText)) {
				searchResult = append(searchResult, copy)
				continue
			}
		case "Members":
			for _, member := range copy.Members {
				if strings.Contains((member), (searchText)) {
					searchResult = append(searchResult, copy)
					break
				}
			}
			continue
		case "Creation Date":
			if strconv.Itoa(int(copy.CreationDate)) == searchText {
				searchResult = append(searchResult, copy)
			}
			continue
		case "First Album":
			if copy.FirstAlbum == searchText {
				searchResult = append(searchResult, copy)
			}
			continue
		case "Location":
			for value := range copy.DatesLocations {
				if strings.Contains(value, searchText) {
					searchResult = append(searchResult, copy)
					break
				}
			}
			continue
		default:
			w.WriteHeader(http.StatusBadRequest)
			ErrorHandler(w, "400: Bad Request", http.StatusBadRequest)
			return
		}
	}

	if searchText == "" && searchOption != "Artist" && searchOption != "Members" && searchOption != "Creation Date" && searchOption != "First Album" && searchOption != "Location" {
		w.WriteHeader(http.StatusBadRequest)
		ErrorHandler(w, "400: Bad Request", http.StatusBadRequest)
	} else if searchResult == nil {
		ErrorHandler(w, "Not Found", http.StatusOK)
	}
	HandleSearch(w, searchResult)
	return
}

//HandleSearch is handling a search
func HandleSearch(w http.ResponseWriter, r []Artist) {
	tpl, err := template.ParseFiles("static/templates/search.html")
	if err != nil {
		fmt.Println(err)
		ErrorHandler(w, "500: Internal Server Error", http.StatusInternalServerError)
	}
	tpl.Execute(w, r)
}

//FaviconHandler is handling for favicon
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/img/favicon.ico")
}

//Openbrowser allows your server open on browser
func Openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		fmt.Println("Can't find your os, please open yourself localhost:8080")
	}
	if err != nil {
		log.Println(err)
	}

}

//FilterHandle is handling for error
func FilterHandle(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/filter.html")
	if err != nil {
		fmt.Println(err)
		ErrorHandler(w, "500: Internal Server Error", http.StatusInternalServerError)
	}
	tmpl.Execute(w, Filter(w,r))
}

//Filter is filtering all data from API
func Filter(w http.ResponseWriter, r *http.Request) []Artist{
	var filterResult []Artist

	var Inputs struct {
		CreationDate struct {
			From string
			To   string
		}
		FirstAlbumDate struct {
			From string
			To   string
		}
		NumberOfMembers struct {
			From string
			To   string
		}
		Location string
	}
	//Check for invalid request data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, "404: Not Found", http.StatusNotFound)
	}

	query, err := url.ParseQuery(string(body))
	if err != nil {
		ErrorHandler(w, "404: Not Found", http.StatusNotFound)
	}
	for word, ArtistStruct := range query {
		switch word {
		case "cd-from":
			Inputs.CreationDate.From = ArtistStruct[0]
		case "cd-to":
			Inputs.CreationDate.To = ArtistStruct[0]
		case "fad-from":
			Inputs.FirstAlbumDate.From = ArtistStruct[0]
		case "fad-to":
			Inputs.FirstAlbumDate.To = ArtistStruct[0]
		case "nom-from":
			Inputs.NumberOfMembers.From = ArtistStruct[0]
		case "nom-to":
			Inputs.NumberOfMembers.To = ArtistStruct[0]
		case "loc":
			Inputs.Location = ArtistStruct[0]
		default:
			ErrorHandler(w, "404: Not Found", http.StatusNotFound)
		}
	}
	//Filtering by user inputs
	for i := 0; i < len(Data); i++ {
		if !compareCreationDate(Inputs.CreationDate.From, Inputs.CreationDate.To, i) {
			continue
		}
		if !compareFirstAlbumDate(Inputs.FirstAlbumDate.From, Inputs.FirstAlbumDate.To, i) {
			continue
		}
		if !compareNumberOfMembers(Inputs.NumberOfMembers.From, Inputs.NumberOfMembers.To, i) {
			continue
		}
		if !compareLocation(Inputs.Location, i) {
			continue
		}
		filterResult = append(filterResult, Data[i])
	}
	if filterResult == nil {
		ErrorHandler(w, "Not Found", http.StatusOK)
		
	}
    return filterResult
}

//Compare users creation dates with Artists creation dates
func compareCreationDate(from string, to string, index int) bool {
	if from == "" && to == "" {
		return true
	}
	if from != "" && to != "" {
		compare := Data[index].CreationDate
		fromN, err := strconv.Atoi(from)
		toN, err2 := strconv.Atoi(to)
		if err != nil || err2 != nil {
			return false
		}
		if compare >= fromN && compare <= toN {
			return true
		}
	}
	if from != "" && to == "" {
		compare := Data[index].CreationDate
		fromN, err := strconv.Atoi(from)
		if err != nil {
			return false
		}
		if compare >= fromN {
			return true
		}
	}
	if from == "" && to != "" {
		compare := Data[index].CreationDate
		toN, err := strconv.Atoi(to)
		if err != nil {
			return false
		}
		if compare <= toN {
			return true
		}
	}
	return false
}

//Compare users first album date with Artists first album date
func compareFirstAlbumDate(from string, to string, index int) bool {
	if from == "" && to == "" {
		return true
	}
	if from != "" && to != "" {
		fullDate := []rune(Data[index].FirstAlbum)
		compare, err := strconv.Atoi(string(fullDate[len(fullDate)-4:]))
		fromN, err1 := strconv.Atoi(from)
		toN, err2 := strconv.Atoi(to)
		if err != nil || err1 != nil || err2 != nil {
			return false
		}
		if compare >= fromN && compare <= toN {
			return true
		}
	}
	if from != "" && to == "" {
		fullDate := Data[index].FirstAlbum
		compare, err := strconv.Atoi(fullDate[len(fullDate)-4:])
		fromN, err := strconv.Atoi(from)
		if err != nil {
			return false
		}
		if compare >= fromN {
			return true
		}
	}
	if from == "" && to != "" {
		fullDate := Data[index].FirstAlbum
		compare, err := strconv.Atoi(fullDate[len(fullDate)-4:])
		toN, err := strconv.Atoi(to)
		if err != nil {
			return false
		}
		if compare <= toN {
			return true
		}
	}
	return false
}

//Compare users number of members with Artists number of members
func compareNumberOfMembers(from string, to string, index int) bool {
	if from == "" && to == "" {
		return true
	}
	if from != "" && to != "" {
		compare := len(Data[index].Members)
		fromN, err := strconv.Atoi(from)
		toN, err2 := strconv.Atoi(to)
		if err != nil || err2 != nil {
			return false
		}
		if compare >= fromN && compare <= toN {
			return true
		}
	}
	if from != "" && to == "" {
		compare := len(Data[index].Members)
		fromN, err := strconv.Atoi(from)
		if err != nil {
			return false
		}
		if compare >= fromN {
			return true
		}
	}
	if from == "" && to != "" {
		compare := len(Data[index].Members)
		toN, err := strconv.Atoi(to)
		if err != nil {
			return false
		}
		if compare <= toN {
			return true
		}
	}
	return false
}

//Compare users location with Artists locations
func compareLocation(location string, index int) bool {
	loc := strings.Split(location, ", ")
	city := strings.ToLower(loc[0])
	location = strings.ToLower(location)
	
	if location == "" {
		return true
	}

	for v := range Data[index].DatesLocations {
		if strings.Index(strings.ToLower(v), location) != -1 {
			return true
		}
		if strings.Index(strings.ToLower(v), city) != -1 {
			return true
		}
	}
	return false
}
