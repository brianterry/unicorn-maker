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
const APIEndpoint = "https://crudcrud.com/api/<Your API ID>/unicorns"

// A Unicorn represents a unicorn.
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
	if err := validateInput(req, currentModel); err != nil {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
		}, nil
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
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
		}, nil
	}
	response := makeRequest(&RequestInput{
		Method: "GET",
		URL:    APIEndpoint + "/" + aws.StringValue(currentModel.UID),
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
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
		}, nil
	}
	reqBody, err := marshal(currentModel)
	if err != nil {
		return handler.ProgressEvent{}, err
	}
	response := makeRequest(&RequestInput{
		Method: "PUT",
		URL:    APIEndpoint + "/" + aws.StringValue(currentModel.UID),
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
		URL:    APIEndpoint + "/" + aws.StringValue(currentModel.UID),
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

func validateInput(req handler.Request, model *Model) error {
	if exist(req, model) {
		return errors.New("Resource exist")
	}
	if model.Name == nil {
		return errors.New("Name required")
	}
	if model.Color == nil {
		return errors.New("Color required")
	}
	return nil
}

func marshal(resource *Model) ([]byte, error) {
	u := Unicorn{}
	if resource.Name != nil {
		u.Name = aws.StringValue(resource.Name)
	}
	if resource.Color != nil {
		u.Color = aws.StringValue(resource.Color)
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
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
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
			HandlerErrorCode: cloudformation.HandlerErrorCodeNetworkFailure,
			Message:          err.Error(),
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 || resp.StatusCode == 400 {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
		}

	}
	return makeReturn(input, resp)
}

func unmarshal(unicorn *Unicorn) *Model {
	m := Model{
		UID:   aws.String(unicorn.ID),
		Name:  aws.String(unicorn.Name),
		Color: aws.String(unicorn.Color),
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

		// The cloudformation service requires that an empty array
		// be return if there are 0 unicorns.
		// Because the model struct has the
		// tag: omitempty, the value is omitted
		// in the JSON return, causing an error.
		// So, make a slice and set it to 1 to
		//return an empty list if there is 0 unicorn
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
