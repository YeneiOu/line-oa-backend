package valueobjects

// LINEProfile represents LINE user profile data
type LINEProfile struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// IsValid validates the LINE profile
func (p *LINEProfile) IsValid() bool {
	return p.UserID != "" && p.DisplayName != ""
}
