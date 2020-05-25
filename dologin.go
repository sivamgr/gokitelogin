package gokitelogin

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func enterCredentials(httpClient *http.Client, referURL string, userid string, password string, pin string) error {
	data := url.Values{
		"user_id":  {userid},
		"password": {password},
	}

	const twofaURL = "https://kite.zerodha.com/api/twofa"
	const loginURL = "https://kite.zerodha.com/api/login"

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", referURL)
	r, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	// Declared an empty interface
	var response map[string]*json.RawMessage
	var respdata map[string]*json.RawMessage
	var result string
	var reqid string

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(string(body)), &response)
	json.Unmarshal(*response["status"], &result)
	if result == "success" {
		json.Unmarshal(*response["data"], &respdata)
		json.Unmarshal(*respdata["request_id"], &reqid)
		return enterPIN(httpClient,
			referURL,
			twofaURL,
			userid,
			reqid,
			pin)

	}
	return errors.New("Login Failed")
}

func enterPIN(httpClient *http.Client, referURL string, twofaURL string, userid string, reqid string, pin string) error {
	data := url.Values{
		"user_id":     {userid},
		"request_id":  {reqid},
		"twofa_value": {pin},
	}

	req, err := http.NewRequest("POST", twofaURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", referURL)
	r, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	// Declared an empty interface
	var response map[string]*json.RawMessage
	var respdata map[string]interface{}
	var result string

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(string(body)), &response)
	json.Unmarshal(*response["status"], &result)
	if result == "success" {
		json.Unmarshal(*response["data"], &respdata)
		r, err := httpClient.Get(referURL + "&skip_session=true")
		if err != nil {
			defer r.Body.Close()
		}
		return err
	}
	return errors.New("Failed to perform login")
}

//Login method uses http client to login into Kite API
func Login(KiteLoginURL string, userid string, password string, pin string) error {
	var httpClient *http.Client
	var jar http.CookieJar

	jar, _ = cookiejar.New(nil)
	httpClient = &http.Client{Jar: jar}
	r, err := httpClient.Get(KiteLoginURL)

	if err != nil {
		return err
	}

	defer r.Body.Close()
	return enterCredentials(httpClient,
		r.Request.URL.String(),
		userid, password, pin)
}
