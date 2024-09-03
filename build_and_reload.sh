#!/bin/bash

# Build the project
make build

# Start the server in the background
./bin/cms &

# Get the PID of the server
SERVER_PID=$!

# Function to check if the server is up
check_server() {
    # Adjust the URL and port as needed
    curl -s http://localhost:4000 > /dev/null
    return $?
}

# Wait for the server to start (timeout after 30 seconds)
COUNTER=0
while ! check_server && [ $COUNTER -lt 30 ]; do
    sleep 1
    let COUNTER=COUNTER+1
done

if [ $COUNTER -lt 30 ]; then
    echo "Server is up, refreshing browser"
    # Run your browser refresh script
    ./script/reloadbrowser.sh
else
    echo "Server didn't start in time"
fi

# Wait for the server process to finish
wait $SERVER_PID