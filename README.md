# File Storage Server

A complete TCP-based file storage system written in pure Go, including both server and client implementations. The system allows clients to perform basic file operations over a network connection using a clean command-line interface.

## Features

- **Send**: Upload files to the server
- **Get**: Download files from the server
- **List**: Retrieve a list of all files stored on the server
- **Delete**: Remove files from the server
- **Environment Configuration**: Configurable server settings via `.env` file
- **Command-line Client**: Easy-to-use CLI client with flag-based commands

## Project Structure

```
filer/
├── cmd/
│   ├── client/
│   │   └── client.go      # Command-line client application
│   └── server/
│       └── server.go      # Server application entry point
├── internal/
│   ├── client/
│   │   ├── client.go      # Client networking implementation
│   │   └── msg.go         # Message handling with goroutines
│   └── server/
│       └── server.go      # Server networking implementation
├── .env                   # Environment configuration
├── go.mod                 # Go module definition
├── README.md             # This file
└── storage/              # Directory for stored files (created automatically)
```

There are also separate branches that contain only server files and only client files.

## Installation

1. Ensure you have [Go](https://golang.org/doc/install) installed (version 1.22+ required)
2. Clone the repository:
   ```bash
   git clone https://github.com/nikaydo/fileStorage.git
   cd fileStorage
   ```

## Usage

### Starting the Server

Run the server with:
```bash
go run cmd/server/server.go
```

The server will:
- Start listening on `localhost:9000`
- Create a `storage/` directory if it doesn't exist
- Accept client connections

### Protocol

The server uses a simple binary protocol over TCP:

1. **Message Format**: Each message is prefixed with an 8-byte length header (big-endian uint64) followed by the data
2. **Commands**: 
   - `get <filename>` - Download a file
   - `list` - Get list of files
   - `delete <filename>` - Delete a file
   - `send <filename>` - Upload a file

## Configuration

The server uses environment variables for configuration, which can be set via the `.env` file:

- **ADDR**: Server address (default: localhost)
- **PORT**: Server port (default: 9000)
- **STORAGEPATH**: Directory for file storage (default: storage)

### Using the Command-Line Client

The project includes a built-in client that you can use to interact with the server:

#### Upload a file to the server:
```bash
go run cmd/client/client.go -mode=send -file=path/to/your/file.txt
```

#### Download a file from the server:
```bash
go run cmd/client/client.go -mode=get -file=filename.txt
```

#### List all files on the server:
```bash
go run cmd/client/client.go -mode=list
```

#### Delete a file from the server:
```bash
go run cmd/client/client.go -mode=delete -file=filename.txt
```

**Note**: The client uses the same environment variables as the server for connection settings. Make sure the `.env` file is properly configured.

### Protocol Details

The server uses a simple binary protocol over TCP:

1. **Message Format**: Each message is prefixed with an 8-byte length header (big-endian uint64) followed by the data
2. **File Transfer**: Files are transferred with filename followed by file size header and raw file data
3. **Commands**: 
   - `send` - Upload a file (client sends filename + file data)
   - `get <filename>` - Download a file from server
   - `list` - Get list of all stored files
   - `delete <filename>` - Remove a file from server

## Development

### Running the Server
```bash
go run cmd/server/server.go
```

### Running the Client
```bash
go run cmd/client/client.go -mode=<command> [-file=<filename>]
```

### Building Binaries
```bash
# Build server
go build -o filer-server cmd/server/server.go

# Build client  
go build -o filer-client cmd/client/client.go
```

## Architecture

- **Server**: Handles TCP connections and implements file storage.
- **Client**: Provides as a software API.
- **Storage**: Files are stored in a separate directory in the program root.
- **Environment Loading**: Loads environment variables using no third-party libraries.
- **Binary Protocol**: Efficient binary communication with length-prefixed messages.

## Technical Features

- **Custom Environment Loader**: Built-in `.env` file parsing without external dependencies
- **Concurrent Message Handling**: Goroutine-based message processing for list/delete operations
- **File Streaming**: Efficient file transfer using 4KB buffer chunks
- **Error Handling**: Comprehensive error handling with proper client feedback
- **Path Safety**: Uses `filepath.Base()` to prevent directory traversal attacks
