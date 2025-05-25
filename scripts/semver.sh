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

# If it is tag - use it
if [ "$CURRENT_BRANCH" = "HEAD" ]; then
    # On tag, output pure semver
    echo $VERSION
else
    # Use distance as patch version on all branches
    MAJOR=$(echo $VERSION | sed -E 's/^v?([0-9]+)\..*/\1/')
    MINOR=$(echo $VERSION | sed -E 's/^v?[0-9]+\.([0-9]+)\..*/\1/')
    PATCH=$DISTANCE
    if [ -z "$PATCH" ] || [ "$PATCH" = "" ]; then
        PATCH=0
    fi
    if [ "$CURRENT_BRANCH" = "$DEFAULT_BRANCH" ]; then
        # If it is default branch - use distance as patch version
        echo "$MAJOR.$MINOR.$PATCH+$SHORT_COMMIT"
    else
        # If it is not default branch - increment patch version by 1 and use branch name as RC
        PATCH=$(($PATCH + 1))
        BRANCH_ALNUM=$(echo "$CURRENT_BRANCH" | tr -cd '[:alnum:]')
        echo "$MAJOR.$MINOR.$PATCH-$BRANCH_ALNUM+$SHORT_COMMIT"
    fi
fi