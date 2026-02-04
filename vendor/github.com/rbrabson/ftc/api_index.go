package ftc

import (
	"encoding/json"
)

// ApiIndex provives information for the FTC server API
type ApiIndex struct {
	Name                    string  `json:"name,omitempty"`
	APIVersion              string  `json:"apiVersion,omitempty"`
	ServiceMainifestName    *string `json:"serviceMainifestName,omitempty"`
	ServiceMainifestVersion *string `json:"serviceMainifestVersion,omitempty"`
	CodePackageName         string  `json:"codePackageName"`
	CodePackageVersion      string  `json:"codePackageVersion"`
	Status                  string  `json:"status,omitempty"`
	CurrentSeason           int     `json:"currentSeason,omitempty"`
	MaxSeason               int     `json:"maxSeason,omitempty"`
}

// GetApiIndex returns the information for the FTC server API
func GetApiIndex() (*ApiIndex, error) {
	url := server

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output ApiIndex
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return &output, nil
}

// String returns a string representation of General. In this case, it is a json string.
func (general *ApiIndex) String() string {
	body, _ := json.Marshal(general)
	return string(body)
}
