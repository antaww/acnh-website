package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)


func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(`{"message": "hello world"}`))
}



func genshin(index int) string {

	url := "http://api.genshin.dev/characters"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var result []string
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}
	return result[index]
}

func handlerCharacters(w http.ResponseWriter, r *http.Request, result string) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"character": "%s"}`, result)))
}

func main(){
	println("Lancement serveur")
	http.HandleFunc("/", handler)
	for i := 0; i < 47; i++ {
		http.HandleFunc(fmt.Sprintf("/%d", i), func(writer http.ResponseWriter, request *http.Request) {
			nbr, err := strconv.Atoi(strings.TrimPrefix(request.URL.Path, "/"))
			if err != nil {
				return
			}
			fmt.Println(request.URL, "=", genshin(nbr))
			handlerCharacters(writer, request, genshin(nbr))
		})
	}
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		return
	}
}