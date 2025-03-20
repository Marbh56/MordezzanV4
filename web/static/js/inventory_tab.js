document.addEventListener('DOMContentLoaded', function() {
    // Setup tab click handler
    const inventoryTab = document.querySelector('.tab-item[data-tab="inventory-tab"]');
    if (inventoryTab) {
        inventoryTab.addEventListener('click', function() {
            console.log('Inventory tab clicked');
            setupInventoryAddItem();
            fetchInventory();
        });
    }
    
    // Setup inventory functionality if tab is active initially
    if (document.getElementById('inventory-tab') && 
        document.getElementById('inventory-tab').classList.contains('active')) {
        console.log('Inventory tab is active by default, initializing');
        setupInventoryAddItem();
        fetchInventory();
    }
});

function setupInventoryAddItem() {
    // First check if we've already set this up to avoid duplicate handlers
    if (window.inventorySetup) return;
    
    console.log('Setting up inventory add item functionality');
    window.inventorySetup = true;
    
    const addItemBtn = document.getElementById('btn-add-item');
    if (!addItemBtn) {
        console.error('Add item button not found');
        return;
    }
    
    if (!document.getElementById('add-item-modal')) {
        createAddItemModal();
    }
    
    const modal = document.getElementById('add-item-modal');
    const closeModal = document.getElementById('close-add-item-modal');
    const itemTypeSelect = document.getElementById('item-type');
    const itemSelect = document.getElementById('item-select');
    const quantityInput = document.getElementById('item-quantity');
    const equippedCheckbox = document.getElementById('item-equipped');
    const notesInput = document.getElementById('item-notes');
    const addItemForm = document.getElementById('add-item-form');
    
    // Remove existing event listeners if any (to prevent duplicates)
    const newAddItemBtn = addItemBtn.cloneNode(true);
    addItemBtn.parentNode.replaceChild(newAddItemBtn, addItemBtn);
    
    newAddItemBtn.addEventListener('click', function() {
        console.log('Add item button clicked');
        modal.style.display = 'block';
        itemTypeSelect.value = '';
        itemSelect.innerHTML = '<option value="">Select an item...</option>';
        quantityInput.value = '1';
        equippedCheckbox.checked = false;
        notesInput.value = '';
        itemSelect.disabled = true;
    });
    
    // Close modal when X is clicked
    closeModal.addEventListener('click', function() {
        modal.style.display = 'none';
    });
    
    // Close modal when clicking outside
    window.addEventListener('click', function(event) {
        if (event.target == modal) {
            modal.style.display = 'none';
        }
    });
    
    // Load items when type is selected
    itemTypeSelect.addEventListener('change', function() {
        const selectedType = this.value;
        console.log('Item type selected:', selectedType);
        
        if (selectedType) {
            fetchItemsByType(selectedType);
        } else {
            itemSelect.innerHTML = '<option value="">Select an item...</option>';
            itemSelect.disabled = true;
        }
    });
    
    // Handle form submission
    addItemForm.addEventListener('submit', function(e) {
        e.preventDefault();
        console.log('Form submitted');
        
        const selectedType = itemTypeSelect.value;
        const selectedItemId = itemSelect.value;
        const quantity = parseInt(quantityInput.value) || 1;
        const isEquipped = equippedCheckbox.checked;
        const notes = notesInput.value;
        
        if (!selectedType || !selectedItemId) {
            alert('Please select an item type and item');
            return;
        }
        
        addItemToInventory(selectedType, selectedItemId, quantity, isEquipped, notes);
    });
}

function createAddItemModal() {
    console.log('Creating add item modal');
    
    const modalHTML = `
    <div id="add-item-modal" class="modal">
        <div class="modal-content">
            <span id="close-add-item-modal" class="close">&times;</span>
            <h2 id="modal-title">Add Item to Inventory</h2>
            <form id="add-item-form">
                <div class="form-group">
                    <label for="item-type">Item Type:</label>
                    <select id="item-type" required>
                        <option value="">Select a type...</option>
                        <option value="weapon">Weapon</option>
                        <option value="armor">Armor</option>
                        <option value="shield">Shield</option>
                        <option value="potion">Potion</option>
                        <option value="magic_item">Magic Item</option>
                        <option value="ring">Ring</option>
                        <option value="ammo">Ammunition</option>
                        <option value="spell_scroll">Spell Scroll</option>
                        <option value="container">Container</option>
                        <option value="equipment">Equipment</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="item-select">Item:</label>
                    <select id="item-select" required disabled>
                        <option value="">Select an item...</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="item-quantity">Quantity:</label>
                    <input type="number" id="item-quantity" min="1" value="1" required>
                </div>
                <div class="form-group checkbox">
                    <input type="checkbox" id="item-equipped">
                    <label for="item-equipped">Equip this item?</label>
                </div>
                <div class="form-group">
                    <label for="item-notes">Notes:</label>
                    <textarea id="item-notes" rows="2"></textarea>
                </div>
                <div class="form-actions">
                    <button type="submit" class="btn btn-primary">Add to Inventory</button>
                </div>
            </form>
        </div>
    </div>
    `;
    
    // Add the modal HTML to the page
    document.body.insertAdjacentHTML('beforeend', modalHTML);
}

async function fetchItemsByType(itemType) {
    try {
        const token = localStorage.getItem('authToken');
        const itemSelect = document.getElementById('item-select');
        
        // Show loading state
        itemSelect.innerHTML = '<option value="">Loading items...</option>';
        itemSelect.disabled = true;
        
        // Determine the correct endpoint based on item type
        let endpoint = '';
        switch (itemType) {
            case 'weapon': endpoint = '/api/weapons'; break;
            case 'armor': endpoint = '/api/armors'; break;
            case 'shield': endpoint = '/api/shields'; break;
            case 'potion': endpoint = '/api/potions'; break;
            case 'magic_item': endpoint = '/api/magic-items'; break;
            case 'ring': endpoint = '/api/rings'; break;
            case 'ammo': endpoint = '/api/ammo'; break;
            case 'spell_scroll': endpoint = '/api/spell-scrolls'; break;
            case 'container': endpoint = '/api/containers'; break;
            case 'equipment': endpoint = '/api/equipment'; break;
            default:
                itemSelect.innerHTML = '<option value="">Invalid item type</option>';
                return;
        }
        
        console.log(`Fetching items from: ${endpoint}`);
        
        const response = await fetch(endpoint, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            }
        });

        if (!response.ok) {
            console.error(`API response error: ${response.status} ${response.statusText}`);
            throw new Error(`Failed to fetch items (${response.status})`);
        }

        const items = await response.json();
        console.log(`Received ${items.length} items for type ${itemType}:`, items);
        
        // Reset and populate the select
        itemSelect.innerHTML = '<option value="">Select an item...</option>';
        
        if (!items || items.length === 0) {
            itemSelect.innerHTML += '<option value="" disabled>No items found</option>';
        } else {
            items.forEach(item => {
                const option = document.createElement('option');
                option.value = item.id;
                option.textContent = item.name || 'Unnamed Item';
                itemSelect.appendChild(option);
            });
        }
        
        itemSelect.disabled = false;
    } catch (error) {
        console.error('Error fetching items:', error);
        const itemSelect = document.getElementById('item-select');
        itemSelect.innerHTML = `<option value="">Error: ${error.message}</option>`;
        itemSelect.disabled = true;
    }
}

async function addItemToInventory(itemType, itemId, quantity, isEquipped, notes) {
    try {
        const token = localStorage.getItem('authToken');
        const inventoryID = await getInventoryID();
        
        if (!inventoryID) {
            throw new Error('Could not find inventory for this character');
        }
        
        console.log(`Adding item to inventory ${inventoryID}:`, {
            item_type: itemType,
            item_id: parseInt(itemId),
            quantity: quantity,
            is_equipped: isEquipped,
            notes: notes
        });
        
        const response = await fetch(`/api/inventories/${inventoryID}/items`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                item_type: itemType,
                item_id: parseInt(itemId),
                quantity: quantity,
                is_equipped: isEquipped,
                notes: notes
            })
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response:', errorText);
            
            try {
                const errorData = JSON.parse(errorText);
                throw new Error(errorData.message || `Failed to add item (${response.status})`);
            } catch (e) {
                throw new Error(`Server error (${response.status}): ${errorText.substring(0, 100)}`);
            }
        }
        
        document.getElementById('add-item-modal').style.display = 'none';
        fetchInventory();
        alert('Item added to inventory successfully!');
    } catch (error) {
        console.error('Error adding item to inventory:', error);
        alert('Failed to add item: ' + error.message);
    }
}

async function getInventoryID() {
    try {
        const token = localStorage.getItem('authToken');
        const characterID = getCharacterIdFromURL();
        
        console.log(`Fetching inventory for character: ${characterID}`);
        
        const response = await fetch(`/api/inventories/character/${characterID}`, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            }
        });
        
        if (!response.ok) {
            console.error(`Failed to fetch inventory: ${response.status} ${response.statusText}`);
            throw new Error(`Failed to fetch inventory (${response.status})`);
        }
        
        const data = await response.json();
        console.log('Inventory data:', data);
        
        if (!data.inventory || !data.inventory.id) {
            console.error('Invalid inventory data structure:', data);
            throw new Error('Invalid inventory data received');
        }
        
        return data.inventory.id;
    } catch (error) {
        console.error('Error getting inventory ID:', error);
        return null;
    }
}

async function fetchInventory() {
    try {
        const token = localStorage.getItem('authToken');
        const characterId = getCharacterIdFromURL();
        
        const inventoryTable = document.getElementById('inventory-table');
        const inventoryItems = document.getElementById('inventory-items');
        const inventoryLoading = document.getElementById('inventory-loading');
        const inventoryEmpty = document.getElementById('inventory-empty');
        
        if (!inventoryTable || !inventoryItems || !inventoryLoading || !inventoryEmpty) {
            console.error('Missing inventory elements');
            return;
        }
        
        inventoryLoading.style.display = 'block';
        inventoryTable.style.display = 'none';
        inventoryEmpty.style.display = 'none';
        
        console.log(`Fetching inventory for character: ${characterId}`);
        
        const response = await fetch(`/api/inventories/character/${characterId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to fetch inventory');
        }
        
        const data = await response.json();
        console.log('Inventory data:', data);
        
        const inventory = data.inventory;
        const enrichedItems = data.items;
        
        inventoryLoading.style.display = 'none';
        
        if (!enrichedItems || enrichedItems.length === 0) {
            inventoryEmpty.style.display = 'block';
            return;
        }
        
        inventoryItems.innerHTML = '';
        
        // Add each item
        enrichedItems.forEach(item => {
            const details = item.item_details;
            const row = document.createElement('tr');
            let itemName = "Unknown Item";
            let itemWeight = "N/A";
            
            if (details) {
                itemName = details.name || "Unnamed Item";
                itemWeight = details.weight ? details.weight + ' lbs' : 'N/A';
            }
            
            row.innerHTML = `
                <td>${itemName}</td>
                <td>${formatItemType(item.item_type)}</td>
                <td>${item.quantity}</td>
                <td>${itemWeight}</td>
                <td>${item.notes || ''}</td>
                <td class="item-actions">
                    <button class="btn btn-item-edit" data-id="${item.id}">Edit</button>
                    <button class="btn btn-item-delete" data-id="${item.id}">Delete</button>
                </td>
            `;
            
            inventoryItems.appendChild(row);
        });
        
        document.querySelectorAll('.btn-item-edit').forEach(btn => {
            btn.addEventListener('click', function() {
                const itemId = this.getAttribute('data-id');
                editInventoryItem(itemId);
            });
        });
        
        document.querySelectorAll('.btn-item-delete').forEach(btn => {
            btn.addEventListener('click', function() {
                const itemId = this.getAttribute('data-id');
                if (confirm('Are you sure you want to remove this item?')) {
                    deleteInventoryItem(itemId);
                }
            });
        });
        
        inventoryTable.style.display = 'table';
    } catch (error) {
        console.error('Error fetching inventory:', error);
        
        const inventoryLoading = document.getElementById('inventory-loading');
        const inventoryEmpty = document.getElementById('inventory-empty');
        
        if (inventoryLoading) {
            inventoryLoading.style.display = 'none';
        }
        
        if (inventoryEmpty) {
            inventoryEmpty.style.display = 'block';
            inventoryEmpty.textContent = 'Failed to load inventory: ' + error.message;
        }
    }
}

function formatItemType(type) {
    if (!type) return 'Unknown';
    return type
        .split('_')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');
}

async function editInventoryItem(itemId) {
    alert('Edit item functionality will be implemented soon');
}

async function deleteInventoryItem(itemId) {
    try {
        const token = localStorage.getItem('authToken');
        const inventoryID = await getInventoryID();
        
        if (!inventoryID) {
            throw new Error('Could not find inventory');
        }
        
        const response = await fetch(`/api/inventories/${inventoryID}/items/${itemId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to delete item');
        }
        
        fetchInventory();
    } catch (error) {
        console.error('Error deleting inventory item:', error);
        alert('Failed to delete item: ' + error.message);
    }
}

// Helper function to get character ID from URL
function getCharacterIdFromURL() {
    const pathParts = window.location.pathname.split('/');
    return pathParts[pathParts.length - 1];
}