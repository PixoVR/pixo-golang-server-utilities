#!/bin/bash

RED="\033[1;31m"
GREEN="\033[1;32m"
NOCOLOR="\033[0m"

COVERAGE_DIR="./coverage"
COVERAGE_FLOOR="70"

VETPKGS="bools,httpresponse,printf,tests,structtag,unreachable,unsafeptr"
COVERPKGS="./src/..."

mkdir -p $COVERAGE_DIR

if [[ $1 ]]; then
  COVERAGE_FLOOR=$1
fi

ginkgo -cover -race -trace -r -covermode atomic -vet $VETPKGS -coverpkg=$COVERPKGS

SUCCESS=$?

if [[ $SUCCESS == 1 ]]; then
  exit 1
fi

echo "Merging coverage reports..."

gocovmerge "$(find . -type f -name "*.coverprofile")" > $COVERAGE_DIR/merged-coverage.out

echo "Coverage reports merged. Calculating package coverage..."

go test -coverprofile=coverage.out $COVERPKGS

SUCCESS=0
go tool cover -func=coverage.out | while read -r line ; do
  if [[ $line == *"total:"* ]]; then
    COVERAGE=$(echo "$line" | awk '{print $NF}' | tr -d '%')
    PKG=$(echo "$line" | awk '{print $1}')
    if [[ $(echo "$COVERAGE < $COVERAGE_FLOOR" | bc -l) == 1 ]]; then
      echo -e "${RED}FAILED:${NOCOLOR} minimum code coverage not met for package ${PKG} - ${COVERAGE} < ${COVERAGE_FLOOR}"
      SUCCESS=1
    else
      echo -e "${GREEN}SUCCESS:${NOCOLOR} minimum code coverage met for package ${PKG} - ${COVERAGE} >= ${COVERAGE_FLOOR}"
    fi
  fi
done

exit $SUCCESS
