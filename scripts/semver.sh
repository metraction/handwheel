#!/bin/bash

DEFAULT_BRANCH=$(git remote show origin | sed -n '/HEAD branch/s/.*: //p')
# Get the latest tag
LATEST_TAG=$(git describe --tags --abbrev=0)

# Get the latest tag's version number
VERSION=$(echo $LATEST_TAG | grep -oP 'v\d+\.\d+\.\d+')

# Get the latest tag's commit hash
COMMIT=$(git rev-list -n 1 $LATEST_TAG)

# Get the number of commits since the latest tag
COMMITS=$(git rev-list --count $COMMIT..HEAD)

# Get the current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Remove non-alphanumeric chars from branch name
BRANCH_ALNUM=$(echo "$CURRENT_BRANCH" | tr -cd '[:alnum:]')

# Print version based on branch
if [ "$CURRENT_BRANCH" = "$DEFAULT_BRANCH" ]; then
    if [ $COMMITS -eq 0 ]; then
        echo $VERSION
    else
        echo $VERSION-$COMMITS.g$(git rev-parse --short HEAD)
    fi
else
    echo ${VERSION}-${BRANCH_ALNUM}-${COMMITS}+$(git rev-parse --short HEAD)
fi