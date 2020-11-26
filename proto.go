package marionette_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func NewDecoderEncoder(protoVersion int32) (DecoderEncoder, error) {
	if protoVersion == MARIONETTE_PROTOCOL_V3 {
		return ProtoV3DecoderEncoder{}, nil
	}

	m := fmt.Sprintf("No DecoderEncoder found for the specified version : %v", protoVersion)

	return nil, errors.New(m)
}

type Decoder interface {
	Decode(buf []byte, r *Response) error
}

type Encoder interface {
	Encode(t Transporter, command string, values interface{}) ([]byte, error)
}

type DecoderEncoder interface {
	Decoder
	Encoder
}

type ProtoV3DecoderEncoder struct{}

func (e ProtoV3DecoderEncoder) Encode(t Transporter, command string, values interface{}) ([]byte, error) {
	message := make([]interface{}, 4)
	message[0] = 0
	message[1] = t.MessageID()
	message[2] = command
	message[3] = values

	b, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	if RunningInDebugMode {
		fmt.Println(string(b))
	}

	return []byte(strconv.Itoa(len(b)) + ":" + string(b)), nil

}

func (e ProtoV3DecoderEncoder) Decode(buf []byte, r *Response) error {
	var v []interface{}
	if err := json.Unmarshal(buf, &v); err != nil {
		return err
	}

	//Debug only
	if RunningInDebugMode {
		if len(buf) >= 512 {
			log.Println(string(buf)[0:256] + " - END - " + string(buf)[len(buf)-256:])
		} else {
			log.Println(string(buf))
		}
	}
	//Debug only end

	r.MessageID = int32(v[1].(float64))
	r.Size = int32(len(buf))

	// error found on response?
	if v[2] != nil {
		re := &DriverError{}
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

		return re
	}

	// It's a JSON Object
	result, found := v[3].(map[string]interface{})
	if found {
		resultBytes, err := json.Marshal(result)
		if err != nil {
			return err
		}

		r.Value = string(resultBytes)
	}

	if !found {
		// JSON Array
		result, found := v[3].([]interface{})
		if found {
			resultBytes, err := json.Marshal(result)
			if err != nil {
				return err
			}

			r.Value = string(resultBytes)
		}

		//TODO: return error?
	}

	return nil
}
