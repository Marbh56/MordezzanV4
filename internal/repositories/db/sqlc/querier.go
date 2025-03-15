// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	AddInventoryItem(ctx context.Context, arg AddInventoryItemParams) (sql.Result, error)
	CreateAmmo(ctx context.Context, arg CreateAmmoParams) (sql.Result, error)
	CreateArmor(ctx context.Context, arg CreateArmorParams) (sql.Result, error)
	CreateCharacter(ctx context.Context, arg CreateCharacterParams) (sql.Result, error)
	CreateContainer(ctx context.Context, arg CreateContainerParams) (sql.Result, error)
	CreateEquipment(ctx context.Context, arg CreateEquipmentParams) (sql.Result, error)
	CreateInventory(ctx context.Context, arg CreateInventoryParams) (sql.Result, error)
	CreateMagicItem(ctx context.Context, arg CreateMagicItemParams) (sql.Result, error)
	CreatePotion(ctx context.Context, arg CreatePotionParams) (sql.Result, error)
	CreateRing(ctx context.Context, arg CreateRingParams) (sql.Result, error)
	CreateShield(ctx context.Context, arg CreateShieldParams) (sql.Result, error)
	CreateSpell(ctx context.Context, arg CreateSpellParams) (sql.Result, error)
	CreateSpellScroll(ctx context.Context, arg CreateSpellScrollParams) (sql.Result, error)
	CreateTreasure(ctx context.Context, arg CreateTreasureParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	CreateWeapon(ctx context.Context, arg CreateWeaponParams) (sql.Result, error)
	DeleteAmmo(ctx context.Context, id int64) (sql.Result, error)
	DeleteArmor(ctx context.Context, id int64) (sql.Result, error)
	DeleteCharacter(ctx context.Context, id int64) (sql.Result, error)
	DeleteContainer(ctx context.Context, id int64) (sql.Result, error)
	DeleteEquipment(ctx context.Context, id int64) (sql.Result, error)
	DeleteInventory(ctx context.Context, id int64) error
	DeleteMagicItem(ctx context.Context, id int64) (sql.Result, error)
	DeletePotion(ctx context.Context, id int64) (sql.Result, error)
	DeleteRing(ctx context.Context, id int64) (sql.Result, error)
	DeleteShield(ctx context.Context, id int64) (sql.Result, error)
	DeleteSpell(ctx context.Context, id int64) (sql.Result, error)
	DeleteSpellScroll(ctx context.Context, id int64) (sql.Result, error)
	DeleteTreasure(ctx context.Context, id int64) (sql.Result, error)
	DeleteUser(ctx context.Context, id int64) (sql.Result, error)
	DeleteWeapon(ctx context.Context, id int64) (sql.Result, error)
	GetAmmo(ctx context.Context, id int64) (Ammo, error)
	GetAmmoByName(ctx context.Context, name string) (Ammo, error)
	GetArmor(ctx context.Context, id int64) (Armor, error)
	GetArmorByName(ctx context.Context, name string) (Armor, error)
	GetCharacter(ctx context.Context, id int64) (Character, error)
	GetCharactersByUser(ctx context.Context, userID int64) ([]Character, error)
	GetContainer(ctx context.Context, id int64) (Container, error)
	GetContainerByName(ctx context.Context, name string) (Container, error)
	GetEquipment(ctx context.Context, id int64) (Equipment, error)
	GetEquipmentByName(ctx context.Context, name string) (Equipment, error)
	GetFullUserByEmail(ctx context.Context, email string) (User, error)
	GetInventory(ctx context.Context, id int64) (Inventory, error)
	GetInventoryByCharacter(ctx context.Context, characterID int64) (Inventory, error)
	GetInventoryItem(ctx context.Context, id int64) (InventoryItem, error)
	GetInventoryItemByTypeAndItemID(ctx context.Context, arg GetInventoryItemByTypeAndItemIDParams) (InventoryItem, error)
	GetInventoryItems(ctx context.Context, inventoryID int64) ([]InventoryItem, error)
	GetInventoryItemsByType(ctx context.Context, arg GetInventoryItemsByTypeParams) ([]InventoryItem, error)
	GetMagicItem(ctx context.Context, id int64) (MagicItem, error)
	GetMagicItemByName(ctx context.Context, name string) (MagicItem, error)
	GetPotion(ctx context.Context, id int64) (Potion, error)
	GetPotionByName(ctx context.Context, name string) (Potion, error)
	GetRing(ctx context.Context, id int64) (Ring, error)
	GetRingByName(ctx context.Context, name string) (Ring, error)
	GetShield(ctx context.Context, id int64) (Shield, error)
	GetShieldByName(ctx context.Context, name string) (Shield, error)
	GetSpell(ctx context.Context, id int64) (Spell, error)
	GetSpellScroll(ctx context.Context, id int64) (GetSpellScrollRow, error)
	GetSpellScrollsBySpell(ctx context.Context, spellID int64) ([]GetSpellScrollsBySpellRow, error)
	GetSpellsByCharacter(ctx context.Context, characterID int64) ([]Spell, error)
	GetTreasure(ctx context.Context, id int64) (Treasure, error)
	GetTreasureByCharacter(ctx context.Context, characterID sql.NullInt64) (Treasure, error)
	GetUser(ctx context.Context, id int64) (GetUserRow, error)
	GetWeapon(ctx context.Context, id int64) (Weapon, error)
	GetWeaponByName(ctx context.Context, name string) (Weapon, error)
	ListAmmo(ctx context.Context) ([]Ammo, error)
	ListArmors(ctx context.Context) ([]Armor, error)
	ListCharacters(ctx context.Context) ([]Character, error)
	ListContainers(ctx context.Context) ([]Container, error)
	ListEquipment(ctx context.Context) ([]Equipment, error)
	ListInventories(ctx context.Context) ([]Inventory, error)
	ListMagicItems(ctx context.Context) ([]MagicItem, error)
	ListMagicItemsByType(ctx context.Context, itemType string) ([]MagicItem, error)
	ListPotions(ctx context.Context) ([]Potion, error)
	ListRings(ctx context.Context) ([]Ring, error)
	ListShields(ctx context.Context) ([]Shield, error)
	ListSpellScrolls(ctx context.Context) ([]ListSpellScrollsRow, error)
	ListSpells(ctx context.Context) ([]Spell, error)
	ListTreasures(ctx context.Context) ([]Treasure, error)
	ListUsers(ctx context.Context) ([]ListUsersRow, error)
	ListWeapons(ctx context.Context) ([]Weapon, error)
	RemoveAllInventoryItems(ctx context.Context, inventoryID int64) error
	RemoveInventoryItem(ctx context.Context, id int64) error
	UpdateAmmo(ctx context.Context, arg UpdateAmmoParams) (sql.Result, error)
	UpdateArmor(ctx context.Context, arg UpdateArmorParams) (sql.Result, error)
	UpdateCharacter(ctx context.Context, arg UpdateCharacterParams) (sql.Result, error)
	UpdateContainer(ctx context.Context, arg UpdateContainerParams) (sql.Result, error)
	UpdateEquipment(ctx context.Context, arg UpdateEquipmentParams) (sql.Result, error)
	UpdateInventory(ctx context.Context, arg UpdateInventoryParams) (sql.Result, error)
	UpdateInventoryItem(ctx context.Context, arg UpdateInventoryItemParams) (sql.Result, error)
	UpdateMagicItem(ctx context.Context, arg UpdateMagicItemParams) (sql.Result, error)
	UpdatePotion(ctx context.Context, arg UpdatePotionParams) (sql.Result, error)
	UpdateRing(ctx context.Context, arg UpdateRingParams) (sql.Result, error)
	UpdateShield(ctx context.Context, arg UpdateShieldParams) (sql.Result, error)
	UpdateSpell(ctx context.Context, arg UpdateSpellParams) (sql.Result, error)
	UpdateSpellScroll(ctx context.Context, arg UpdateSpellScrollParams) (sql.Result, error)
	UpdateTreasure(ctx context.Context, arg UpdateTreasureParams) (sql.Result, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error)
	UpdateWeapon(ctx context.Context, arg UpdateWeaponParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
