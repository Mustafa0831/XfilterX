package groupietracker

/*Artist is a struct to hold artists information*/
type Artist struct {
	ID             int64    `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	DatesLocations map[string][]string
}

//Assistant helps get relations
type Assistant struct {
	Indx []DatesLocations `json:"index"`
}

//DatesLocations gets Dates and Locations
type DatesLocations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

