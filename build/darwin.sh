#!/bin/zsh

# Exit on error
set -e

# Check if fyne is installed
if ! command -v fyne &> /dev/null; then
    echo "Error: fyne is not installed. Please install it with: go install fyne.io/fyne/v2/cmd/fyne@latest"
    exit 1
fi

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo "Error: go is not installed. Please install Go first."
    exit 1
fi

echo "Creating dist directory..."
mkdir -p "dist"

echo "Building executable..."
go build -o "dist/hosts-manager"

echo "Packaging application..."
fyne package \
    --os darwin \
    --icon "resources/icon.png" \
    --name "Hosts Profiles Manager" \
    --appID "io.github.jonaskahn.hostsmanager" \
    --use-raw-icon \
    --executable "dist/hosts-manager" \
    --release

echo "Moving .app to dist folder..."
mv "Hosts Profiles Manager.app" "dist/Hosts Profiles Manager.app"

echo "Cleaning up..."
rm "dist/hosts-manager"

echo "Build completed successfully!"
