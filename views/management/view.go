package management

import (
	"bufio"
	"fmt"
	"strings"

	"hosts-manager/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ManagementView struct {
	hostsManager    *models.HostsManager
	window          fyne.Window
	profileList     *widget.List
	addBtn          *widget.Button
	editBtn         *widget.Button
	deleteBtn       *widget.Button
	selectedProfile *models.HostProfile
	selectedIndex   int
	onProfileChange func()
}

func NewManagementView(hostsManager *models.HostsManager, window fyne.Window, onProfileChange func()) *ManagementView {
	view := &ManagementView{
		hostsManager:    hostsManager,
		window:          window,
		onProfileChange: onProfileChange,
	}
	view.initialize()
	return view
}

func (v *ManagementView) initialize() {
	// Create styled buttons with icons
	v.addBtn = widget.NewButtonWithIcon("Add New", theme.ContentAddIcon(), v.handleAdd)
	v.editBtn = widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), v.handleEdit)
	v.deleteBtn = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), v.handleDelete)

	v.addBtn.Importance = widget.HighImportance
	v.editBtn.Importance = widget.MediumImportance
	v.deleteBtn.Importance = widget.DangerImportance

	v.editBtn.Disable()
	v.deleteBtn.Disable()

	v.profileList = widget.NewList(
		func() int { return len(v.hostsManager.Profiles) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Template"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			profile := v.hostsManager.Profiles[id]
			activeProfile := v.hostsManager.GetActiveProfile()
			box := item.(*fyne.Container)
			label := box.Objects[1].(*widget.Label)
			icon := box.Objects[0].(*widget.Icon)

			if profile.Name == activeProfile {
				label.SetText(profile.Name + " (active)")
			} else {
				label.SetText(profile.Name)
			}
			icon.SetResource(theme.DocumentIcon())
		},
	)

	v.profileList.OnSelected = func(id widget.ListItemID) {
		v.selectedIndex = id
		v.selectedProfile = &v.hostsManager.Profiles[id]
		activeProfile := v.hostsManager.GetActiveProfile()

		if v.selectedProfile.Name == activeProfile {
			v.editBtn.Disable()
			v.deleteBtn.Disable()
			return
		}

		v.editBtn.Enable()
		v.deleteBtn.Enable()
	}
}

func (v *ManagementView) handleAdd() {
	nameEntry := widget.NewEntry()
	entriesEntry := widget.NewMultiLineEntry()
	entriesEntry.SetPlaceHolder("Enter host entries (one per line)\nFormat: IP_ADDRESS HOSTNAME")

	form := dialog.NewForm("New Profile", "Create", "Cancel",
		[]*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Entries", Widget: entriesEntry},
		},
		func(ok bool) {
			if !ok || nameEntry.Text == "" {
				return
			}

			entries := []string{}
			scanner := bufio.NewScanner(strings.NewReader(entriesEntry.Text))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					entries = append(entries, line)
				}
			}

			newProfile := models.HostProfile{
				Name:    nameEntry.Text,
				Entries: entries,
			}

			v.hostsManager.Profiles = append(v.hostsManager.Profiles, newProfile)
			if err := v.hostsManager.SaveProfiles(); err != nil {
				dialog.ShowError(err, v.window)
				return
			}

			v.profileList.Refresh()
			v.onProfileChange()
		}, v.window)

	form.Resize(fyne.NewSize(400, 300))
	form.Show()
}

func (v *ManagementView) handleEdit() {
	if v.selectedProfile == nil {
		return
	}

	if v.selectedProfile.Name == v.hostsManager.GetActiveProfile() {
		dialog.ShowError(fmt.Errorf("cannot edit active profile - deactivate it first"), v.window)
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(v.selectedProfile.Name)

	entriesEntry := widget.NewMultiLineEntry()
	entriesEntry.SetText(strings.Join(v.selectedProfile.Entries, "\n"))

	form := dialog.NewForm("Edit Profile", "Save", "Cancel",
		[]*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Entries", Widget: entriesEntry},
		},
		func(ok bool) {
			if !ok || nameEntry.Text == "" {
				return
			}

			entries := []string{}
			scanner := bufio.NewScanner(strings.NewReader(entriesEntry.Text))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					entries = append(entries, line)
				}
			}

			v.hostsManager.Profiles[v.selectedIndex] = models.HostProfile{
				Name:    nameEntry.Text,
				Entries: entries,
			}

			if err := v.hostsManager.SaveProfiles(); err != nil {
				dialog.ShowError(err, v.window)
				return
			}

			v.profileList.Refresh()
			v.onProfileChange()
		}, v.window)

	form.Resize(fyne.NewSize(400, 300))
	form.Show()
}

func (v *ManagementView) handleDelete() {
	if v.selectedProfile == nil {
		return
	}

	if v.selectedProfile.Name == v.hostsManager.GetActiveProfile() {
		dialog.ShowError(fmt.Errorf("cannot delete active profile - deactivate it first"), v.window)
		return
	}

	dialog.ShowConfirm("Confirm", "Are you sure you want to delete this profile?",
		func(ok bool) {
			if !ok {
				return
			}

			// Remove the profile from the slice
			if v.selectedIndex >= 0 && v.selectedIndex < len(v.hostsManager.Profiles) {
				v.hostsManager.Profiles = append(v.hostsManager.Profiles[:v.selectedIndex], v.hostsManager.Profiles[v.selectedIndex+1:]...)
				if err := v.hostsManager.SaveProfiles(); err != nil {
					dialog.ShowError(err, v.window)
					return
				}
			}

			// Reset selection and disable buttons
			v.selectedProfile = nil
			v.selectedIndex = -1
			v.editBtn.Disable()
			v.deleteBtn.Disable()

			// Clear the list selection and refresh
			v.profileList.UnselectAll()
			v.profileList.Refresh()
			v.onProfileChange()
		}, v.window)
}

func (v *ManagementView) refreshActiveStatus() {
	activeProfile := v.hostsManager.GetActiveProfile()
	v.profileList.Refresh()

	// Update button states if a profile is selected
	if v.selectedProfile != nil {
		if v.selectedProfile.Name == activeProfile {
			v.editBtn.Disable()
			v.deleteBtn.Disable()
		} else {
			v.editBtn.Enable()
			v.deleteBtn.Enable()
		}
	}
}

func (v *ManagementView) Container() *fyne.Container {
	// Create a button container with padding and centered
	buttons := container.NewCenter(
		container.NewHBox(
			layout.NewSpacer(),
			v.addBtn,
			widget.NewSeparator(),
			v.editBtn,
			widget.NewSeparator(),
			v.deleteBtn,
			layout.NewSpacer(),
		),
	)

	// Create a list container with padding
	listContainer := container.NewPadded(v.profileList)

	// Combine all elements with proper spacing
	return container.NewBorder(
		nil,
		buttons,
		nil,
		nil,
		container.NewPadded(listContainer),
	)
}

func (v *ManagementView) RefreshList() {
	v.refreshActiveStatus()
}
