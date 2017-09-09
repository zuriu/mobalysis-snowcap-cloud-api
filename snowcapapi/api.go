// Package snowcapapi is the beta version of the Mobalysis® SnowCAP Public API. For more information
// or to obtain keys, please visit http://snowcapalerts.com. Pull requests are always welcome! :)
package snowcapapi

import (
	jwt "github.com/dgrijalva/jwt-go"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(url string, publicKey string, privateKey string) (*Response, error) {

	// Assemble signature claims...
	t := time.Now()
	claims := map[string]interface{}{}
	claims["iat"] = t.Unix()
	claims["exp"] = t.Add(1 * time.Minute).Unix()

	// Assemble the signature...
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims(claims)
	signedToken, err := token.SignedString([]byte(privateKey))

	client := &http.Client{}

	// Assemble the request...
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s %s", publicKey, signedToken))

	// Make request...
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the entity body from the response...
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Check the response signature...
	signature := res.Header.Get("Authorization")
	if signature != "" {
		// TODO - Verify response signature from header
	}

	// Unmarshal the response and finish up...
	response := &Response{StatusCode: res.StatusCode, Body: string(body)}
	json.Unmarshal(body, &response)

	return response, nil
}
