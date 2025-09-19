<p align="center">
  <img src="https://github.com/coolapso/dygma-indicator/blob/main/img/Logo.png" width="200" >
</p>

# Dygma Indicator

[![Release](https://github.com/coolapso/dygma-indicator/actions/workflows/release.yaml/badge.svg?branch=main)](https://github.com/coolapso/dygma-indicator/actions/workflows/release.yaml)
![GitHub Tag](https://img.shields.io/github/v/tag/coolapso/dygma-indicator?logo=semver&label=semver&labelColor=gray&color=green)
[![Go Report Card](https://goreportcard.com/badge/github.com/coolapso/dygma-indicator)](https://goreportcard.com/report/github.com/coolapso/dygma-indicator)
![GitHub Sponsors](https://img.shields.io/github/sponsors/coolapso?style=flat&logo=githubsponsors)

A simple CLI utility to get the battery level of Dygma keyboards.

Because only one process can use the serial port at a time, this app is designed to run, get the battery level, print it to standard output, and exit immediately.

It is designed to work with status bars like `waybar`, but its simple JSON output makes it easy to integrate with other tools and scripts.

## Installation

### AUR

On Arch linux you can use the AUR `dygma-indicator-bin`

### Go Install

#### Latest version

`go install github.com/coolapso/dygma-indicator`

#### Specific version

`go install github.com/coolapso/dygma-indicator@v1.0.0`

### Linux Script

It is also impossible to install on any linux distro with the installation script

#### Latest version

```
curl -L https://dygma-indicator.coolapso.sh/install.sh | bash
```

#### Specific version

```
curl -L https://dygma-indicator.coolapso.sh/install.sh | VERSION="v1.1.0" bash
```

### Manual install

* Grab the binary from the [releases page](https://github.com/coolapso/dygma-indicator/releases).
* Extract the binary
* Execute it

## Usage

At the moment there's nothing speciall about it, just execute it and it will print the battery level to standard output.

```bash
dygma-indicator
```

### output


The application outputs a JSON object with the battery percentage and a corresponding text status.

```jsom
{"text":"L:50% R:70%","tooltip":"Left side: 50%\rRight side: 70%","percentage":50}
```

### Waybar

<p align="center">
  <img src="https://github.com/coolapso/dygma-indicator/blob/main/img/waybar.jpg">
</p>

To use this with `waybar`, add the following configuration to your `config` file.

```json
  "custom/keyboard": {
    "format": "{icon}   {text}",
    "return-type": "json",
    "interval": 3600,
    "format-icons": [
      "󰂃",
      "󰁻",
      "󰁾",
      "󰂀",
      "󰁹"
    ],
    "max-length": 40,
    "escape": true,
    "exec": "dygma-indicator"
  },
```

> [!IMPORTANT]
> The serial protocol used by the keyboard can only be used by one application at a time. It is recommended to set a reasonably high `interval` (e.g., 3600 seconds) to avoid blocking other applications like Bazecor when you need to use them.
