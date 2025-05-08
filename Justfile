
test *args='-f dots-v2 -- ./...': generate
  go run gotest.tools/gotestsum@v1.12.1 {{args}}

generate:
  go generate ./examples/...

build:
  go build -o oapi-resty-codegen ./main.go


lint:
  pre-commit run --all-files --show-diff-on-failure --color=always

[doc("""
  Use nektos/act to run GHA workflows locally

  Example: `just act --matrix=version:1.24
""")]
act *args:
  # See .actrc for default flags
  act {{args}}
