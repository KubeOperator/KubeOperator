# Verify Generated Modules

Pre-commit hook for verifying that generated library modules match
their EXPECTED content. Library modules are generated from fragments
under the `roles/lib_(openshift|utils)/src/` directories.

If the attempted commit modified files under the
`roles/lib_(openshift|utils)/` directories this script will run the
`generate.py --verify` command.

This script will **NOT RUN** if module source fragments are modified
but *not part of the commit*. I.e., you can still make commits if you
modified module fragments AND other files but are *not comitting the
the module fragments*.

# Setup Instructions

Standard installation procedure. Copy the hook to the `.git/hooks/`
directory and ensure it is executable.
