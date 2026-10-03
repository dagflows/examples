import { Workflow } from "@dagflows/sdk/authoring";

// Interfaces are reflected into the manifest as JSON Schema, read from this source
// by the typescript package when the manifest is built.
interface Order {
  id: number;
  customer_id: string;
  amount_cents: number;
  currency: string;
}

interface Orders {
  orders: Order[];
}

export const wf = new Workflow("order-pipeline");

// Node handlers must be exported so Dagflows can discover them.
// The node key defaults to the function name unless specified in options.
// The return type becomes the node's output schema.
export const ingestOrders = wf.node(
  function ingestOrders({ ctx }): Orders {
    const orders: Order[] = [
      { id: 1001, customer_id: "cus_alpha", amount_cents: 1250, currency: "usd" },
      { id: 1002, customer_id: "cus_beta", amount_cents: 34000, currency: "USD" },
      { id: 1003, customer_id: "cus_gamma", amount_cents: 9999, currency: "usd" },
      { id: 1004, customer_id: "cus_delta", amount_cents: 500, currency: "eur" },
      { id: 1005, customer_id: "cus_epsilon", amount_cents: 0, currency: "usd" },
    ];

    // ctx provides run metadata and structured logging sent to the run console.
    ctx.log.info(`run ${ctx.run.workflowRunId}: ingested ${orders.length} orders`);

    return { orders };
  },
  { key: "ingest_orders" },
);
