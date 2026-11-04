# Dagflows Examples

These examples build one small order pipeline, a step at a time. Step 1 is a single node that
returns five orders. By step 8 the pipeline runs six nodes across Python, Go and TypeScript, starts
when a webhook delivers the orders, and signs its report with a secret key.

Each step adds one idea to the step before it, so start at step 1.

| Step | What it adds | Languages |
|---|---|---|
| 1. [Your first node](01-first-node) | A workflow with one node, run by hand | Python, Go or Node |
| 2. [Typed node](02-typed-node) | The node's output declared as a type | Python, Go or TypeScript |
| 3. [Several nodes](03-several-nodes) | A chain of nodes, each fed its parent's output | Python, Go or TypeScript |
| 4. [Two languages](04-two-languages) | One workflow across a Python and a Go project | Python and Go |
| 5. [Fan-out and join](05-fan-out-and-join) | Two branches that run side by side and meet in one node | Python, Go and TypeScript |
| 6. [Trigger](06-trigger) | The orders arrive as a typed event that starts the run | Python, Go and TypeScript |
| 7. [Webhook](07-webhook) | The same trigger, started by a signed request from outside | Python, Go and TypeScript |
| 8. [Variables and secrets](08-variables-and-secrets) | Values and a key kept on Dagflows instead of in the code | Python, Go and TypeScript |

## A Few Words First

- A **workflow** is your code: a set of **nodes**, functions that Dagflows runs for you in order.
- A node's **output** is what it returns. Dagflows hands it to the nodes that depend on it.
- To **deploy** is to have Dagflows build your repository into a **version** of the workflow.
- A **run** is one execution of that version. Each node of a run has its own status, logs and
  output.

## Deploy an Example on Dagflows

Every example deploys the same way. You need a GitHub account and a Dagflows account.

### 1. Sign Up

Create an account at [app.dagflows.com](https://app.dagflows.com) and verify your email. Your
account comes with an organization, where your projects live.

### 2. Put the Example in a Repository

Dagflows builds a workflow from the root of a git repository, so an example's folder becomes the
root of a repository of your own. Create an empty repository on GitHub, then:

```bash
git clone https://github.com/dagflows/examples.git
mkdir my-workflow
cp -r examples/01-first-node/python/. my-workflow/
cd my-workflow
git init -b main
git add .
git commit -m "Step 1"
git remote add origin https://github.com/YOUR-NAME/my-workflow.git
git push -u origin main
```

Steps 1 to 3 have a folder per language, so copy the one you want (`python`, `go` or `node`). From
step 4 on, an example's folder holds all of its projects, and the whole folder is copied.

A public repository builds as it is. To build a private one, connect GitHub first, in
**Settings → Connected accounts → Connect GitHub**.

### 3. Create a Project

In the sidebar, open **Projects** and click **New project**. Name it, for example `Examples`.

![The Create project dialog](images/new-project.png)

### 4. Create the Workflow

In the project, click **New workflow**. Name it `order-pipeline`, then point it at your repository:
pick it from your GitHub repositories, or click **Enter a URL instead** and give its URL and the
branch, `main`. Click **Create**.

![The Create workflow dialog with the repository's URL and branch](images/new-workflow.png)

### 5. Deploy It

On the workflow's page, click **Deploy**. Leave **Commit hash** empty, which builds the newest
commit of the branch, and click **Deploy**.

![The Deploy workflow dialog](images/deploy.png)

The **Builds** tab follows the build while Dagflows clones the repository, installs what each
project needs and builds the nodes. Once it finishes, that build is the version that runs. If it
fails, its log says why.

![A finished build on the Builds tab](images/build.png)

### 6. Run It

On the **Runs** tab, click **Run**. Steps 1 to 5 take no payload, so leave the box empty and click
**Run**. Steps 6 to 8 say what to send.

![The Run workflow dialog](images/run.png)

Open the run to follow it: each node, its status, its logs and its output.

![A finished run with its node, its log and its output](images/run-detail.png)

### Moving On to the Next Step

Keep the same repository and the same workflow for every step. Replace the repository's files with
the next step's, push, and click **Deploy** again:

```bash
cd my-workflow
git rm -rq .
cp -r ../examples/02-typed-node/python/. .
git add .
git commit -m "Step 2"
git push
```

This way the whole series fits the Free plan, which holds two workflows.
