# Step 2: Typed Node (Node)

A minimal workflow with one root node (`ingest_orders`) that returns the same mock order data as step 1, now typed with TypeScript. This node takes no trigger payload and is invoked manually.

## What's New

Compared with [step 1](https://github.com/dagflows/examples/tree/main/01-first-node/node):

- The workflow is TypeScript (`workflow.ts`): `Order` and `Orders` are interfaces, and `ingestOrders` returns `Orders`.
- `typescript` is a dev dependency. The SDK uses it to read the types when it builds the manifest, so the step 1 "no types" warning is gone. Node runs the `.ts` file directly, with no compile step.
- `dagflows.yaml` names the entry module, since Dagflows only discovers `workflow.js` (or `app/workflow.js`) on its own.
- The manifest now carries the node's output schema. The data the node returns is unchanged.

## Files

```
workflow.ts         workflow definition, types and node handler
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

# 1. Build the workflow manifest and print the node's output schema
npx dagflows-sdk build manifest workflow.ts -o dagflows-manifest.json
jq '.nodes[0].io.output.schema' dagflows-manifest.json

# 2. Execute the node locally
npx dagflows-sdk dev run workflow.ts:ingestOrders
```

- `build manifest`: Reflects the `Orders` return type into the node's output schema:

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
          "id": { "type": "number" },
          "customer_id": { "type": "string" },
          "amount_cents": { "type": "number" },
          "currency": { "type": "string" }
        },
        "required": ["id", "customer_id", "amount_cents", "currency"]
      }
    }
  },
  "required": ["orders"]
}
```

A TypeScript `number` is a JSON number, where Python's `int` and Go's `int64` are integers. That only matters once nodes in different languages meet, in step 5.

- `dev run`: Executes the handler as in step 1, logging the same line and returning the same JSON:

```
INFO run local: ingested 5 orders node=ingestOrders

workflow.ts:ingestOrders -> SUCCESS
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

*Note: The local runner identifies the node by its export (`ingestOrders`). On the Dagflows platform, it is tracked by its registered key (`ingest_orders`).*

## Run on Dagflows

Push this directory to your repository, then deploy and trigger a run as shown in [Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows).

The workflow requires no input payload. Once execution completes, the run log will show:

```
INFO run <workflow_run_id>: ingested 5 orders node=ingest_orders
```
