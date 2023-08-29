# Blimp

Customizable terminal UI for monitoring weather information, application status, network latency, and more

_I built this to run on a RaspberryPi Zero inside a Macintosh 1 on my shelf to quickly see weather and uptime info._

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
# On smaller screens you probably don't want to cramp everything in, so remove the views that you don't need
layout = [
    ["weather", "weather"],
    ["weather", "weather"],
    ["feeds", "feeds"]
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
# You openweathermap.org API key
owm_api_key = ""
# The latitude and longitude of the weather location
owm_lat = 51.33
owm_lon = 12.37

[[views.feeds.targets]]
url = "https://hnrss.org/newest"
```

Look at the [default config](internal/config/default.config.toml) for other values you can override.
