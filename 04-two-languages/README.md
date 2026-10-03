# Step 4: Two Languages

The chain from step 3 split across two languages: `ingest_orders` and `build_report` stay in Python, and `validate_orders` moves to Go. A `dagflows.yaml` workspace builds both projects as one workflow. The workflow takes no trigger payload and is invoked manually.

```
ingest_orders (Python) → validate_orders (Go) → build_report (Python)
```

### What you'll learn
- How a `dagflows.yaml` workspace builds several projects as one workflow.
- How a node depends on a node in another project, by its key.
- How Dagflows checks the types on every edge between languages when it builds the workflow.

## What's New

Compared with step 3 in [Python](https://github.com/dagflows/examples/tree/main/03-several-nodes/python) and [Go](https://github.com/dagflows/examples/tree/main/03-several-nodes/go):

- `dagflows.yaml` lists the two projects and names the workflow. Dagflows builds and deploys them together.
- `py-orders` owns the workflow with a named `Workflow`, and `go-orders` contributes its node through an unnamed one.
- A node in another project is named by its key: `wf.external_node("validate_orders")` in Python, `wf.External[Orders]("ingest_orders")` in Go. Node keys are unique across the workspace.
- Each language declares its own types for an edge, and Dagflows checks that they agree (see [Types Across Languages](#types-across-languages)).

## Files

```
dagflows.yaml          workspace: the workflow name and its projects
py-orders/
  workflow.py          the workflow, ingest_orders and build_report
  requirements.txt     pinned Dagflows SDK dependency
go-orders/
  main.go              validate_orders
  go.mod               module declaration requiring the Dagflows SDK
```

## Code

[py-orders/workflow.py](py-orders/workflow.py) and [go-orders/main.go](go-orders/main.go) contain the complete workflow implementation.

## Run Locally

```bash
# 1. Build each project's manifest and print the nodes it needs from the other project
cd py-orders
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -r requirements.txt
python -m dagflows build manifest workflow -o dagflows-manifest.json
jq '.nodes[] | {key, external_depends}' dagflows-manifest.json

cd ../go-orders
go run . build manifest -o dagflows-manifest.json
jq '.nodes[] | {key, external_depends}' dagflows-manifest.json

# 2. Run each node on its parent's output, across the two languages
cd ../py-orders
python -m dagflows dev fixture workflow:ingest_orders -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node ingest_orders
jq .output.data /tmp/out.json > /tmp/orders.json

cd ../go-orders
go run . dev run validate_orders --input ingest_orders=/tmp/orders.json --json | jq .result.output.data > /tmp/validated.json

cd ../py-orders
python -m dagflows dev fixture workflow:build_report --input validate_orders=/tmp/validated.json -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node build_report
jq .output.data /tmp/out.json
```

- `build manifest`: Each project's manifest lists the node it needs from the other project under `external_depends`: `build_report` needs `validate_orders`, and `validate_orders` needs `ingest_orders`.
- `--input <parent>=<file>`: Hands each node its parent's output, whichever language produced it. The report is the same as in step 3:

```
INFO run local: ingested 5 orders
INFO report: 3 orders, 2 rejected, 45249 cents gross
```

```json
{
  "total_orders": 3,
  "rejected": 2,
  "gross_cents": 45249
}
```

## Types Across Languages

Python and Go each declare the types of the edges they read. When Dagflows builds the workflow it compares every node's input schema with its parent's output schema, across projects and languages. If the Go `Order` struct declared `AmountCents string`, the deployment would fail with:

```
node "build_report" expects a ValidatedOrders from "validate_orders", but "validate_orders" produces a ValidatedOrders: orders[].amount_cents is a string, but integer is expected
```

## Run on Dagflows

Push this directory to your repository, with `dagflows.yaml` at its root, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Dagflows builds both projects, then runs the three nodes in order, passing each output to the next. Once execution completes, the run logs will show:

```
INFO run <workflow_run_id>: ingested 5 orders
time=<timestamp> level=INFO msg="validated orders" kept=3 rejected=2
INFO report: 3 orders, 2 rejected, 45249 cents gross
```
