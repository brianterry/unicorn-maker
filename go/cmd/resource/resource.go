package resource

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/encoding"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
)

const APIEndpoint = "https://crudcrud.com/api/64e82541a53643fdac7fdb7576feff19/unicorns"

type unicorn struct {
	ID    string `json:"_id,omitempty"`
	Name  string `json:"name,omitempty"`
	Age   int64  `json:"age,omitempty"`
	Color string `json:"color,omitempty"`
}

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Create complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: Create")
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	resp, err := http.Get(APIEndpoint + "/" + *currentModel.UID.Value())

	if err != nil {
		return handler.ProgressEvent{}, err
	}
	u := unicorn{}

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return handler.ProgressEvent{}, err
	}
	unmarshal(&u)

	// Construct a new handler.ProgressEvent and return it
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read complete",
		ResourceModel:   &u,
	}

	return response, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Update complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: Update")
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Delete complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: Delete")
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "List complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: List")
}

func unmarshal(unicorn *unicorn) *Model {
	m := Model{}

	m.UID = encoding.NewString(unicorn.ID)
	m.Name = encoding.NewString(unicorn.Name)
	m.Age = encoding.NewInt(unicorn.Age)
	m.Color = encoding.NewString(unicorn.Color)
	return &m
}
