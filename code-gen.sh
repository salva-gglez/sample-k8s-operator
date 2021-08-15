#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(
  cd "${SCRIPT_ROOT}"
  ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator
)}
GO_PKG="github.com/mmontes11/echoperator/pkg"

bash "${CODEGEN_PKG}"/generate-groups.sh "all" \
  ${GO_PKG}/echo/v1alpha1/apis \
  ${GO_PKG} \
  echo:v1alpha1 \
  --go-header-file "${SCRIPT_ROOT}"/codegen/boilerplate.go.txt