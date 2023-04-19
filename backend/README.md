# Arma Horus

Arma Horus is a tool that controls Arma 3 Server instances. It allows you to start, stop, restart, update, and monitor your server with ease.

## Features

- Supports multiple server instances with different configurations and mods
- Automatically updates the server and mods when needed
- Provides a web interface for managing and monitoring the server
- Sends notifications via email or Discord when the server status changes
- Logs server events and performance metrics
- Supports custom commands and scripts

## Installation

To install Arma Horus, you need to have Go installed on your system. You can download it from [here](https://golang.org/dl/).

Then, you can clone this repository and build the executable:

```bash
git clone https://github.com/arma-horus/arma-horus.git
cd arma-horus
go build
Alternatively, you can download the latest release from here.

Usage
To use Arma Horus, you need to create a configuration file for each server instance you want to control. The configuration file is a JSON file that specifies the server settings, mods, and other options. You can find an example configuration file here.

Then, you can run Arma Horus with the following command:

./arma-horus -config config.json
This will start the server and the web interface. You can access the web interface at http://localhost:8080 by default.

You can also use the following flags to customize the behavior of Arma Horus:

-port: The port number for the web interface (default: 8080)
-log: The log level for Arma Horus (default: info)
-debug: Enable debug mode for Arma Horus (default: false)
For more information, you can run:

./arma-horus -help
License
Arma Horus is licensed under the MIT License. See LICENSE for details.
