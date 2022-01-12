## groupie-tracker

### Objectives

Groupie Trackers consists on receiving a given API and manipulate the data contained in it, in order to create a site, displaying the information.

- It will be given an [API](https://groupietrackers.herokuapp.com/api), that consists in four parts:

  - The first one, `artists`, containing information about some bands and artists like their name(s), image, in which year they began their activity, the date of their first album and the members.

  - The second one, `locations`, consists in their last and/or upcoming concert locations.

  - The third one, `dates`, consists in their last and/or upcoming concert dates.

  - And the last one, `relation`, does the link between all the other parts, `artists`, `dates` and `locations`.

- You must use at least the `artists` and the `relation` parts given in the API.

- Given all this you should build a user friendly website where you can display the bands info through several data visualizations (examples : blocks, cards, tables, list, pages, graphics, etc). It is up to you to decide which info you will present and how you will display it.

- This project also focuses on the creation of events and on their visualization.

  - An event consists in a system that responds to some kind of action triggered by the client, time, or any other factor.

### Instructions

- The backend must be written in **Go**.
- The code must respect the [**good practices**](../good-practices/README.md).
- It is recommended that the code should present a **test file**.

### Allowed packages

- Only the [standard Go](https://golang.org/pkg/) packages are allowed.

-------

### Authors
- mus11110<br>
- zhangir11 <br>

### How to run
- Execute command:  `go run .` or `go build && ./groupietracker`
- Follow the link shown at the terminal: http://localhost:8080/

### Implemention details
- Two endpoints main `"/"` and `"/artist/#"` for artist individual page.
- Both endpoints supports only `GET` method.
- Server is running only after success parsing from API.