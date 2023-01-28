#!/usr/bin/env bash
#
#   Copyright The OpenTelemetry Authors.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
# Adds code owners without write access as reviewers on a PR. Note that
# the code owners must still be a member of the `open-telemetry`
# organization.
#
# Note that since this script is considered a requirement for PRs,
# it should never fail.

set -euo pipefail

if [[ -z "${REPO:-}" || -z "${PR:-}" ]]; then
    echo "One or more of REPO and PR have not been set, please ensure each is set."
    exit 0
fi

main () {
    CUR_DIRECTORY=$(dirname "$0")
    # Reviews may have comments that need to be cleaned up for jq,
    # so restrict output to only printable characters and ensure escape
    # sequences are removed.
    # The latestReviews key only returns the latest review for each reviewer,
    # cutting out any other reviews. We use that instead of requestedReviews
    # since we need to get the list of users eligible for requesting another
    # review. The GitHub CLI does not offer a list of all reviewers, which
    # is only available through the API. To cut down on API calls to GitHub,
    # we use the latest reviews to determine which users to filter out.
    JSON=$(gh pr view "${PR}" --json "files,author,latestReviews" | tr -dc '[:print:]' | sed -E 's/\\[a-z]//g')
    AUTHOR=$(printf "${JSON}"| jq -r '.author.login')
    FILES=$(printf "${JSON}"| jq -r '.files[].path')
    REVIEW_LOGINS=$(printf "${JSON}"| jq -r '.latestReviews[].author.login')
    COMPONENTS=$(bash "${CUR_DIRECTORY}/get-components.sh")
    REVIEWERS=""
    LABELS=""
    declare -A PROCESSED_COMPONENTS
    declare -A REVIEWED

    for REVIEWER in ${REVIEW_LOGINS}; do
        REVIEWED["@${REVIEWER}"]=true
    done

    for COMPONENT in ${COMPONENTS}; do
        # Files will be in alphabetical order and there are many files to
        # a component, so loop through files in an inner loop. This allows
        # us to remove all files for a component from the list so they
        # won't be checked against the remaining components in the components
        # list. This provides a meaningful speedup in practice.
        for FILE in ${FILES}; do
            MATCH=$(echo "${FILE}" | grep -E "^${COMPONENT}" || true)

            if [[ -z "${MATCH}" ]]; then
                continue
            fi

            # If we match a file with a component we don't need to process the file again.
            FILES=$(printf "${FILES}" | grep -v "${FILE}")

            if [[ -v PROCESSED_COMPONENTS["${COMPONENT}"] ]]; then
                continue
            fi

            PROCESSED_COMPONENTS["${COMPONENT}"]=true

            OWNERS=$(COMPONENT="${COMPONENT}" bash "${CUR_DIRECTORY}/get-codeowners.sh")

            for OWNER in ${OWNERS}; do
                # Users that leave reviews are removed from the "requested reviewers"
                # list and are eligible to have another review requested. We only want
                # to request a review once, so remove them from the list.
                if [[ -v REVIEWED["${OWNER}"] || "${OWNER}" = "@${AUTHOR}" ]]; then
                    continue
                fi

                if [[ -n "${REVIEWERS}" ]]; then
                    REVIEWERS+=","
                fi
                REVIEWERS+=$(echo "${OWNER}" | sed -E 's/@(.+)/"\1"/')
            done

            # Convert the CODEOWNERS entry to a label
            COMPONENT_NAME=$(echo "${COMPONENT}" | sed -E 's%^(.+)/(.+)\1%\1/\2%')

            if (( "${#COMPONENT_NAME}" > 50 )); then
                echo "'${COMPONENT_NAME}' exceeds GitHub's 50-character limit on labels, skipping adding label"
                continue
            fi

            if [[ -n "${LABELS}" ]]; then
                LABELS+=","
            fi
            LABELS+="${COMPONENT_NAME}"
        done
    done

    if [[ -n "${LABELS}" ]]; then
        gh pr edit "${PR}" --add-label "${LABELS}" || echo "Failed to add labels to #${PR}"
    else
        echo "No labels found"
    fi

    # Note that adding the labels above will not trigger any other workflows to
    # add code owners, so we have to do it here.
    #
    # We have to use the GitHub API directly due to an issue with how the CLI
    # handles PR updates that causes it require access to organization teams,
    # and the GitHub token doesn't provide that permission.
    # For more: https://github.com/cli/cli/issues/4844
    #
    # The GitHub API validates that authors are not requested to review, but
    # accepts duplicate logins and logins that are already reviewers.
    if [[ -n "${REVIEWERS}" ]]; then
        curl \
            -X POST \
            -H "Accept: application/vnd.github+json" \
            -H "Authorization: Bearer ${GITHUB_TOKEN}" \
            "https://api.github.com/repos/${REPO}/pulls/${PR}/requested_reviewers" \
            -d "{\"reviewers\":[${REVIEWERS}]}" \
            | jq ".message" \
            || echo "Failed to add reviewers to #${PR}"
    else
        echo "No code owners found"
    fi
}

main || echo "Failed to run $0"
