package marionette_client

import (
	"encoding/json"
	"log"
)

func makeProto2Response(buf []byte) (*response, error) {
	r := &response{}
	r.Size = int32(len(buf))
	r.Value = string(buf)

	return r, nil
}

func makeProto3Response(buf []byte) (*response, error) {
	//Debug only
	if RunningInDebugMode {
		if (len(buf) >= 512){
			log.Println(string(buf)[0:512] + " - END - " + string(buf)[len(buf) - 512:])
		} else {
			log.Println(string(buf))
		}
	}
	//Debug only end

	var v []interface{}

	if err := json.Unmarshal(buf, &v); err != nil {
		return nil, err
	}

	r := &response{}
	r.MessageID = int32(v[1].(float64))
	r.Size = int32(len(buf))

	if v[2] != nil {
		re := &DriverError{}
		// JSON Object?
		for key, value := range v[2].(map[string]interface{}) {
			if key == "error" {
				re.ErrorType = value.(string)
			}

			if key == "message" {
				re.Message = value.(string)
			}

			if key == "stacktrace" && value != nil {
				theTrace := value.(string)
				re.Stacktrace = &theTrace
			}
		}

		return nil, re
	}

	// It's a JSON Object
	result, found := v[3].(map[string]interface{})
	if found {
		resultBytes, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		r.Value = string(resultBytes)
	}

	if !found {
		// JSON Array
		result, found := v[3].([]interface{})
		if found {
			resultBytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}

			r.Value = string(resultBytes)
		}

		//TODO: return error?
	}

	return r, nil
}
