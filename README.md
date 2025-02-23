# Hosts Profiles Manager

A simple and efficient GUI application to manage multiple hosts file profiles on your system. Switch between different hosts configurations easily without manual editing.

## Features

- 🔄 Switch between different hosts profiles with one click
- ✏️ Edit hosts files with a built-in editor
- 💾 Save and manage multiple hosts profiles
- 🔒 Automatic backup of original hosts file
- 🎨 Modern and intuitive user interface

## Prerequisites

To build from source, you need:

- Go 1.20 or later
- Fyne toolkit
- Git

### Installing Prerequisites

#### macOS
```bash
# Install Go
brew install go

# Install Fyne
go install fyne.io/fyne/v2/cmd/fyne@latest
```

#### Linux
```bash
# Install Go (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install golang-go

# Install Fyne
go install fyne.io/fyne/v2/cmd/fyne@latest
```

#### Windows
1. Download and install Go from [golang.org](https://golang.org)
2. Install Fyne:
```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```

## Building from Source

1. Clone the repository:
```bash
git clone https://github.com/jonaskahn/hosts-manager.git
cd hosts-manager
```

2. Build the application:

### macOS
```bash
chmod +x build/darwin.sh
./build/darwin.sh
```

The built application will be available in the `dist` folder as `Hosts Profiles Manager.app`.

### Linux (coming soon)
```bash
chmod +x build/linux.sh
./build/linux.sh
```

### Windows (coming soon)
```bash
build\windows.bat
```

## Running the Application

Since the application modifies system files, it requires administrative privileges to run.

### macOS
1. Right-click on `Hosts Profiles Manager.app`
2. Select "Open" from the context menu
3. Click "Open" in the security dialog
4. Enter your administrator password when prompted

Alternatively, from terminal:
```bash
sudo open "dist/Hosts Profiles Manager.app"
```

### Linux
Run the application with sudo:
```bash
sudo ./hosts-manager
```

### Windows
1. Right-click on `Hosts Profiles Manager.exe`
2. Select "Run as administrator"
3. Click "Yes" in the User Account Control (UAC) dialog

## Usage

1. Launch the application with administrative privileges (see above)
2. The app will automatically backup your current hosts file on first run
3. Create a new profile:
   - Click "New Profile"
   - Enter a profile name
   - Edit the hosts entries
   - Click "Save"

4. Switch between profiles:
   - Select a profile from the list
   - Click "Activate Profile"

5. Edit existing profiles:
   - Select a profile
   - Click "Edit"
   - Make your changes
   - Click "Save"

## File Locations

- macOS: Profiles are stored in `~/Library/Application Support/hosts-manager/`
- Linux: Profiles are stored in `~/.config/hosts-manager/`
- Windows: Profiles are stored in `%APPDATA%\hosts-manager\`

The system hosts file is located at:
- macOS & Linux: `/etc/hosts`
- Windows: `C:\Windows\System32\drivers\etc\hosts`

## Backup and Recovery

The application automatically creates a backup of your original hosts file before making any changes. The backup is stored in the application's data directory with the name `hosts.backup`.

To restore the original hosts file:
1. Open the application with administrative privileges
2. Click on "Settings"
3. Click "Restore Original Hosts"

## Troubleshooting

### Permission Denied Errors
- Make sure you're running the application with administrative privileges
- Check that your hosts file is not set to read-only
- Verify that no other application is currently accessing the hosts file

### macOS Security Warning
If you see "App can't be opened because it is from an unidentified developer":
1. Go to System Settings > Privacy & Security
2. Scroll down to Security
3. Click "Open Anyway"
4. Enter your administrator password when prompted

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please file an issue on the GitHub repository.
