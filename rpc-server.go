package pcommon

import (
	"encoding/json"
	"reflect"
)

func HandleAction(a []byte, service interface{}) Response {

	//check if the action is a valid action
	action := Action{}
	if err := json.Unmarshal(a, &action); err != nil {
		return Response{Data: nil, Error: err.Error(), Id: "not found"}
	}

	if action.Method == "" {
		return Response{Data: nil, Error: "method not found", Id: action.Id}
	}

	// Obtain the reflection Value of the interface
	val := reflect.ValueOf(service)

	// Get the method by name
	method := val.MethodByName(action.Method)
	if !method.IsValid() {
		return Response{Data: nil, Error: "method not found", Id: action.Id}
	}

	// Prepare input arguments for reflection call
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(action.Payload)

	ret := method.Call(in)

	var errStr = ""
	// Check and convert the error value if it's not nil
	if errInter := ret[1].Interface(); errInter != nil {
		if err, ok := errInter.(error); ok {
			errStr = err.Error()
		} else {
			errStr = "error asserting type"
		}
	} else {
		errStr = ""
	}

	return Response{Data: ret[0].Interface(), Error: errStr, Id: action.Id}
}
