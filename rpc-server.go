package pcommon

import (
	"encoding/json"
	"reflect"
)

type rpcServer struct{}

var RPCServer = rpcServer{}

func (rpc rpcServer) HandleRPCRequest(a []byte, service interface{}) RPCResponse {

	//check if the action is a valid action
	action := RPCAction{}
	if err := json.Unmarshal(a, &action); err != nil {
		return RPCResponse{Data: nil, Error: err.Error(), Id: "not found"}
	}

	if action.Method == "" {
		return RPCResponse{Data: nil, Error: "method not found", Id: action.Id}
	}

	// Obtain the reflection Value of the interface
	val := reflect.ValueOf(service)

	// Get the method by name
	method := val.MethodByName(action.Method)
	if !method.IsValid() {
		return RPCResponse{Data: nil, Error: "method not found", Id: action.Id}
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

	return RPCResponse{Data: ret[0].Interface(), Error: errStr, Id: action.Id}
}
