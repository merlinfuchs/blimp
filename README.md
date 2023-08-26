# Blimp

Customizable terminal UI for monitoring weather information, application status, network latency, and more

![Example](example.png)

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
owm_api_key = "1dee9412f62ba03c40a23c6aa436710e"
owm_lat = 51.33
owm_lon = 12.37
owm_location = "Leipzig"
```
