# OpenTelemetry Connector Test

A test suite for OpenTelemetry Collector's **example connector** that demonstrates traces-to-metrics conversion functionality.

## Overview

This repository contains test code to validate the behavior of an OpenTelemetry Collector custom connector. The connector is designed to:

- **Consume traces** from OTLP receiver
- **Filter spans** based on specific attributes
- **Generate metrics** when target attributes are found
- **Export metrics** through debug exporter

## What is the Example Connector?

The example connector is a **traces-to-metrics** connector that:

1. Monitors incoming trace data for spans containing a specific attribute (`request.n` by default)
2. When a span with the target attribute is detected, it generates a metric
3. Passes the generated metric to the next component in the pipeline (debug exporter)

### Connector Behavior

```
Trace with "request.n" attribute → Connector detects → Generates metric → Debug output
Trace without target attribute   → Connector ignores → No metric generated
```

## Prerequisites

- **Go 1.21+**
- **OpenTelemetry Collector** with example connector built-in
- **Git** for cloning the repository

## Installation

1. Clone this repository:

```bash
git clone https://github.com/jaehanbyun/connector-test.git
cd connector-test
```

2. Install dependencies:

```bash
go mod tidy
```

## Usage

### Step 1: Start OpenTelemetry Collector

First, you need to have an OpenTelemetry Collector running with the example connector configured. The collector should be configured with:

- **OTLP Receiver** (listening on port 4318 for HTTP)
- **Example Connector** (with `attribute_name: "request.n"`)
- **Debug Exporter** (for viewing generated metrics)

Example collector configuration:

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318

connectors:
  example:
    attribute_name: "request.n"

exporters:
  debug:

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [example]
    metrics:
      receivers: [example]
      exporters: [debug]
```

### Step 2: Run the Test

Execute the test script to send trace data to the collector:

```bash
go run test_connector.go
```

## Test Scenarios

The test script sends two types of traces to validate connector behavior:

### Test Case 1: Trigger Connector

- **Span Name**: `test-operation-with-trigger`
- **Attributes**:
  - `request.n: "some-value"` ← **This triggers the connector**
  - `http.method: "GET"`
  - `http.url: "/api/test"`
- **Expected Result**: Connector detects the `request.n` attribute and generates a metric

### Test Case 2: No Trigger

- **Span Name**: `test-operation-without-trigger`
- **Attributes**:
  - `http.method: "POST"`
  - `http.url: "/api/other"`
  - `user.id: "12345"`
- **Expected Result**: Connector ignores this span (no `request.n` attribute)

## Expected Output

### Test Script Output

```
🧪 Test 1: Sending trace with 'request.n' attribute...
🧪 Test 2: Sending trace without 'request.n' attribute...
✅ Test completed! Please check the Collector logs.
```

### Collector Logs

When the test runs successfully, you should see:

1. **Trace Reception**: OTLP receiver logs showing incoming trace data
2. **Connector Activation**: Connector building logs
3. **Metric Generation**: Debug exporter showing metric output (only for Test Case 1)

Example collector log output:

```
2025-07-15T22:19:07.123+0900    info    exampleconnector@v0.129.0/connector.go:26    Building exampleconnector connector
2025-07-15T22:19:07.456+0900    info    Metrics {"resource metrics": 0, "metrics": 0, "data points": 0}
```

> **Note**: The current example connector generates empty metrics as per the official OpenTelemetry tutorial. The appearance of metric logs indicates successful detection and processing.

## Understanding the Results

### ✅ Success Indicators

- Test script completes without errors
- Collector shows metric generation logs for Test Case 1
- No metric logs appear for Test Case 2

### ❌ Troubleshooting

- **No metric logs**: Check if collector is running and configured correctly
- **Connection errors**: Verify collector is listening on `localhost:4318`
- **Import errors**: Run `go mod tidy` to resolve dependencies

## Architecture

```
[Test Script] → [OTLP Receiver] → [Example Connector] → [Debug Exporter]
     |                                    |                     |
  Generates                          Filters by              Outputs
  trace data                      "request.n" attr           metrics
```

## Customization

### Change Target Attribute

To test with a different attribute, modify the collector configuration:

```yaml
connectors:
  example:
    attribute_name: "your.custom.attribute" # Change this
```

Then update the test script to include your custom attribute in the span.

### Add More Test Cases

You can extend `test_connector.go` to include additional test scenarios:

```go
// Add more spans with different attributes
span3.SetAttributes(
    attribute.String("your.custom.attribute", "test-value"),
    // ... other attributes
)
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Resources

- [OpenTelemetry Collector Documentation](https://opentelemetry.io/docs/collector/)
- [Building Custom Connectors](https://opentelemetry.io/docs/collector/building/connector/)
- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/)
