# Gocom

A management engine for remote machines. Uses a client server model.

## TODO

### Phase 1

- Return json instead of a string from the server
- Run a single, static command on the client
- Have client POST back response to the server

### Phase 2

- Support main command and args
- Be able to queue commands for a client
- Be able to handle many clients
- Be able to specify which client runs a command

### Phase 3

- Refactor server and client into seperate files
+ Refactor args parse into seperate files (TODO) <<<<<<<<<<<<
- Ensure client never dies

### Phase 4

- Support sending commands to a client via the c&c server
- Support getting client information from the API
- Test with multiple clients
+ Why is the jitter the same across both clients (TODO) <<<<<<<<<<<<


### Phase 5

+ Hook up c&c server to a MySQL DB
    - Get migrations working
    - Record client information
    - Record commands sent to a client
    - Record responses from those clients
    + Update test script to updated endpoint

### Phase 6

+ Run different types of commands: scripts, commands, new c&c, delete yourself
+ Refactor server to controller methods
+ modes: client, server, proxy, beacon
+ add authentication to admin endpoints
+ Use uuid for commands

## Features

+ Allow client to run fileless
+ Client<->Server authentication
    + tokens
    + mTLS
+ Support client groups
+ Allow the server to act as an install server (given a key)
+ Proxy functionality for a server
+ CLI that allows easy interation at a server level for sending commands to a single or groups of clients
+ Output the clients and their last communication
+ Configurable client heartbeats
+ Support various methods for client to server authentication
+ Persistance mechanisms for Linux, Mac, Windows
+ Each time compiled, new file hash
