package cmd

func New(id uint8, data []byte) *Cmd {
	return &Cmd{
		ID:   id,
		Data: data,
	}
}

func (c *Cmd) Pack() []byte {
	buf := make([]byte, 2+len(c.Data))

	buf[0] = c.ID
	buf[1] = uint8(len(c.Data))
	copy(buf[2:], c.Data)

	return buf
}

func (c *Cmd) Unpack(buf []byte) *Cmd {
	if len(buf) < 2 {
		return nil
	}

	size := int(buf[1])

	if len(buf) < 2+size {
		return nil
	}

	return &Cmd{
		ID:   buf[0],
		Data: buf[2 : 2+size],
	}
}
