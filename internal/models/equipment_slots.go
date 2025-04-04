package models

import "strings"

type EquipmentSlot string

const (
	SlotHead      EquipmentSlot = "head"
	SlotBody      EquipmentSlot = "body"
	SlotMainHand  EquipmentSlot = "main_hand"
	SlotOffHand   EquipmentSlot = "off_hand"
	SlotRingLeft  EquipmentSlot = "ring_left"
	SlotRingRight EquipmentSlot = "ring_right"
	SlotNeck      EquipmentSlot = "neck"
	SlotBack      EquipmentSlot = "back"
	SlotBelt      EquipmentSlot = "belt"
	SlotFeet      EquipmentSlot = "feet"
	SlotHands     EquipmentSlot = "hands"
)

// GetItemTypeSlots returns the slots a given item type can occupy
func GetItemTypeSlots(itemType string) []EquipmentSlot {
	switch itemType {
	case "weapon":
		return []EquipmentSlot{SlotMainHand, SlotOffHand}
	case "armor":
		return []EquipmentSlot{SlotBody}
	case "shield":
		return []EquipmentSlot{SlotOffHand}
	case "ring":
		return []EquipmentSlot{SlotRingLeft, SlotRingRight}
	case "helmet", "headgear":
		return []EquipmentSlot{SlotHead}
	case "boots":
		return []EquipmentSlot{SlotFeet}
	case "gloves":
		return []EquipmentSlot{SlotHands}
	case "amulet", "necklace":
		return []EquipmentSlot{SlotNeck}
	case "cloak":
		return []EquipmentSlot{SlotBack}
	default:
		return []EquipmentSlot{}
	}
}

// IsTwoHanded checks if a weapon's properties indicate it's two-handed
func IsTwoHanded(properties string) bool {
	return strings.Contains(strings.ToLower(properties), "two-handed")
}
