package xjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// PrettyPrintJSON pretty prints JSON
func PrettyPrintJSON(body []byte, writer io.Writer) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(writer, prettyJSON.String())
	return nil
}
