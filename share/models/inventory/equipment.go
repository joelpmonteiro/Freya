package inventory

import (
	"bytes"
	"encoding/binary"
	"sort"
)

type Equipment struct {
	Equip map[int]Item
}

// Initializes Equipment
func (e *Equipment) Init() {
	e.Equip = make(map[int]Item)
	var item = Item{}

	for key, _ := range eqTypes {
		item.Slot = uint16(key)
		e.Equip[key] = item
	}
}

// Sets equipment item by slot
func (e *Equipment) Set(slot uint16, item Item) {
	e.Equip[int(slot)] = item
}

// Returns equipment item by slot
func (e *Equipment) Get(slot uint16) Item {
	if value, ok := e.Equip[int(slot)]; ok {
		return value
	}

	return Item{}
}

// Removes equipment item by slot
func (e *Equipment) Remove(slot uint16) bool {
	if _, ok := e.Equip[int(slot)]; ok {
		delete(e.Equip, int(slot))
		return true
	}

	return false
}

// Serializes equipment into byte array
func (e *Equipment) Serialize() ([]byte, int) {
	// collect keys for sorted iteration
	var keys []int
	for k := range e.Equip {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var length = 0

	var equip bytes.Buffer
	for _, value := range keys {
		if e.Equip[value].Kind > 0 {
			binary.Write(&equip, binary.LittleEndian, e.Equip[value])
			length++
		}
	}

	return equip.Bytes(), length
}

// Serializes equipment kind_idx into byte array
func (e *Equipment) SerializeKind() []byte {
	// collect keys for sorted iteration
	var keys []int
	for k := range e.Equip {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var equip bytes.Buffer
	for _, value := range keys {
		binary.Write(&equip, binary.LittleEndian, e.Equip[value].Kind)
	}

	return equip.Bytes()
}

// SerializeEx serializes equipment with kind and option into byte array
func (e *Equipment) SerializeEx() ([]byte, int) {
	// collect keys for sorted iteration
	var keys []int
	for k := range e.Equip {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var length = 0

	var equip bytes.Buffer
	for _, value := range keys {
		if e.Equip[value].Kind > 0 {
			binary.Write(&equip, binary.LittleEndian, byte(e.Equip[value].Slot))
			binary.Write(&equip, binary.LittleEndian, e.Equip[value].Kind)
			binary.Write(&equip, binary.LittleEndian, e.Equip[value].Option)
			length++
		}
	}

	return equip.Bytes(), length
}
