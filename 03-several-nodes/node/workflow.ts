import { Workflow } from "@dagflows/sdk/authoring";
import type { Ctx } from "@dagflows/sdk/runtime";

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

interface ValidatedOrders {
  orders: Order[];
  rejected: number;
}

interface Report {
  total_orders: number;
  rejected: number;
  gross_cents: number;
}

export const wf = new Workflow("order-pipeline");

// Node handlers must be exported so Dagflows can discover them.
// The node key defaults to the function name unless specified in options.
// The return type becomes the node's output schema.
export const ingestOrders = wf.node(
  function ingestOrders ({ ctx }): Orders {
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

// depends names the node's parents. With exactly one parent, the handler receives
// that parent's output as input. Annotating input puts it in the manifest as the
// node's input schema, which Dagflows checks against the parent's output schema.
export const validateOrders = wf.node(
  function validateOrders ({ ctx, input }: { ctx: Ctx; input: Orders }): ValidatedOrders {
    const orders = input.orders
      .filter((order) => order.currency.toLowerCase() === "usd" && order.amount_cents > 0)
      .map((order) => ({ ...order, currency: "usd" }));

    const rejected = input.orders.length - orders.length;

    ctx.log.info(`kept ${orders.length} orders, rejected ${rejected}`);

    return { orders, rejected };
  },
  { key: "validate_orders", depends: [ingestOrders] },
);

export const buildReport = wf.node(
  function buildReport ({ ctx, input }: { ctx: Ctx; input: ValidatedOrders }): Report {
    const report: Report = {
      total_orders: input.orders.length,
      rejected: input.rejected,
      gross_cents: input.orders.reduce((sum, order) => sum + order.amount_cents, 0),
    };

    ctx.log.info(
      `report: ${report.total_orders} orders, ${report.rejected} rejected, ${report.gross_cents} cents gross`,
    );

    return report;
  },
  { key: "build_report", depends: [validateOrders] },
);
