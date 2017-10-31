package lib

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
    "time"

	"github.com/andybalholm/cascadia"
)

/* GET Method
 *  - returns html for "Hello World" page
 */

func Hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1 id='header' class='something'>Hello, World! <small>inner tag </small></h1><h2> subhead </h2><h1 class='something'>Second Title</h1>")
}


/* POST Method
 *  - returns all html that matches the selector in the given url
 */
type parsePageBody struct {
	Url string `json:"url"`
	Selector string `json:"selector"`
}

type foundElementsResponse struct {
    NumMatches int `json:"numMatches"`
    Matches []string `json:"matches"`
}

func getSelectorMatches(node *html.Node, cssSelector string) ([]*html.Node, error) {

	c, err := cascadia.Compile(cssSelector)

	if err != nil {
		log.Println(err)
		return nil, ErrInvalidSelector(cssSelector)
	}

	selectedNodes := c.MatchAll(node)

	return selectedNodes, nil
}

func ParsePage(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	var toLookUp parsePageBody
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&toLookUp)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
	}
	defer req.Body.Close()

    netClient := &http.Client{
      Timeout: time.Second * 10,
    }

	resp, err := netClient.Get(toLookUp.Url)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid Url", http.StatusBadRequest)
	}
	defer resp.Body.Close()

	temp, err := html.Parse(resp.Body)

	if err != nil {
		log.Println(err)
	}

	selected, err := getSelectorMatches(temp, toLookUp.Selector)
    foundElementsResponse := &foundElementsResponse{
        NumMatches: len(selected),
        Matches: make([]string, 0),
    }

	for _, n := range selected {
		buf := bytes.NewBufferString("")
		if err := html.Render(buf, n); err != nil {
			log.Println(err)
		} else {
            stringHolder := buf.String()
            foundElementsResponse.Matches = append(foundElementsResponse.Matches, stringHolder)
        }
	}

	responseEncoder := json.NewEncoder(w)
    responseEncoder.SetEscapeHTML(false)
    responseEncoder.Encode(foundElementsResponse)
}
