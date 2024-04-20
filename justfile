_tmp_run_dir := "test/.run"
gobin        := absolute_path(".gobin")
nmbin        := absolute_path("node_modules/.bin")
version      := env_var_or_default("VERSION", "dev")

import ".scripts/just/tools.justfile"
import ".scripts/just/release.justfile"
import ".scripts/just/db.justfile"
import ".scripts/just/go.justfile"
import ".scripts/just/oci.justfile"
import ".scripts/just/ui.justfile"

_default:
    @just --list

fmt: _go-fmt _ui-fmt

lint: _go-lint _ui-lint

lint-ci:  _go-lint-ci _ui-lint-ci

test *flags:
    just _go-test {{ flags }}

test-ci: _go-test-ci

alias tw := test-watch
test-watch *flags:
    just _go-test-watch {{ flags }}

build:
    just _ui-build
    just _go-build

run: _ui-build
    @just _go-run

watch: (_install-go-tool "github.com/bokwoon95/wgo")
    mkdir -p test/tmp/.run
    {{ gobin }}/wgo \
        -file '.*\.go' \
        -xfile '.*_test\.go' \
        just _go-run :: \
    wgo \
        -file package.json \
       just _ui-watch-styles  :: \
    wgo \
        -file package.json \
        just _ui-watch-browser :: \
    wgo \
        -file 'ui/src/.*\.tsx?' \
        just _ui-build-server

clean:
    rm -rf node_modules build .gobin .tmp {{ _tmp_run_dir }} ui/build
    go clean -cache
