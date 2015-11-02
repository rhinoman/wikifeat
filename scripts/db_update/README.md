Database Update Scripts
=======================

This directory contains scripts to update CouchDB database design documents in pre-existing Wikifeat databases.
The scripts require Python 3 to be installed on the target system.

## Naming Convention

Scripts are named thusly: update_{FROM_VERSION}_to_{TO_VERSION}.py

## Usage

Simply run the script that corresponds to the Wikifeat version(s) you are upgrading from/to.  These scripts update
the database(s) one 'step' at a time, so if you haven't upgraded in a while and missed a few updates, you will need
to run each upgrade script sequentially.

