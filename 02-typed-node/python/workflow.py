from dataclasses import dataclass

from dagflows.authoring import Workflow
from dagflows.runtime import Ctx

wf = Workflow("order-pipeline")


# Dataclasses are reflected into the manifest as JSON Schema.
@dataclass
class Order:
    id: int
    customer_id: str
    amount_cents: int
    currency: str


@dataclass
class Orders:
    orders: list[Order]


# Node keys default to the decorated function name unless overridden in @wf.node().
# The return annotation becomes the node's output schema.
@wf.node()
def ingest_orders(ctx: Ctx) -> Orders:
    orders = [
        Order(id=1001, customer_id="cus_alpha", amount_cents=1250, currency="usd"),
        Order(id=1002, customer_id="cus_beta", amount_cents=34000, currency="USD"),
        Order(id=1003, customer_id="cus_gamma", amount_cents=9999, currency="usd"),
        Order(id=1004, customer_id="cus_delta", amount_cents=500, currency="eur"),
        Order(id=1005, customer_id="cus_epsilon", amount_cents=0, currency="usd"),
    ]

    # ctx provides structured logging attached to this workflow run.
    ctx.log.info("run %s: ingested %d orders", ctx.run.workflow_run_id, len(orders))

    return Orders(orders=orders)
