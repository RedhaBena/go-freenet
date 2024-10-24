# Freenet Client - Practical Work MA_ParaDis

## Overview

This project implements a **Freenet Client** for the **MA_ParaDis** course. The application provides a simulation of a distributed storage system, allowing clients to interact with a warehouse, search for files, and log network activity. This client is built using **Go** with a **Text-Based User Interface (TUI)** using the `tview` library.

## Demo

To quickly start the application, you can navigate to the `demo` folder and run the provided scripts:

```bash
cd demo
```

- **Windows**: Run `start_freenet.bat`:
```bash
start_freenet.bat
```
- **Linux/Mac**: Run `start_freenet.sh`:
```bash
chmod 777 ./start_freenet.sh
./start_freenet.sh
```

These scripts will start multiple Freenet clients with their configurations.

## Usage

### Command-Line Interface (CLI)

You can run the Freenet client using the following command:

```bash
go run .
```

To see the available options and commands, use:

```bash
go run . -h
```

This will display:

```bash
NAME:
   freenet - Freenet client for practical work MA_ParaDis

USAGE:
   freenet [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

   LOGS

   --debug  debug logs (default: false) [$DEBUG]

   NETWORK

   --address value  network address (default: "127.0.0.1") [$ADDRESS]
   --port value     network port (default: 43210) [$PORT]

   WAREHOUSE

   --warehouse value  warehouse file path (default: "warehouse.yaml") [$WAREHOUSE]
```

## Example Usage

1. Run the client with default settings:

```bash
go run .
```

2. Specify a custom network address and port:

```bash
go run . --address 192.168.1.10 --port 12345
```

3. Enable debug logs:

```bash
go run . --debug
```

4. Use a custom warehouse file:
```bash
go run . --warehouse mywarehouse.yaml
```

## Warehouse File

The **warehouse.yaml** file stores the mapping between file IDs and their locations. Each file is represented by a key (file ID) and its corresponding location (local or a network address). Here's an example of a typical `warehouse.yaml` file:

```yaml
files:
  55: local              # File 55 is stored locally
  10: 127.0.0.1:43214    # File 10 is located at the node running on address 127.0.0.1, port 43214
  20: 127.0.0.1:43212    # File 20 is located at the node running on address 127.0.0.1, port 43212
```

In this warehouse file:

- **File 55** is stored locally on the current node.
- **File 10** and **File 20** are located at different nodes within the network, specified by their respective IP addresses and ports.

The warehouse file is essential for keeping track of the file locations and ensuring efficient file retrieval when a search request is made.

### Global Options

- **--debug**: Enable detailed logging for debugging purposes.
- **--address**: Set the network address for communication (default is `127.0.0.1`).
- **--port**: Set the network port for communication (default is `43210`).
- **--warehouse**: Specify the path to the warehouse file (default is `warehouse.yaml`).
