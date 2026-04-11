#!/bin/bash

# Arguments:
# $1: Additional go test flags (e.g., -short)

TEST_FLAGS=$1
COVERAGE_OUT="coverage.out"
FILTERED_OUT="coverage_filtered.out"

# 1. Run tests with coverage
# We exclude /mocks and /migrations from the package list to speed up testing
echo "Running tests with flags: $TEST_FLAGS"
go test $TEST_FLAGS -coverprofile=$COVERAGE_OUT $(go list ./... | grep -vE "/mocks|/migrations") > /dev/null

if [ ! -f $COVERAGE_OUT ]; then
    echo "Error: coverage.out not generated."
    exit 1
fi

# 2. Filter out unwanted files from the profile
# - *_mock.go (individual mock files)
# - cmd/api/main.go (entry point/wiring)
# - internal/adapters/db/connection.go (database initialization/boilerplate)
grep -vE "_mock.go|cmd/api/main.go|internal/adapters/db/connection.go" $COVERAGE_OUT > $FILTERED_OUT

# 3. Calculate total coverage percentage
TOTAL_COVERAGE=$(go tool cover -func=$FILTERED_OUT | grep total | awk '{print substr($3, 1, length($3)-1)}')

# 4. Clean up
# Keep $FILTERED_OUT if the user wants to see the detailed report
# rm -f $COVERAGE_OUT

# Output the result
echo "Total coverage: $TOTAL_COVERAGE%"

# Check against threshold (e.g., 80%)
THRESHOLD=80.0
is_below=$(echo "$TOTAL_COVERAGE < $THRESHOLD" | bc -l)

if [ "$is_below" -eq 1 ]; then
    echo "Error: Test coverage is $TOTAL_COVERAGE%, which is below the $THRESHOLD% threshold."
    rm -f $FILTERED_OUT $COVERAGE_OUT
    exit 1
fi

# If we are just showing the summary, we can remove the files
# but if the user ran task test:coverage-out, they expect coverage.out to exist.
# For now, let's rename FILTERED_OUT to coverage.out so 'go tool cover' works on it
mv $FILTERED_OUT coverage.out
