# Step 1: your first node (Python)

A workflow with one node, `ingest_orders`. It returns a batch of five orders and logs how
many there are. It has no types and no trigger: you run it by hand.

## The files

```
workflow.py        the workflow and its one node
requirements.txt   the Dagflows SDK, pinned
```

Dagflows finds the workflow without any configuration: a `requirements.txt` naming
`dagflows` makes this a Python project, and `workflow.py` is where the workflow is looked
for (`app/workflow.py` works too).

## The code

```python
wf = Workflow("order-pipeline")


@wf.node()
def ingest_orders(ctx):
    orders = [
        {"id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd"},
        ...
    ]

    ctx.log.info("run %s: ingested %d orders", ctx.run.workflow_run_id, len(orders))
    return {"orders": orders}
```

- `Workflow("order-pipeline")` names the workflow. Every step of these examples builds on it.
- `@wf.node()` makes the function a node. Its key is the function's name, `ingest_orders`.
- `ctx` is the node's context. `ctx.log` writes to the run's logs, and
  `ctx.run.workflow_run_id` is the id of the run.
- What the function returns is the node's output, plain JSON here.

## Run it on your machine

```bash
python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
python -m dagflows build manifest workflow -o dagflows-manifest.json
python -m dagflows dev fixture workflow:ingest_orders -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/output.json python -m dagflows invoke --node ingest_orders
jq . /tmp/output.json
```

`build manifest` writes what Dagflows reads from your code, one node called
`ingest_orders`. Dagflows builds the manifest itself when it deploys, so it is not
committed. `dev fixture` writes the input a node receives, and `invoke` runs the node on
it the way Dagflows does. The node logs:

```
INFO run local: ingested 5 orders
```

and `/tmp/output.json` holds its output:

```json
{
  "status": "SUCCESS",
  "output": {
    "type": "INLINE",
    "content_type": "application/json",
    "data": {
      "orders": [
        { "id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd" },
        ...
      ]
    }
  }
}
```

## Run it on Dagflows

Copy this directory into a repository of your own, then deploy and run it as
[Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows)
shows. The run takes no body. When it has finished, the node's log includes this line:

```
INFO run <the run's id>: ingested 5 orders
```
