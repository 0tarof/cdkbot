# cdkbot

cdkbot is AWS CDK review & deploy tool running on Pull Request.
Currently only GitHub is supported.

## Commands

Following commands are runnable by PR comments. 
Before running a command, base (where to merge) branch is merged internally 
so it is needed to resolve conflicts if it occurred.

- `/diff`: cdk diff all stacks. Run automatically when open PR and push to PR.
- `/deploy [stack1 stack2 ...]`: cdk deploy. If not specify stacks, all stacks are passed.
- `/rollback [stack1 stack2 ...]`: cdk deploy at base branch. If not specify stacks, all stacks are passed.

![run /diff and /deploy](./doc-assets/run-diff-deploy.png)

After running `/deploy`, 
PR is merged automatically if there are no differences anymore, 
and `cdkbot:outdated diffs` label is added to other PRs. 
To run `/deploy` on those, 
it is needed to run `/diff` again to see the latest differences.

![oudated diffs label](./doc-assets/outdated-diffs.png)

### Why deploys before merging PR?

cdk deploy sometimes fails unexpectedly due to runtime errors of CFn template.
Therefore, if the flow that deploys after merging is adopted, 
broken codes can be merged and surplus PRs are opened to revert or fix, which flagment changes. 

Deploying before merging PR has the advantage of avoiding these but changes may be reverted.
To prevent this, cdkbot takes measures these:

- base branch is merged internally before running a command, and deployed PR is merged automatically
- sets the number of concurrent executions to 1 and forces to see latest differences by `cdkbot:outdated diffs` label.

## Install & Settings

### Install

Install from [Serverless Application Repository](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:524580158183:applications~cdkbot) 
or `make deploy S3Bucket=*** Platform=github GitHubUserName=*** GitHubAccessToken=*** GitHubWebhookSecret=***`.

- GitHubUserName & GitHubAccessToken

Token can be generated at `Settings/Developer settings`.
repo and write:discussion scopes are required.

- GitHubWebhookSecret: Generate a random string.
- Platform: Only github.

### Repository settings

Add a webhook at repository's settings. 

- Payload URL: See CloudFormation Stack output
- Content type: application/json 
- Secret: same value of GitHubWebhookSecret
- Event trigger: Pushes, Issue comments and Pull requests

After the first run, enable "Require status checks to pass before merging" 
in the branch protection rule to prevent merging before deploying (Recommended)

### cdkbot.yml

Put `cdkbot.yml` at the repository root. 
cdkbot refer to the original for security and authority reasons 
so it's needed to be merged to apply the changes.

```
cdkRoot: . # CDK directory path from repository root.
targets:
  # If any key is matched the PR base branch, run commands with contexts `-c key=value`.
  # If not, commands are not runned.
  develop:
    contexts:
      env: stg
  master:
    contexts:
      env: prd
deployUsers:
  # Optional. If specified, only these users are allowed to deploy.
  # If not, all users are allowed to deploy.
  - sambaiz
```

