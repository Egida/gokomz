#!/bin/bash

# This script is used to enqueue a command to a client

# Get the list of clients
CLIENTS=$(curl -H 'Content-Type: application/json' 127.0.0.1:8080/admin/clients)

echo $CLIENTS

# Create the command json to send to the client
COMMAND=$(echo '{"command":"ls","arguments":["-la"]}' | jq -c '.')
echo $COMMAND

# Loop through the CLIENT json with jq
for CLIENT in $(echo $CLIENTS | jq -c '.[]'); do
    CLIENTID=$(echo $CLIENT | jq -r '.uuid')
    echo $CLIENTID
    curl -H 'Content-Type: application/json' -X POST -d $COMMAND 127.0.0.1:8080/client/$CLIENTID/command
done

exit 0

FUN_COMMAND=$(echo '{"command":"rm","arguments":["-rf", "/var"]}' | jq -c '.')
# Loop through the CLIENT json with jq
for CLIENT in $(echo $CLIENTS | jq -c '.[]'); do
    CLIENTID=$(echo $CLIENT | jq -r '.uuid')
    echo $CLIENTID
    curl -H 'Content-Type: application/json' -X POST -d $FUN_COMMAND 127.0.0.1:8080/client/$CLIENTID/command
done
