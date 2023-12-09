package net

/*
*	
*	Handles network traffic.
*	
 */

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/buffermet/epoxy/log"
	"github.com/buffermet/epoxy/session"
)

var (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
)

func SendRequest(url string, s *session.SessionConfig) ([]byte, string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		log.Error("malformed request packet for " + url + " (" + err.Error() + ")")
		return []byte(""), ""
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := client.Do(req)
	if err != nil {
		log.Error("cannot retrieve resource at " + url + " (" + err.Error() + ")")
		return []byte(""), ""
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("cannot read body of response from " + url + " (" + err.Error() + ")")
		return []byte(""), ""
	}

	return body, res.Header.Get("Content-Type")
}
