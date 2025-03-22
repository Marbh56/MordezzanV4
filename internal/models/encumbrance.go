package models

// EncumbranceThresholds defines weights at which various encumbrance effects occur
type EncumbranceThresholds struct {
	BaseEncumbered      float64 `json:"base_encumbered"`       // When movement speed is reduced
	BaseHeavyEncumbered float64 `json:"base_heavy_encumbered"` // When heavily encumbered
	MaximumCapacity     float64 `json:"maximum_capacity"`      // Cannot carry more than this
}

// EncumbranceStatus represents a character's current encumbrance state
type EncumbranceStatus struct {
	Encumbered      bool    `json:"encumbered"`       // Basic encumbrance - reduced movement
	HeavyEncumbered bool    `json:"heavy_encumbered"` // Heavy encumbrance - further penalties
	Overloaded      bool    `json:"overloaded"`       // Cannot move normally
	CurrentWeight   float64 `json:"current_weight"`
	MaximumCapacity float64 `json:"maximum_capacity"`
	WeightRemaining float64 `json:"weight_remaining"`
	PercentFull     int     `json:"percent_full"` // How full is inventory (0-100)
}

// InventoryWeightDetails provides a detailed breakdown of inventory weight
type InventoryWeightDetails struct {
	TotalWeight   float64                 `json:"total_weight"`
	WeightByType  map[string]float64      `json:"weight_by_type"`
	Thresholds    EncumbranceThresholds   `json:"thresholds"`
	Status        EncumbranceStatus       `json:"status"`
	HeaviestItems []WeightedInventoryItem `json:"heaviest_items"`
}

// WeightedInventoryItem represents an inventory item with weight information
type WeightedInventoryItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	ItemType    string  `json:"item_type"`
	Weight      float64 `json:"weight"`
	TotalWeight float64 `json:"total_weight"` // Weight × quantity
	Quantity    int     `json:"quantity"`
}

// CalculateEncumbranceThresholds calculates the encumbrance thresholds for a character
// based on their strength and constitution
func CalculateEncumbranceThresholds(strength, constitution int) EncumbranceThresholds {
	baseThresholds := EncumbranceThresholds{
		BaseEncumbered:      75,
		BaseHeavyEncumbered: 150,
		MaximumCapacity:     300, // Base maximum capacity
	}

	// Calculate strength modifier (in pounds)
	strMod := 0
	maxMod := 0
	switch {
	case strength <= 6:
		strMod = -25
		maxMod = -100
	case strength >= 7 && strength <= 8:
		strMod = -15
		maxMod = -50
	case strength >= 13 && strength <= 14:
		strMod = 15
		maxMod = 50
	case strength >= 15 && strength <= 16:
		strMod = 25
		maxMod = 100
	case strength == 17:
		strMod = 35
		maxMod = 150
	case strength == 18:
		strMod = 50
		maxMod = 200
	}

	// Calculate constitution modifier (in pounds)
	conMod := 0
	conMaxMod := 0
	switch {
	case constitution <= 6:
		conMod = -10
		conMaxMod = -25
	case constitution >= 7 && constitution <= 8:
		conMod = -5
		conMaxMod = -15
	case constitution >= 13 && constitution <= 14:
		conMod = 5
		conMaxMod = 15
	case constitution >= 15 && constitution <= 16:
		conMod = 10
		conMaxMod = 25
	case constitution >= 17:
		conMod = 15
		conMaxMod = 35
	}

	// Apply modifiers
	baseThresholds.BaseEncumbered += float64(strMod + conMod)
	baseThresholds.BaseHeavyEncumbered += float64((strMod + conMod) * 2)
	baseThresholds.MaximumCapacity += float64(maxMod + conMaxMod)

	// Ensure minimum thresholds
	if baseThresholds.BaseEncumbered < 40 {
		baseThresholds.BaseEncumbered = 40
	}
	if baseThresholds.BaseHeavyEncumbered < 60 {
		baseThresholds.BaseHeavyEncumbered = 60
	}
	if baseThresholds.MaximumCapacity < 100 {
		baseThresholds.MaximumCapacity = 100
	}

	return baseThresholds
}

// CalculateEncumbranceStatus determines the character's encumbrance status
// based on current weight and thresholds
func CalculateEncumbranceStatus(currentWeight float64, thresholds EncumbranceThresholds) EncumbranceStatus {
	status := EncumbranceStatus{
		CurrentWeight:   currentWeight,
		MaximumCapacity: thresholds.MaximumCapacity,
		WeightRemaining: thresholds.MaximumCapacity - currentWeight,
	}

	// Calculate percent full (0-100)
	if thresholds.MaximumCapacity > 0 {
		status.PercentFull = int((currentWeight / thresholds.MaximumCapacity) * 100)
		if status.PercentFull > 100 {
			status.PercentFull = 100
		}
	}

	// Determine encumbrance state
	status.Encumbered = currentWeight > thresholds.BaseEncumbered
	status.HeavyEncumbered = currentWeight > thresholds.BaseHeavyEncumbered
	status.Overloaded = currentWeight > thresholds.MaximumCapacity

	return status
}
