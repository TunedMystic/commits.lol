#!/usr/bin/env bash

COVERAGE_AMOUNT=$(go tool cover -func coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')

if [ 1 -eq "$(echo "${COVERAGE_AMOUNT} < ${COVERAGE_THRESHOLD}" | bc)" ]
then
    echo -e "🚫 Coverage (${COVERAGE_AMOUNT}%) is below ${COVERAGE_THRESHOLD}%"
    exit 1
else
    echo -e "✅ Coverage (${COVERAGE_AMOUNT}%)"
fi
