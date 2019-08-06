package backend

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Details(w http.ResponseWriter, r *http.Request) {

	httpClient := &http.Client{Timeout: 2 * time.Second}

	url := "https://api.github.com/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintln(w, string(body))
}
