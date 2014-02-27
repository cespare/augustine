package augustine

type Doc struct {
	Text []byte
}

func (d *Doc) MarshalBinary() ([]byte, error) {
	return d.Text, nil
}

func (d *Doc) UnmarshalBinary(b []byte) error {
	d.Text = make([]byte, len(b))
	copy(d.Text, b)
	return nil
}
