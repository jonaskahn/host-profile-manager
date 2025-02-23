@echo off

REM Create dist directory if it doesn't exist
if not exist dist mkdir dist

REM First build the executable
go build -o dist\hosts-manager.exe

REM Package the application
%USERPROFILE%\go\bin\fyne package ^
  --os windows ^
  --icon resources\icon.png ^
  --name "Hosts Manager" ^
  --appID "com.jonas.hostsmanager" ^
  --executable dist\hosts-manager.exe ^
  --release

REM Move the packaged application to dist folder
move "Hosts Manager.exe" "dist\Hosts Manager.exe"

REM Clean up the temporary executable
del dist\hosts-manager.exe

echo Build completed successfully!
