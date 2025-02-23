package main

import (
	"hosts-manager/models"
	"hosts-manager/views/management"
	"hosts-manager/views/profile"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

func getPrivilegeInstructions() string {
	switch runtime.GOOS {
	case "windows":
		return "This application needs administrator privileges to modify the hosts file.\n\n" +
			"Please right-click on the application and select 'Run as administrator'."
	case "darwin":
		return "This application needs administrator privileges to modify the hosts file.\n\n" +
			"Please run the application with sudo:\n" +
			"sudo /Applications/hosts-manager.app/Contents/MacOS/hosts-manager"
	case "linux":
		return "This application needs administrator privileges to modify the hosts file.\n\n" +
			"Please run the application with sudo:\n" +
			"sudo ./hosts-manager"
	default:
		return "This application needs administrator privileges to modify the hosts file."
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hosts Profiles Manager")

	// Set application icon
	if icon, err := fyne.LoadResourceFromPath("Icon.png"); err == nil {
		myWindow.SetIcon(icon)
		myApp.SetIcon(icon)
	}

	hostsManager := models.NewHostsManager()

	// Check initial permissions and show warning if needed
	if err := models.CheckWritePermission(); err != nil {
		dialog.ShowInformation("Administrator Privileges Required",
			getPrivilegeInstructions(),
			myWindow)
	}

	// Create views
	profileView := profile.NewProfileView(hostsManager, myWindow)
	managementView := management.NewManagementView(hostsManager, myWindow, profileView.RefreshList)

	// Create tabs with custom styling
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Profiles", theme.HomeIcon(), profileView.Container()),
		container.NewTabItemWithIcon("Management", theme.SettingsIcon(), managementView.Container()),
	)

	// Set tab location and style
	tabs.SetTabLocation(container.TabLocationTop)

	// Add tab change listener
	tabs.OnChanged = func(tab *container.TabItem) {
		if tab.Text == "Management" {
			managementView.RefreshList()
		}
	}

	// Create a container that centers the content horizontally
	mainContent := container.NewGridWithColumns(3,
		layout.NewSpacer(),
		tabs,
		layout.NewSpacer(),
	)

	myWindow.SetContent(mainContent)
	myWindow.Resize(fyne.NewSize(400, 500))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
