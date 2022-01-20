package checkptclient

import (
	"encoding/json"
	"fmt"
	"log"
)

type ErrHandler struct{}

func (eh ErrHandler) Handle(code int, data []byte) error {

	if code == 200 {
		return nil
	}
	//no special handling based on body
	if len(data) < 1 {
		return nil
	}
	e := ErrResponse{}
	json.Unmarshal([]byte(data), &e)

	f := e.Message

	log.Printf("error handler: %+v", e)
	if len(e.Errors) > 0 || len(e.Blocking) > 0 {

		for _, ef := range e.Errors {
			f = fmt.Sprintf("%s%s\n", f, ef.Message)
		}
		for _, ef := range e.Blocking {
			f = fmt.Sprintf("%s%s\n", f, ef.Message)
		}
	}
	if f != "" {
		return fmt.Errorf(f)
	}
	return nil
	// log.Printf("Error [%v], Message [%s] --> %d", e.Error, e.Message, len(e.Error))

	// if len(e.Error) > 0 {
	// 	return fmt.Errorf("%d - %s : %s", code, e.Error, e.Message)
	// }
	// return nil
}
