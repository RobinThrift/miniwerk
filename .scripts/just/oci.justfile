oci_repo       := env_var_or_default("OCI_REPO", "ghcr.io/robinthrift/miniwerk")

_oci_tool := env_var_or_default("OCI_TOOL", "docker")

_oci_build_cmd := env_var_or_default("OCI_BUILD_CMD", "build")
_oci-build:
    {{ _oci_tool }} {{ _oci_build_cmd }} --build-arg="VERSION={{ version }}" -f ./deployment/Dockerfile  -t {{ oci_repo }}:{{ version }} .

_oci-run: _oci-build
    {{ _oci_tool }} run --rm \
        -p 8080:8080 \
        {{ oci_repo }}:{{ version }}


