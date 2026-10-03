# Step 6: Trigger

The workflow from step 5, started by a trigger. The orders are no longer in the code: they arrive as the run's payload, Dagflows checks them against `Orders` before the run starts, and every node can read what started the run.

```
               orders_received (trigger)
                          |
                  ingest_orders (Python)
                   /                  \
    validate_orders (Go)        score_risk (TypeScript)
            |                            |
   enrich_customers (TypeScript)   apply_pricing (Go)
                   \                  /
                   build_report (Python)
```

### What you'll learn
- How to declare a trigger and the type of its event.
- How a node receives the event as its input.
- How a node reads `ctx.trigger`, the delivery that started the run.
- How to start a run with a payload.

## What's New

Compared with [step 5](https://github.com/dagflows/examples/tree/main/05-fan-out-and-join):

- `orders_received`, a trigger whose event is an `Orders`. Dagflows checks every payload against it and refuses one that does not match, so no node ever sees a malformed batch.
- `ingest_orders` depends on the trigger and takes the event as its input. The five orders move out of the code into `orders.json`.
- `ingest_orders` logs `ctx.trigger`: how the run was started, the event's id and when it was received.
- The other five nodes are unchanged.

## Files

```
dagflows.yaml          workspace: the workflow name, its projects and the TypeScript entry module
orders.json            the five orders, sent as the run's payload
py-orders/
  workflow.py          the workflow, the orders_received trigger, ingest_orders and build_report
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

[py-orders/workflow.py](py-orders/workflow.py) declares the trigger and `ingest_orders`. [go-orders/main.go](go-orders/main.go) and [node-orders/workflow.ts](node-orders/workflow.ts) are the same as in step 5.

## Run Locally

```bash
# 1. Install the Python and TypeScript dependencies
cd py-orders
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -r requirements.txt
cd ../node-orders
npm ci

# 2. Hand ingest_orders the event, with the delivery a manual run would describe
cd ../py-orders
python -m dagflows dev fixture workflow:ingest_orders --input orders_received=../orders.json -o /tmp/fixture.json
jq '.ctx.trigger = {kind: "manual", id: "local-event", received_at: "2026-10-03T09:00:00Z"}' /tmp/fixture.json > /tmp/triggered.json
DAGFLOWS_INPUT=/tmp/triggered.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node ingest_orders
jq .output.data /tmp/out.json > /tmp/orders.json

# 3. Run each node on its parents' outputs, branch by branch
cd ../go-orders
go run . dev run validate_orders --input ingest_orders=/tmp/orders.json --json | jq .result.output.data > /tmp/validated.json

cd ../node-orders
npx dagflows-sdk dev run workflow.ts:enrichCustomers --input validate_orders=/tmp/validated.json --json | jq .result.output.data > /tmp/enriched.json
npx dagflows-sdk dev run workflow.ts:scoreRisk --input ingest_orders=/tmp/orders.json --json | jq .result.output.data > /tmp/scores.json

cd ../go-orders
go run . dev run apply_pricing --input score_risk=/tmp/scores.json --json | jq .result.output.data > /tmp/priced.json

# 4. Join the two branches
cd ../py-orders
python -m dagflows dev fixture workflow:build_report --input enrich_customers=/tmp/enriched.json --input apply_pricing=/tmp/priced.json -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node build_report
jq .output.data /tmp/out.json
```

- `--input orders_received=../orders.json`: Hands `ingest_orders` the event, the way a parent's output is handed to a node.
- `jq '.ctx.trigger = ...'`: A fixture describes no delivery, so `ctx.trigger` is `None` locally. Setting it shows what a triggered run logs:

```
INFO manual event local-event received at 2026-10-03T09:00:00Z
INFO run local: ingested 5 orders
```

The report is the same as in step 5.

## Reading the Trigger

Every node of a triggered run sees the same trigger, in any language:

| | Python | Go | TypeScript |
|---|---|---|---|
| The trigger | `ctx.trigger`, `None` when no trigger started the run | `ctx.Trigger()`, `nil` when no trigger started the run | `ctx.trigger`, `undefined` when no trigger started the run |
| Its fields | `kind`, `id`, `received_at`, `attributes` | `Kind`, `ID`, `ReceivedAt`, `Attributes` | `kind`, `id`, `receivedAt`, `attributes` |

- `kind`: `manual` for a run started as shown below, `webhook` for a webhook, `cron` for a schedule and `api` for an event emitted through the events API.
- `id`: The event's id. Each run of this workflow starts from one event.
- `received_at`: When Dagflows received the event, in RFC 3339.
- `attributes`: Empty for a manual run. [Step 7](https://github.com/dagflows/examples/tree/main/07-webhook-and-schedule) shows what a webhook and a schedule add.

## Run on Dagflows

Push this directory to your repository, with `dagflows.yaml` at its root, then deploy it as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The run takes `orders.json` as its payload, so start it with this command instead of the guide's:

```bash
RUN=$(jq '{payload: .}' orders.json | curl -sS -X POST "$API/api/v1/workflows/$WORKFLOW/runs" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" -H "Content-Type: application/json" \
  --data-binary @- | jq -r .data.id)

curl -sS "$API/api/v1/workflows/$WORKFLOW/runs/$RUN" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" | jq .data.trigger
```

- `{payload: ...}`: The event. With one trigger declared, the run goes to it. A workflow with several names one with `"trigger_key"`.
- A payload that does not match `Orders` is refused with `422` and the run does not start, so `RUN` is `null`. Without a payload, the event is `{}`, which has no `orders` and is refused the same way.

The run names the trigger and the event that started it:

```json
{
  "kind": "manual",
  "key": "orders_received",
  "event_id": "<event_id>",
  "delivery_id": "<delivery_id>"
}
```

Read the logs as the guide shows. `ingest_orders` logs the same event id:

```
INFO manual event <event_id> received at <received_at>
INFO run <workflow_run_id>: ingested 5 orders
time=<timestamp> level=INFO msg="validated orders" kept=3 rejected=2
INFO enriched 3 orders node=enrich_customers
INFO scored 5 orders node=score_risk
time=<timestamp> level=INFO msg="priced orders" scored=5 discounted=3
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
```
