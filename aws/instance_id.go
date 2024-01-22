package aws

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetInstanceId it work in aws ec2
func GetInstanceId() (string, error) {
	client := http.Client{Timeout: time.Second}
	res, err := client.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		b, err := io.ReadAll(res.Body)
		return string(b), err
	} else {
		return "", fmt.Errorf("StatusCode != 200, response: %d", res.StatusCode)
	}

}
