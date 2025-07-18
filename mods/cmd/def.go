package cmd

const (
	CMD_EXPOSE         uint8 = 1
	CMD_DISPOSE        uint8 = 2
	CMD_SHUTDOWN       uint8 = 3
	CMD_OPEN_CTRL_CHAN uint8 = 4
	CMD_CONNECT        uint8 = 5
	CMD_STATUS_OK      uint8 = 6
	CMD_STATUS_UNKNOWN uint8 = 7
	CMD_INVALID_CMD    uint8 = 8
	CMD_MALFORMED_DATA uint8 = 9
	CMD_PORT_UNAVAIL   uint8 = 10
	CMD_OP_FAILED      uint8 = 11
)

type Cmd struct {
	ID   uint8
	Data []byte
}
