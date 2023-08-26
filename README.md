# Blimp

Customizable terminal UI for monitoring weather information, application status, network latency, and more

![Example](example.png)

## Features

- **Weather**: Display current weather information and forecast
- **Network Latency**: Display a live chart of the network latency / ping
- **Application Status**: Monitor the status of web applications

## Installation

```shell
# Install from GitHub
go install github.com/merlinfuchs/blimp
# Run blimp
blimp
```

## Configuration

The app will look for a configuration file called `blimp.toml`. Here is an example configuration for the example above:

```toml
# The layout is based on a grid, you can add rows and columns or remove some widgets
layout = [
    ["weather", "weather"],
    ["weather", "weather"],
    ["latency", "status"]
]

[[views.status.targets]]
name = "Xenon Bot"
type = "https"
host = "xenon.bot"

[[views.status.targets]]
name = "Embed Generator"
type = "https"
host = "message.style"

[[views.status.targets]]
name = "Friendly Captcha API"
type = "https"
host = "eu-api.friendlycaptcha.eu"

[views.weather]
owm_api_key = ""
owm_lat = 51.33
owm_lon = 12.37
owm_location = "Leipzig"
```
