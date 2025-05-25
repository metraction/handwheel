#!/bin/bash

DEFAULT_BRANCH=main
# Get describe info
DESCRIBE=$(git describe --tags --long --always)
# Example: v1.2.3-5-gabc1234
VERSION=$(echo "$DESCRIBE" | sed -E 's/^v([0-9]+\.[0-9]+\.[0-9]+).*/\1/')
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  VERSION="0.0.0"
fi
DISTANCE=$(echo "$DESCRIBE" | sed -E 's/^v?[0-9]+\.[0-9]+\.[0-9]+-([0-9]+)-g[0-9a-f]+$/\1/')
SHORT_COMMIT=$(echo "$DESCRIBE" | sed -E 's/.*-g([0-9a-f]+)$/\1/')

# Get the current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Remove non-alphanumeric chars from branch name
BRANCH_ALNUM=$(echo "$CURRENT_BRANCH" | tr -cd '[:alnum:]')

# Print version based on branch or tag
if [ "$CURRENT_BRANCH" = "HEAD" ]; then
    # On tag, output pure semver
    echo $VERSION
elif [ "$CURRENT_BRANCH" = "$DEFAULT_BRANCH" ]; then
    if [ "$DISTANCE" = "0" ] || [ -z "$DISTANCE" ]; then
        echo $VERSION
    else
        echo $VERSION-$DISTANCE+$SHORT_COMMIT
    fi
else
    # Increment patch version by 1
    PATCH=$(echo $VERSION | sed -E 's/^v?([0-9]+)\.([0-9]+)\.([0-9]+)$/\3/')
    PATCH=$((PATCH + 1))
    MAJOR=$(echo $VERSION | sed -E 's/^v?([0-9]+)\..*/\1/')
    MINOR=$(echo $VERSION | sed -E 's/^v?[0-9]+\.([0-9]+)\..*/\1/')
    NEW_VERSION="$MAJOR.$MINOR.$PATCH"
    echo ${NEW_VERSION}-${BRANCH_ALNUM}-${DISTANCE}+$SHORT_COMMIT
fi