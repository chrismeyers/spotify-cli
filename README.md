# spotify-cli

A command-line interface (CLI) tool for interacting with the Spotify Web API.
Use it to search Spotify from your terminal.

## Getting Started

### 1. Get Spotify API Credentials

To use this tool, you need a Spotify Developer account:

1. Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications).
2. Log in and create a new application.
3. Copy your **Client ID** and **Client Secret**.

### 2. Create `config.json`

The program will look for your `config.json` in several common locations.
You can place it in any of the following paths:

- `$XDG_CONFIG_HOME/spotify-cli/config.json`
- `$HOME/.config/spotify-cli/config.json`
- `./config.json` (current directory)

Create the file with the following structure:

```json
{
  "api": {
    "clientId": "YOUR_CLIENT_ID",
    "clientSecret": "YOUR_CLIENT_SECRET"
  }
}
```

Replace `YOUR_CLIENT_ID` and `YOUR_CLIENT_SECRET` with your actual credentials.

### 3. Build and Run

Make sure you have Go installed (1.24+ recommended):

```sh
go build -o spotify-cli
./spotify-cli
```

or

```sh
go run .
```
