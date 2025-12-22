package links

import "time"

type Link struct {
	ID          int
	UserID      int
	URL         string
	Title       string
	Description string
	CreatedAt   time.Time
}
