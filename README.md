# Data Store Server

An application written in Go that will take a JSON file and serve the contents over APIs.

## Prerequisites
If you intend to compile and run the raw code, you must have Go installed.

## Usage

### Running Pre-Compiled Binary (tested on MacOS only)
```
./bin/data-store-server <optional_flags>
```

### Manual Compile and Run
```
go run main.go <optional_flags>
```

### Compile Into New Binary (tested on MacOS)
```
go build -o bin/data-store-server main.go
```

## Optional Flags

You can combine any of the following flags in order to customize your application usage.

Name | Description | Usage | Default Value
---- | ----------- | ----- | -------------
port | port the HTTP server runs on | --port=8081 | 8081
file | the path to the file you want to load | --file=data/somefile.json | data/data.json

## HTTP APIs

*Assumes HTTP server is running on port 8081*

### Retrieve the full list of key/values
http://localhost:8081

### Retrieve a single key's value
http://localhost:8081?name=<key>