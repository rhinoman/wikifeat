#!/bin/sh

# Find all running Wikifeat services and kill them
pgrep wikifeat-wikis | xargs kill
pgrep wikifeat-users | xargs kill
pgrep wikifeat-notifications | xargs kill
pgrep wikifeat-frontend | xargs kill
pgrep wikifeat-auth | xargs kill
