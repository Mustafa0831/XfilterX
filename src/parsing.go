package groupietracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/*Data is a main variable for unmarshaling*/
var Data []Artist

/*Parse is calling GetJSon() function and copying data to Artists struct
and returns nil error in case of success*/
func Parse() error {
	for {
		jsn, err := GetJSON("https://groupietrackers.herokuapp.com/api/artists")
		if err != nil {
			return errors.New(err.Error() + " GetJSON()1st")
		}
		err = json.Unmarshal(jsn, &Data)
		if err != nil {
			return errors.New(err.Error() + " Unmarchal()1st")
		}

		jsn, err = GetJSON("https://groupietrackers.herokuapp.com/api/relation")
		if err != nil {
			return errors.New(err.Error() + " GetJSON()2nd")
		}
		Assist := Assistant{}
		err = json.Unmarshal(jsn, &Assist)
		if err != nil {
			return errors.New(err.Error() + " Unmarchal()2nd")
		}
		for index, DatesLocs := range Assist.Indx {
			Data[index].DatesLocations = DatesLocs.DatesLocations
		}
		fmt.Println("Data Updated")

		time.Sleep(time.Second * 1)
		return err
	}
}

//GetJSON is getting json from API
func GetJSON(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {

		return nil, err
	}

	json, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	return json, err
}
