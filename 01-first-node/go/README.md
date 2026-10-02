# Step 1: your first node (Go)

A workflow with one node, `ingest_orders`. It returns a batch of five orders and logs how many there are. It has no types and no trigger: you run it by hand.

## The files

```
main.go   the workflow and its one node
go.mod    the module, requiring the Dagflows SDK
go.sum    the SDK's checksums
```

Dagflows finds the workflow without any configuration: a `go.mod` at the root makes this a
Go project, and Dagflows builds the module root, where `package main` lives.

## The code

```go
func ingestOrders(ctx *df.Ctx, _ df.None) (map[string]any, error) {
	orders := []map[string]any{
		{"id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd"},
		...
	}

	ctx.Log().Info("ingested orders", "run", ctx.Run().WorkflowRunID, "count", len(orders))
	return map[string]any{"orders": orders}, nil
}

func main() {
	wf := df.NewWorkflow("order-pipeline", df.WorkflowOptions{})
	wf.Node(ingestOrders, df.Root, df.NodeOptions{Key: "ingest_orders"})
	df.Main()
}
```

- `df.NewWorkflow("order-pipeline", ...)` names the workflow. Every step of these examples builds on it.
- A node's handler always has the shape `func(ctx *df.Ctx, in In) (Out, error)`. Here `In`
  is `df.None`, nothing, and `Out` is `map[string]any`, plain JSON.
- `wf.Node` registers the handler with where its input comes from, `df.Root` for a node
  with no parents, and its key, `ingest_orders`.
- `ctx` is the node's context. `ctx.Log()` is a `*slog.Logger` that writes to the run's
  logs, and `ctx.Run().WorkflowRunID` is the id of the run.
- `df.Main()` answers Dagflows: it writes the manifest when Dagflows builds the code, and
  runs the node when Dagflows invokes it.

## Run it on your machine

```bash
go run . build manifest -o dagflows-manifest.json
go run . dev run ingest_orders
```

`build manifest` writes what Dagflows reads from your code, one node called
`ingest_orders`. Dagflows builds the manifest itself when it deploys, so it is not
committed. `dev run` runs the node the way Dagflows does, and prints what it logged and
returned:

```
time=2026-10-02T15:26:51.839Z level=INFO msg="ingested orders" run=local count=5

ingest_orders -> SUCCESS
  {
    "orders": [
      {
        "amount_cents": 1250,
        "currency": "usd",
        "customer_id": "cus_alpha",
        "id": 1001
      },
      ...
    ]
  }
```

## Run it on Dagflows

Copy this directory into a repository of your own, then deploy and run it as
[Running an example](https://github.com/dagflows/examples#running-an-example-on-dagflows)
shows. The run takes no body. When it has finished, the node's log includes this line:

```
time=<when it ran> level=INFO msg="ingested orders" run=<the run's id> count=5
```
