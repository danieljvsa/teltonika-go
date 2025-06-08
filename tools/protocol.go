package tools

import (
	"encoding/hex"
	"fmt"

	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
)

func GetProtocol(header []byte) (*tool_domain.ProtocolData, error) {
	if len(header) < 4 {
		return nil, fmt.Errorf("header is too small")
	}

	if hex.EncodeToString(header[:4]) == "00000000" {
		return &tool_domain.ProtocolData{Protocol: "TCP"}, nil
	}

	return &tool_domain.ProtocolData{Protocol: "UDP"}, nil
}
