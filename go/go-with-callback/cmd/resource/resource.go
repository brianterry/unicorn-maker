package resource

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// APIEndpoint is the crudcrud endpoint.
// You can  obtain an endpoint by going to https://crudcrud.com.
var APIEndpoint string = fmt.Sprintf("https://crudcrud.com/api/%s/unicorns", "ed9e3091188f4f3e82258b3990b5b0b5")

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
	// First, we check if callback context is set. When your handler is reinvoked you can use this property to identify
	// the current context. in this example, we just use a string to see if the resource is stable.

	_, ok := req.CallbackContext["status"]

	if ok {
		// We call the read handler to see if the resource is stable. You can use the value of the Callback Context
		// to determine the state. You can compare it to the a prop in the model, etc...
		// But since this is a simple example, if the read handler returns an error, we will
		// ask AWS CloudFormation to reinvoke the handler because we assume the resource is not stable.
		_, err := Read(req, prevModel, currentModel)
		if err != nil {
			r := handler.ProgressEvent{}
			r.CallbackContext = map[string]interface{}{"status": "stabilizing"}
			r.CallbackDelaySeconds = 60
			r.Message = "Create in progress.... "
			r.OperationStatus = handler.InProgress
			r.ResourceModel = currentModel
		}
		// In this example, we assume that if the read was successful, the resource is stable so we return success.
		return handler.ProgressEvent{
			OperationStatus: handler.Success,
			Message:         "Update Complete",
		}, nil

	}

	if err := validateInput(req, currentModel); err != nil {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
		}, nil
	}

	b, err := marshal(currentModel)
	if err != nil {
		return handler.ProgressEvent{}, err
	}
	r := makeRequest(&RequestInput{
		Method: "POST",
		URL:    APIEndpoint,
		Body:   bytes.NewBuffer(b),
		Action: "Create",
	})
	return r, nil
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
	// Each handler must return a ProgressEvent
	result := handler.ProgressEvent{
		OperationStatus: handler.Success,
	}
	switch input.Action {
	case "Create":
		u := Unicorn{}
		if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
			return handler.NewFailedEvent(err)
		}
		// Notice that we set the CallbackContext. The CallbackContext is a map[string]interface{},and
		// you can store any type.
		result.CallbackContext = map[string]interface{}{"status": "stabilizing"}
		// Setting CallbackDelaySeconds tells AWS CloudFormation how long to wait before reinvoking the handler.
		result.CallbackDelaySeconds = 60
		result.Message = "Create in progress.... "
		// We return an InProgress state
		result.OperationStatus = handler.InProgress
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

		// The AWS CloudFormation service requires that an empty array
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
