package message

import (
	"io"
	"net/http"
)

var (
	forceBase64File = false
)

// SetForceBase64File read or download a file and convert it to base64:// before sending.
func SetForceBase64File(x bool) {
	forceBase64File = x
}

func dl(file string) ([]byte, error) {
	resp, err := http.Get(file)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
