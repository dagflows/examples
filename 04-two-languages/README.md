# Step 4: Two Languages

The chain from [step 3](../03-several-nodes), split across two languages: `ingest_orders` and
`build_report` stay in Python, and `validate_orders` moves to Go. Dagflows builds both projects and
runs them as one workflow, handing each node its parent's output whatever language produced it.

```
ingest_orders (Python) → validate_orders (Go) → build_report (Python)
```

This folder is the whole example. Copy all of it:

```
dagflows.yaml          the workspace: the workflow's name and its two projects
py-orders/             Python: ingest_orders and build_report
go-orders/             Go: validate_orders
```

## The Code

**[dagflows.yaml](dagflows.yaml)** makes the repository a workspace, several projects built as one
workflow:

```yaml
v: 1

workflow: order-pipeline

projects:
  - path: py-orders
  - path: go-orders
```

**Python**, [py-orders/workflow.py](py-orders/workflow.py), owns the workflow, so it names it. Its
`build_report` depends on a node of the Go project, named by its key:

```python
wf = Workflow("order-pipeline")


@wf.node(depends=[wf.external_node("validate_orders")])
def build_report(validated: ValidatedOrders, ctx: Ctx) -> Report:
    ...
```

**Go**, [go-orders/main.go](go-orders/main.go), adds a node to that workflow, so its workflow has no
name. Its `validate_orders` depends on `ingest_orders` in the Python project:

```go
wf := df.NewWorkflow("", df.WorkflowOptions{})

wf.Node(validateOrders, wf.External[Orders]("ingest_orders"), df.NodeOptions{Key: "validate_orders"})
```

- **One project names the workflow**, and every other one adds its nodes without a name.
- **A node of another project is named by its key**: `wf.external_node("validate_orders")` in
  Python, `wf.External[Orders]("ingest_orders")` in Go. Keys are unique across the workspace.
- **Each language declares the types it reads.** Python and Go both declare `Orders` and
  `ValidatedOrders`. When Dagflows builds the workflow, it checks that both sides of every edge
  agree. Had the Go `Order` declared `AmountCents string`, the build would fail with:

  ```
  node "build_report" expects a ValidatedOrders from "validate_orders", but "validate_orders" produces a ValidatedOrders: orders[].amount_cents is a string, but integer is expected
  ```

## Deploy It

Follow [Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows) with this whole
folder, so `dagflows.yaml` is at the root of the repository. Dagflows builds the Python and the Go
project in the same deploy. The run takes no payload.

## What You'll See

The same three nodes as in step 3, the middle one in Go:

![The run with its Python and Go nodes](../images/04-run.png)

```
INFO run <run id>: ingested 5 orders
time=<time> level=INFO msg="validated orders" kept=3 rejected=2
INFO report: 3 orders, 2 rejected, 45249 cents gross
```

## Try It on Your Machine

Each project runs on its own, and you hand its node what the other project's node would. The Go
node, from `go-orders`:

```bash
go run . dev run validate_orders --input ingest_orders='{"orders": [{"id": 1, "customer_id": "c1", "amount_cents": 1250, "currency": "usd"}, {"id": 2, "customer_id": "c2", "amount_cents": 500, "currency": "eur"}]}'
```

Then the Python node on what Go returned, from `py-orders`:

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
python -m dagflows dev fixture workflow:build_report -o fixture.json --input validate_orders='{"orders": [{"id": 1, "customer_id": "c1", "amount_cents": 1250, "currency": "usd"}], "rejected": 1}'
DAGFLOWS_INPUT=fixture.json DAGFLOWS_OUTPUT=out.json python -m dagflows invoke --node build_report
```

It logs `INFO report: 1 orders, 1 rejected, 1250 cents gross`.

## Next

[Step 5](../05-fan-out-and-join) adds a second branch in TypeScript and Go, which runs beside the
first and meets it in `build_report`.
