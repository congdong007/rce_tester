package worker

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	cfgpkg "rce_tester/config"
)

type Task struct {
	URL      string
	Data     string
	Payload  string
	Headers  map[string]string
	Keywords []string
}

func RunWorkersWithData(cfg *cfgpkg.Config, payloads []string, keywords []string, dataTemplate string) (chan string, error) {
	results := make(chan string, len(payloads))
	tasks := make(chan Task, len(payloads))

	var wg sync.WaitGroup

	// default headers
	defaultHeaders := map[string]string{
		"Host":                      cfg.Host,
		"Content-Type":              cfg.ContentT,
		"Content-Length":            "84",
		"Cache-Control":             "max-age=0",
		"Accept-Language":           "zh-CN,zh;q=0.9",
		"Origin":                    cfg.HostUrl,
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		"Upgrade-Insecure-Requests": "1",
		"Referer":                   cfg.URL,
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Connection":                "keep-alive",
		"Cookie":                    cfg.CookieT, //"session=b60715bf40a4404192ad1e54a80acb83; security=low; PHPSESSID=j90tj5i17agigrlsrev44df8c6",
	}

	for i := 0; i < cfg.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				sendRequest(task, results)
			}
		}()
	}

	for _, payload := range payloads {
		data := strings.ReplaceAll(dataTemplate, "FUZZ", payload)
		tasks <- Task{
			URL:      cfg.URL,
			Data:     data,
			Payload:  payload,
			Headers:  defaultHeaders,
			Keywords: keywords,
		}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil
}

func sendRequest(task Task, results chan string) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", task.URL, bytes.NewBufferString(task.Data))
	if err != nil {
		return
	}

	// 设置 headers
	for k, v := range task.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	bodyStr := string(body)
	if len(task.Keywords) > 0 {
		for _, k := range task.Keywords {
			if strings.Contains(bodyStr, k) {
				results <- fmt.Sprintf("[Vulnerable] %s -> payload: %s", task.URL, task.Payload)
				fmt.Printf("\033[31m[ =-= Vulnerable =-= ] payload=%s status=%d\033[0m\n", task.Payload, resp.StatusCode)
				return
			}
		}
	} else {
		results <- fmt.Sprintf("[RESULT] %s -> payload: %s", task.URL, task.Payload)
	}
}
