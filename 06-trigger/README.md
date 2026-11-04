# Step 6: Trigger

The workflow from [step 5](../05-fan-out-and-join), with the orders taken out of the code. They now
arrive with the run, as an event of a trigger the code declares. Dagflows checks every event
against the trigger's type before the run starts, so no node ever sees a malformed batch.

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

This folder is the whole example. Copy all of it:

```
dagflows.yaml          the workspace: the workflow's name and its three projects
orders.json            the five orders, the event a run starts with
py-orders/             Python: the trigger, ingest_orders and build_report
go-orders/             Go: validate_orders and apply_pricing
node-orders/           TypeScript: enrich_customers and score_risk
```

## The Code

Only `ingest_orders` changes, in [py-orders/workflow.py](py-orders/workflow.py):

```python
orders_received = wf.trigger("orders_received", event=Orders)


@wf.node(depends=[orders_received])
def ingest_orders(batch: Orders, ctx: Ctx) -> Orders:
    if ctx.trigger:
        ctx.log.info("%s event %s received at %s", ctx.trigger.kind, ctx.trigger.id, ctx.trigger.received_at)

    ctx.log.info("run %s: ingested %d orders", ctx.run.workflow_run_id, len(batch.orders))

    return batch
```

- **`wf.trigger`** declares a trigger, `orders_received`, and the type of its event, `Orders`.
- **A node that depends on a trigger** receives the event as its input, the way a node receives
  its parent's output.
- **`ctx.trigger`** describes what started the run, in every node and every language: its `kind`,
  the event's `id`, when it was `received_at`, and its `attributes`. It is empty when no trigger
  started the run. In Go it is `ctx.Trigger()`, in TypeScript `ctx.trigger`.
- The other five nodes are the same as in step 5.

## Deploy It

Follow [Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows) with this whole
folder, or replace your repository's files with it, push, and click **Deploy**.

This time the run takes the orders. On the **Runs** tab, click **Run**. The **Payload** box opens
with an example built from `Orders`. Replace it with the contents of [orders.json](orders.json) and
click **Run**:

![The Run workflow dialog with the orders as the payload](../images/06-run-payload.png)

A payload that does not match `Orders`, say an order without its `amount_cents`, is refused and no
run starts.

The **Triggers** tab now lists `orders_received`. [Step 7](../07-webhook) connects it to a webhook.

## What You'll See

The run is marked as started manually, and `ingest_orders` logs the event it was handed:

![The run started by the orders_received event](../images/06-run.png)

```
INFO manual event <event id> received at <time>
INFO run <run id>: ingested 5 orders
time=<time> level=INFO msg="validated orders" kept=3 rejected=2
INFO enriched 3 orders node=enrich_customers
INFO scored 5 orders node=score_risk
time=<time> level=INFO msg="priced orders" scored=5 discounted=3
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
```

The report is the same as in step 5, from the same five orders, which now came with the run.

## Try It on Your Machine

Locally, the event is handed to `ingest_orders` like a parent's output, by the trigger's key. From
`py-orders`:

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
python -m dagflows dev run workflow:ingest_orders --input orders_received=../orders.json
```

It logs `INFO run local: ingested 5 orders`. Nothing started the run locally, so `ctx.trigger` is
empty and the first line is not logged.

## Next

[Step 7](../07-webhook) starts the same workflow from outside, through a webhook, without changing
its code.
