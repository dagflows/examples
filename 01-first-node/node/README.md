# Step 1: your first node (Node)

A workflow with one node, `ingest_orders`. It returns a batch of five orders and logs how many there are. It has no types and no trigger: you run it by hand.

## The files

```
workflow.js         the workflow and its one node
package.json        the project, depending on the Dagflows SDK
package-lock.json   the exact versions installed
```

Dagflows finds the workflow without any configuration: a `package.json` depending on
`@dagflows/sdk` makes this a Node project, and `workflow.js` is where the workflow is
looked for (`app/workflow.js` works too). The lockfile is required, since Dagflows
installs with `npm ci`.

## The code

```js
export const wf = new Workflow("order-pipeline");

export const ingestOrders = wf.node(
  function ingestOrders({ ctx }) {
    const orders = [
      { id: 1001, customer_id: "cus_alpha", amount_cents: 1250, currency: "usd" },
      ...
    ];

    ctx.log.info(`run ${ctx.run.workflowRunId}: ingested ${orders.length} orders`);
    return { orders };
  },
  { key: "ingest_orders" },
);
```

- `new Workflow("order-pipeline")` names the workflow. Every step of these examples builds on it.
- `wf.node` makes the function a node, keyed `ingest_orders`. Without `key`, the key
  would be the function's name. A node is found by its export, so it is exported.
- The handler receives `{ ctx }`, the node's context. `ctx.log` writes to the run's logs,
  and `ctx.run.workflowRunId` is the id of the run.
- What the function returns is the node's output, plain JSON here.

## Run it on your machine

```bash
npm ci
npx dagflows-sdk build manifest workflow.js -o dagflows-manifest.json
npx dagflows-sdk dev run workflow.js:ingestOrders
```

`build manifest` writes what Dagflows reads from your code, one node called
`ingest_orders`. Dagflows builds the manifest itself when it deploys, so it is not
committed. It also warns:

```
warning: workflow.js has no TypeScript source beside it, so the manifest carries no types, the platform will check nothing on its edges
```

That is expected here, since this node has no types. Step 2 adds them, with TypeScript.

`dev run` runs the node the way Dagflows does, and prints what it logged and returned:

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

On your machine the runner names the node after its export, `ingestOrders`. On Dagflows
it is the key, `ingest_orders`.

## Run it on Dagflows

Copy this directory into a repository of your own, then deploy and run it as
[Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows)
shows. The run takes no body. When it has finished, the node's log includes this line:

```
INFO run <the run's id>: ingested 5 orders node=ingest_orders
```
