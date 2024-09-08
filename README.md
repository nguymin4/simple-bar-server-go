# <img src="https://github.com/Jean-Tinland/simple-bar-server/raw/main/images/logo-simple-bar-server.png" width="200" alt="simple-bar-server" />

## Overview

This is the server part of [https://github.com/nguymin4/simple-bar](https://github.com/nguymin4/simple-bar)

Architecture diagram:
```mermaid
flowchart LR
  Client-->http_server
  subgraph simple-bar-server
    http_server[HTTP Server]-->Websocket
    scheduled_task[Scheduled Task]-->Websocket
  end
  Websocket<-->uebersicht[Ubersicht Widget: simple-bar]
```

## Features

- Refresh, toggle, enable or disable simple-bar widgets
- Refresh AeroSpace spaces, windows and displays simple-bar widgets


## Installation

Clone this project to `~/.config/uebersicht`. This is important as currently we need to use python script to get app badges

```bash
git clone https://github.com/nguymin4/simple-bar-server-go.git ~/.config/uebersicht/
```

```
go run .
```

The process named `simple-bar-server`, this can be checked via `ps`, `pgrep` etc.
