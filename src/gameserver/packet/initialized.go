package packet

import (
    "share/network"
    "share/models/character"
    "share/rpc"
)

// Initialized Packet
func Initialized(session *network.Session, reader *network.Reader) {
    var charId = reader.ReadInt32()

    if !session.Data.Verified || !session.Data.LoggedIn {
        log.Errorf("User is not verified (char: %d)", charId)
        return
    }

    // verify char id
    if (charId >> 3) != session.Data.AccountId {
        log.Errorf("User is using invalid character id (id: %d, char: %d)",
            session.Data.AccountId, charId)
        return
    }

    var c = character.Character{}

    // fetch character
    for _, data := range session.Data.CharacterList {
        if data.Id == charId {
            c = data
            break
        }
    }

    // check if character exists
    if c.Id != charId {
        log.Errorf("User is using invalid character id (id: %d, char: %d)",
            session.Data.AccountId, charId)
        return
    }

    // load additional character data
    var req = character.DataReq{byte(g_ServerSettings.ServerId), c.Id}
    var res = character.DataRes{}
    g_RPCHandler.Call(rpc.LoadCharacterData, req, &res)

    // serialize data
    var eq, eqlen   = c.Equipment.Serialize()
    var inv, invlen = res.Inventory.Serialize()
    var sk, sklen   = res.Skills.Serialize()
    var sl, sllen   = res.Links.Serialize()

    var packet = network.NewWriter(INITIALIZED)
    packet.WriteBytes(make([]byte, 57))
    packet.WriteByte(0x00)
    packet.WriteByte(0x14)
    packet.WriteByte(g_ServerSettings.ChannelId)
    packet.WriteBytes(make([]byte, 23))
    packet.WriteByte(0xFF)
    packet.WriteUint16(g_ServerConfig.MaxUsers)
    packet.WriteUint32(0x8501A8C0)
    packet.WriteUint16(0x985A)
    packet.WriteInt32(0x01)
    packet.WriteInt32(0x0100001F)

    packet.WriteInt32(c.World)
    packet.WriteInt32(0x00)
    packet.WriteUint16(c.X)
    packet.WriteUint16(c.Y)
    packet.WriteUint64(c.Exp)
    packet.WriteUint64(c.Alz)
    packet.WriteUint64(c.WarExp)
    packet.WriteUint32(c.Level)
    packet.WriteInt32(0x00)

    packet.WriteUint32(c.STR)
    packet.WriteUint32(c.DEX)
    packet.WriteUint32(c.INT)
    packet.WriteUint32(c.PNT)
    packet.WriteByte(c.SwordRank)
    packet.WriteByte(c.MagicRank)
    packet.WriteUint16(0x00) // padding for skillrank
    packet.WriteUint32(0x00)
    packet.WriteUint16(c.MaxHP)
    packet.WriteUint16(c.CurrentHP)
    packet.WriteUint16(c.MaxMP)
    packet.WriteUint16(c.CurrentMP)
    packet.WriteUint16(c.MaxSP)
    packet.WriteUint16(c.CurrentSP)
    packet.WriteUint16(0x00)//stats.DungeonPoints)
    packet.WriteUint16(0x00)
    packet.WriteInt32(0x2A30)
    packet.WriteInt32(0x01)
    packet.WriteUint16(0x00)//stats.SwordExp)
    packet.WriteUint16(0x00)//stats.SwordPoint)
    packet.WriteUint16(0x00)//stats.MagicExp)
    packet.WriteUint16(0x00)//stats.MagicPoint)
    packet.WriteUint16(0x00)//stats.SwordExpPoint)
    packet.WriteUint16(0x00)//stats.MagicExpPoint)
    packet.WriteInt32(0x00)
    packet.WriteInt32(0x00)
    packet.WriteInt32(0x00)           // honour pnt
    packet.WriteUint64(0x00)         // death penalty exp
    packet.WriteUint64(0x00)         // death hp
    packet.WriteUint64(0x00)         // death mp
    packet.WriteUint16(0x00)       // pk penalty // pk pna

    packet.WriteUint32(0x8501A8C0)    // chat ip
    packet.WriteUint16(0x9858)      // chat port

    packet.WriteUint32(0x8501A8C0)    // ah ip
    packet.WriteUint16(0x9859)      // ah port

    packet.WriteByte(c.Nation)
    packet.WriteInt32(0x00)
    packet.WriteInt32(0x07)           // warp code
    packet.WriteInt32(0x07)           // map code
    packet.WriteUint32(c.Style.Get())
    packet.WriteBytes(make([]byte, 39))

    packet.WriteUint16(eqlen)
    packet.WriteUint16(invlen)
    packet.WriteUint16(sklen)
    packet.WriteUint16(sllen)

    packet.WriteBytes(make([]byte, 6))
    packet.WriteUint16(0x00)        // ap
    packet.WriteUint32(0x00)          // ap exp
    packet.WriteInt16(0x00)
    packet.WriteByte(0x00)         // blessing bead count
    packet.WriteByte(0x00)          // active quest count
    packet.WriteUint16(0x00)        // period item count
    packet.WriteBytes(make([]byte, 1023))

    packet.WriteBytes(make([]byte, 128))       // quest dungeon flags
    packet.WriteBytes(make([]byte, 128))       // mission dungeon flags

    packet.WriteByte(0x00)          // Craft Lv 0
    packet.WriteByte(0x00)          // Craft Lv 1
    packet.WriteByte(0x00)          // Craft Lv 2
    packet.WriteByte(0x00)          // Craft Lv 3
    packet.WriteByte(0x00)          // Craft Lv 4
    packet.WriteUint16(0x00)        // Craft Exp 0
    packet.WriteUint16(0x00)        // Craft Exp 1
    packet.WriteUint16(0x00)        // Craft Exp 2
    packet.WriteUint16(0x00)        // Craft Exp 3
    packet.WriteUint16(0x00)        // Craft Exp 4
    packet.WriteBytes(make([]byte, 16))        // Craft Flags
    packet.WriteUint32(0x00)          // Craft Type

    packet.WriteInt32(0x00)           // Help Window Index
    packet.WriteBytes(make([]byte, 163))

    packet.WriteUint32(0x00)          // TotalPoints
    packet.WriteUint32(0x00)          // GeneralPoints
    packet.WriteUint32(0x00)          // QuestPoints
    packet.WriteUint32(0x00)          // DungeonPoints
    packet.WriteUint32(0x00)          // ItemPoints
    packet.WriteUint32(0x00)          // PVPPoints
    packet.WriteUint32(0x00)          // MissionWarPoints
    packet.WriteUint32(0x00)          // HuntingPoints
    packet.WriteUint32(0x00)          // CraftingPoints
    packet.WriteUint32(0x00)          // CommunityPoints
    packet.WriteUint32(0x00)          // SharedAchievments
    packet.WriteUint32(0x00)          // SpecialPoints

    packet.WriteUint32(0x00)
    packet.WriteUint32(0x00)          // QuestsCount
    packet.WriteUint32(0x00)          // QuestFlagsCount
    packet.WriteUint32(0x00)

    packet.WriteByte(len(c.Name) + 1)
    packet.WriteString(c.Name)

    packet.WriteBytes(eq)
    packet.WriteBytes(inv)
    packet.WriteBytes(sk)
    packet.WriteBytes(sl)

    session.Send(packet)
}