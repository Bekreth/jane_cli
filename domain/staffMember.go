package domain

type StaffMember struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	LocationIDs []int  `json:"location_ids"`
}
