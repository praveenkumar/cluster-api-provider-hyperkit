#!/bin/sh
# Example:  ./hack/go-lint.sh installer/... pkg/... tests/smoke

if [ "$IS_CONTAINER" != "" ]; then
  golint -set_exit_status "${@}"
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/praveenkumar/cluster-api-provider-hyperkit:z" \
    --workdir /go/src/github.com/praveenkumar/cluster-api-provider-hyperkit \
    openshift/origin-release:golang-1.10 \
    ./hack/go-lint.sh "${@}"
fi