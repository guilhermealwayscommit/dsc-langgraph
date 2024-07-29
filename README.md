# dsc-langgraph

## Structure

- [logs](logs)
  - folders by date
- [.env](.env)
- [.env.example](.env.example)
- go.mod (go package manager)
- go.sum (go package manager)
- [langsmith_proxy.go](langsmith_proxy.go)
  - file for langsmith server proxy (logs langsmith requests to langchain)
  - [installation process](#go-installation)
- [langgraph.go](langgraph.go)

## Go Installation

- How to install Go
  - <https://golang.org/doc/install>
- How to run the file
  - `go run langsmith_proxy.go`
