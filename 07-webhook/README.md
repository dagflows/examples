# Step 7: Webhook

The workflow from [step 6](../06-trigger), started from outside. The code does not change. On
Dagflows, the `orders_received` trigger is connected to a webhook, and anything that can send a
signed HTTP request, a shop or a payment service say, starts a run by posting the orders to its URL.

```
              POST <the webhook's URL>
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

This folder is the whole example, the same as step 6's. Copy all of it:

```
dagflows.yaml          the workspace: the workflow's name and its three projects
orders.json            the five orders, posted to the webhook
py-orders/             Python: the trigger, ingest_orders and build_report
go-orders/             Go: validate_orders and apply_pricing
node-orders/           TypeScript: enrich_customers and score_risk
```

## The Code

Nothing changes. The code declares what starts a run, the `orders_received` trigger and its
`Orders` event. Where the events come from is set on Dagflows, so the same deployment listens to a
webhook without a new build.

What a webhook adds is in `ctx.trigger`: its `kind` is `webhook`, and its `attributes` name the
webhook, `source` and `source_id`. In Python that is `ctx.trigger.attributes["source"]`, in Go
`ctx.Trigger().Attributes["source"]`, in TypeScript `ctx.trigger?.attributes["source"]`.

## Deploy It

If your workflow runs step 6, there is nothing to deploy. Otherwise, follow
[Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows) with this whole folder.

### Connect the Trigger to a Webhook

On the workflow's **Triggers** tab, click **Connect** on `orders_received`. Under **Listens to**,
choose **Create a webhook for this**, leave **Require a signature** on, and click **Connect**:

![Connecting orders_received to a new webhook](../images/07-connect.png)

Dagflows shows the webhook's URL and its signing secret. The secret is shown only this once, so copy
both before you click **I've saved it**:

![The webhook's URL and signing secret](../images/07-secret.png)

### Send the Orders

**From the dashboard**, to check the connection: in **Event sources**, open `orders_received`, click
**Send test event**, paste the contents of [orders.json](orders.json) and click **Send event**. It
starts a run as an event would, from inside Dagflows, so it needs no signature:

![Sending a test event](../images/07-test-event.png)

**As a service would**: sign the orders with the secret and post them to the URL. From this folder,
with `openssl` and `curl`:

```bash
HOOK='the webhook URL'
SECRET='the signing secret'

TS=$(date +%s)
SIGNATURE=$( (printf '%s.' "$TS"; cat orders.json) | openssl dgst -sha256 -hmac "$SECRET" | awk '{print $NF}')
curl -sS -X POST "$HOOK" \
  -H "Content-Type: application/json" \
  -H "X-Dagflows-Signature: t=$TS,v1=$SIGNATURE" \
  --data-binary @orders.json
```

- **The signature** is an HMAC-SHA256 of the time and the body joined by a dot, keyed by the
  secret. A request whose time is more than five minutes off is refused, so a captured request
  cannot be sent again later.
- **The webhook answers** `202` once it has recorded the event, and the run starts right after. A
  wrong signature is answered `401`, and orders that do not match `Orders` `422`, naming the field
  that failed.

## What You'll See

A run started by the signed request is marked as started by a webhook, and `ingest_orders` logs the
webhook's event. The test event's run is marked as started by an event, and logs `api event`:

![The run started by the webhook](../images/07-run.png)

```
INFO webhook event <event id> received at <time>
INFO run <run id>: ingested 5 orders
time=<time> level=INFO msg="validated orders" kept=3 rejected=2
INFO enriched 3 orders node=enrich_customers
INFO scored 5 orders node=score_risk
time=<time> level=INFO msg="priced orders" scored=5 discounted=3
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
```

The webhook's page in **Event sources** lists every event it received, and the run each one
started.

## Try It on Your Machine

Locally, the code runs as in [step 6](../06-trigger#try-it-on-your-machine). The webhook lives on
Dagflows.

## Next

[Step 8](../08-variables-and-secrets) moves the discount and the gold tier out of the code, and signs
the report with a key that stays out of the repository.
