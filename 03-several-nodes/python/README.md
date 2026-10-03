# Step 3: Several Nodes (Python)

A workflow of three typed nodes in a chain: `ingest_orders` → `validate_orders` → `build_report`. Each node receives its parent's output, decoded into its declared type. The workflow takes no trigger payload and is invoked manually.

## What's New

Compared with [step 2](https://github.com/dagflows/examples/tree/main/02-typed-node/python):

- `validate_orders` keeps the USD orders worth more than zero, normalizes their currency, and counts the rest as rejected.
- `build_report` totals the validated orders into a `Report`.
- `depends=[...]` names each node's parent. With one parent, the parameter that is not `ctx` receives the parent's output, decoded into its annotation.
- Each node's input is declared in the manifest, and Dagflows checks it against the parent's output when it builds the workflow.

## Files

```
workflow.py        workflow definition, types and node handlers
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

# 1. Build the workflow manifest and print each node's parents
python -m dagflows build manifest workflow -o dagflows-manifest.json
jq '.nodes[] | {key, depends}' dagflows-manifest.json

# 2. Run each node on its parent's output
python -m dagflows dev fixture workflow:ingest_orders -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node ingest_orders
jq .output.data /tmp/out.json > /tmp/orders.json

python -m dagflows dev fixture workflow:validate_orders --input ingest_orders=/tmp/orders.json -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node validate_orders
jq .output.data /tmp/out.json > /tmp/validated.json

python -m dagflows dev fixture workflow:build_report --input validate_orders=/tmp/validated.json -o /tmp/fixture.json
DAGFLOWS_INPUT=/tmp/fixture.json DAGFLOWS_OUTPUT=/tmp/out.json python -m dagflows invoke --node build_report
jq .output.data /tmp/out.json
```

- `build manifest`: Lists each node with its parents: `ingest_orders` has none, `validate_orders` depends on `ingest_orders`, `build_report` on `validate_orders`.
- `dev fixture --input <parent>=<file>`: Hands the node its parent's output, keyed by the parent, as Dagflows does during a run.
- `invoke`: Runs each node in turn. Order 1004 is in euros and 1005 is worth nothing, so two are rejected:

```
INFO run local: ingested 5 orders
INFO kept 3 orders, rejected 2
INFO report: 3 orders, 2 rejected, 45249 cents gross
```

The report:

```json
{
  "total_orders": 3,
  "rejected": 2,
  "gross_cents": 45249
}
```

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Dagflows runs the three nodes in order, passing each output to the next. Once execution completes, the run logs will show:

```
INFO run <workflow_run_id>: ingested 5 orders
INFO kept 3 orders, rejected 2
INFO report: 3 orders, 2 rejected, 45249 cents gross
```
