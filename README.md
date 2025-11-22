# MLR Desktop Tool

A user-friendly desktop application for transforming data using the [Miller (mlr)](https://github.com/johnkerl/miller) command-line tool. This GUI application makes it easy to experiment with mlr transformations interactively.

## Features

- **Interactive Preview**: See real-time output as you build your transformation pipeline
- **Sample Data**: Comes pre-loaded with sample grocery data to help you get started
- **Multiple Input Formats**: Support for CSV, TSV, JSON, and NDJSON
  - CSV/TSV-specific options: Ragged, Headerless, Custom field separator
- **Multiple Output Formats**: Pretty Print, CSV, TSV, JSON, NDJSON
- **Verb Pipeline Builder**: Chain multiple mlr verbs with drag-and-drop reordering
- **Quick Add Shortcuts**: Common transformation patterns available with one click
  - Head 5 lines
  - Clean headers (replace spaces with underscores)
  - Filter by column value
  - Label columns
  - Cut columns
  - Add computed columns (split/extract)
- **Command Preview**: See the exact mlr command that will be executed
- **Save Output**: Export transformed data to a file
- **Auto-save**: Your work is automatically saved between sessions
- **File Input**: Load data from files or paste it directly

## Prerequisites

Before building or running this application, you need:

1. **Go** (1.19 or later) - [Install Go](https://golang.org/doc/install)
2. **Node.js** (14 or later) - [Install Node.js](https://nodejs.org/)
3. **Wails CLI** - Install with: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
4. **Linux Development Dependencies**:
   - WebKit2GTK: `sudo apt install libwebkit2gtk-4.1-dev` (Ubuntu 22.04+) or `libwebkit2gtk-4.0-dev` (older versions)
   - Build tools: `sudo apt install build-essential`

> [!NOTE]
> The Miller data transformation library is embedded directly in the application - no external `mlr` binary installation is required!

## Building

### Development Build

To build the application for development and testing:

```bash
# For Ubuntu 22.04+ (WebKit 4.1)
wails build -tags webkit2_41

# For older Ubuntu/Debian (WebKit 4.0)
wails build -tags webkit2
```

The compiled binary will be available at `build/bin/mlr-desktop`.

### Production Build

For a production build with optimizations:

```bash
wails build -tags webkit2_41
```

## Running

After building, run the application:

```bash
./build/bin/mlr-desktop
```

## Development

To run in live development mode with hot reload:

```bash
wails dev -tags webkit2_41
```

This will start a development server with your Go backend and a Vite frontend server for fast iteration.

## Usage

1. **Input Data**: 
   - Use the pre-loaded sample data, or
   - Paste your own data in the text area, or
   - Click "File Path" and select a file from your system

2. **Configure Input Format**: 
   - Select the input format (CSV, TSV, JSON, NDJSON) from the dropdown
   - For CSV/TSV, enable optional flags like "Ragged" or "Headerless" if needed

3. **Add Transformations**:
   - Use the "Quick Add" shortcuts for common operations, or
   - Type custom mlr verbs in the text field (e.g., `cut -f SKU,Price`)
   - Reorder verbs by clicking the ▲/▼ buttons
   - Enable/disable verbs with checkboxes

4. **Configure Output Format**:
   - Choose your desired output format from the dropdown

5. **Preview & Save**:
   - The output updates automatically as you make changes
   - Click "Save to File" to export the transformed data
   - Copy the generated command to use mlr directly in your terminal

## CI/CD and Downloads

This project uses GitHub Actions to automatically build executables for multiple platforms:

- **Linux** (x86_64)
- **macOS** (Intel and Apple Silicon)
- **Windows** (x86_64)

Builds are automatically triggered on:
- Push to `main` branch
- Pull requests
- Version tags (e.g., `v1.0.0`)

### Downloading Pre-built Binaries

1. Go to the [Releases](../../releases) page
2. Download the appropriate binary for your platform:
   - `mlr-desktop-linux` - Linux x86_64
   - `mlr-desktop-darwin-amd64` - macOS Intel
   - `mlr-desktop-darwin-arm64` - macOS Apple Silicon (M1/M2/M3)
   - `mlr-desktop-windows` - Windows x86_64
3. Extract and run the executable

### Creating a Release

To create a new release with pre-built binaries:

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

The GitHub Actions workflow will automatically build for all platforms and create a GitHub release with the binaries attached.

## Project Structure

- `app.go` - Go backend using Miller library directly for data transformations
- `frontend/src/App.jsx` - Main React application
- `frontend/src/components/` - React components
  - `InputSection.jsx` - Input data and format configuration
  - `VerbBuilder.jsx` - Transformation pipeline builder
  - `OutputPreview.jsx` - Output display and export
- `.github/workflows/build.yml` - GitHub Actions CI/CD for multi-platform builds

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Third-Party Dependencies

This project uses the following tools and frameworks:
- [Wails](https://wails.io/) - MIT License
- [Miller (mlr)](https://github.com/johnkerl/miller) - BSD 2-Clause License

Please see their respective licenses for details.
