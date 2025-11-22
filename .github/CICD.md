# GitHub Actions CI/CD Setup

## Overview

GitHub Actions workflow configured in `.github/workflows/build.yml` to automatically build mlr-desktop for multiple platforms.

## Supported Platforms

The workflow builds for the following platforms:

1. **Linux** (x86_64)
   - Runner: `ubuntu-latest`
   - Dependencies: libwebkit2gtk-4.1-dev, libgtk-3-dev, build-essential
   - Build flags: `-tags webkit2_41`

2. **macOS Intel** (x86_64)
   - Runner: `macos-latest`
   - Platform: `darwin/amd64`
   - No additional dependencies (uses system WebKit)

3. **macOS Apple Silicon** (ARM64)
   - Runner: `macos-latest`
   - Platform: `darwin/arm64`
   - No additional dependencies (uses system WebKit)

4. **Windows** (x86_64)
   - Runner: `windows-latest`
   - Platform: `windows/amd64`
   - No additional dependencies (uses WebView2)

## Workflow Triggers

The workflow runs on:
- **Push to main branch** - Builds all platforms, uploads artifacts
- **Pull requests** - Builds all platforms for testing
- **Tag push** (v*) - Builds all platforms and creates a GitHub Release
- **Manual trigger** - Can be triggered manually from Actions tab

## Build Process

Each platform build follows these steps:

1. **Checkout code** - Get the latest source
2. **Setup Go** (v1.21) - Install Go compiler
3. **Setup Node.js** (v20) - Install Node for frontend build
4. **Install platform dependencies** - OS-specific libraries
5. **Install Wails CLI** - Build tool for the desktop app
6. **Install frontend dependencies** - npm install
7. **Build application** - Run `wails build`
8. **Upload artifacts** - Store binaries for 30 days

## Release Process

When you push a version tag:

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

The workflow will:
1. Build for all platforms (Linux, macOS Intel, macOS ARM, Windows)
2. Create a GitHub Release
3. Attach all platform binaries to the release
4. Generate release notes automatically

## Artifacts

Build artifacts are stored for 30 days and include:

- `mlr-desktop-linux` - Linux executable
- `mlr-desktop-darwin-amd64.zip` - macOS Intel .app bundle (zipped)
- `mlr-desktop-darwin-arm64.zip` - macOS ARM .app bundle (zipped)
- `mlr-desktop-windows` - Windows executable (.exe)

**Note**: macOS builds are packaged as .zip files containing the `.app` bundle, which is the standard macOS application format. Users should extract the .zip and drag the .app to their Applications folder.

## Binary Sizes

Expected binary sizes (approximate):
- Linux: ~40 MB (includes Miller library)
- macOS: ~40 MB (includes Miller library)
- Windows: ~40 MB (includes Miller library and WebView2 runtime)

The larger size is due to the embedded Miller library, making the application completely self-contained with no external dependencies.

## Downloading Binaries

Users can download pre-built binaries from:
1. The [Releases](https://github.com/USERNAME/mlr-desktop/releases) page (for tagged releases)
2. The [Actions](https://github.com/USERNAME/mlr-desktop/actions) tab (for development builds)

## Troubleshooting

### macOS Security

On macOS, downloaded binaries may be blocked by Gatekeeper. Users should:
```bash
xattr -cr mlr-desktop
```

### Linux Permissions

On Linux, make the binary executable:
```bash
chmod +x mlr-desktop
```

### Windows SmartScreen

Windows may show a SmartScreen warning for unsigned executables. Users can click "More info" → "Run anyway".

## Future Improvements

Potential enhancements:
- Code signing for macOS and Windows
- Notarization for macOS
- Creating installers (.dmg for macOS, .msi for Windows, .deb/.rpm for Linux)
- Automated testing on each platform
- Performance benchmarks
