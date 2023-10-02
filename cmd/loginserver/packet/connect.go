package packet

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
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

// CheckVersion Packet
func CheckVersion(session *network.Session, reader *network.Reader) {
	var version1 = reader.ReadInt32()

	targetVersion := int32(g_ServerConfig.Version)

	if g_ServerConfig.IgnoreVersionCheck {
		targetVersion = version1
	}

	session.Data.Verified = true

	if version1 != targetVersion {
		log.Errorf("Client version mismatch (Client: %d, server: %d, src: %s)",
			version1, targetVersion, session.GetEndPnt(),
		)

		session.Data.Verified = false
	}

	var packet = network.NewWriter(CHECKVERSION)
	packet.WriteInt32(targetVersion)
	packet.WriteInt32(0x00) // debug
	packet.WriteInt32(0x00) // reserved
	packet.WriteInt32(0x00) // reserved

	session.Send(packet)
}

// FDisconnect Packet
func FDisconnect(session *network.Session, reader *network.Reader) {
	var idx = reader.ReadInt32()

	var packet = network.NewWriter(FDISCONNECT)
	if idx != session.Data.AccountId {
		// wooops invalid account id
		packet.WriteByte(0x00) // failed
	} else {
		var res = account.OnlineRes{}
		g_RPCHandler.Call(rpc.ForceDisconnect, account.OnlineReq{idx, true}, &res)
		if res.Result {
			packet.WriteByte(0x01) // success
		} else {
			packet.WriteByte(0x00) // failed
		}
	}

	session.Send(packet)
}
