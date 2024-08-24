#!/bin/sh

PUID=${PUID:-1000}
PGID=${PGID:-1000}

# Create the group and user with the specified PUID and PGID
addgroup -g "$PGID" gitgroup
adduser -u "$PUID" -G gitgroup -s /bin/sh -D gituser

# Set ownership of the directories
chown -R gituser:gitgroup /backups /git-sync

# Switch to the new user and execute the command
exec su-exec gituser "$@"
