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
	AddKnownSpell(ctx context.Context, arg AddKnownSpellParams) (sql.Result, error)
	AddWeaponMastery(ctx context.Context, arg AddWeaponMasteryParams) error
	ClearPreparedSpells(ctx context.Context, characterID int64) error
	CountPreparedSpellsByLevelAndClass(ctx context.Context, arg CountPreparedSpellsByLevelAndClassParams) (int64, error)
	CountWeaponMasteries(ctx context.Context, arg CountWeaponMasteriesParams) (int64, error)
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
	DeleteWeaponMastery(ctx context.Context, arg DeleteWeaponMasteryParams) error
	GetAllClassData(ctx context.Context, className string) ([]ClassDatum, error)
	GetAmmo(ctx context.Context, id int64) (Ammo, error)
	GetAmmoByName(ctx context.Context, name string) (Ammo, error)
	GetArmor(ctx context.Context, id int64) (Armor, error)
	GetArmorByName(ctx context.Context, name string) (Armor, error)
	// Gets all barbarian abilities available to a character based on their level
	GetBarbarianAbilities(ctx context.Context, characterLevel int64) ([]BarbarianAbility, error)
	// Gets all bard abilities available to a character based on their level
	GetBardAbilities(ctx context.Context, characterLevel int64) ([]BardAbility, error)
	GetBardDruidSpells(ctx context.Context, level int64) (BardDruidSpell, error)
	GetBardIllusionistSpells(ctx context.Context, level int64) (BardIllusionistSpell, error)
	// Gets all berserker abilities available to a character based on their level
	GetBerserkerAbilities(ctx context.Context, characterLevel int64) ([]BerserkerAbility, error)
	GetBerserkerNaturalAC(ctx context.Context, arg GetBerserkerNaturalACParams) (int64, error)
	// Gets all cataphract abilities available to a character based on their level
	GetCataphractAbilities(ctx context.Context, characterLevel int64) ([]CataphractAbility, error)
	GetCharacter(ctx context.Context, id int64) (GetCharacterRow, error)
	GetCharacterForSpellcasting(ctx context.Context, id int64) (Character, error)
	GetCharactersByUser(ctx context.Context, userID int64) ([]GetCharactersByUserRow, error)
	GetClassAbilities(ctx context.Context, className string) ([]GetClassAbilitiesRow, error)
	GetClassAbilitiesByLevel(ctx context.Context, arg GetClassAbilitiesByLevelParams) ([]GetClassAbilitiesByLevelRow, error)
	GetClassData(ctx context.Context, arg GetClassDataParams) (ClassDatum, error)
	GetClassDataForSpellcasting(ctx context.Context, arg GetClassDataForSpellcastingParams) (ClassDatum, error)
	// Gets all cleric abilities available to a character based on their level
	GetClericAbilities(ctx context.Context, characterLevel int64) ([]ClericAbility, error)
	GetClericTurningAbility(ctx context.Context, level int64) (int64, error)
	GetContainer(ctx context.Context, id int64) (Container, error)
	GetContainerByName(ctx context.Context, name string) (Container, error)
	// Gets all cryomancer abilities available to a character based on their level
	GetCryomancerAbilities(ctx context.Context, characterLevel int64) ([]CryomancerAbility, error)
	// Gets all druid abilities available to a character based on their level
	GetDruidAbilities(ctx context.Context, characterLevel int64) ([]DruidAbility, error)
	GetEquipment(ctx context.Context, id int64) (Equipment, error)
	GetEquipmentByName(ctx context.Context, name string) (Equipment, error)
	GetEquippedItems(ctx context.Context, inventoryID int64) ([]InventoryItem, error)
	// Gets all fighter abilities available to a character based on their level
	GetFighterAbilities(ctx context.Context, characterLevel int64) ([]FighterAbility, error)
	GetFullUserByEmail(ctx context.Context, email string) (User, error)
	// Gets all huntsman abilities available to a character based on their level
	GetHuntsmanAbilities(ctx context.Context, characterLevel int64) ([]HuntsmanAbility, error)
	// Gets all illusionist abilities available to a character based on their level
	GetIllusionistAbilities(ctx context.Context, characterLevel int64) ([]IllusionistAbility, error)
	GetInventory(ctx context.Context, id int64) (Inventory, error)
	GetInventoryByCharacter(ctx context.Context, characterID int64) (Inventory, error)
	GetInventoryItem(ctx context.Context, id int64) (InventoryItem, error)
	GetInventoryItemByTypeAndItemID(ctx context.Context, arg GetInventoryItemByTypeAndItemIDParams) (InventoryItem, error)
	GetInventoryItems(ctx context.Context, inventoryID int64) ([]InventoryItem, error)
	GetInventoryItemsByType(ctx context.Context, arg GetInventoryItemsByTypeParams) ([]InventoryItem, error)
	GetItemsBySlot(ctx context.Context, arg GetItemsBySlotParams) ([]InventoryItem, error)
	GetKnownSpellByCharacterAndSpell(ctx context.Context, arg GetKnownSpellByCharacterAndSpellParams) (KnownSpell, error)
	GetKnownSpells(ctx context.Context, characterID int64) ([]KnownSpell, error)
	GetKnownSpellsByClass(ctx context.Context, arg GetKnownSpellsByClassParams) ([]KnownSpell, error)
	// Gets all legerdemainist abilities available to a character based on their level
	GetLegerdemainistAbilities(ctx context.Context, characterLevel int64) ([]LegerdemainistAbility, error)
	GetMagicItem(ctx context.Context, id int64) (MagicItem, error)
	GetMagicItemByName(ctx context.Context, name string) (MagicItem, error)
	// Gets all magician abilities available to a character based on their level
	GetMagicianAbilities(ctx context.Context, characterLevel int64) ([]MagicianAbility, error)
	GetMonkACBonus(ctx context.Context, level int64) (int64, error)
	// Gets all monk abilities available to a character based on their level
	GetMonkAbilities(ctx context.Context, characterLevel int64) ([]MonkAbility, error)
	GetMonkEmptyHandDamage(ctx context.Context, level int64) (string, error)
	// Gets all necromancer abilities available to a character based on their level
	GetNecromancerAbilities(ctx context.Context, characterLevel int64) ([]NecromancerAbility, error)
	GetNecromancerTurningAbility(ctx context.Context, level int64) (int64, error)
	GetNextAvailableSlotIndex(ctx context.Context, arg GetNextAvailableSlotIndexParams) (int64, error)
	GetNextLevelData(ctx context.Context, arg GetNextLevelDataParams) (ClassDatum, error)
	// Gets all paladin abilities available to a character based on their level
	GetPaladinAbilities(ctx context.Context, characterLevel int64) ([]PaladinAbility, error)
	GetPaladinTurningAbility(ctx context.Context, level int64) (int64, error)
	GetPotion(ctx context.Context, id int64) (Potion, error)
	GetPotionByName(ctx context.Context, name string) (Potion, error)
	GetPreparedSpellByCharacterAndSpell(ctx context.Context, arg GetPreparedSpellByCharacterAndSpellParams) (PreparedSpell, error)
	GetPreparedSpells(ctx context.Context, characterID int64) ([]PreparedSpell, error)
	GetPreparedSpellsByClass(ctx context.Context, arg GetPreparedSpellsByClassParams) ([]PreparedSpell, error)
	// Gets all priest abilities available to a character based on their level
	GetPriestAbilities(ctx context.Context, characterLevel int64) ([]PriestAbility, error)
	// Gets all purloiner abilities available to a character based on their level
	GetPurloinerAbilities(ctx context.Context, characterLevel int64) ([]PurloinerAbility, error)
	// Gets all pyromancer abilities available to a character based on their level
	GetPyromancerAbilities(ctx context.Context, characterLevel int64) ([]PyromancerAbility, error)
	// Gets all ranger abilities available to a character based on their level
	GetRangerAbilities(ctx context.Context, characterLevel int64) ([]RangerAbility, error)
	GetRangerDruidSpellSlots(ctx context.Context, classLevel int64) ([]GetRangerDruidSpellSlotsRow, error)
	GetRangerMagicianSpellSlots(ctx context.Context, classLevel int64) ([]GetRangerMagicianSpellSlotsRow, error)
	GetRing(ctx context.Context, id int64) (Ring, error)
	GetRingByName(ctx context.Context, name string) (Ring, error)
	// Gets all runegraver abilities available to a character based on their level
	GetRunegraverAbilities(ctx context.Context, characterLevel int64) ([]RunegraverAbility, error)
	GetRunesPerDay(ctx context.Context, level int64) (GetRunesPerDayRow, error)
	// Gets all scout abilities available to a character based on their level
	GetScoutAbilities(ctx context.Context, characterLevel int64) ([]ScoutAbility, error)
	// Gets all shaman abilities available to a character based on their level
	GetShamanAbilities(ctx context.Context, characterLevel int64) ([]ShamanAbility, error)
	GetShamanArcaneSpells(ctx context.Context, level int64) (ShamanArcaneSpell, error)
	GetShamanDivineSpells(ctx context.Context, level int64) (ShamanDivineSpell, error)
	GetShield(ctx context.Context, id int64) (Shield, error)
	GetShieldByName(ctx context.Context, name string) (Shield, error)
	GetSpell(ctx context.Context, id int64) (Spell, error)
	GetSpellForSpellcasting(ctx context.Context, id int64) (Spell, error)
	GetSpellScroll(ctx context.Context, id int64) (GetSpellScrollRow, error)
	GetSpellScrollsBySpell(ctx context.Context, spellID int64) ([]GetSpellScrollsBySpellRow, error)
	GetSpellsByClassLevel(ctx context.Context, arg GetSpellsByClassLevelParams) ([]Spell, error)
	// Gets all thief abilities available to a character based on their level
	GetThiefAbilities(ctx context.Context, characterLevel int64) ([]ThiefAbility, error)
	GetThiefSkillsByLevel(ctx context.Context, level int64) ([]ThiefSkill, error)
	GetTreasure(ctx context.Context, id int64) (Treasure, error)
	GetTreasureByCharacter(ctx context.Context, characterID sql.NullInt64) (Treasure, error)
	GetUser(ctx context.Context, id int64) (GetUserRow, error)
	// Gets all warlock abilities available to a character based on their level
	GetWarlockAbilities(ctx context.Context, characterLevel int64) ([]WarlockAbility, error)
	GetWeapon(ctx context.Context, id int64) (Weapon, error)
	GetWeaponByName(ctx context.Context, name string) (Weapon, error)
	GetWeaponMasteriesByCharacter(ctx context.Context, characterID int64) ([]WeaponMastery, error)
	GetWeaponMasteryByBaseName(ctx context.Context, arg GetWeaponMasteryByBaseNameParams) (WeaponMastery, error)
	GetWeaponMasteryByID(ctx context.Context, id int64) (WeaponMastery, error)
	// Gets all witch abilities available to a character based on their level
	GetWitchAbilities(ctx context.Context, characterLevel int64) ([]WitchAbility, error)
	ListAmmo(ctx context.Context) ([]Ammo, error)
	ListArmors(ctx context.Context) ([]Armor, error)
	ListCharacters(ctx context.Context) ([]ListCharactersRow, error)
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
	MarkSpellAsMemorized(ctx context.Context, arg MarkSpellAsMemorizedParams) error
	MarkSpellAsMemorizedBySpellID(ctx context.Context, arg MarkSpellAsMemorizedBySpellIDParams) error
	PrepareSpell(ctx context.Context, arg PrepareSpellParams) (sql.Result, error)
	RecalculateInventoryWeight(ctx context.Context, id int64) error
	RemoveAllInventoryItems(ctx context.Context, inventoryID int64) error
	RemoveInventoryItem(ctx context.Context, id int64) error
	RemoveKnownSpell(ctx context.Context, id int64) error
	ResetAllMemorizedSpells(ctx context.Context, characterID int64) error
	UnprepareSpell(ctx context.Context, id int64) error
	UpdateAmmo(ctx context.Context, arg UpdateAmmoParams) (sql.Result, error)
	UpdateArmor(ctx context.Context, arg UpdateArmorParams) (sql.Result, error)
	UpdateCharacter(ctx context.Context, arg UpdateCharacterParams) (sql.Result, error)
	UpdateContainer(ctx context.Context, arg UpdateContainerParams) (sql.Result, error)
	UpdateEquipment(ctx context.Context, arg UpdateEquipmentParams) (sql.Result, error)
	UpdateInventory(ctx context.Context, arg UpdateInventoryParams) (sql.Result, error)
	UpdateInventoryItem(ctx context.Context, arg UpdateInventoryItemParams) (sql.Result, error)
	UpdateInventoryWeight(ctx context.Context, arg UpdateInventoryWeightParams) error
	UpdateMagicItem(ctx context.Context, arg UpdateMagicItemParams) (sql.Result, error)
	UpdatePotion(ctx context.Context, arg UpdatePotionParams) (sql.Result, error)
	UpdateRing(ctx context.Context, arg UpdateRingParams) (sql.Result, error)
	UpdateShield(ctx context.Context, arg UpdateShieldParams) (sql.Result, error)
	UpdateSpell(ctx context.Context, arg UpdateSpellParams) (sql.Result, error)
	UpdateSpellScroll(ctx context.Context, arg UpdateSpellScrollParams) (sql.Result, error)
	UpdateTreasure(ctx context.Context, arg UpdateTreasureParams) (sql.Result, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	UpdateWeapon(ctx context.Context, arg UpdateWeaponParams) (sql.Result, error)
	UpdateWeaponMasteryLevel(ctx context.Context, arg UpdateWeaponMasteryLevelParams) error
}

var _ Querier = (*Queries)(nil)
