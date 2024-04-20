go_ldflgas := env_var_or_default("GO_LDFLGAS", "") + " -X 'github.com/RobinThrift/miniwerk/miniwerk.Version=" + version + "'"

_go-fmt:
    go fmt ./...

_go-lint: (_install-go-tool "honnef.co/go/tools/cmd/staticcheck") (_install-go-tool "github.com/golangci/golangci-lint/cmd/golangci-lint")
    {{ gobin }}/staticcheck ./...
    {{ gobin }}/golangci-lint run ./...

_go-lint-ci: (_install-go-tool "honnef.co/go/tools/cmd/staticcheck") (_install-go-tool "github.com/golangci/golangci-lint/cmd/golangci-lint")
    {{ gobin }}/golangci-lint run --timeout 5m --out-format=junit-xml ./... > lint.junit.xml
    {{ gobin }}/staticcheck ./...

_go-test *flags="-failfast -v -timeout 5m": (_install-go-tool "gotest.tools/gotestsum")
    {{ gobin }}/gotestsum --format short-verbose -- {{ flags }} ./...

_go-test-ci: (_install-go-tool "gotest.tools/gotestsum")
    {{ gobin }}/gotestsum --format short-verbose --junitfile=test.junit.xml -- -timeout 10m ./...

_go-test-watch *flags="-failfast -timeout 5m": (_install-go-tool "gotest.tools/gotestsum")
    {{ gobin }}/gotestsum --watch --format short-verbose -- {{ flags }}

_go-run:
    @-mkdir -p {{ _tmp_run_dir }}
    @MINIWERK_LOG_LEVEL="debug" MINIWERK_LOG_FORMAT="console" \
        MINIWERK_ADDRESS="localhost:8080" \
        go run -tags dev -race ./miniwerk/bin/miniwerk

_go-build:
    go build -ldflags="{{go_ldflgas}}" -o build/miniwerk ./miniwerk/bin/miniwerk
