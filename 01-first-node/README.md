# Step 1: Your First Node

A minimal workflow with a single untyped, manually triggered root node (`ingest_orders`) returning mock order data.

Pick your language - each directory is a standalone runnable project:

- [Python](python)
- [Go](go)
- [Node](node)

### What you'll learn
- How to define a workflow and register a root node.
- What the execution context (`ctx`) provides (logging and run metadata).
- How Dagflows auto-discovers projects from standard package files (`package.json`, `go.mod`, `requirements.txt`).
- How to inspect the workflow manifest and run the node locally before deploying.
