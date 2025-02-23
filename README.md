# Hosts Profiles Manager

A GUI tool to manage and switch between different hosts file configurations.

## Install Requirements

- Go 1.20+
- Fyne: `go install fyne.io/fyne/v2/cmd/fyne@latest`

## Build

```bash
# macOS
./build/darwin.sh

# Linux
./build/linux.sh

# Windows
build\windows.bat
```

Built app will be in `dist` folder.

## Run

The app needs admin rights to modify hosts file.

### macOS
```bash
sudo open "dist/Hosts Profiles Manager.app"
```

### Linux
```bash
sudo ./hosts-manager
```

### Windows
Right-click > Run as administrator

## Hosts File Location

- macOS & Linux: `/etc/hosts`
- Windows: `C:\Windows\System32\drivers\etc\hosts`

## License

MIT License - see [LICENSE](LICENSE) file
