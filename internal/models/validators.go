package models

// Validate checks if the Rate model has valid data.
func (d *Depth) Validate() error {
	// Check if the depth data has both asks and bids
	// and if the timestamp is valid.
	if len(d.Asks) == 0 || len(d.Bids) == 0 {
		return ErrIncompleteDepthData
	}
	if d.Timestamp <= 0 {
		return ErrInvalidTimestamp
	}

	return nil
}
