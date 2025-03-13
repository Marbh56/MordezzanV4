package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type WeaponRepository interface {
	GetWeapon(ctx context.Context, id int64) (*models.Weapon, error)
	GetWeaponByName(ctx context.Context, name string) (*models.Weapon, error)
	ListWeapons(ctx context.Context) ([]*models.Weapon, error)
	CreateWeapon(ctx context.Context, input *models.CreateWeaponInput) (int64, error)
	UpdateWeapon(ctx context.Context, id int64, input *models.UpdateWeaponInput) error
	DeleteWeapon(ctx context.Context, id int64) error
}

type SQLCWeaponRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCWeaponRepository(db *sql.DB) *SQLCWeaponRepository {
	return &SQLCWeaponRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCWeaponRepository) GetWeapon(ctx context.Context, id int64) (*models.Weapon, error) {
	weapon, err := r.q.GetWeapon(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("weapon", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbWeaponToModel(weapon), nil
}

func (r *SQLCWeaponRepository) GetWeaponByName(ctx context.Context, name string) (*models.Weapon, error) {
	weapon, err := r.q.GetWeaponByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("weapon", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbWeaponToModel(weapon), nil
}

func (r *SQLCWeaponRepository) ListWeapons(ctx context.Context) ([]*models.Weapon, error) {
	weapons, err := r.q.ListWeapons(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Weapon, len(weapons))
	for i, weapon := range weapons {
		result[i] = mapDbWeaponToModel(weapon)
	}
	return result, nil
}

func (r *SQLCWeaponRepository) CreateWeapon(ctx context.Context, input *models.CreateWeaponInput) (int64, error) {
	var rangeShort, rangeMedium, rangeLong sql.NullInt64

	if input.RangeShort != nil {
		rangeShort.Int64 = int64(*input.RangeShort)
		rangeShort.Valid = true
	}

	if input.RangeMedium != nil {
		rangeMedium.Int64 = int64(*input.RangeMedium)
		rangeMedium.Valid = true
	}

	if input.RangeLong != nil {
		rangeLong.Int64 = int64(*input.RangeLong)
		rangeLong.Valid = true
	}

	var rateOfFire, damageTwoHanded, properties sql.NullString

	if input.RateOfFire != "" {
		rateOfFire.String = input.RateOfFire
		rateOfFire.Valid = true
	}

	if input.DamageTwoHanded != "" {
		damageTwoHanded.String = input.DamageTwoHanded
		damageTwoHanded.Valid = true
	}

	if input.Properties != "" {
		properties.String = input.Properties
		properties.Valid = true
	}

	result, err := r.q.CreateWeapon(ctx, sqlcdb.CreateWeaponParams{
		Name:            input.Name,
		Category:        input.Category,
		WeaponClass:     int64(input.WeaponClass),
		Cost:            input.Cost,
		Weight:          int64(input.Weight),
		RangeShort:      rangeShort,
		RangeMedium:     rangeMedium,
		RangeLong:       rangeLong,
		RateOfFire:      rateOfFire,
		Damage:          input.Damage,
		DamageTwoHanded: damageTwoHanded,
		Properties:      properties,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	return id, nil
}

func (r *SQLCWeaponRepository) UpdateWeapon(ctx context.Context, id int64, input *models.UpdateWeaponInput) error {
	_, err := r.GetWeapon(ctx, id)
	if err != nil {
		return err
	}

	var rangeShort, rangeMedium, rangeLong sql.NullInt64

	if input.RangeShort != nil {
		rangeShort.Int64 = int64(*input.RangeShort)
		rangeShort.Valid = true
	}

	if input.RangeMedium != nil {
		rangeMedium.Int64 = int64(*input.RangeMedium)
		rangeMedium.Valid = true
	}

	if input.RangeLong != nil {
		rangeLong.Int64 = int64(*input.RangeLong)
		rangeLong.Valid = true
	}

	var rateOfFire, damageTwoHanded, properties sql.NullString

	if input.RateOfFire != "" {
		rateOfFire.String = input.RateOfFire
		rateOfFire.Valid = true
	}

	if input.DamageTwoHanded != "" {
		damageTwoHanded.String = input.DamageTwoHanded
		damageTwoHanded.Valid = true
	}

	if input.Properties != "" {
		properties.String = input.Properties
		properties.Valid = true
	}

	_, err = r.q.UpdateWeapon(ctx, sqlcdb.UpdateWeaponParams{
		Name:            input.Name,
		Category:        input.Category,
		WeaponClass:     int64(input.WeaponClass),
		Cost:            input.Cost,
		Weight:          int64(input.Weight),
		RangeShort:      rangeShort,
		RangeMedium:     rangeMedium,
		RangeLong:       rangeLong,
		RateOfFire:      rateOfFire,
		Damage:          input.Damage,
		DamageTwoHanded: damageTwoHanded,
		Properties:      properties,
		ID:              id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCWeaponRepository) DeleteWeapon(ctx context.Context, id int64) error {
	_, err := r.GetWeapon(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteWeapon(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbWeaponToModel(weapon sqlcdb.Weapon) *models.Weapon {
	var rangeShort, rangeMedium, rangeLong *int

	if weapon.RangeShort.Valid {
		shortVal := int(weapon.RangeShort.Int64)
		rangeShort = &shortVal
	}

	if weapon.RangeMedium.Valid {
		medVal := int(weapon.RangeMedium.Int64)
		rangeMedium = &medVal
	}

	if weapon.RangeLong.Valid {
		longVal := int(weapon.RangeLong.Int64)
		rangeLong = &longVal
	}

	var rateOfFire, damageTwoHanded, properties string

	if weapon.RateOfFire.Valid {
		rateOfFire = weapon.RateOfFire.String
	}

	if weapon.DamageTwoHanded.Valid {
		damageTwoHanded = weapon.DamageTwoHanded.String
	}

	if weapon.Properties.Valid {
		properties = weapon.Properties.String
	}

	return &models.Weapon{
		ID:              weapon.ID,
		Name:            weapon.Name,
		Category:        weapon.Category,
		WeaponClass:     int(weapon.WeaponClass),
		Cost:            weapon.Cost,
		Weight:          int(weapon.Weight),
		RangeShort:      rangeShort,
		RangeMedium:     rangeMedium,
		RangeLong:       rangeLong,
		RateOfFire:      rateOfFire,
		Damage:          weapon.Damage,
		DamageTwoHanded: damageTwoHanded,
		Properties:      properties,
		CreatedAt:       weapon.CreatedAt,
		UpdatedAt:       weapon.UpdatedAt,
	}
}
