# Step 5: Fan-Out and Join

The workflow from [step 4](../04-two-languages) grows a second branch. `ingest_orders` now feeds
two nodes, which start two branches that run side by side, and `build_report` waits for both and
joins them. Six nodes in three languages:

```
                  ingest_orders (Python)
                   /                  \
    validate_orders (Go)        score_risk (TypeScript)
            |                            |
   enrich_customers (TypeScript)   apply_pricing (Go)
                   \                  /
                   build_report (Python)
```

- `enrich_customers` puts each valid order in a tier, gold from 10,000 cents.
- `score_risk` scores every order, and `apply_pricing` discounts by 500 cents each order scored
  below 0.34.
- `build_report` reports the tiers, the discounts and the total after them.

This folder is the whole example. Copy all of it:

```
dagflows.yaml          the workspace: the workflow's name and its three projects
py-orders/             Python: ingest_orders and build_report
go-orders/             Go: validate_orders and apply_pricing
node-orders/           TypeScript: enrich_customers and score_risk
```

## The Code

**Fan-out needs nothing new.** Two nodes name the same parent: `validate_orders` in Go as before,
and `score_risk` in [node-orders/workflow.ts](node-orders/workflow.ts):

```ts
const ingestOrders = wf.externalNode<Orders>("ingest_orders");

export const scoreRisk = wf.node(
  function scoreRisk({ ctx, input }: { ctx: Ctx; input: Orders }): RiskScores {
    const scores = input.orders.map((order) => ({ order_id: order.id, score: (order.id % 10) / 10 }));
    return { scores };
  },
  { key: "score_risk", depends: [ingestOrders] },
);
```

**A join names several parents.** `build_report`, in [py-orders/workflow.py](py-orders/workflow.py),
depends on the last node of each branch. With more than one parent, it takes `inputs` and reads each
parent through its handle:

```python
enriched = wf.external_node("enrich_customers", EnrichedOrders)
priced = wf.external_node("apply_pricing", PricedOrders)


@wf.node(depends=[enriched, priced])
def build_report(inputs: Inputs, ctx: Ctx) -> Report:
    orders = inputs[enriched].value()
    adjustments = inputs[priced].value()
    ...
```

**Integers in TypeScript are `Int`.** A TypeScript `number` may hold a fraction, so it never matches
the integers Python's `int` and Go's `int64` declare. Every whole number that crosses to Python or
Go is declared `Int`, a plain `number` at run time:

```ts
import type { Ctx, Int } from "@dagflows/sdk/runtime";

interface Order {
  id: Int;
  customer_id: string;
  amount_cents: Int;
  currency: string;
}
```

Had `score_risk` declared `order_id: number`, the build would fail with:

```
node "apply_pricing" expects a RiskScores from "score_risk", but "score_risk" produces a RiskScores: scores[].order_id is a number, but integer is expected
```

**[dagflows.yaml](dagflows.yaml)** lists the third project, and names its TypeScript entry:

```yaml
  - path: node-orders
    module: workflow.ts
```

## Deploy It

Follow [Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows) with this whole
folder, or replace your repository's files with it, push, and click **Deploy**. Dagflows builds the
three projects in the same deploy. The run takes no payload.

## What You'll See

The run fans out after `ingest_orders` and joins at `build_report`. The two branches run at the same
time as far as your plan allows, one node at a time on the Free plan, so their lines can come in
either order:

![The run with its two branches](../images/05-run.png)

```
INFO run <run id>: ingested 5 orders
time=<time> level=INFO msg="validated orders" kept=3 rejected=2
INFO enriched 3 orders node=enrich_customers
INFO scored 5 orders node=score_risk
time=<time> level=INFO msg="priced orders" scored=5 discounted=3
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
```

Orders 1001 to 1003 score below 0.34, so each takes the 500 cent discount, and 1002, at 34,000
cents, is the one gold order.

## Try It on Your Machine

Run the new branch one node at a time. `score_risk`, from `node-orders`:

```bash
npm ci
npx dagflows-sdk dev run workflow.ts:scoreRisk --input ingest_orders='{"orders": [{"id": 1001, "customer_id": "c1", "amount_cents": 1250, "currency": "usd"}, {"id": 1004, "customer_id": "c2", "amount_cents": 500, "currency": "eur"}]}'
```

Then `apply_pricing` on its scores, from `go-orders`:

```bash
go run . dev run apply_pricing --input score_risk='{"scores": [{"order_id": 1001, "score": 0.1}, {"order_id": 1004, "score": 0.4}]}'
```

It discounts order 1001, scored 0.1, by 500 cents, and leaves 1004, scored 0.4, as it is.

## Next

[Step 6](../06-trigger) takes the orders out of the code: they arrive with the run, as a typed event.
