# 🛠️ Go OTLP Trace Demo

A fully self-contained Go example that produces a rich, multi-step trace and exports it via OTLP/HTTP to a local OpenTelemetry Collector.

---

## 🚀 Requirements

- Go 1.18+
- Docker (or a local OTEL Collector setup)
- Internet access (for the external HTTP call in the demo)

---

## ⚙️ Quick Start (with Docker Collector)

1. **Clone this repo**  
   ```bash
   git clone https://your-url/otel-go-demo.git
   cd otel-go-demo

2. **Install dependencies**
    ```bash
    go mod tidy

3. Add your collector url in [otel-go-demo/go.sum, l:23](https://github.com/IrisDyr/traces-playground/blob/752d31295319f470e8b027af2a7924ddfa28c0d6/otel-go-demo/main.go#L23)

4. Run the go app
    ```bash
    go run main.go   
