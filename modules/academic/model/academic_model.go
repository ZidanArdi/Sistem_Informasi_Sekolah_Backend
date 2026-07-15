package model

// TeacherContext is a lightweight struct to hold identity for the current request.
// It DOES NOT store any cached academic assignments to avoid stale data.
type TeacherContext struct {
	UserID    uint
	GuruID    uint
	Role      string
	RequestID string
	SchoolID  uint
}
