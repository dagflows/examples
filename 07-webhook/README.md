# Step 7: Webhook

The workflow from step 6, started by a webhook. The code does not change: on Dagflows, the `orders_received` trigger is bound to a webhook, and anything that can send a signed HTTP request starts a run by posting orders to its URL.

```
     POST https://hooks.dagflows.com/hooks/<token>
                          |
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
- How to bind a trigger to a webhook, outside the code.
- How to sign a delivery.
- How to follow a delivery to the run it started.
- What a webhook adds to `ctx.trigger`.

## What's New

Compared with [step 6](https://github.com/dagflows/examples/tree/main/06-trigger):

- No code changes. The code declares the trigger and its type, Dagflows binds it to a source, so the same deployment can listen to a webhook without a new build.
- Binding a webhook creates a source with a URL and a signing secret.
- A delivery is checked against `Orders` like a manual run's payload, recorded as an event, and starts the run on its own.

## Files

```
dagflows.yaml          workspace: the workflow name, its projects and the TypeScript entry module
orders.json            the five orders, posted to the webhook
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

The same as in [step 6](https://github.com/dagflows/examples/tree/main/06-trigger#code).

## Run Locally

As in [step 6](https://github.com/dagflows/examples/tree/main/06-trigger#run-locally). A webhook is bound on Dagflows, so locally the fixture's `ctx.trigger` stands in for it.

## Run on Dagflows

Push this directory to your repository, with `dagflows.yaml` at its root, then deploy it as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows). Signing needs `openssl` besides the guide's tools.

```bash
# 1. Bind a webhook to the trigger
BINDING=$(curl -sS -X POST "$API/api/v1/workflows/$WORKFLOW/triggers" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" -H "Content-Type: application/json" \
  -d '{"trigger_key": "orders_received", "name": "Orders webhook", "webhook": {}}')
HOOK=$(printf '%s' "$BINDING" | jq -r .data.source.url)
SECRET=$(printf '%s' "$BINDING" | jq -r .data.secret)

# 2. Sign the orders and post them
TS=$(date +%s)
SIGNATURE=$( (printf '%s.' "$TS"; cat orders.json) | openssl dgst -sha256 -hmac "$SECRET" | awk '{print $NF}')
EVENT=$(curl -sS -X POST "$HOOK" \
  -H "Content-Type: application/json" -H "X-Dagflows-Signature: t=$TS,v1=$SIGNATURE" \
  --data-binary @orders.json | jq -r .data.event_id)

# 3. Follow the event to its run
curl -sS "$API/api/v1/events/$EVENT" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" \
  | jq '.data | {source_kind, attributes, deliveries: [.deliveries[] | {status, run_id}]}'
```

- `"webhook": {}`: Creates a webhook source named after the trigger key and binds the trigger to it. The answer carries the source's URL and its signing secret. The secret is shown only this once, so keep it like a password.
- `X-Dagflows-Signature`: An HMAC-SHA256 of the time and the body joined by a dot, keyed by the secret. A delivery whose time is more than five minutes away is refused, so a captured request cannot be sent again later. `--data-binary` posts the file byte for byte, as it was signed.
- The webhook answers `202` once the event is recorded, and the run starts right after. A wrong signature is answered `401`, and orders that do not match `Orders` are answered `422` with the field that failed.
- A sender that retries can add an `Idempotency-Key` header. A repeat of a key is answered `200` with the first event and starts nothing.

The event names its source, and its delivery the run it started:

```json
{
  "source_kind": "webhook",
  "attributes": {
    "source": "orders_received",
    "source_id": "<source_id>"
  },
  "deliveries": [
    {
      "status": "STARTED",
      "run_id": "<run_id>"
    }
  ]
}
```

While the delivery is `PENDING`, repeat the command. Then keep its run, `RUN=the-run-id`, and read the logs as the guide shows:

```
INFO webhook event <event_id> received at <received_at>
INFO run <workflow_run_id>: ingested 5 orders
time=<timestamp> level=INFO msg="validated orders" kept=3 rejected=2
INFO enriched 3 orders node=enrich_customers
INFO scored 5 orders node=score_risk
time=<timestamp> level=INFO msg="priced orders" scored=5 discounted=3
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
```

## What a Webhook Adds to the Trigger

Every node of the run sees `ctx.trigger.kind` as `webhook`, and `ctx.trigger.attributes` as the event's attributes above: `source`, the source's name, and `source_id`. The source's name reads as:

- Python: `ctx.trigger.attributes["source"]`
- Go: `ctx.Trigger().Attributes["source"]`
- TypeScript: `ctx.trigger?.attributes["source"]`
