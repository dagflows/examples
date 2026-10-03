import { Workflow } from "@dagflows/sdk/authoring";
import type { Ctx } from "@dagflows/sdk/runtime";

// Interfaces are reflected into the manifest as JSON Schema, read from this source
// by the typescript package when the manifest is built. A bigint is reflected as an
// int64 integer, which the Python and Go integers on the other side of each edge need.
// A number is a JSON number, which never satisfies an integer.
interface Order {
  id: bigint;
  customer_id: string;
  amount_cents: bigint;
  currency: string;
}

interface Orders {
  orders: Order[];
}

interface ValidatedOrders {
  orders: Order[];
  rejected: bigint;
}

interface EnrichedOrder extends Order {
  tier: string;
}

interface EnrichedOrders {
  orders: EnrichedOrder[];
  rejected: bigint;
}

interface RiskScores {
  scores: { order_id: bigint; score: number }[];
}

// A project contributing nodes to a workflow another project owns uses an unnamed workflow.
export const wf = new Workflow();

// externalNode names a node of another project by its key, typed with what this
// project expects of it, which Dagflows checks against what that project declares.
const validateOrders = wf.externalNode<ValidatedOrders>("validate_orders");
const ingestOrders = wf.externalNode<Orders>("ingest_orders");

// Node handlers must be exported so Dagflows can discover them. Annotating input
// puts it in the manifest as the node's input schema.
export const enrichCustomers = wf.node(
  function enrichCustomers ({ ctx, input }: { ctx: Ctx; input: ValidatedOrders }): EnrichedOrders {
    const orders = input.orders.map((order) => ({
      ...order,
      tier: order.amount_cents >= 10_000n ? "gold" : "standard",
    }));

    ctx.log.info(`enriched ${orders.length} orders`);

    return { orders, rejected: input.rejected };
  },
  { key: "enrich_customers", depends: [validateOrders] },
);

export const scoreRisk = wf.node(
  function scoreRisk ({ ctx, input }: { ctx: Ctx; input: Orders }): RiskScores {
    const scores = input.orders.map((order) => {
      // An integer within 2^53 arrives as a number whatever its declared type.
      // BigInt() makes it the bigint it is declared as before bigint arithmetic.
      const id = BigInt(order.id);

      return { order_id: id, score: Number(id % 10n) / 10 };
    });

    ctx.log.info(`scored ${scores.length} orders`);

    return { scores };
  },
  { key: "score_risk", depends: [ingestOrders] },
);
