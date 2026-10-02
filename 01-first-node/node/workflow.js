// A workflow with one node, run by hand.

import { Workflow } from "@dagflows/sdk/authoring";

// The workflow every step of these examples builds on.
export const wf = new Workflow("order-pipeline");

/**
 * Returns a batch of orders as plain data, and logs how many there are.
 * `ctx` is the node's context: `ctx.log` writes to the run's logs, and `ctx.run`
 * says which run this is. A node is found by its export, so it is exported.
 */
export const ingestOrders = wf.node(
  function ingestOrders({ ctx }) {
    const orders = [
      { id: 1001, customer_id: "cus_alpha", amount_cents: 1250, currency: "usd" },
      { id: 1002, customer_id: "cus_beta", amount_cents: 34000, currency: "USD" },
      { id: 1003, customer_id: "cus_gamma", amount_cents: 9999, currency: "usd" },
      { id: 1004, customer_id: "cus_delta", amount_cents: 500, currency: "eur" },
      { id: 1005, customer_id: "cus_epsilon", amount_cents: 0, currency: "usd" },
    ];

    ctx.log.info(`run ${ctx.run.workflowRunId}: ingested ${orders.length} orders`);
    return { orders };
  },
  { key: "ingest_orders" },
);
