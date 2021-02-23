package report

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func PrintResultAsJSON(violations interface{}) {
	resultJSON, err := json.Marshal(violations)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(resultJSON)) // nolint:forbidigo
}
