package modules

import (
	"errors"
	"time"
)

type Shop struct {
	Id          string                         `json:"user_id"`
	Title       string                         `json:"title"`
	KioskItems  map[string]map[string]ShopItem `json:"items"`
	LastUpdated int64                          `json:"last_updated"`
}

type ShopItem struct {
	Image       string `json:",omitempty"`
	Stock       int64  `json:"stock"`
	Price       int64  `json:"price"`
	LastUpdated int64  `json:"last_updated"`
}

// ShopManager provides methods for managing the kiosk and its items.
type ShopManager struct {
	kiosk *Shop
}

// NewShopManager creates a new instance of KioskManager with an empty kiosk.
func NewShopManager(id, title string) *ShopManager {
	return &ShopManager{
		kiosk: &Shop{
			Id:          id,
			Title:       title,
			KioskItems:  make(map[string]map[string]ShopItem),
			LastUpdated: time.Now().Unix(),
		},
	}
}

// AddItem adds a new item to the kiosk.
func (km *ShopManager) AddItem(name, itemType string, item ShopItem) {
	km.kiosk.LastUpdated = time.Now().Unix()

	// Check if the item type already exists in the kiosk.
	if km.kiosk.KioskItems[itemType] == nil {
		// If not, create a new map for that item type.
		km.kiosk.KioskItems[itemType] = make(map[string]ShopItem)
	}

	// Check if the item with the given name already exists for the item type.
	if _, ok := km.kiosk.KioskItems[itemType][name]; !ok {
		// If not, set the stock to a default value of 1 if it's less than or equal to 0.
		if item.Stock <= 0 {
			item.Stock = 1
		}
	}

	// Add or update the item in the kiosk.
	km.kiosk.KioskItems[itemType][name] = item
}

// GetItem retrieves an item from the kiosk based on its name.
func (km *ShopManager) GetItem(name string) (ShopItem, error) {
	// Iterate over each item type in the kiosk.
	for _, items := range km.kiosk.KioskItems {
		// Check if the item with the given name exists for any item type.
		if item, ok := items[name]; ok {
			// If found, return the item.
			return item, nil
		}
	}

	// If the item is not found, return an error.
	return ShopItem{}, errors.New("item not found")
}

// RemoveItem removes an item from the kiosk based on its name and item type.
func (km *ShopManager) RemoveItem(name, itemType string) error {
	// Check if the item type exists in the kiosk.
	if items, ok := km.kiosk.KioskItems[itemType]; ok {
		// Check if the item with the given name exists for the item type.
		delete(items, name)

		// If there are no more items for the item type, remove the item type from the kiosk.
		if len(items) == 0 {
			delete(km.kiosk.KioskItems, itemType)
		}

		km.kiosk.LastUpdated = time.Now().Unix()
		return nil
	}

	// If the item is not found, return an error.
	return errors.New("item not found")
}

// GetLastUpdated returns the last updated timestamp of the kiosk.
func (km *ShopManager) GetLastUpdated() int64 {
	return km.kiosk.LastUpdated
}

// GetKiosk returns a copy of the kiosk.
func (km *ShopManager) GetKiosk() Shop {
	// Create a deep copy of the kiosk to prevent direct modifications.
	// This ensures that modifications can only be made through the KioskManager methods.
	kioskCopy := *km.kiosk
	return kioskCopy
}
