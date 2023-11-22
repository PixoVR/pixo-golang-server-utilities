#!/bin/bash

RED="\033[1;31m"
GREEN="\033[1;32m"
NOCOLOR="\033[0m"

VETPKGS="bools,httpresponse,printf,tests,structtag,unreachable,unsafeptr"

COVERAGE_OUTPUT=$(ginkgo -r -race -trace -v --cover --covermode atomic -vet $VETPKGS | grep "composite coverage:")

COVERAGE_NUMBER=$(echo "$COVERAGE_OUTPUT" | awk '{print $3}' | tr -d '%')


if [[ $1 ]]; then
  COVERAGE_FLOOR=$1
fi

if [[ -z $COVERAGE_FLOOR ]]; then
  COVERAGE_FLOOR=80
fi

if [[ $(echo "$COVERAGE_FLOOR > $COVERAGE_NUMBER" | awk '{print ($1 > $2)}') == 1 ]]; then
  echo -e "${RED}FAILED:${NOCOLOR} minimum code coverage not met for project - ${COVERAGE_NUMBER} < ${COVERAGE_FLOOR}"
  exit 1
fi

echo -e "${GREEN}SUCCESS:${NOCOLOR} minimum code coverage met for project - ${COVERAGE_NUMBER} >= ${COVERAGE_FLOOR}"
