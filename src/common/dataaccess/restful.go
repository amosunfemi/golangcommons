package dataaccess

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	utility "common/utility"
	"github.com/parnurzeal/gorequest"
	log "gopkg.in/inconshreveable/log15.v2"
)

//RESTFul ... Restful data source
type RESTFul struct {
	ServiceURI string
	Logger     log.Logger
}

var request = gorequest.New()

//SaveObject ... This method saves any object calling the POST method
func (restful *RESTFul) SaveObject(entity interface{}, endpoint string) (interface{}, int, error) {
	var restErr *RESTError
	resp, body, _ := request.Post(restful.ServiceURI + "/" + endpoint).Send(entity).End()
	fmt.Println(reflect.TypeOf(entity))
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		//err := convertBodyToInterface(body, &entity)
		if err := json.Unmarshal([]byte(body), &entity); err != nil {
			restErr := restful.composeCustomError(body, resp.StatusCode)
			return nil, 1, restErr
		}
		restful.Logger.Debug(fmt.Sprintf("Saved %v", entity))
		return entity, 0, nil
	}
	restErr = restful.composeCustomError(body, resp.StatusCode)
	return nil, 1, restErr
}

//UpdateObject ... This method saves any object calling the POST method
func (restful *RESTFul) UpdateObject(reqbody string, endpoint string) (interface{}, int, error) {
	var restErr *RESTError
	resp, body, _ := request.Patch(restful.ServiceURI + "/" + endpoint).Send(reqbody).End()
	entity := map[string]interface{}{}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		//err := convertBodyToInterface(body, &entity)
		if err := json.Unmarshal([]byte(body), &entity); err != nil {
			restErr := restful.composeCustomError(body, resp.StatusCode)
			return nil, 1, restErr
		}
		return entity, 0, nil
	}
	restErr = restful.composeCustomError(body, resp.StatusCode)
	return nil, 1, restErr
}

//UpdateObjectAuthToken ... Update entites with token
func (restful *RESTFul) UpdateObjectAuthToken(reqbody string, endpoint string, token string) (interface{}, int, error) {
	var restErr *RESTError
	resp, body, _ := request.Patch(restful.ServiceURI+"/"+endpoint).Set("Authorization", token).Send(reqbody).End()
	entity := map[string]interface{}{}
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		if err := json.Unmarshal([]byte(body), &entity); err != nil {
			restErr := restful.composeCustomError(body, resp.StatusCode)
			return nil, 1, restErr
		}
		return entity, 0, nil
	}
	restErr = restful.composeCustomError(body, resp.StatusCode)
	return nil, 1, restErr
}

//Set("Authorization", r.Header.Get("Authorization"))

//RetrieveObject ...
func (restful *RESTFul) RetrieveObject(endpoint string) ([]interface{}, error) {
	var restErr *RESTError
	var entities []interface{}
	resp, _, reqErrs := request.Get(restful.ServiceURI + "/" + endpoint).End()
	if len(reqErrs) == 0 {
		err := getJSON(resp, entities)
		if err != nil {
			return nil, err
		}
	}
	return entities, restErr
}

//DeleteObject ...
func (restful *RESTFul) DeleteObject(endpoint string) error {
	var restErr *RESTError
	resp, body, _ := request.Delete(restful.ServiceURI + "/" + endpoint).End()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		restErr = restful.composeCustomError(body, resp.StatusCode)
	}
	return restErr
}

//SaveObjectAuthToken ... This method saves any object calling the POST method
func (restful *RESTFul) SaveObjectAuthToken(entity interface{}, endpoint string, token string) (interface{}, int, error) {
	var restErr *RESTError
	resp, body, _ := request.Post(restful.ServiceURI+"/"+endpoint).Set("Authorization", token).Send(entity).End()
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		//err := convertBodyToInterface(body, &entity)
		restful.Logger.Debug(fmt.Sprintf("Saved %v", resp))
		if err := json.Unmarshal([]byte(body), &entity); err != nil {
			restErr := restful.composeCustomError(body, resp.StatusCode)
			return nil, 1, restErr
		}
		restful.Logger.Debug(fmt.Sprintf("Saved %v", entity))
		return entity, 0, nil
	}
	restErr = restful.composeCustomError(body, resp.StatusCode)
	return nil, 1, restErr
}

//RetrieveObjectAuthToken ...
func (restful *RESTFul) RetrieveObjectAuthToken(endpoint string, token string) ([]interface{}, error) {
	var err error
	var entities []interface{}
	resp, _, reqErrs := request.Get(restful.ServiceURI + endpoint).End()
	if len(reqErrs) == 0 {
		err := getJSON(resp, entities)
		if err != nil {
			return nil, err
		}
	}
	return entities, err
}

func getJSON(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func convertBodyToInterface(reqbody string, entity interface{}) error {
	if err := json.Unmarshal([]byte(reqbody), &entity); err != nil {
		return err
	}
	return nil
}

//RESTError ...
type RESTError struct {
	utility.CustomErr // Err is an embedded struct - ErrWidget_A inherits it's data and behavior
}

// Error ...
func (resterr *RESTError) Error() string {
	return fmt.Sprintf("%s [%d] %s", resterr.When, resterr.ErrNo, resterr.Msg)
}

//composeCustomError
func (restful *RESTFul) composeCustomError(body string, errInt int) *RESTError {
	errCustom := new(RESTError)
	errCustom.ErrNo = errInt
	errCustom.Msg = body
	errCustom.When = time.Now()
	return errCustom
}
