_ui_src_files := "ui/src/*.{ts,tsx} ui/src/**/*.{ts,tsx}"

_ui-install:
    @if [ ! -d node_modules ] || [ package.json -nt node_modules/.package-lock.json ]; then \
        npm i --no-audit --no-fund ; \
    fi

_ui-fmt: _ui-install
    {{ nmbin }}/biome format --write {{ _ui_src_files }}

_ui-lint: _ui-install
    {{ nmbin }}/biome check {{ _ui_src_files }}

_ui-lint-ci: _ui-install
    {{ nmbin }}/biome check {{ _ui_src_files }}

_ui-build: _ui-build-browser _ui-build-server


_ui-build-styles: _ui-install
    NODE_ENV=production {{ nmbin }}/postcss -c frontend/postcss.config.js ./frontend/src/styles.css -o ./frontend/build/styles.css --no-map

_ui-watch-styles: _ui-install
    {{ nmbin }}/postcss --watch --verbose ./ui/src/styles/index.css -o ./ui/build/browser/styles.css

_ui_browser_bundle := "ui/build/browser/index.min.js"
_ui-build-browser: _ui-install
    {{ nmbin }}/esbuild ui/src/browser/index.ts --format=esm --target=es2020 --minify --bundle --outfile={{ _ui_browser_bundle }}

_ui-watch-browser: _ui-install
    {{ nmbin }}/esbuild ui/src/browser/index.ts --format=esm --target=es2020 --bundle --outfile={{ _ui_browser_bundle }} --watch=forever

_ui_sever_bundle := "ui/build/server/index.min.js"
_ui-build-server: _ui-bundle-server _install_javy
    {{ gobin }}/javy compile {{ _ui_sever_bundle }} -o {{ replace(_ui_sever_bundle, ".min.js", ".wasm") }} --no-source-compression

_ui-bundle-server: _ui-install
    {{ nmbin }}/esbuild ui/src/server/index.ts --format=esm --target=es2020 --minify --bundle --outfile={{ _ui_sever_bundle }}
