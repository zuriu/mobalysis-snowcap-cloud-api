package snowcapapi

import (
	jwt "github.com/dgrijalva/jwt-go"

	"fmt"
	"reflect"
	"time"
)

type Response struct {
	StatusCode int    `json:"-"`
	Body       string `json:"-"`
	Token      string `json:"token"`
	Error      string `json:"error"`
}

func (r *Response) ParseEvents(privateKey string) (*EventListResults, error) {

	list, err := r.ParseList(privateKey)
	if err != nil {
		return nil, err
	}

	listResults, ok := list.Results.([]interface{})
	if !ok {
		return nil, fmt.Errorf("listResults not of type []interface{}")
	}

	events := []Event{}
	for _, result := range listResults {

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("resultMap not of type map[string]interface{}")
		}

		event := Event{}
		typ := reflect.TypeOf(event)
		val := reflect.ValueOf(&event).Elem()
		for i := 0; i < typ.NumField(); i++ {

			fieldType := typ.Field(i)
			fieldValue := val.Field(i)

			switch fieldValue.Kind() {
			case reflect.Float64:
				mapValue, ok := resultMap[fieldType.Tag.Get("json")].(float64)
				if !ok {
					continue
				}
				fieldValue.SetFloat(mapValue)

			case reflect.String:
				mapValue, ok := resultMap[fieldType.Tag.Get("json")].(string)
				if !ok {
					continue
				}
				fieldValue.SetString(mapValue)

			case reflect.Struct:
				switch fieldValue.Type().Name() {
				case "Time":

					mapValue, ok := resultMap[fieldType.Tag.Get("json")].(string)
					if !ok {
						continue
					}

					mapValueTime, err := time.Parse("2006-01-02T15:04:05Z", mapValue)
					if err != nil {
						continue
					}

					fieldValue.Set(reflect.ValueOf(mapValueTime))
				}
			}

		}

		events = append(events, event)
	}

	eventList := &EventListResults{}
	eventList.TotalCount = list.TotalCount
	eventList.Results = events
	return eventList, nil
}

func (r *Response) ParseList(privateKey string) (*ListResults, error) {

	result, err := r.ParseResult(privateKey)
	if err != nil {
		return nil, err
	}

	totalCount, ok := result["totalCount"].(float64)
	if !ok {
		return nil, fmt.Errorf("totalCount not of type int64")
	}

	results, ok := result["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("results not of type []interface{}")
	}

	list := &ListResults{}
	list.TotalCount = int64(totalCount)
	list.Results = results
	return list, nil
}

func (r *Response) ParseResult(privateKey string) (map[string]interface{}, error) {

	token, err := jwt.Parse(r.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("response token is invalid")
	}

	claims := token.Claims.(jwt.MapClaims)
	result, ok := claims["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response not of type map[string]interface{}")
	}

	return result, nil
}
