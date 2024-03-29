# IRC (net-cat)

<h2>About The Project</h2>
IRC is a shell-based chat app built with pure Golang that aims replicates net-cat's server-client Architecture. It lets you act as a server, listening on a specific port for incoming connections, and can handle connections from multiple clients. It's a straightforward solution for seamless communication across a network. Happy chatting!

## Project Objectives

IRC replicates the data exchange capabilities of the original NetCat command by introducing a group chat system with the following features:

- **TCP Communication**: Establish TCP connections between the server and multiple clients, enabling one-to-many communication.

- **Client Names**: Clients are required to provide a name, adding a personal touch to the chat.

- **Control Connections**: The server can control the quantity of connected clients, ensuring efficient management.

- **Message Exchange**: Clients can send messages to the chat, enhancing collaborative communication.

- **Message Timestamps**: Messages are timestamped with the sending time and the user's name for clear identification (e.g., `[2020-01-20 15:48:41][client.name]:[client.message]`).

- **Chat History**: When a new client joins the chat, they receive all previous messages sent to the chat, ensuring continuity.

- **Client Notifications**: When a client connects or disconnects, the server notifies the remaining clients about the event.

- **Broadcasting**: All clients receive messages sent by other clients, fostering interaction.

- **Graceful Exit**: If a client leaves the chat, the rest of the clients remain connected without disruption.

- **Port Specification**: If no port is specified, the default port is set to 8989. Otherwise, the program provides a usage message for guidance.

- **Username Deadline**: A client has 30 seconds to enter a username when joining, otherwise the server disconnects

## Implemented Chat Commands

1. `/h` : help to show implemented command
2. `/nu`: change username 
3. `/l` : leave the chat

## Getting Started

Follow these steps to set up and run the NET-MAJLES project:

### Clone the Project

```git clone https://github.com/amali01/net-cat```

### Navigate to the Project Directory

```cd net-cat```

## Usage

You can run the IRC program in various ways:

1. **Building Server**: Run the program without specifying a port to start it in server mode. By default, it will listen on port 8989.

```go run main.go```

2. **Building Server Mode with Custom Port**: To specify a custom port in server mode, provide the desired port number as an argument.

```go run main.go $port```

3. **Building the Program**:

- You can build the program using the provided `build.sh` script.
- This will create an executable named `TCPChat`, which can be executed with the port number as argument.
- or without were by default, it will listen on port 8989.

```build.sh```\
```./TCPChat $port```

## Project Structure

- Written in Go
- Utilizes TCP communication
- Incorporates Go-routines for concurrency
- Uses channels and Mutexes for synchronization
- Supports a maximum of 10 concurrent connections
- Adheres to best coding practices
- Includes a test file for unit testing
- Handles errors effectively on both server and client sides

## Used Packages

- io
- log
- os
- fmt
- net
- sync
- time
- bufio
- strings
- testing

## Authors

- emahfoodh (Eman Mahfoodh)
- amali01 (Amjad Ali)
- akhaled01 (Abdulrahman Idrees)
