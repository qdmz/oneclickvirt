package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// 测试数据库连接 API
	testDBAPI()
}

func testDBAPI() {
	// 等待服务启动
	time.Sleep(2 * time.Second)

	url := "http://localhost:8890/api/v1/public/test-db-connection"

	payload := map[string]interface{}{
		"type":     "mysql",
		"host":     "localhost",
		"port":     "3306",
		"database": "test",
		"username": "root",
		"password": "",
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
}
