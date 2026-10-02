# Order Pipeline

A polyglot Dagflows workflow executed across three languages as a single directed acyclic graph (DAG).

Six nodes coordinate processing: Python handles order ingestion and final reporting, Go executes validation and pricing, and TypeScript manages customer enrichment and risk scoring. Every edge in the workflow traverses a language boundary. A run starts from the `orders_received` trigger, whose event is the batch of orders to process.

```
                    orders_received (trigger)
                               |
                               v
                    ingest_orders (Python)
                      /                \
                     v                  v
       validate_orders (Go)        score_risk (TypeScript)
             |                           |
             v                           v
   enrich_customers (TypeScript)   apply_pricing (Go)
             \                          /
              \                        /
               v                      v
                    build_report (Python)
```

---

## Workspace Layout

```
dagflows.yaml        Workspace configuration registering project paths
orders.json          A sample batch of orders, the event to start a run with
py-orders/           Python package: the orders_received trigger, ingest_orders, build_report
  app/workflow.py
  requirements.txt
go-orders/           Go module: validate_orders, apply_pricing
  main.go
  go.mod
node-orders/         TypeScript package: enrich_customers, score_risk
  app/workflow.ts
  package.json
```

Each subdirectory is a standard language project depending on the respective Dagflows SDK. The build engine installs project dependencies and extracts JSON Schema node manifests automatically.

---

## Core Rules for Multi-Project Workflows

> [!IMPORTANT]
> **Workflow Ownership and Naming**
> Exactly one project in the workspace must define global workflow settings by declaring a named workflow (such as `Workflow("order-pipeline")` in `py-orders`). All other projects must instantiate unnamed workflows (`NewWorkflow("")` in Go or `new Workflow()` in TypeScript). The `workflow:` name declared in [dagflows.yaml](file:///p:/go%20projects/dagflows/examples/order-pipeline/dagflows.yaml) must match the owner project workflow name.

> [!WARNING]
> **Flat Workspace Namespace**
> Node keys form a flat, workspace-wide namespace. Duplicate node keys across different projects are rejected during cross-project plan resolution at build time. Ensure every node key in the workspace is unique.

> [!NOTE]
> **Triggers Live Beside the Nodes They Start**
> A node can depend only on a trigger its own project declares, so `py-orders` declares `orders_received` with `wf.trigger(...)` beside `ingest_orders`. Trigger keys share the workspace-wide namespace with node keys.

### External Dependencies

Dependencies residing in peer projects are linked by string key and typed schema:
* **Python**: `@wf.node(depends=[wf.external_node("enrich_customers")])`
* **Go**: `wf.Node(applyPricing, wf.External[RiskScores]("score_risk"), ...)`
* **TypeScript**: `const validateOrders = wf.externalNode<ValidatedOrders>("validate_orders")`

The builder validates all external node references against peer manifests before deployment.

---

## Cross-Language Type Safety

Each SDK reflects handler data structures into JSON Schema manifests. Dagflows validates input and output schema compatibility across every edge at build time.

> [!IMPORTANT]
> **Numeric Type Compatibility: `bigint` vs `number`**
> In JSON Schema, an integer satisfies a number, but a number **never** satisfies an integer.
>
> In TypeScript, `number` maps to a JSON floating-point number. To emit an integer compatible with Go (`int64`) and Python (`int`), fields must use `bigint`:
> ```ts
> interface Order {
>   id: bigint // Emitted as int64 integer, matching Go int64 and Python int
>   amount_cents: number
> }
> ```

---

## Pipeline Execution Summary

1. `ingest_orders` (Python) receives the batch the run was started with, through the `orders_received` trigger, and hands it to both branches. The figures below are for [orders.json](orders.json), a batch of 5 orders.
2. `validate_orders` (Go) drops non-USD and non-positive orders (3 survive: 1004 is in euros and 1005 is worth nothing).
3. `enrich_customers` (TypeScript) classifies surviving orders into customer tiers.
4. `score_risk` (TypeScript) computes risk scores directly from ingested orders in parallel.
5. `apply_pricing` (Go) computes discount deductions based on risk scores.
6. `build_report` (Python) merges enriched orders and pricing adjustments by order ID to produce the final report:

```json
{
  "total_orders": 3,
  "gross_cents": 45249,
  "discount_cents": 1500,
  "net_cents": 43749,
  "by_tier": { "standard": 2, "gold": 1 }
}
```

---

## Running the Workflow

A manual run starts from the trigger, and the request's `payload` is the event. To run [orders.json](orders.json):

```bash
jq '{payload: .}' orders.json |
  curl -sS -X POST "$DAGFLOWS_API/api/v1/workflows/$WORKFLOW_ID/runs" \
    -H "Authorization: Bearer $TOKEN" \
    -H "X-Organization-Id: $ORG_ID" \
    -H "Content-Type: application/json" \
    --data-binary @-
```

* The workflow declares one trigger, so a run takes it without naming it. `"trigger_key": "orders_received"` beside `payload` names it explicitly.
* The payload is checked against `Orders` before the run starts. A run without one, or with an order missing a field or holding the wrong type, is refused with `422`.
* The `201` response is the run, and its `trigger` gives the kind (`manual`), the key (`orders_received`), the event and the delivery.

---

## Local Development & Manifest Inspection

You can test nodes locally and inspect emitted manifests using the CLI. A node takes its parent's output, or the event of the trigger it depends on, keyed by that parent's key:

```bash
# Go
cd go-orders
go run . build manifest -o /tmp/manifest.json
go run . dev run validate_orders --input ingest_orders=../orders.json

# TypeScript
cd ../node-orders
npm install
npx tsc --noEmit
./node_modules/.bin/dagflows-sdk build manifest app/workflow.ts -o /tmp/manifest.json
./node_modules/.bin/dagflows-sdk dev run app/workflow.ts:scoreRisk --input ingest_orders=../orders.json

# Python
cd ../py-orders
pip install -r requirements.txt
python -m dagflows build manifest app.workflow -o /tmp/manifest.json
python -m dagflows dev fixture app.workflow:ingest_orders --input orders_received=../orders.json -o /tmp/m.json
DAGFLOWS_INPUT=/tmp/m.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node ingest_orders
```

> [!TIP]
> Running `build manifest` locally validates node schemas and external dependencies without deploying to the platform.
