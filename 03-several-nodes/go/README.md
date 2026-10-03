# Step 3: Several Nodes (Go)

A workflow of three typed nodes in a chain: `ingest_orders` → `validate_orders` → `build_report`. Each node receives its parent's output, decoded into its declared type. The workflow takes no trigger payload and is invoked manually.

## What's New

Compared with [step 2](https://github.com/dagflows/examples/tree/main/02-typed-node/go):

- `validateOrders` keeps the USD orders worth more than zero, normalizes their currency, and counts the rest as rejected.
- `buildReport` totals the validated orders into a `Report`.
- `wf.Node` returns the node's handle. Passing a handle as another node's edge makes it that node's parent, and the compiler checks the parent's output type is the child's input type.
- Each node's input is declared in the manifest, and Dagflows checks it against the parent's output when it builds the workflow.

## Files

```
main.go   workflow definition, types and node handlers
```

Dagflows automatically discovers the project by finding `go.mod` at the root and builds the module where `package main` lives.

## Code

[main.go](main.go) contains the complete workflow implementation.

## Run Locally

```bash
# 1. Build the workflow manifest and print each node's parents
go run . build manifest -o dagflows-manifest.json
jq '.nodes[] | {key, depends}' dagflows-manifest.json

# 2. Run each node on its parent's output
go run . dev run ingest_orders --json | jq .result.output.data > /tmp/orders.json
go run . dev run validate_orders --input ingest_orders=/tmp/orders.json --json | jq .result.output.data > /tmp/validated.json
go run . dev run build_report --input validate_orders=/tmp/validated.json
```

- `build manifest`: Lists each node with its parents: `ingest_orders` has none, `validate_orders` depends on `ingest_orders`, `build_report` on `validate_orders`.
- `dev run --json`: Prints the result as JSON, with the node's log in a `logs` field, so `jq` can pass its output to the next node.
- `dev run --input <parent>=<file>`: Hands the node its parent's output, keyed by the parent, as Dagflows does during a run. Order 1004 is in euros and 1005 is worth nothing, so two are rejected and the report reads:

```
time=2026-10-03T12:39:24.183Z level=INFO msg="built report" orders=3 rejected=2 gross_cents=45249

build_report -> SUCCESS
  {
    "gross_cents": 45249,
    "rejected": 2,
    "total_orders": 3
  }
```

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Dagflows runs the three nodes in order, passing each output to the next. Once execution completes, the run logs will show:

```
time=<timestamp> level=INFO msg="ingested orders" run=<workflow_run_id> count=5
time=<timestamp> level=INFO msg="validated orders" kept=3 rejected=2
time=<timestamp> level=INFO msg="built report" orders=3 rejected=2 gross_cents=45249
```
