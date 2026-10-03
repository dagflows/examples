# Step 5: Fan-Out and Join

The workflow from step 4 grows a second branch in TypeScript and Go, and the two branches meet in `build_report`. Six nodes in three languages, one workspace. The workflow takes no trigger payload and is invoked manually.

```
                  ingest_orders (Python)
                   /                  \
    validate_orders (Go)        score_risk (TypeScript)
            |                            |
   enrich_customers (TypeScript)   apply_pricing (Go)
                   \                  /
                   build_report (Python)
```

### What you'll learn
- How one node's output fans out to several children.
- How a node with several parents joins them with `inputs`.
- How integers cross between TypeScript, Go and Python.

## What's New

Compared with [step 4](https://github.com/dagflows/examples/tree/main/04-two-languages):

- `node-orders`, a TypeScript project: `enrich_customers` puts each validated order in a tier (gold from 10,000 cents), and `score_risk` scores every ingested order.
- `apply_pricing` in Go discounts every order scored below 0.34 by 500 cents.
- `ingest_orders` feeds two children that do not depend on each other, so Dagflows can run the two branches at the same time, up to the workflow's concurrent node limit (10 by default).
- `build_report` has two parents. It takes `inputs` and reads each parent through its typed handle, then reports tiers, discounts and the net total.

## Files

```
dagflows.yaml          workspace: the workflow name, its projects and the TypeScript entry module
py-orders/
  workflow.py          the workflow, ingest_orders and build_report
  requirements.txt     pinned Dagflows SDK dependency
go-orders/
  main.go              validate_orders and apply_pricing
  go.mod               module declaration requiring the Dagflows SDK
node-orders/
  workflow.ts          enrich_customers and score_risk
  package.json         project manifest depending on @dagflows/sdk, with typescript for development
  package-lock.json    pinned dependency tree (required for npm ci)
```

## Code

[py-orders/workflow.py](py-orders/workflow.py), [go-orders/main.go](go-orders/main.go) and [node-orders/workflow.ts](node-orders/workflow.ts) contain the complete workflow implementation.

## Run Locally

```bash
# 1. Install the Python and TypeScript dependencies
cd py-orders
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -r requirements.txt
cd ../node-orders
npm ci

# 2. Run each node on its parents' outputs, branch by branch
cd ../py-orders
python -m dagflows dev fixture workflow:ingest_orders -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node ingest_orders
jq .output.data /tmp/out.json > /tmp/orders.json

cd ../go-orders
go run . dev run validate_orders --input ingest_orders=/tmp/orders.json --json | jq .result.output.data > /tmp/validated.json

cd ../node-orders
npx dagflows-sdk dev run workflow.ts:enrichCustomers --input validate_orders=/tmp/validated.json --json | jq .result.output.data > /tmp/enriched.json
npx dagflows-sdk dev run workflow.ts:scoreRisk --input ingest_orders=/tmp/orders.json --json | jq .result.output.data > /tmp/scores.json

cd ../go-orders
go run . dev run apply_pricing --input score_risk=/tmp/scores.json --json | jq .result.output.data > /tmp/priced.json

# 3. Join the two branches
cd ../py-orders
python -m dagflows dev fixture workflow:build_report --input enrich_customers=/tmp/enriched.json --input apply_pricing=/tmp/priced.json -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node build_report
jq .output.data /tmp/out.json
```

- `--input <parent>=<file>`: Hands a node one parent's output. `build_report` takes two, one per parent, as Dagflows hands it both branches during a run.
- Orders 1001 to 1003 score below 0.34, so each gets the 500 cent discount, and 1002 (34,000 cents) is the one gold order:

```
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
```

```json
{
  "total_orders": 3,
  "rejected": 2,
  "gross_cents": 45249,
  "discount_cents": 1500,
  "net_cents": 43749,
  "by_tier": {
    "standard": 2,
    "gold": 1
  }
}
```

## Integers Across Languages

TypeScript has two numeric types. A `number` is reflected as a JSON number, which never satisfies an integer, so every field that Go's `int64` or Python's `int` reads is declared `bigint`, which is reflected as an int64 integer. If `score_risk` declared `order_id: number`, the deployment would fail with:

```
node "apply_pricing" expects a RiskScores from "score_risk", but "score_risk" produces a RiskScores: scores[].order_id is a number, but integer is expected
```

At run time, an integer within 2^53 arrives in TypeScript as a `number`, whatever its declared type. Comparing it with a `bigint` works either way, as the tier check in `enrich_customers` does. Before bigint arithmetic, convert it with `BigInt()`, as `score_risk` does.

## Run on Dagflows

Push this directory to your repository, with `dagflows.yaml` at its root, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Dagflows builds the three projects, runs both branches from `ingest_orders`, and runs `build_report` once both have finished. Once execution completes, the run logs will show (the two branches may run in either order):

```
INFO run <workflow_run_id>: ingested 5 orders
time=<timestamp> level=INFO msg="validated orders" kept=3 rejected=2
INFO enriched 3 orders node=enrich_customers
INFO scored 5 orders node=score_risk
time=<timestamp> level=INFO msg="priced orders" scored=5 discounted=3
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
```
