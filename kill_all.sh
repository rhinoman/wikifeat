#!/bin/sh

# Find all running Wikifeat services and kill them
pgrep wikifeat | xargs kill
