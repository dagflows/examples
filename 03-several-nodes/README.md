# Step 3: Several Nodes

The typed node from [step 2](../02-typed-node), followed by two more: `validate_orders` keeps the
orders it can price and rejects the rest, and `build_report` totals what is left. Each node gets the
output of the node before it.

```
ingest_orders → validate_orders → build_report
```

Pick a language. Each folder is a complete project:

- [python](python): `workflow.py` and `requirements.txt`
- [go](go): `main.go` and `go.mod`
- [node](node), in TypeScript: `workflow.ts`, `dagflows.yaml` and `package.json`

## The Code

A node names its parent, and receives the parent's output as its input, typed.

**Python**, [workflow.py](python/workflow.py). `depends` names the parent, and the parameter that
is not `ctx` receives its output:

```python
@wf.node(depends=[ingest_orders])
def validate_orders(batch: Orders, ctx: Ctx) -> ValidatedOrders:
    valid = [
        replace(order, currency="usd")
        for order in batch.orders
        if order.currency.lower() == "usd" and order.amount_cents > 0
    ]
    return ValidatedOrders(orders=valid, rejected=len(batch.orders) - len(valid))


@wf.node(depends=[validate_orders])
def build_report(validated: ValidatedOrders, ctx: Ctx) -> Report:
    ...
```

**Go**, [main.go](go/main.go). A node's input type is its parent's output type, and registering a
node returns a handle that the next node takes as its parent:

```go
func validateOrders(ctx *df.Ctx, in Orders) (ValidatedOrders, error) { ... }

func buildReport(ctx *df.Ctx, in ValidatedOrders) (Report, error) { ... }

func main() {
	wf := df.NewWorkflow("order-pipeline", df.WorkflowOptions{})

	ingested := wf.Node(ingestOrders, df.Root, df.NodeOptions{Key: "ingest_orders"})
	validated := wf.Node(validateOrders, ingested, df.NodeOptions{Key: "validate_orders"})
	wf.Node(buildReport, validated, df.NodeOptions{Key: "build_report"})

	df.Main()
}
```

**TypeScript**, [workflow.ts](node/workflow.ts). `depends` names the parent, and `input` is its
output:

```ts
export const validateOrders = wf.node(
  function validateOrders({ ctx, input }: { ctx: Ctx; input: Orders }): ValidatedOrders {
    const orders = input.orders
      .filter((order) => order.currency.toLowerCase() === "usd" && order.amount_cents > 0)
      .map((order) => ({ ...order, currency: "usd" }));

    return { orders, rejected: input.orders.length - orders.length };
  },
  { key: "validate_orders", depends: [ingestOrders] },
);
```

- **Dagflows runs the nodes in order**, each once its parent has finished, and hands each its
  parent's output.
- **The types must agree.** `validate_orders` returns a `ValidatedOrders`, and `build_report` takes
  one. When Dagflows builds the workflow, it checks every such pair, and a mismatch fails the build
  with the field that differs, before anything runs. In Go, the compiler checks it too.
- Of the five orders, three are in US dollars with an amount, so `validate_orders` keeps three and
  rejects two.

## Deploy It

Follow [Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows) with the folder
of your language, or replace your repository's files with this step's, push, and click **Deploy**.
The run takes no payload.

## What You'll See

The run shows the three nodes, one after the other:

![The run with its three nodes](../images/03-run.png)

```
Python   INFO run <run id>: ingested 5 orders
         INFO kept 3 orders, rejected 2
         INFO report: 3 orders, 2 rejected, 45249 cents gross

Go       time=<time> level=INFO msg="ingested orders" run=<run id> count=5
         time=<time> level=INFO msg="validated orders" kept=3 rejected=2
         time=<time> level=INFO msg="built report" orders=3 rejected=2 gross_cents=45249

Node     INFO run <run id>: ingested 5 orders node=ingest_orders
         INFO kept 3 orders, rejected 2 node=validate_orders
         INFO report: 3 orders, 2 rejected, 45249 cents gross node=build_report
```

## Try It on Your Machine

Locally, you hand a node its parent's output yourself. This runs `validate_orders` on two orders,
one it keeps and one in euros it rejects:

**Python**:

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
python -m dagflows dev run workflow:validate_orders --input ingest_orders='{"orders": [{"id": 1, "customer_id": "c1", "amount_cents": 1250, "currency": "usd"}, {"id": 2, "customer_id": "c2", "amount_cents": 500, "currency": "eur"}]}'
```

**Go**:

```bash
go run . dev run validate_orders --input ingest_orders='{"orders": [{"id": 1, "customer_id": "c1", "amount_cents": 1250, "currency": "usd"}, {"id": 2, "customer_id": "c2", "amount_cents": 500, "currency": "eur"}]}'
```

**TypeScript**:

```bash
npm ci
npx dagflows-sdk dev run workflow.ts:validateOrders --input ingest_orders='{"orders": [{"id": 1, "customer_id": "c1", "amount_cents": 1250, "currency": "usd"}, {"id": 2, "customer_id": "c2", "amount_cents": 500, "currency": "eur"}]}'
```

`--input ingest_orders=...` is the parent's output, by the parent's key. Each logs that it kept 1
order and rejected 1.

## Next

[Step 4](../04-two-languages) moves `validate_orders` to Go, so one workflow spans two languages.
