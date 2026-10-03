# Step 3: Several Nodes (Node)

A workflow of three typed nodes in a chain: `ingest_orders` → `validate_orders` → `build_report`. Each node receives its parent's output as its typed `input`. The workflow takes no trigger payload and is invoked manually.

## What's New

Compared with [step 2](https://github.com/dagflows/examples/tree/main/02-typed-node/node):

- `validateOrders` keeps the USD orders worth more than zero, normalizes their currency, and counts the rest as rejected.
- `buildReport` totals the validated orders into a `Report`.
- `depends: [...]` names each node's parent. With one parent, the handler receives the parent's output as `input`.
- `input` is annotated with the parent's type. That declares the node's input in the manifest, and Dagflows checks it against the parent's output when it builds the workflow. An `input` left to inference is checked by your editor, but not declared.

## Files

```
workflow.ts         workflow definition, types and node handlers
dagflows.yaml       names workflow.ts as the entry module
package.json        project manifest depending on @dagflows/sdk, with typescript for development
package-lock.json   pinned dependency tree (required for npm ci)
```

## Code

[workflow.ts](workflow.ts) contains the complete workflow implementation.

## Run Locally

Node 22.18 or newer runs the TypeScript entry directly.

```bash
npm ci

# 1. Build the workflow manifest and print each node's parents
npx dagflows-sdk build manifest workflow.ts -o dagflows-manifest.json
jq '.nodes[] | {key, depends}' dagflows-manifest.json

# 2. Run each node on its parent's output
npx dagflows-sdk dev run workflow.ts:ingestOrders --json | jq .result.output.data > /tmp/orders.json
npx dagflows-sdk dev run workflow.ts:validateOrders --input ingest_orders=/tmp/orders.json --json | jq .result.output.data > /tmp/validated.json
npx dagflows-sdk dev run workflow.ts:buildReport --input validate_orders=/tmp/validated.json
```

- `build manifest`: Lists each node with its parents: `ingest_orders` has none, `validate_orders` depends on `ingest_orders`, `build_report` on `validate_orders`.
- `dev run --json`: Prints the result as JSON, with the node's log in a `logs` field, so `jq` can pass its output to the next node.
- `dev run --input <parent>=<file>`: Hands the node its parent's output, keyed by the parent's key, as Dagflows does during a run. Order 1004 is in euros and 1005 is worth nothing, so two are rejected and the report reads:

```
INFO report: 3 orders, 2 rejected, 45249 cents gross node=buildReport

workflow.ts:buildReport -> SUCCESS
  {
    "total_orders": 3,
    "rejected": 2,
    "gross_cents": 45249
  }
```

*Note: The local runner identifies the node by its export (`buildReport`). On the Dagflows platform, it is tracked by its registered key (`build_report`).*

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Dagflows runs the three nodes in order, passing each output to the next. Once execution completes, the run logs will show:

```
INFO run <workflow_run_id>: ingested 5 orders node=ingest_orders
INFO kept 3 orders, rejected 2 node=validate_orders
INFO report: 3 orders, 2 rejected, 45249 cents gross node=build_report
```
