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

Create a file named `config.json` in the project root with the following structure:

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
go run .
```

or

```sh
go build -o spotify-cli
./spotify-cli
```
