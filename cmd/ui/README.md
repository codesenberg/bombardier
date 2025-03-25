# Bombardier UI

A web-based user interface for Bombardier HTTP benchmarking tool.

## Requirements

- Go 1.18 or higher
- Bombardier command-line tool installed and available in PATH

## Building the UI

### Windows

```
build.bat
```

### Linux/macOS

```
./build.sh
```

If the build script is not executable, make it executable first:

```
chmod +x build.sh
```

## Running the UI

After building the UI, you can run it using:

### Windows

```
..\..\bombardier-ui.exe
```

### Linux/macOS

```
../../bombardier-ui
```

## Using the UI

1. Open your browser and navigate to http://localhost:8080
2. Configure your test parameters using the form
3. Click "Run Test" to execute the benchmark
4. View results in real-time
5. Save configurations for later use

## Features

- Configure all Bombardier options through a simple web interface
- Run tests and view results in real-time
- Save and load test configurations using browser local storage
- User-friendly visualization of test results 