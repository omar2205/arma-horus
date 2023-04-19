# Arma Horus

Arma Horus is a tool that controls Arma 3 Server instances. It allows you to start, stop, restart, update, and monitor your server with ease.

## Features

    [ ] Supports multiple server instances with different configurations and mods
    [ ] Automatically updates the server and mods when needed
    [ ] Provides a web interface for managing and monitoring the server
    [ ] Sends notifications via email or Discord when the server status changes
    [ ] Logs server events and performance metrics
    [ ] Supports custom commands and scripts

## Installation

To install Arma Horus, you need to have Go installed on your system. You can download it from [here](https://golang.org/dl/).

Then, you can clone this repository and build the executable:

```bash
git clone https://github.com/omar2205/arma-horus.git
cd arma-horus
make build/api
```

## Usage

To use Arma Horus, you need to create a configuration file for each server instance you want to control. The configuration file is a JSON file that specifies the server settings, mods, and other options.

```json
{
  "server_folder": "change_me",
  "server_script": "change_me",
  "server_pid_file": "change_me",
  "db_file": "./db.sqlite",
  "admins_email": ["change_me"]
}
```

Arma Horus also support these command line args:
```
  -cors-trusted-origins value
        Trust CORS origins (space seperated)
  -env string
        Environment (dev|staging|prod) (default "dev")
  -limiter-burst int
        Rate limiter maximum burst (default 4)
  -limiter-enabled
        Enable rate limiter (default true)
  -limiter-rps float
        Rate limiter maximum requests per second (default 2)
  -port int
        API server port (default 3000)
  -version
        Display version and exit
```

## License

Arma Horus is licensed under the MIT License. See LICENSE for details.
