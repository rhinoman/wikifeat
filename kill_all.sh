#!/bin/sh

# Find all running Wikifeat services and kill them
pgrep wikis | xargs kill
pgrep users | xargs kill
pgrep notifications | xargs kill
pgrep frontend | xargs kill
