#!/bin/bash

PLAYERS=$(docker exec xenodochial_spence rcon /players)
ONLINE_COUNT=$(echo "$PLAYERS" | grep -c " (online)$")

if [[ "$ONLINE_COUNT" -gt "0" ]]; then
    echo "$PLAYERS"
    # exit with 75 (EX_TEMPFAIL) for watchtower
    # https://containrrr.dev/watchtower/lifecycle-hooks/
    exit 75
fi
