package profile

import (
	"hosts-manager/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ProfileView struct {
	hostsManager    *models.HostsManager
	window          fyne.Window
	statusLabel     *widget.Label
	profileList     *widget.List
	activateBtn     *widget.Button
	deactivateBtn   *widget.Button
	selectedProfile *models.HostProfile
}

func NewProfileView(hostsManager *models.HostsManager, window fyne.Window) *ProfileView {
	view := &ProfileView{
		hostsManager: hostsManager,
		window:       window,
	}
	view.initialize()
	return view
}

func (v *ProfileView) initialize() {
	// Create a styled status label
	v.statusLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	v.updateStatus()

	// Create styled buttons with icons
	v.activateBtn = widget.NewButtonWithIcon("Activate", theme.ConfirmIcon(), v.handleActivate)
	v.deactivateBtn = widget.NewButtonWithIcon("Deactivate", theme.CancelIcon(), v.handleDeactivate)
	v.activateBtn.Importance = widget.HighImportance
	v.deactivateBtn.Importance = widget.WarningImportance
	v.activateBtn.Disable()
	v.deactivateBtn.Disable()

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
				label.SetText(profile.Name)
				icon.SetResource(theme.RadioButtonCheckedIcon())
			} else {
				label.SetText(profile.Name)
				icon.SetResource(theme.RadioButtonIcon())
			}
		},
	)

	v.profileList.OnSelected = func(id widget.ListItemID) {
		v.selectedProfile = &v.hostsManager.Profiles[id]
		activeProfile := v.hostsManager.GetActiveProfile()
		if v.selectedProfile.Name == activeProfile {
			v.activateBtn.Disable()
			v.deactivateBtn.Enable()
		} else {
			v.activateBtn.Enable()
			v.deactivateBtn.Disable()
		}
	}
}

func (v *ProfileView) updateStatus() {
	activeProfile := v.hostsManager.GetActiveProfile()
	if activeProfile == "" {
		v.statusLabel.SetText("No profile active")
	} else {
		v.statusLabel.SetText("Active profile: " + activeProfile)
	}
}

func (v *ProfileView) handleActivate() {
	if v.selectedProfile == nil {
		return
	}
	dialog.ShowConfirm("Confirm", "Are you sure you want to activate this profile?",
		func(ok bool) {
			if ok {
				err := v.hostsManager.ApplyProfile(v.selectedProfile)
				if err != nil {
					dialog.ShowError(err, v.window)
					return
				}
				v.updateStatus()
				v.profileList.Refresh()
				v.activateBtn.Disable()
				v.deactivateBtn.Enable()
				dialog.ShowInformation("Success", "Profile activated successfully", v.window)
			}
		}, v.window)
}

func (v *ProfileView) handleDeactivate() {
	dialog.ShowConfirm("Confirm", "Are you sure you want to deactivate the current profile?",
		func(ok bool) {
			if ok {
				err := v.hostsManager.RemoveCurrentProfile()
				if err != nil {
					dialog.ShowError(err, v.window)
					return
				}
				v.updateStatus()
				v.profileList.Refresh()
				v.activateBtn.Enable()
				v.deactivateBtn.Disable()
				dialog.ShowInformation("Success", "Profile deactivated successfully", v.window)
			}
		}, v.window)
}

func (v *ProfileView) Container() *fyne.Container {
	// Create a button container with padding and centered
	buttons := container.NewCenter(
		container.NewHBox(
			layout.NewSpacer(),
			v.activateBtn,
			widget.NewSeparator(),
			v.deactivateBtn,
			layout.NewSpacer(),
		),
	)

	// Create a list container with padding and status label
	listContainer := container.NewVBox(
		container.NewCenter(v.statusLabel),
		container.NewPadded(v.profileList),
	)

	// Combine all elements with proper spacing
	return container.NewBorder(
		nil,
		buttons,
		nil,
		nil,
		listContainer,
	)
}

func (v *ProfileView) RefreshList() {
	v.profileList.Refresh()
}
