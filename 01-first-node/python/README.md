# Step 1: Your First Node (Python)

A minimal workflow with one root node (`ingest_orders`) that returns mock order data. This node is untyped, takes no trigger payload, and is invoked manually.

## Files

```
workflow.py        workflow definition and node handler
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

# 1. Build the workflow manifest
python -m dagflows build manifest workflow -o dagflows-manifest.json

# 2. Generate a local test fixture and invoke the node
python -m dagflows dev fixture workflow:ingest_orders -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/output.json python -m dagflows invoke --node ingest_orders
jq . /tmp/output.json
```

- `build manifest`: Inspects the workflow module and generates `dagflows-manifest.json` describing `ingest_orders`. (Dagflows generates this automatically during deployment, but running it locally verifies discovery).
- `dev fixture`: Generates a mock input payload matching what Dagflows supplies to root nodes.
- `invoke`: Runs the node handler against the fixture file, outputting structured logs and writing the result:

```
INFO run local: ingested 5 orders
```

Output in `/tmp/output.json`:

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

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Once execution completes, the run log will show:

```
INFO run <workflow_run_id>: ingested 5 orders
```
