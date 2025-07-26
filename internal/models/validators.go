package models

func (d *Depth) Validate() error {
	if len(d.Asks) == 0 || len(d.Bids) == 0 {
		return ErrIncompleteDepthData
	}
	if d.Timestamp <= 0 {
		return ErrInvalidTimestamp
	}

	return nil
}
