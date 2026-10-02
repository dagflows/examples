# Step 1: Your First Node (Go)

A minimal workflow with one root node (`ingest_orders`) that returns mock order data. This node is untyped, takes no trigger payload, and is invoked manually.

## Files

```
main.go   the workflow and its node handler
go.mod    module declaration requiring the Dagflows SDK
```

Dagflows automatically discovers the project by finding `go.mod` at the root and builds the module where `package main` lives.

## Code

[main.go](main.go) contains the complete workflow implementation.

## Run Locally

```bash
# 1. Build the workflow manifest
go run . build manifest -o dagflows-manifest.json

# 2. Execute the node locally
go run . dev run ingest_orders
```

- `build manifest`: Generates `dagflows-manifest.json`, the execution graph Dagflows uses to schedule nodes. (Dagflows builds this automatically during deployment, but running it locally lets you inspect your node metadata).
- `dev run`: Executes the handler directly on your machine without a microVM or server, printing execution status, structured logs, and JSON output:

```
time=2026-10-02T15:26:51.839Z level=INFO msg="ingested orders" run=local count=5

ingest_orders -> SUCCESS
  {
    "orders": [
      {
        "amount_cents": 1250,
        "currency": "usd",
        "customer_id": "cus_alpha",
        "id": 1001
      },
      ...
    ]
  }
```

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Once execution completes, the node's structured log will show:

```
time=<timestamp> level=INFO msg="ingested orders" run=<workflow_run_id> count=5
```
