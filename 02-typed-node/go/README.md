# Step 2: Typed Node (Go)

A minimal workflow with one root node (`ingest_orders`) that returns the same mock order data as step 1, now typed. This node takes no trigger payload and is invoked manually.

## What's New

Compared with [step 1](https://github.com/dagflows/examples/tree/main/01-first-node/go):

- `Order` and `Orders` are structs with `json` tags, and `ingestOrders` returns `Orders` instead of `map[string]any`.
- The manifest now carries the node's output schema. The data the node returns is unchanged.

## Files

```
main.go   workflow definition, types and node handler
go.mod    module declaration requiring the Dagflows SDK
```

Dagflows automatically discovers the project by finding `go.mod` at the root and builds the module where `package main` lives.

## Code

[main.go](main.go) contains the complete workflow implementation.

## Run Locally

```bash
# 1. Build the workflow manifest and print the node's output schema
go run . build manifest -o dagflows-manifest.json
jq '.nodes[0].io.output.schema' dagflows-manifest.json

# 2. Execute the node locally
go run . dev run ingest_orders
```

- `build manifest`: Reflects the `Orders` struct into the node's output schema:

```json
{
  "type": "object",
  "title": "Orders",
  "properties": {
    "orders": {
      "type": "array",
      "items": {
        "type": "object",
        "title": "Order",
        "properties": {
          "id": { "type": "integer", "format": "int64" },
          "customer_id": { "type": "string" },
          "amount_cents": { "type": "integer", "format": "int64" },
          "currency": { "type": "string" }
        },
        "required": ["id", "customer_id", "amount_cents", "currency"]
      }
    }
  },
  "required": ["orders"]
}
```

- `dev run`: Executes the handler as in step 1, logging the same line and returning the same JSON:

```
time=2026-10-03T11:54:19.380Z level=INFO msg="ingested orders" run=local count=5

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
