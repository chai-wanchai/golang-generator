## Install gitlint

```sh
brew install gitlint       # Homebrew (macOS)
sudo port install gitlint  # Macports (macOS)
apt-get install gitlint    # Ubuntu
```

## View Test Coverage Each Line in vs code

-   Add .vscode/settings.json

```json
{
    "go.coverOnSave": true,
    "go.coverOnSingleTest": true,
    "go.coverageDecorator": {
        "type": "gutter",
        "coveredHighlightColor": "rgba(64,128,128,0.5)",
        "uncoveredHighlightColor": "rgba(128,64,64,0.25)",
        "coveredGutterStyle": "blockgreen",
        "uncoveredGutterStyle": "blockred"
    }
}
```

## Run Application local develop

```sh
    make watch
```

## Run Application with docker

```sh
    docker compose -f docker-compose.yml up -d
```

