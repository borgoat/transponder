# Terraform Transponder

[![Go Report Card](https://goreportcard.com/badge/github.com/transponder-tf/transponder)](https://goreportcard.com/report/github.com/transponder-tf/transponder)

###### IMPORTANT
*Transponder* is compatible with __Terraform 0.12__
(currently in beta).

Transponder is meant to offer 2 features:
- a remote state backend for Terraform w/ support for locking
- a generic, language-independent, API to query the current (and possibly past) state

The code is organised as follows:
- `server`: the listener answering Terraform
  when changes are planned or applied.
  For the time being, this is compatible with the
  [http](https://www.terraform.io/docs/backends/types/http.html)
  backend. In the future, a specific client implementation
  (i.e. on Terraform side), could be created.
- `statemgrmap`: store and retrieve the .tfstate file,
  with support for versions and multitenancy
- `transformer`: the external API is implemented here

---

###### TODO

- Implement auth for the HTTP backend
- Support namespaces and workspaces
- Develop a custom client+backend
- GraphQL access
- "Raw" Terraform address access
- Distributed implementation of statemgr.Full
