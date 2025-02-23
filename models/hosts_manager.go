package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type HostProfile struct {
	Name    string   `json:"name"`
	Entries []string `json:"entries"`
}

type HostsManager struct {
	Profiles     []HostProfile
	ProfilesPath string
	HostsPath    string
}

func getHostsPath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

func NewHostsManager() *HostsManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	profilesPath := filepath.Join(homeDir, ".hosts-manager", "profiles.json")
	if err := os.MkdirAll(filepath.Dir(profilesPath), 0755); err != nil {
		panic(err)
	}

	hm := &HostsManager{
		ProfilesPath: profilesPath,
		HostsPath:    getHostsPath(),
	}

	hm.loadProfiles()
	return hm
}

func (hm *HostsManager) loadProfiles() {
	data, err := ioutil.ReadFile(hm.ProfilesPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error reading profiles: %v\n", err)
		}
		return
	}

	if err := json.Unmarshal(data, &hm.Profiles); err != nil {
		fmt.Printf("Error parsing profiles: %v\n", err)
	}
}

func (hm *HostsManager) SaveProfiles() error {
	data, err := json.MarshalIndent(hm.Profiles, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(hm.ProfilesPath, data, 0644)
}

func CheckWritePermission() error {
	file, err := os.OpenFile(getHostsPath(), os.O_WRONLY, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: this application needs administrator privileges to modify the hosts file.\n\nPlease run the application with sudo")
		}
		return err
	}
	file.Close()
	return nil
}

func (hm *HostsManager) ApplyProfile(profile *HostProfile) error {
	if err := CheckWritePermission(); err != nil {
		return err
	}

	content, err := ioutil.ReadFile(hm.HostsPath)
	if err != nil {
		return err
	}

	backupPath := hm.HostsPath + ".backup"
	if err := ioutil.WriteFile(backupPath, content, 0644); err != nil {
		return err
	}

	var existingEntries []string
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	inManagedSection := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "### BEGIN HOST CHANGER") {
			inManagedSection = true
			continue
		}
		if strings.Contains(line, "### END HOST CHANGER ###") {
			inManagedSection = false
			continue
		}
		if !inManagedSection {
			existingEntries = append(existingEntries, line)
		}
	}

	var newContent strings.Builder
	for _, entry := range existingEntries {
		newContent.WriteString(entry + "\n")
	}

	newContent.WriteString("\n### BEGIN HOST CHANGER - [" + profile.Name + "] ###\n")
	for _, entry := range profile.Entries {
		newContent.WriteString(entry + "\n")
	}
	newContent.WriteString("### END HOST CHANGER ###\n")

	return ioutil.WriteFile(hm.HostsPath, []byte(newContent.String()), 0644)
}

func (hm *HostsManager) RemoveCurrentProfile() error {
	if err := CheckWritePermission(); err != nil {
		return err
	}

	content, err := ioutil.ReadFile(hm.HostsPath)
	if err != nil {
		return err
	}

	var newContent strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	inManagedSection := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "### BEGIN HOST CHANGER") {
			inManagedSection = true
			continue
		}
		if strings.Contains(line, "### END HOST CHANGER ###") {
			inManagedSection = false
			continue
		}
		if !inManagedSection {
			newContent.WriteString(line + "\n")
		}
	}

	content = []byte(strings.TrimRight(newContent.String(), "\n"))
	return ioutil.WriteFile(hm.HostsPath, content, 0644)
}

func (hm *HostsManager) GetActiveProfile() string {
	content, err := ioutil.ReadFile(hm.HostsPath)
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "### BEGIN HOST CHANGER - [") {
			start := strings.Index(line, "[") + 1
			end := strings.Index(line, "]")
			if start > 0 && end > start {
				return line[start:end]
			}
		}
	}
	return ""
}
