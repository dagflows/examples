# Step 1: Your First Node (Node)

A minimal workflow with one root node (`ingest_orders`) that returns mock order data. This node is untyped, takes no trigger payload, and is invoked manually.

## Files

```
workflow.js         workflow definition and node handler
package.json        project manifest depending on @dagflows/sdk
package-lock.json   pinned dependency tree (required for npm ci)
```

Dagflows automatically discovers the project from `package.json` and loads `workflow.js` (or `app/workflow.js`). The lockfile is required because Dagflows builds dependencies using `npm ci`.

## Code

[workflow.js](workflow.js) contains the complete workflow implementation.

## Run Locally

```bash
npm ci

# 1. Build the workflow manifest
npx dagflows-sdk build manifest workflow.js -o dagflows-manifest.json

# 2. Execute the node locally
npx dagflows-sdk dev run workflow.js:ingestOrders
```

`build manifest` inspects the exported workflow and extracts the node metadata into `dagflows-manifest.json`. You'll see this expected warning:

```
warning: workflow.js has no TypeScript source beside it, so the manifest carries no types, the platform will check nothing on its edges
```

This is expected for untyped JavaScript nodes. (Step 2 adds build-time edge type checking with TypeScript).

`dev run` executes the handler locally, printing the run log and output:

```
INFO run local: ingested 5 orders node=ingestOrders

workflow.js:ingestOrders -> SUCCESS
  {
    "orders": [
      {
        "id": 1001,
        "customer_id": "cus_alpha",
        "amount_cents": 1250,
        "currency": "usd"
      },
      ...
    ]
  }
```

*Note: The local runner identifies the node by its JavaScript export (`ingestOrders`). On the Dagflows platform, it is tracked by its registered key (`ingest_orders`).*

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Once execution completes, the run log will show:

```
INFO run <workflow_run_id>: ingested 5 orders node=ingest_orders
```
