# Step 2: Typed Node (Python)

A minimal workflow with one root node (`ingest_orders`) that returns the same mock order data as step 1, now typed. This node takes no trigger payload and is invoked manually.

## What's New

Compared with [step 1](https://github.com/dagflows/examples/tree/main/01-first-node/python):

- `Order` and `Orders` are dataclasses, and `ingest_orders` returns `Orders`.
- `ctx` is annotated with `Ctx` from `dagflows.runtime`.
- The manifest now carries the node's output schema. The data the node returns is unchanged.

## Files

```
workflow.py        workflow definition, types and node handler
requirements.txt   pinned Dagflows SDK dependency
```

Dagflows automatically discovers the project from `requirements.txt` and loads `workflow.py` (or `app/workflow.py`).

## Code

[workflow.py](workflow.py) contains the complete workflow implementation.

## Run Locally

```bash
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -r requirements.txt

# 1. Build the workflow manifest and print the node's output schema
python -m dagflows build manifest workflow -o dagflows-manifest.json
jq '.nodes[0].io.output.schema' dagflows-manifest.json

# 2. Generate a local test fixture and invoke the node
python -m dagflows dev fixture workflow:ingest_orders -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/output.json python -m dagflows invoke --node ingest_orders
jq . /tmp/output.json
```

- `build manifest`: Reflects the `-> Orders` annotation into the node's output schema:

```json
{
  "type": "object",
  "title": "Orders",
  "properties": {
    "orders": {
      "type": "array",
      "items": {
        "type": "object",
        "title": "Order",
        "properties": {
          "id": { "type": "integer", "format": "int64" },
          "customer_id": { "type": "string" },
          "amount_cents": { "type": "integer", "format": "int64" },
          "currency": { "type": "string" }
        },
        "required": ["id", "customer_id", "amount_cents", "currency"]
      }
    }
  },
  "required": ["orders"]
}
```

- `dev fixture` and `invoke`: Run the node as in step 1. It logs the same line and returns the same JSON, encoded from the dataclasses:

```
INFO run local: ingested 5 orders
```

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Once execution completes, the run log will show:

```
INFO run <workflow_run_id>: ingested 5 orders
```
