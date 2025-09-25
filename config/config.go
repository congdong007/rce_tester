package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

// Config holds runtime configuration
type Config struct {
	URL        string
	Uf         string
	Host       string
	HostUrl    string
	Ip         string
	Port       int
	Data       string
	Threads    int
	Output     string
	Payloads   string
	Keywords   string
	TimeoutSec int
	ContentT   string
	CookieT    string
	MaxUf      int
}

// ParseFlags parses CLI flags into Config
func ParseFlags() *Config {
	url := flag.String("u", "", "Target URL, e.g. http://dvwa/vulnerabilities/exec/")
	uf := flag.String("uf", "", "File with multiple target URLs containing FUZZ placeholder, one per line")
	data := flag.String("d", "", "POST data template; use FUZZ as placeholder to be replaced by payload")
	threads := flag.Int("t", 10, "Number of concurrent threads")
	maxUf := flag.Int("max-uf", 0, "Max number of URL concurrently tested in batch mode (default = CPU cores)")
	output := flag.String("o", "results.txt", "Output file to record hits/possible results")
	payloadFile := flag.String("pf", "", "Payload file path, one payload per line")
	keywordFile := flag.String("kf", "", "Keywords file path, one keyword per line; match any to consider a hit")
	timeoutS := flag.Int("timeout", 10, "HTTP client timeout (seconds)")
	contentT := flag.String("content-type", "application/x-www-form-urlencoded", "Content-Type header")
	cookieT := flag.String("cookie", "application/x-www-form-urlencoded", "Content-Type header")

	flag.Parse()

	if *payloadFile == "" || *keywordFile == "" || *url == "" && *uf == "" {
		fmt.Println("Usage: rce_tester -u <url> -d <data template> -pf <payload file> [-kf <keywords file>] [-t <threads>] [-o <output>]")
		os.Exit(1)
	}

	if *maxUf <= 0 {
		*maxUf = runtime.NumCPU()
	}

	return &Config{
		URL:        *url,
		Uf:         *uf,
		Data:       *data,
		Threads:    *threads,
		Output:     *output,
		Payloads:   *payloadFile,
		Keywords:   *keywordFile,
		TimeoutSec: *timeoutS,
		ContentT:   *contentT,
		CookieT:    *cookieT,
		MaxUf:      *maxUf,
	}
}
