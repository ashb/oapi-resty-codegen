
generate:
  go generate ./examples/...

build:
  go build -o oapi-resty-codegen ./main.go

test: generate

lint:
  pre-commit run --all-files --show-diff-on-failure --color=always
