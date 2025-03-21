const inventoryState = {
    items: [],
    inventory: null,
    itemTypes: [
        { value: 'weapon', label: 'Weapon' },
        { value: 'armor', label: 'Armor' },
        { value: 'shield', label: 'Shield' },
        { value: 'potion', label: 'Potion' },
        { value: 'magic_item', label: 'Magic Item' },
        { value: 'ring', label: 'Ring' },
        { value: 'ammo', label: 'Ammunition' },
        { value: 'spell_scroll', label: 'Spell Scroll' },
        { value: 'container', label: 'Container' },
        { value: 'equipment', label: 'Equipment' }
    ]
};

// Main function to fetch inventory data
async function fetchInventory(characterId = null) {
    try {
        if (!characterId) {
            characterId = getCharacterIdFromURL();
        }
        
        console.log(`Fetching inventory for character: ${characterId}`);
        
        const token = localStorage.getItem('authToken');
        if (!token) {
            throw new Error('Authentication token not found');
        }
        
        const inventoryTable = document.getElementById('inventory-table');
        const inventoryItems = document.getElementById('inventory-items');
        const inventoryLoading = document.getElementById('inventory-loading');
        const inventoryEmpty = document.getElementById('inventory-empty');
        
        if (inventoryLoading) {
            inventoryLoading.style.display = 'block';
        }
        if (inventoryTable) {
            inventoryTable.style.display = 'none';
        }
        if (inventoryEmpty) {
            inventoryEmpty.style.display = 'none';
        }
        
        const response = await fetch(`/api/inventories/character/${characterId}`, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            }
        });
        
        if (!response.ok) {
            throw new Error(`Failed to fetch inventory: ${response.statusText}`);
        }
        
        const data = await response.json();
        console.log('Inventory data received:', data);
        
        inventoryState.inventory = data.inventory;
        inventoryState.items = data.items || [];
        
        // Update the UI
        renderInventoryItems();
        
        return data;
    } catch (error) {
        console.error('Error fetching inventory:', error);
        displayInventoryError(error.message);
        return null;
    }
}

function renderInventoryItems() {
    const inventoryTable = document.getElementById('inventory-table');
    const inventoryItems = document.getElementById('inventory-items');
    const inventoryLoading = document.getElementById('inventory-loading');
    const inventoryEmpty = document.getElementById('inventory-empty');
    
    if (inventoryLoading) {
        inventoryLoading.style.display = 'none';
    }
    
    if (!inventoryState.items || inventoryState.items.length === 0) {
        if (inventoryTable) {
            inventoryTable.style.display = 'none';
        }
        if (inventoryEmpty) {
            inventoryEmpty.style.display = 'block';
        }
        return;
    }
    
    if (inventoryItems) {
        inventoryItems.innerHTML = '';
        // Add each item to the table
        inventoryState.items.forEach(item => {
            const details = item.item_details || {};
            const row = document.createElement('tr');
            const weight = details.weight || 0;
            const weightStr = weight > 0 ? `${weight} lbs` : 'N/A';
            
            // Add an equipped status toggle
            const equippedStatus = item.is_equipped ? 'Equipped' : 'Not equipped';
            const equippedButtonText = item.is_equipped ? 'Unequip' : 'Equip';
            const equippedButtonClass = item.is_equipped ? 'btn-unequip' : 'btn-equip';
            
            row.innerHTML = `
                <td>${details.name || 'Unknown'}</td>
                <td>${formatItemType(item.item_type)}</td>
                <td>${item.quantity}</td>
                <td>${weightStr}</td>
                <td>${equippedStatus}</td>
                <td>${item.notes || ''}</td>
                <td class="item-actions">
                    <button class="btn ${equippedButtonClass}" data-id="${item.id}" data-equipped="${item.is_equipped}">
                        ${equippedButtonText}
                    </button>
                    <button class="btn btn-item-edit" data-id="${item.id}">Edit</button>
                    <button class="btn btn-item-delete" data-id="${item.id}">Delete</button>
                </td>
            `;
            
            inventoryItems.appendChild(row);
        });
    }
    
    if (inventoryTable) {
        inventoryTable.style.display = 'table';
    }
    if (inventoryEmpty) {
        inventoryEmpty.style.display = 'none';
    }
    
    // Add event listeners for all buttons
    addInventoryButtonListeners();
}

async function toggleEquippedStatus(itemId, currentStatus) {
    try {
        const token = localStorage.getItem('authToken');
        const characterId = getCharacterIdFromURL();
        const inventoryID = inventoryState.inventory ? inventoryState.inventory.id : null;
        
        if (!inventoryID) {
            throw new Error('Could not find inventory');
        }
        
        console.log(`Toggling equipped status for item ${itemId}: ${currentStatus} → ${!currentStatus}`);
        
        const response = await fetch(`/api/inventories/${inventoryID}/items/${itemId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                is_equipped: !currentStatus
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to update item equipped status');
        }
        
        // Refresh the inventory display
        await fetchInventory(characterId);
        
        // If the combat tab exists and is active, refresh that too
        const combatTab = document.getElementById('combat-tab');
        if (combatTab && combatTab.classList.contains('active') && typeof loadCombatData === 'function') {
            loadCombatData();
        }
        
    } catch (error) {
        console.error('Error toggling equipped status:', error);
        alert('Failed to update item: ' + error.message);
    }
}

// Format item type for display
function formatItemType(type) {
    if (!type) return 'Unknown';
    return type.replace(/_/g, ' ')
        .split(' ')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');
}

// Add event listeners to the edit and delete buttons
function addItemActionListeners() {
    // Add listeners for equip/unequip buttons
    document.querySelectorAll('.btn-equip, .btn-unequip').forEach(btn => {
        btn.addEventListener('click', function() {
            const itemId = this.getAttribute('data-id');
            const isEquipped = this.getAttribute('data-equipped') === 'true';
            toggleEquippedStatus(itemId, isEquipped);
        });
    });
    
    // Existing code for edit buttons
    document.querySelectorAll('.btn-item-edit').forEach(btn => {
        btn.addEventListener('click', function() {
            const itemId = this.getAttribute('data-id');
            editInventoryItem(itemId);
        });
    });
    
    // Existing code for delete buttons
    document.querySelectorAll('.btn-item-delete').forEach(btn => {
        btn.addEventListener('click', function() {
            const itemId = this.getAttribute('data-id');
            confirmDeleteItem(itemId);
        });
    });
}

// Create the modal for adding/editing inventory items
function createInventoryModal() {
    if (document.getElementById('inventory-modal')) {
        return; // Modal already exists
    }
    
    const modalHTML = `
    <div id="inventory-modal" class="modal">
        <div class="modal-content">
            <span class="close" id="close-inventory-modal">&times;</span>
            <h2 id="modal-title">Add Item to Inventory</h2>
            <form id="inventory-form">
                <div class="form-group">
                    <label for="item-type">Item Type:</label>
                    <select id="item-type" required>
                        <option value="">Select Type</option>
                        ${inventoryState.itemTypes.map(type =>
                            `<option value="${type.value}">${type.label}</option>`
                        ).join('')}
                    </select>
                </div>
                <div class="form-group">
                    <label for="item-select">Item:</label>
                    <select id="item-select" required disabled>
                        <option value="">Select an Item</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="item-quantity">Quantity:</label>
                    <input type="number" id="item-quantity" min="1" value="1" required>
                </div>
                <div class="form-group">
                    <label for="item-equipped">Equipped:</label>
                    <input type="checkbox" id="item-equipped">
                </div>
                <div class="form-group">
                    <label for="item-notes">Notes:</label>
                    <textarea id="item-notes" rows="3"></textarea>
                </div>
                <div class="form-actions">
                    <button type="submit" id="save-item" class="btn btn-primary">Save</button>
                    <button type="button" id="cancel-item" class="btn">Cancel</button>
                </div>
            </form>
        </div>
    </div>
    `;
    
    const modalContainer = document.createElement('div');
    modalContainer.innerHTML = modalHTML;
    document.body.appendChild(modalContainer);
    
    setupModalEventListeners();
}

// Set up event listeners for the modal
function setupModalEventListeners() {
    const modal = document.getElementById('inventory-modal');
    const closeBtn = document.getElementById('close-inventory-modal');
    const cancelBtn = document.getElementById('cancel-item');
    const itemTypeSelect = document.getElementById('item-type');
    const itemSelect = document.getElementById('item-select');
    const inventoryForm = document.getElementById('inventory-form');
    
    closeBtn.addEventListener('click', closeModal);
    cancelBtn.addEventListener('click', closeModal);
    
    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            closeModal();
        }
    });
    
    itemTypeSelect.addEventListener('change', function() {
        const selectedType = this.value;
        console.log(`Item type selected: ${selectedType}`);
        
        if (selectedType) {
            loadItemsByType(selectedType);
        } else {
            itemSelect.disabled = true;
            itemSelect.innerHTML = '<option value="">Select an Item</option>';
        }
    });
    
    inventoryForm.addEventListener('submit', function(event) {
        event.preventDefault();
        saveInventoryItem();
    });
}

// Load items of the selected type
async function loadItemsByType(itemType) {
    try {
        const token = localStorage.getItem('authToken');
        const itemSelect = document.getElementById('item-select');
        
        console.log(`Loading items of type: ${itemType}`);
        itemSelect.innerHTML = '<option value="">Loading items...</option>';
        itemSelect.disabled = true;
        
        const endpoint = `/api/${getEndpointForType(itemType)}`;
        console.log(`Fetching from: ${endpoint}`);
        
        const response = await fetch(endpoint, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            }
        });
        
        if (!response.ok) {
            console.error(`API response error: ${response.status} ${response.statusText}`);
            throw new Error(`Failed to load items (${response.status})`);
        }
        
        const items = await response.json();
        console.log(`Received ${items.length} items:`, items);
        
        itemSelect.innerHTML = '<option value="">Select an Item</option>';
        
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
        console.error(`Error loading ${itemType} items:`, error);
        const itemSelect = document.getElementById('item-select');
        itemSelect.innerHTML = `<option value="">Error: ${error.message}</option>`;
        itemSelect.disabled = true;
    }
}

// Get the API endpoint for a given item type
function getEndpointForType(type) {
    const endpoints = {
        'weapon': 'weapons',
        'armor': 'armors',
        'shield': 'shields',
        'potion': 'potions',
        'magic_item': 'magic-items',
        'ring': 'rings',
        'ammo': 'ammo',
        'spell_scroll': 'spell-scrolls',
        'container': 'containers',
        'equipment': 'equipment'
    };
    
    return endpoints[type] || type;
}

// Open the modal for adding a new item
function openAddItemModal() {
    resetFormState();
    document.getElementById('modal-title').textContent = 'Add Item to Inventory';
    const modal = document.getElementById('inventory-modal');
    modal.style.display = 'block';
}

// Edit an existing inventory item
function editInventoryItem(itemId) {
    const item = inventoryState.items.find(i => i.id == itemId);
    if (!item) {
        console.error(`Item with ID ${itemId} not found`);
        return;
    }
    
    console.log(`Editing item: ${itemId}`, item);
    
    document.getElementById('modal-title').textContent = 'Edit Inventory Item';
    document.getElementById('item-quantity').value = item.quantity;
    document.getElementById('item-equipped').checked = item.is_equipped;
    document.getElementById('item-notes').value = item.notes || '';
    
    // Set the item type and disable it
    const typeSelect = document.getElementById('item-type');
    typeSelect.value = item.item_type;
    typeSelect.disabled = true;
    
    // Need to load the items of this type and select the current one
    loadItemsByType(item.item_type).then(() => {
        const itemSelect = document.getElementById('item-select');
        itemSelect.value = item.item_id;
        itemSelect.disabled = true;
    });
    
    // Store the item ID for the save function
    document.getElementById('inventory-form').setAttribute('data-editing-id', itemId);
    
    const modal = document.getElementById('inventory-modal');
    modal.style.display = 'block';
}

// Reset the form state
function resetFormState() {
    const form = document.getElementById('inventory-form');
    form.reset();
    form.removeAttribute('data-editing-id');
    
    const typeSelect = document.getElementById('item-type');
    const itemSelect = document.getElementById('item-select');
    
    typeSelect.disabled = false;
    itemSelect.disabled = true;
    itemSelect.innerHTML = '<option value="">Select an Item</option>';
}

// Close the modal
function closeModal() {
    const modal = document.getElementById('inventory-modal');
    if (modal) {
        modal.style.display = 'none';
        resetFormState();
    }
}

// Save an inventory item (add or update)
async function saveInventoryItem() {
    try {
        const form = document.getElementById('inventory-form');
        const editingItemId = form.getAttribute('data-editing-id');
        const isEditing = !!editingItemId;
        
        const token = localStorage.getItem('authToken');
        const itemType = document.getElementById('item-type').value;
        const itemId = document.getElementById('item-select').value;
        const quantity = parseInt(document.getElementById('item-quantity').value);
        const isEquipped = document.getElementById('item-equipped').checked;
        const notes = document.getElementById('item-notes').value;
        
        if (!isEditing && (!itemType || !itemId)) {
            alert('Please select an item type and item');
            return;
        }
        
        if (!quantity || quantity < 1) {
            alert('Quantity must be at least 1');
            return;
        }
        
        console.log(`Saving inventory item (${isEditing ? 'update' : 'new'}):`, {
            type: itemType,
            id: itemId,
            quantity,
            isEquipped,
            notes
        });
        
        let url, method, body;
        
        if (isEditing) {
            url = `/api/inventories/${inventoryState.inventory.id}/items/${editingItemId}`;
            method = 'PUT';
            body = {
                quantity: quantity,
                is_equipped: isEquipped,
                notes: notes
            };
        } else {
            url = `/api/inventories/${inventoryState.inventory.id}/items`;
            method = 'POST';
            body = {
                item_type: itemType,
                item_id: parseInt(itemId),
                quantity: quantity,
                is_equipped: isEquipped,
                notes: notes
            };
        }
        
        const response = await fetch(url, {
            method: method,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify(body)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response:', errorText);
            
            try {
                const errorData = JSON.parse(errorText);
                throw new Error(errorData.message || `Failed to save item (${response.status})`);
            } catch (e) {
                throw new Error(`Server error (${response.status}): ${errorText.substring(0, 100)}`);
            }
        }
        
        console.log('Item saved successfully');
        closeModal();
        await fetchInventory();
        
    } catch (error) {
        console.error('Error saving inventory item:', error);
        alert(`Error saving item: ${error.message}`);
    }
}

// Confirm deletion of an inventory item
function confirmDeleteItem(itemId) {
    if (confirm('Are you sure you want to remove this item from your inventory?')) {
        deleteInventoryItem(itemId);
    }
}

// Delete an inventory item
async function deleteInventoryItem(itemId) {
    try {
        const token = localStorage.getItem('authToken');
        
        console.log(`Deleting inventory item: ${itemId}`);
        
        const response = await fetch(`/api/inventories/${inventoryState.inventory.id}/items/${itemId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error(`Failed to delete item: ${response.statusText}`);
        }
        
        console.log('Item deleted successfully');
        await fetchInventory();
        
    } catch (error) {
        console.error('Error deleting inventory item:', error);
        alert(`Error deleting item: ${error.message}`);
    }
}

// Get the character ID from the URL
function getCharacterIdFromURL() {
    const pathParts = window.location.pathname.split('/');
    return pathParts[pathParts.length - 1];
}

// Display an error in the inventory section
function displayInventoryError(message) {
    const inventoryTable = document.getElementById('inventory-table');
    const inventoryLoading = document.getElementById('inventory-loading');
    const inventoryEmpty = document.getElementById('inventory-empty');
    
    if (inventoryLoading) {
        inventoryLoading.style.display = 'none';
    }
    
    if (inventoryTable) {
        inventoryTable.style.display = 'none';
    }
    
    if (inventoryEmpty) {
        inventoryEmpty.style.display = 'block';
        inventoryEmpty.innerHTML = `<p class="error">Error: ${message}</p><p>Please try refreshing the page.</p>`;
    }
}

// Initialize the inventory functionality
function initInventory() {
    const characterId = getCharacterIdFromURL();
    if (!characterId) return;
    
    console.log('Initializing inventory for character:', characterId);
    
    createInventoryModal();
    
    const addItemBtn = document.getElementById('btn-add-item');
    if (addItemBtn) {
        addItemBtn.addEventListener('click', openAddItemModal);
    }
    
    fetchInventory(characterId);
}

async function toggleEquippedStatus(itemId, currentStatus) {
    try {
        const token = localStorage.getItem('authToken');
        const inventoryID = await getInventoryID();
        
        if (!inventoryID) {
            throw new Error('Could not find inventory');
        }
        
        const response = await fetch(`/api/inventories/${inventoryID}/items/${itemId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                is_equipped: !currentStatus
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to update item equipped status');
        }
        
        // Refresh the inventory display
        fetchInventory();
        
        // If we're on the combat tab, refresh that too
        if (document.getElementById('combat-tab').classList.contains('active')) {
            loadCombatData();
        }
        
    } catch (error) {
        console.error('Error toggling equipped status:', error);
        alert('Failed to update item: ' + error.message);
    }
}

// Call this when the inventory tab is clicked
document.addEventListener('DOMContentLoaded', function() {
    const inventoryTab = document.querySelector('.tab-item[data-tab="inventory-tab"]');
    if (inventoryTab) {
        console.log('Setting up inventory tab click handler');
        inventoryTab.addEventListener('click', function() {
            console.log('Inventory tab clicked');
            initInventory();
        });
    }
    
    // Initialize if inventory tab is active by default
    if (document.getElementById('inventory-tab') && 
        document.getElementById('inventory-tab').classList.contains('active')) {
        console.log('Inventory tab is active by default');
        initInventory();
    }
});