"""A workflow with one node, run by hand."""

from dagflows.authoring import Workflow

# The workflow every step of these examples builds on.
wf = Workflow("order-pipeline")


@wf.node()
def ingest_orders(ctx):
    """Returns a batch of orders as plain data, and logs how many there are.

    `ctx` is the node's context: `ctx.log` writes to the run's logs, and `ctx.run`
    says which run this is.
    """
    orders = [
        {"id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd"},
        {"id": 1002, "customer_id": "cus_beta", "amount_cents": 34000, "currency": "USD"},
        {"id": 1003, "customer_id": "cus_gamma", "amount_cents": 9999, "currency": "usd"},
        {"id": 1004, "customer_id": "cus_delta", "amount_cents": 500, "currency": "eur"},
        {"id": 1005, "customer_id": "cus_epsilon", "amount_cents": 0, "currency": "usd"},
    ]

    ctx.log.info("run %s: ingested %d orders", ctx.run.workflow_run_id, len(orders))
    return {"orders": orders}
