_tmp_db_file := "./tmp/_miniwerk.db"

new-migration name: (_install-go-tool "github.com/pressly/goose/v3/cmd/goose")
    @- mkdir -p {{ parent_directory(_tmp_db_file) }}
    @rm -rf {{ _tmp_db_file }}
    {{ gobin }}/goose -table migrations -dir ./miniwerk/storage/database/sqlite/migrations sqlite3 {{ _tmp_db_file }} create {{name}} sql
    @rm -rf {{ _tmp_db_file }}

alias gen := generate
generate: (_install-go-tool "github.com/pressly/goose/v3/cmd/goose") (_install-go-tool "github.com/stephenafamo/bob/gen/bobgen-sqlite")
    @- mkdir -p {{ parent_directory(_tmp_db_file) }}
    @rm -rf {{ _tmp_db_file }}
    {{ gobin }}/goose -table migrations -dir ./miniwerk/storage/database/sqlite/migrations sqlite3 {{ _tmp_db_file }} up
    {{ gobin }}/bobgen-sqlite -c ./miniwerk/storage/database/sqlite/bob.yaml
    go fmt ./...
    @rm -rf {{ _tmp_db_file }}

