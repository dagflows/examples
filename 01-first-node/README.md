# Step 1: Your First Node

A workflow with one node, `ingest_orders`, which returns five orders and logs how many it found.
It is the smallest thing Dagflows can run.

Pick a language. Each folder is a complete project:

- [python](python): `workflow.py` and `requirements.txt`
- [go](go): `main.go` and `go.mod`
- [node](node): `workflow.js` and `package.json`

## The Code

**Python**, [workflow.py](python/workflow.py):

```python
from dagflows.authoring import Workflow

wf = Workflow("order-pipeline")


@wf.node()
def ingest_orders(ctx):
    orders = [
        {"id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd"},
        # ... four more orders
    ]

    ctx.log.info("run %s: ingested %d orders", ctx.run.workflow_run_id, len(orders))

    return {"orders": orders}
```

**Go**, [main.go](go/main.go):

```go
func ingestOrders(ctx *df.Ctx, _ df.None) (map[string]any, error) {
	orders := []map[string]any{
		{"id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd"},
		// ... four more orders
	}

	ctx.Log().Info("ingested orders", "run", ctx.Run().WorkflowRunID, "count", len(orders))

	return map[string]any{"orders": orders}, nil
}

func main() {
	wf := df.NewWorkflow("order-pipeline", df.WorkflowOptions{})
	wf.Node(ingestOrders, df.Root, df.NodeOptions{Key: "ingest_orders"})
	df.Main()
}
```

**Node**, [workflow.js](node/workflow.js):

```js
export const wf = new Workflow("order-pipeline");

export const ingestOrders = wf.node(
  function ingestOrders({ ctx }) {
    const orders = [
      { id: 1001, customer_id: "cus_alpha", amount_cents: 1250, currency: "usd" },
      // ... four more orders
    ];

    ctx.log.info(`run ${ctx.run.workflowRunId}: ingested ${orders.length} orders`);

    return { orders };
  },
  { key: "ingest_orders" },
);
```

In every language:

- **The workflow** is named `order-pipeline`, and every step keeps that name.
- **A node** is a function registered on the workflow. Its key, `ingest_orders`, is its name in
  runs and logs.
- **`ctx`** is the node's view of the run: `ctx.log` writes to the run's logs, and `ctx.run` says
  which run it is.
- **What a node returns** is its output. From step 3 on, Dagflows hands it to the next node.
- Dagflows recognizes the project by its `requirements.txt`, `go.mod` or `package.json`.

And per language:

- **Go**: the first node of a workflow is a root, `df.Root`, and takes `df.None` since nothing
  comes before it. `df.Main()` lets Dagflows read the workflow and run its nodes, so `main` ends
  with it.
- **Node**: a node is exported, which is how Dagflows finds it.

## Deploy It

Follow [Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows) with the folder
of your language. The run takes no payload.

## What You'll See

The run has one node, and its log shows the orders it ingested:

![The run with ingest_orders and its log](../images/run-detail.png)

```
Python   INFO run <run id>: ingested 5 orders
Go       time=<time> level=INFO msg="ingested orders" run=<run id> count=5
Node     INFO run <run id>: ingested 5 orders node=ingest_orders
```

## Try It on Your Machine

You can run the node on your own machine, without Dagflows.

**Python**, 3.11 or newer:

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
python -m dagflows dev fixture workflow:ingest_orders -o fixture.json
DAGFLOWS_INPUT=fixture.json DAGFLOWS_OUTPUT=out.json python -m dagflows invoke --node ingest_orders
```

`dev fixture` writes the input Dagflows would hand the node, and `invoke` runs the node on it. It
logs `INFO run local: ingested 5 orders` and writes the output to `out.json`.

**Go**, 1.27 or newer:

```bash
go run . dev run ingest_orders
```

**Node**, 22.18 or newer:

```bash
npm ci
npx dagflows-sdk dev run workflow.js:ingestOrders
```

Go and Node print the log line and the output.

## Next

[Step 2](../02-typed-node) gives the orders a type, so Dagflows knows what the node returns.
