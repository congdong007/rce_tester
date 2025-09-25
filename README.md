# rce_fuzzer / rce_tester

## 1️⃣ Project Overview

`rce_fuzzer` is a **multi-threaded HTTP Remote Command Execution (RCE) testing tool**. Its main features include:

- Support **single URL testing** or **batch URL testing**
- Supports **POST data templates** with automatic `FUZZ` placeholder replacement
- Supports **payload files**, testing requests one by one
- Supports **keyword matching** to quickly locate potential vulnerable responses
- Supports **multi-threading**, processing payloads in parallel per URL
- Supports **batch URL concurrency**, testing multiple URLs simultaneously
- Supports **custom HTTP headers**, retaining User-Agent, Cookie, Referer, etc.

> Suitable for local testing environments (e.g., DVWA, WebGoat) to quickly identify RCE vulnerabilities.

------

## 2️⃣ Directory Structure

```
rce_fuzzer/
├── cmd/
│   └── main.go           # Entry point
├── config/
│   └── config.go         # Configuration and flag parsing
├── core/
│   └── core.go           # Core logic for single or batch URL testing
├── util/
│   ├── loader.go         # File loading utilities (payloads, keywords, URLs)
│   └── writer.go         # Result writing utilities
├── worker/
│   └── worker.go         # Multi-threaded request execution
├── go.mod
└── README.md             # Project documentation
```

------

## 3️⃣ Installation & Dependencies

1. Install Go 1.21+
2. Clone the repository:

```
git clone https://github.com/yourusername/rce_fuzzer.git
cd rce_fuzzer
```

1. Build the project:

```
go build -o rce_fuzzer ./cmd
```

Or run directly:

```
go run ./cmd
```

------

## 4️⃣ Command-line Flags

| Flag            | Description                                                  |
| --------------- | ------------------------------------------------------------ |
| `-u`            | Single URL for testing, e.g., `http://dvwa/vulnerabilities/exec/?ip=FUZZ` |
| `-uf`           | File with multiple target URLs, each must contain `FUZZ` placeholder |
| `-t`            | Number of concurrent threads per URL (default 10)            |
| `-max-uf`       | Maximum number of URLs tested concurrently (default = CPU cores) |
| `-pf`           | Payload file path, one payload per line                      |
| `-kf`           | Keyword file path, one keyword per line                      |
| `-o`            | Output file (default `results.txt`)                          |
| `-timeout`      | HTTP request timeout in seconds (default 10)                 |
| `-content-type` | HTTP Content-Type header (default `application/x-www-form-urlencoded`) |

------

## 5️⃣ Usage

### Single URL Testing

```
./rce_fuzzer -u "http://dvwa/vulnerabilities/exec/?ip=FUZZ" -pf payloads.txt -kf keywords.txt -t 10 -o results.txt
```

- `FUZZ` will be replaced with each payload
- Keywords are used to detect potential vulnerabilities
- Results are saved to `results.txt`

------

### Batch URL Testing

```
./rce_fuzzer -uf urls.txt -pf payloads.txt -kf keywords.txt -t 10 -max-uf 8 -o results.txt
```

- Each line in `urls.txt` must contain a `FUZZ` placeholder:

```
http://dvwa/vulnerabilities/exec/?ip=FUZZ&Submit=Submit&user_token=xxxx
http://example.com/test.php?cmd=FUZZ
```

- `-max-uf 8` → Max 8 URLs tested concurrently
- Each URL internally uses `-t` threads for payload testing

------

## 6️⃣ Example Files

### payloads.txt

```
whoami
id
ls
```

### keywords.txt

```
desktop-
root
uid=
```

### urls.txt

```
http://dvwa/vulnerabilities/exec/?ip=FUZZ&Submit=Submit&user_token=xxxx
http://example.com/test.php?cmd=FUZZ
```

------

## 7️⃣ Sample Output

```
[HIT] http://dvwa/vulnerabilities/exec/ -> payload: whoami
[RESULT] http://example.com/test.php -> payload: id
```

- `[HIT]` → response contains keyword, potentially vulnerable
- `[RESULT]` → normal response without keyword match

------

## 8️⃣ Multi-threading & Concurrency

- `-t` → Threads per URL (payload processing)
- `-max-uf` → Concurrent batch URL count
- Default `-max-uf` = number of CPU cores
- Adjust according to system performance for faster scanning

------

## 9️⃣ HTTP Headers

Default headers simulate a real browser:

```
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...
Content-Type: application/x-www-form-urlencoded
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8
Referer: <URL>
Connection: keep-alive
```

- Custom headers (Cookie, Origin, Host) can be added in `worker/worker.go`

------

## 10️⃣ Notes

1. **Use in local testing environments only** (DVWA, WebGoat, etc.)
2. Ensure payloads are safe and legal; avoid attacking real servers
3. Customize `keywords.txt` for more accurate detection
