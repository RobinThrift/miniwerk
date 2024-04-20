
_install-go-tool cmd:
    #!/bin/sh
    set -o nounset 
    set -o pipefail
    set -o errexit

    cmd="{{ file_name(cmd) }}"

    if [ ! -x "{{  gobin  }}/$cmd" ]; then
        GOBIN={{ gobin }} go install -mod=readonly {{ cmd }}
        exit 0;
    fi;

    version_file="{{ gobin }}/.$cmd.version"
    if [ -f "$version_file" ]; then\
        curr_version=`cat $version_file`
        req_version=`go list -m all | grep $cmd | cut -d\  -f2`
        req_version=${req_version:1}
        if [ `just _semver-matches $curr_version ">=$req_version"` = "true" ]; then
            exit 0;
        fi;
    fi

    curr_version=`go list -m all | grep $cmd | cut -d\  -f2`
    curr_version=${curr_version:1}
    echo "$curr_version" > "$version_file"

_semver-matches version requirement:
    @echo {{ semver_matches(version, requirement) }}


_javy_version := "1.4.0"
_javy_arch := if arch() == "aarch64" { "arm" } else { "x86_64" }
_javy_archive := "javy-" + _javy_arch + "-" + os() + "-v" + _javy_version 
_install_javy:
    @([ -x "{{ gobin }}/javy" ] && [ "$({{ gobin }}/javy -V)" = "javy {{ _javy_version }}" ]) || just _download_and_install_javy

_download_and_install_javy:
    @mkdir -p .tmp
    curl -L "https://github.com/bytecodealliance/javy/releases/download/v{{ _javy_version }}/{{ _javy_archive }}.gz" -o .tmp/{{ _javy_archive }}.gz
    @gzip -d .tmp/{{ _javy_archive }}.gz
    @- rm -rf {{ gobin }}/javy
    @mv .tmp/{{ _javy_archive }} {{ gobin }}/javy
    @chmod +x {{ gobin }}/javy
    @echo "Installed" `{{ gobin }}/javy -V`
