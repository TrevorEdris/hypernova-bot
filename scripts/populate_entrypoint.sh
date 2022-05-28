#!/usr/bin/env bash

# Wait for the localstack container's dynamodb service to become available
cmd="aws dynamodb list-tables --endpoint=${DYNAMODB_ENDPOINT}"

total_attempts=10
attempts_left=${total_attempts}
while [ ${attempts_left} -gt 0 ]; do
    echo "Attempting to connect to localstack... ${attempts_left} attempts left"

    # Attempt to run the command
    $cmd
    cmd_status=$?
    if [ ! $cmd_status -eq 0 ]; then
        ((attempts_left--))
        if [ ${attempts_left} -eq 0 ]; then
            echo "Unable to connect after ${total_attempts} attempts. Exiting."
            exit 1
        fi
    else
        echo "Connected!"
        attempts_left=0
    fi
done

# If the populate_data.sh script fails, we want the entire script to fail.
# This is not set at the top because we are attempting to run $cmd
# multiple times, expecting it to fail until the localstack container
# is actually ready.
set -euo pipefail

sh /src/scripts/populate_data.sh
tail -f /dev/null