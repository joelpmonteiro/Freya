package packet

import (
	"github.com/ubis/Freya/share/network"
)

// Connect2Svr Packet
func Connect2Svr(session *network.Session, reader *network.Reader) {
	var packet = network.NewWriter(CONNECT2SVR)
	packet.WriteUint32(session.Encryption.Key.Seed2nd)
	packet.WriteUint32(session.AuthKey)
	packet.WriteUint16(session.UserIdx)
	packet.WriteUint16(session.Encryption.RecvXorKeyIdx)

	session.Send(packet)
}
