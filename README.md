# Dagflows examples

Each example is a complete project. Read it here, then copy its directory into your own
repository to run it on Dagflows (see [Running an example](#running-an-example-on-dagflows)).

## Building workflows, step by step

Every step is the previous one with one idea added, so start at step 1.

| Step | What it adds | Python | Go | Node |
|---|---|---|---|---|
| 1. [Your first node](01-first-node) | One node, no types, run by hand | [python](01-first-node/python) | [go](01-first-node/go) | [node](01-first-node/node) |
| 2. [Typed node](02-typed-node) | The node's output declared as a type | [python](02-typed-node/python) | [go](02-typed-node/go) | [node](02-typed-node/node) (TypeScript) |
| 3. [Several nodes](03-several-nodes) | A chain of typed nodes, each fed its parent's output | [python](03-several-nodes/python) | [go](03-several-nodes/go) | [node](03-several-nodes/node) (TypeScript) |
| 4. [Two languages](04-two-languages) | One workflow across a Python and a Go project | [workspace](04-two-languages) | [workspace](04-two-languages) | |
| 5. [Fan-out and join](05-fan-out-and-join) | Two branches in parallel, joined in one node, across three languages | [workspace](05-fan-out-and-join) | [workspace](05-fan-out-and-join) | [workspace](05-fan-out-and-join) |
| 6. [Trigger](06-trigger) | The orders arrive as a typed event that starts the run | [workspace](06-trigger) | [workspace](06-trigger) | [workspace](06-trigger) |
| 7. [Webhook](07-webhook) | The same trigger bound to a webhook, started by a signed request | [workspace](07-webhook) | [workspace](07-webhook) | [workspace](07-webhook) |

---

## Running an example on Dagflows

Every example runs the same way. You need a Dagflows account with a verified email, and
`git`, `curl` and `jq`. The commands use a POSIX shell.

### 1. Put the example in a repository of your own

Dagflows builds a workflow from the root of a git repository, so copy one example's
directory, not this whole repository:

```bash
git clone https://github.com/dagflows/examples.git
cp -r examples/01-first-node/python my-workflow
cd my-workflow
git init -b main
git add .
git commit -m "My first Dagflows workflow"
git remote add origin https://github.com/YOUR-NAME/my-workflow.git
git push -u origin main
```

A public repository is cloned as it is. A private one needs a git account connected to
your Dagflows organization.

### 2. Log in

```bash
API=https://api.dagflows.com
TOKEN=$(curl -sS -X POST "$API/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "password": "your-password"}' | jq -r .data.token.access_token)
```

The token expires. When a call answers `401`, log in again.

### 3. Choose an organization and a project

```bash
curl -sS "$API/api/v1/organizations" -H "Authorization: Bearer $TOKEN" | jq '.data[] | {id, name}'
ORG=the-organization-id
```

If you have no organization yet, create one:

```bash
ORG=$(curl -sS -X POST "$API/api/v1/organizations" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name": "My team", "slug": "my-team"}' | jq -r .data.id)
```

Then a project to keep the examples in:

```bash
PROJECT=$(curl -sS -X POST "$API/api/v1/projects" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" -H "Content-Type: application/json" \
  -d '{"name": "Examples", "slug": "examples"}' | jq -r .data.id)
```

### 4. Create the workflow from your repository

```bash
WORKFLOW=$(curl -sS -X POST "$API/api/v1/projects/$PROJECT/workflows" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" -H "Content-Type: application/json" \
  -d '{"name": "order-pipeline", "git_url": "https://github.com/YOUR-NAME/my-workflow", "git_branch": "main"}' \
  | jq -r .data.id)
```

### 5. Deploy it

A deployment clones the branch, builds the code and makes it the version that runs.

```bash
DEPLOYMENT=$(curl -sS -X POST "$API/api/v1/workflows/$WORKFLOW/deployments" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" -H "Content-Type: application/json" \
  -d '{}' | jq -r .data.id)

curl -sS "$API/api/v1/workflows/$WORKFLOW/deployments/$DEPLOYMENT/progress" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" | jq '.data | {phase, pct, msg}'
```

Repeat the second command until `phase` is `DONE`. If it is `FAILED`, the build log says
why: `GET $API/api/v1/workflows/$WORKFLOW/deployments/$DEPLOYMENT/logs` holds its `lines`
while it builds and an `archive_url` once it has finished.

### 6. Run it and read what it logged

```bash
RUN=$(curl -sS -X POST "$API/api/v1/workflows/$WORKFLOW/runs" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" | jq -r .data.id)

curl -sS "$API/api/v1/workflows/$WORKFLOW/runs/$RUN" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" \
  | jq '.data | {status, nodes: [.nodes[] | {node_key, status}]}'
```

A run answered with `404` and "workflow version not found" means the deployment has not
finished yet. Repeat the second command until `status` is `SUCCESS` or `FAILED`.

From step 6 on, a workflow declares a trigger and its run carries a payload. The example's
README gives the command that starts it.

What a node prints goes to its log. Once a node has finished, its log is stored and the
logs endpoint gives a link to it, a gzipped file with one JSON line per printed line:

```bash
curl -sS "$API/api/v1/workflows/$WORKFLOW/runs/$RUN/logs" \
  -H "Authorization: Bearer $TOKEN" -H "X-Organization-Id: $ORG" \
  | jq -r '.data.nodes[] | .archive[]?.url' \
  | while read -r url; do curl -sS "$url" | gunzip | jq -r .line; done
```

While a node is still running, its lines come back in the same response instead, under
`.data.nodes[].lines[].line`.
