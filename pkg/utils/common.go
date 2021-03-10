package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func SendDingMsgWithUrl(robot, content string, IsAtAll bool) error {
	if robot == "" {
		return errors.New("robot url miss")
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content},
		"at":      map[string]bool{"isAtAll": IsAtAll},
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(robot, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
