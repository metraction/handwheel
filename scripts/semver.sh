#!/bin/bash

DEFAULT_BRANCH=main
# Get describe info
DESCRIBE=$(git describe --tags --long --always)
# Example: v1.2.3-5-gabc1234
VERSION=$(echo "$DESCRIBE" | sed -E 's/^((v[0-9]+\.[0-9]+\.[0-9]+)).*/\1/')
COMMITS=$(echo "$DESCRIBE" | sed -E 's/^v[0-9]+\.[0-9]+\.[0-9]+-([0-9]+)-g[0-9a-f]+$/\1/')
SHORT_COMMIT=$(echo "$DESCRIBE" | sed -E 's/.*-g([0-9a-f]+)$/\1/')

# Get the current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Remove non-alphanumeric chars from branch name
BRANCH_ALNUM=$(echo "$CURRENT_BRANCH" | tr -cd '[:alnum:]')

# Print version based on branch
if [ "$CURRENT_BRANCH" = "$DEFAULT_BRANCH" ]; then
    if [ "$COMMITS" = "0" ] || [ -z "$COMMITS" ]; then
        echo $VERSION
    else
        echo $VERSION-$COMMITS.g$SHORT_COMMIT
    fi
else
    echo ${VERSION}-${BRANCH_ALNUM}-${COMMITS}+${SHORT_COMMIT}
fi