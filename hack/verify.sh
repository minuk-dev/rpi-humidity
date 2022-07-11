#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail
ORI_DIR=$(pwd)
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

DIFFROOT="${SCRIPT_ROOT}/pkg"
TMP_DIFFROOT="${SCRIPT_ROOT}/_tmp/pkg"
_tmp="${SCRIPT_ROOT}/_tmp"

cleanup() {
  rm -rf "${_tmp}"
}
trap "cleanup" EXIT SIGINT

cleanup

mkdir -p "${TMP_DIFFROOT}"
cp -a "${DIFFROOT}"/* "${TMP_DIFFROOT}"
cp go.mod ${_tmp}
cd ${_tmp}
"../${SCRIPT_ROOT}/hack/update.sh"
cd ${ORI_DIR}

ret1=0
ret2=0
diff -Naup "${SCRIPT_ROOT}/go.mod" "${_tmp}/go.mod" || ret1=$?
diff -Naup "${SCRIPT_ROOT}/go.sum" "${_tmp}/go.sum" || ret2=$?
if [ $ret1 -eq 0 ] && [ $ret2 -eq 0 ]
then
  echo "${DIFFROOT} up to date."
else
  echo "${DIFFROOT} is out of date. Please run hack/update.sh"
  exit 1
fi
