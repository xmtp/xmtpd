---
name: git
description: >-
  Use when creating branches, committing changes, pushing, or opening pull
  requests. Triggers on "make a PR", "open a PR", "create a branch", "push
  this", "commit and push", or any request to ship code to GitHub.
---

# Creating Branches and Pull Requests

## Branch Naming Convention

**Always** use `<github-username>/<branch-description>`:

- GitHub username: !`gh api user --jq '.login'`

```
<username>/payer-id-lru-cache
<username>/fix-fee-calculation
<username>/add-congestion-metrics
```

- Prefix is always the result of `gh api user --jq '.login'`
- Description is kebab-case, concise, imperative
- Never use `main`, a bare description, or any other prefix