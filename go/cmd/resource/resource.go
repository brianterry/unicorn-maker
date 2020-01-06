package resource

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/encoding"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
)

// APIEndpoint is the crudcrud endpoint.
// You can  obtain an endpoint by going to https://crudcrud.com.
const APIEndpoint = "https://crudcrud.com/api/<YOUR ENDPOINT ID>/unicorns"

// A Unicorn represents unicorn.
type Unicorn struct {
	// ID is the ID of the unicorn.
	ID string `json:"_id,omitempty"`
	// Name is the name of the unicorn.
	Name string `json:"name,omitempty"`
	// Color is the color of the unicorn.
	Color string `json:"color,omitempty"`
}

//RequestInput represents the input when making the HTTP request.
type RequestInput struct {
	// Method is the the HTTP request method.
	Method string
	// URL is the request URL
	URL string
	// Body is the body of the request.
	Body io.Reader
	// Action is the Cloudformation resource action
	// Create, Read, etc...
	Action string
	// Model is the resource model.
	Model *Model
}

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	if exist(req, currentModel) {
		return handler.ProgressEvent{}, errors.New("Resource exist")
	}
	reqBody, err := marshal(currentModel)
	if err != nil {
		return handler.ProgressEvent{}, err
	}
	response := makeRequest(&RequestInput{
		Method: "POST",
		URL:    APIEndpoint,
		Body:   bytes.NewBuffer(reqBody),
		Action: "Create",
	})
	return response, nil
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	if currentModel.UID == nil {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource not found",
			HandlerErrorCode: handler.NotFound,
		}, nil
	}
	response := makeRequest(&RequestInput{
		Method: "GET",
		URL:    APIEndpoint + "/" + *currentModel.UID.Value(),
		Action: "Read",
	})
	return response, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	if !exist(req, currentModel) {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource not found",
			HandlerErrorCode: handler.NotFound,
		}, nil
	}
	reqBody, err := marshal(currentModel)
	if err != nil {
		return handler.ProgressEvent{}, err
	}
	response := makeRequest(&RequestInput{
		Method: "PUT",
		URL:    APIEndpoint + "/" + *currentModel.UID.Value(),
		Body:   bytes.NewBuffer(reqBody),
		Action: "Update",
		Model:  currentModel,
	})
	return response, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	response := makeRequest(&RequestInput{
		Method: "DELETE",
		URL:    APIEndpoint + "/" + *currentModel.UID.Value(),
		Body:   nil,
		Action: "Delete",
	})
	return response, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	response := makeRequest(&RequestInput{
		Method: "GET",
		URL:    APIEndpoint,
		Action: "List",
	})
	return response, nil
}

func exist(req handler.Request, model *Model) bool {
	event, _ := Read(req, nil, model)
	if event.OperationStatus != handler.Failed {
		return true
	}
	return false
}

func marshal(resource *Model) ([]byte, error) {
	u := Unicorn{
		Name:  *resource.Name.Value(),
		Color: *resource.Color.Value(),
	}
	body, err := json.Marshal(&u)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func makeRequest(input *RequestInput) handler.ProgressEvent {
	// Create client
	client := &http.Client{}

	// Create request
	re, err := http.NewRequest(input.Method, input.URL, input.Body)
	if err != nil {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			HandlerErrorCode: handler.InvalidRequest,
			Message:          err.Error(),
		}
	}

	// If the body is not nil, we set the content header
	if input.Body != nil {
		re.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	// Fetch Request
	resp, err := client.Do(re)
	if err != nil {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			HandlerErrorCode: handler.NetworkFailure,
			Message:          err.Error(),
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			HandlerErrorCode: handler.NotFound,
		}

	}
	return makeReturn(input, resp)
}

func unmarshal(unicorn *Unicorn) *Model {
	m := Model{
		UID:   encoding.NewString(unicorn.ID),
		Name:  encoding.NewString(unicorn.Name),
		Color: encoding.NewString(unicorn.Color),
	}
	return &m
}

func makeReturn(input *RequestInput, resp *http.Response) handler.ProgressEvent {
	result := handler.ProgressEvent{
		OperationStatus: handler.Success,
	}
	switch input.Action {
	case "Create":
		u := Unicorn{}
		if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
			return handler.NewFailedEvent(err)
		}
		result.Message = "Create Complete"
		result.ResourceModel = unmarshal(&u)

	case "Read":
		u := Unicorn{}
		if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
			return handler.NewFailedEvent(err)
		}
		result.Message = "Read Complete"
		result.ResourceModel = unmarshal(&u)

	case "Update":
		result.Message = "Update Complete"
		result.ResourceModel = input.Model

	case "Delete":
		result.Message = "Delete Complete"

	case "List":
		var unicorns []Unicorn

		//setting to 1 to return an empty list if there is 0 unicorn
		models := make([]interface{}, 1)
		if err := json.NewDecoder(resp.Body).Decode(&unicorns); err != nil {
			return handler.NewFailedEvent(err)
		}
		for _, unicorn := range unicorns {
			models = append(models, unmarshal(&unicorn))
		}
		result.Message = "List Complete"
		result.ResourceModels = models
	}
	return result
}
