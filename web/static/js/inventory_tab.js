document.addEventListener('DOMContentLoaded', function() {
    const inventoryButton = document.querySelector('.tab-button[data-tab="inventory"]');
    const pathSegments = window.location.pathname.split('/');
    if (!pathSegments.includes('view')) {
        console.error("Not on a character detail page");
        return;
    }
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];
    if (!characterId || isNaN(parseInt(characterId))) {
        console.error("Invalid character ID:", characterId);
        return;
    }
    let inventoryId = null;
    let equipmentStatus = null;
    
    if (inventoryButton) {
        inventoryButton.addEventListener('click', () => {
            const inventoryTab = document.getElementById('inventory');
            if (inventoryTab.querySelector('.loading-container') || inventoryTab.innerHTML.trim() === '') {
                fetchInventory();
            }
        });
        // Load inventory if it's the active tab on page load
        if (inventoryButton.classList.contains('active')) {
            fetchInventory();
        }
    }
    
    
    async function fetchInventory() {
        const inventoryTab = document.getElementById('inventory');
        inventoryTab.innerHTML = `
            <div class="loading-container">
                <div class="loading-spinner"></div>
                <p>Loading inventory...</p>
            </div>
        `;
        try {
            console.log(`Fetching inventory data from /api/inventories/character/${characterId}`);
            const response = await fetch(`/api/inventories/character/${characterId}`);
            console.log(`Response status: ${response.status}`);
            if (!response.ok) {
                let errorMessage = 'Failed to fetch inventory';
                try {
                    const errorData = await response.json();
                    errorMessage = errorData.message || errorMessage;
                } catch (e) {
                    const errorText = await response.text();
                    if (errorText) errorMessage = errorText;
                }
                throw new Error(errorMessage);
            }
            const responseData = await response.json();
            console.log("API Response Structure:", responseData);
            
            const inventory = responseData.inventory || {};
            const enrichedItems = responseData.items || [];
            const encumbrance = responseData.encumbrance || {
                total_weight: inventory.current_weight || 0,
                status: {
                    current_weight: inventory.current_weight || 0,
                    maximum_capacity: inventory.max_weight || 100,
                    weight_remaining: (inventory.max_weight || 100) - (inventory.current_weight || 0),
                    percent_full: Math.min(100, ((inventory.current_weight || 0) / (inventory.max_weight || 100)) * 100),
                    encumbered: false,
                    heavy_encumbered: false,
                    overloaded: false
                },
                thresholds: {
                    maximum_capacity: inventory.max_weight || 100,
                    base_encumbered: 75,
                    base_heavy_encumbered: 150
                }
            };
            
            inventoryId = inventory.id;
            if (inventoryId) {
                try {
                    console.log(`Fetching equipment status from /api/characters/${characterId}/equipment-status`);
                    const statusResponse = await fetch(`/api/characters/${characterId}/equipment-status`);
                    if (statusResponse.ok) {
                        equipmentStatus = await statusResponse.json();
                        console.log("Equipment status:", equipmentStatus);
                    } else {
                        console.warn(`Failed to fetch equipment status: ${statusResponse.status}`);
                        equipmentStatus = { equipped_slots: {}, available_slots: [] };
                    }
                } catch (statusError) {
                    console.error("Error fetching equipment status:", statusError);
                    equipmentStatus = { equipped_slots: {}, available_slots: [] };
                }
            } else {
                console.warn("No inventory ID available, skipping equipment status fetch");
                equipmentStatus = { equipped_slots: {}, available_slots: [] };
            }
            
            const processedInventory = {
                weapons: [],
                armor: [],
                shields: [],
                potions: [],
                magic_items: [],
                rings: [],
                ammunition: [],
                spell_scrolls: [],
                containers: [],
                equipment: [],
                treasure: inventory.treasure ? [inventory.treasure] : []
            };
            
            for (const item of enrichedItems) {
                const details = item.item_details || {};
                const typeMap = {
                    'weapon': 'weapons',
                    'armor': 'armor',
                    'shield': 'shields',
                    'potion': 'potions',
                    'magic_item': 'magic_items',
                    'ring': 'rings',
                    'ammo': 'ammunition',
                    'spell_scroll': 'spell_scrolls',
                    'container': 'containers',
                    'equipment': 'equipment'
                };
                
                const category = typeMap[item.item_type] || 'equipment';
                
                // Use fallback values if details are missing
                const displayName = details.name || `Unknown ${item.item_type} (ID: ${item.item_id})`;
                
                processedInventory[category].push({
                    id: item.id,
                    item_id: item.item_id,
                    name: displayName,
                    description: details.description || '',
                    weight: details.weight || 0,
                    cost: details.cost || 0,
                    quantity: item.quantity,
                    is_equipped: item.is_equipped,
                    slot: item.slot,
                    notes: item.notes,
                    damage: details.damage || '',
                    properties: details.properties || '',
                    ac: details.ac || 0,
                    defense_modifier: details.defense_modifier || 0,
                    capacity: details.capacity || 0,
                    value: details.value || 0
                });
            }
            
            console.log("Inventory passed to renderInventory:", processedInventory);
            console.log("Equipment items to render:", processedInventory.equipment || []);
            console.log("Number of equipment items:", (processedInventory.equipment || []).length);
            
            renderInventory(processedInventory, encumbrance);
        } catch (error) {
            console.error('Error:', error);
            inventoryTab.innerHTML = `
                <div class="error-message">
                    Failed to load inventory: ${error.message}
                </div>
            `;
        }
    }
    
    function renderInventory(inventory, encumbrance) {
        const inventoryTab = document.getElementById('inventory');
        const percentFull = Math.min(100, (encumbrance.status.current_weight / encumbrance.thresholds.maximum_capacity) * 100);
        const encumbranceHTML = `
            <div class="encumbrance-summary">
                <div class="encumbrance-bar">
                    <div class="encumbrance-fill" style="width: ${percentFull}%"></div>
                </div>
                <div class="encumbrance-stats">
                    <span>Weight: ${encumbrance.status.current_weight.toFixed(1)} / ${encumbrance.thresholds.maximum_capacity.toFixed(1)}</span>
                    <span class="encumbrance-status">Status: ${getEncumbranceStatus(encumbrance)}</span>
                    ${getEncumbrancePenalties(encumbrance)}
                </div>
            </div>
        `;
        
        let equipmentDisplayHTML = '';
        if (equipmentStatus && Object.keys(equipmentStatus.equipped_slots).length > 0) {
            const slotDisplayNames = {
                'head': 'Head',
                'body': 'Body',
                'main_hand': 'Main Hand',
                'off_hand': 'Off Hand',
                'ring_left': 'Left Ring',
                'ring_right': 'Right Ring',
                'neck': 'Neck',
                'back': 'Back',
                'belt': 'Belt',
                'feet': 'Feet',
                'hands': 'Hands'
            };
            equipmentDisplayHTML = `
                <div class="equipment-display">
                    <h3>Equipped Items</h3>
                    <div class="equipment-slots">
                        ${Object.entries(equipmentStatus.equipped_slots).map(([slot, item]) => `
                            <div class="equipment-slot">
                                <div class="slot-name">${slotDisplayNames[slot] || slot}</div>
                                <div class="slot-item">${item.name}</div>
                            </div>
                        `).join('')}
                    </div>
                </div>
            `;
        }
        
        let sectionsHTML = '';
        // Define sections to display based on API response structure
        const sections = [
            { title: "Weapons", items: inventory.weapons || [] },
            { title: "Armor", items: inventory.armor || [] },
            { title: "Shields", items: inventory.shields || [] },
            { title: "Potions", items: inventory.potions || [] },
            { title: "Magic Items", items: inventory.magic_items || [] },
            { title: "Rings", items: inventory.rings || [] },
            { title: "Ammunition", items: inventory.ammunition || [] },
            { title: "Spell Scrolls", items: inventory.spell_scrolls || [] },
            { title: "Containers", items: inventory.containers || [] },
            { title: "Equipment", items: inventory.equipment || [] },
            { title: "Treasure", items: inventory.treasure || [] }
        ];
        
        console.log("Sections to be rendered:", sections.map(s => ({
            title: s.title,
            count: s.items.length
        })));
        
        for (const section of sections) {
            if (section.items.length > 0) {
                sectionsHTML += `
                    <div class="inventory-section">
                        <h3>${section.title}</h3>
                        <div class="inventory-items">
                            ${renderItems(section.items, section.title.toLowerCase())}
                        </div>
                    </div>
                `;
            }
        }
        
        // If no items found
        if (sectionsHTML === '') {
            sectionsHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Items Found</h3>
                    <p class="empty-description">This character doesn't have any items in their inventory yet.</p>
                    <button id="addItemBtn" class="btn btn-primary">Add Item</button>
                </div>
            `;
        } else {
            // Add button to add more items
            sectionsHTML += `
                <div class="inventory-actions">
                    <button id="addItemBtn" class="btn btn-primary">Add Item</button>
                </div>
            `;
        }
        
        // Combine all components
        inventoryTab.innerHTML = `
            <h2>Inventory</h2>
            ${encumbranceHTML}
            ${equipmentDisplayHTML}
            ${sectionsHTML}
        `;
        
        // Add event listener for Add Item button
        const addItemBtn = document.getElementById('addItemBtn');
        if (addItemBtn) {
            addItemBtn.addEventListener('click', () => {
                document.getElementById('inventoryModal').style.display = 'flex';
                populateItemTypes();
            });
        }
        
        // Add event listeners for item actions
        setupItemActionListeners();
    }
    
    function getEncumbranceStatus(encumbrance) {
        if (encumbrance.status.overloaded) {
            return "Overloaded";
        } else if (encumbrance.status.heavy_encumbered) {
            return "Heavily Encumbered";
        } else if (encumbrance.status.encumbered) {
            return "Encumbered";
        } else {
            return "Unencumbered";
        }
    }
    
    function getEncumbrancePenalties(encumbrance) {
        if (encumbrance.status.overloaded) {
            return `<span class="encumbrance-penalties">Penalties: Cannot move</span>`;
        } else if (encumbrance.status.heavy_encumbered) {
            return `<span class="encumbrance-penalties">Penalties: Move at 1/2 speed, -2 AC</span>`;
        } else if (encumbrance.status.encumbered) {
            return `<span class="encumbrance-penalties">Penalties: Move at 3/4 speed</span>`;
        }
        return '';
    }
    
    function renderItems(items, type) {
        let itemsHTML = '';
        for (const item of items) {
            let itemDetails = '';
            // Common properties
            if (item.weight) {
                itemDetails += `<span class="item-weight">${item.weight} lb</span>`;
            }
            if (item.cost) {
                itemDetails += `<span class="item-cost">${item.cost} gp</span>`;
            }
            // Type-specific properties
            switch (type) {
                case 'weapons':
                    itemDetails += `<span class="item-damage">Damage: ${item.damage || 'N/A'}</span>`;
                    if (item.properties) {
                        itemDetails += `<span class="item-properties">${item.properties}</span>`;
                    }
                    break;
                case 'armor':
                    itemDetails += `<span class="item-ac">AC: ${item.ac || 'unknown'}</span>`;
                    break;
                case 'shields':
                    itemDetails += `<span class="item-defense">Defense: +${item.defense_modifier}</span>`;
                    break;
                case 'potions':
                case 'magic_items':
                case 'rings':
                case 'spell_scrolls':
                    if (item.description) {
                        itemDetails += `<span class="item-description">${item.description}</span>`;
                    }
                    break;
                case 'ammunition':
                    itemDetails += `<span class="item-quantity">Qty: ${item.quantity}</span>`;
                    break;
                case 'containers':
                    itemDetails += `<span class="item-capacity">Capacity: ${item.capacity || 'N/A'} lb</span>`;
                    break;
                case 'equipment':
                    // Add specific handling for equipment items
                    if (item.description) {
                        itemDetails += `<span class="item-description">${item.description}</span>`;
                    }
                    break;
                case 'treasure':
                    itemDetails += `<span class="item-value">Value: ${item.value} gp</span>`;
                    break;
            }
            
            // If quantity is greater than 1, display it (for all item types)
            if (item.quantity > 1 && type !== 'ammunition') {
                itemDetails += `<span class="item-quantity">Qty: ${item.quantity}</span>`;
            }
            
            // Determine if this item type can be equipped
            const equippableTypes = ['weapons', 'armor', 'shields', 'rings', 'equipment'];
            const canEquip = equippableTypes.includes(type);
            
            // Create equipped toggle if appropriate
            const equippedToggle = canEquip ? `
                <label class="equip-toggle">
                    <input type="checkbox" class="toggle-equipped" ${item.is_equipped ? 'checked' : ''}>
                    <span class="toggle-label">${item.is_equipped ? 'Equipped' : 'Equip'}</span>
                </label>
            ` : '';
            
            itemsHTML += `
                <div class="inventory-item ${item.is_equipped ? 'item-equipped' : ''}" data-id="${item.id}" data-type="${type}" data-item-id="${item.item_id}">
                    <div class="item-header">
                        <h4 class="item-name">${item.name}</h4>
                        <div class="item-actions">
                            ${equippedToggle}
                            <button class="item-action edit-item" title="Edit Item">✏️</button>
                            <button class="item-action remove-item" title="Remove Item">❌</button>
                        </div>
                    </div>
                    <div class="item-details">
                        ${itemDetails}
                    </div>
                    ${item.slot ? `<div class="item-slot">Slot: ${item.slot}</div>` : ''}
                    ${item.notes ? `<div class="item-notes">Notes: ${item.notes}</div>` : ''}
                </div>
            `;
        }
        return itemsHTML;
    }
    
    function setupItemActionListeners() {
        // Add event listeners for edit and remove buttons
        document.querySelectorAll('.edit-item').forEach(button => {
            button.addEventListener('click', (e) => {
                const itemElement = e.target.closest('.inventory-item');
                const itemId = itemElement.dataset.id;
                const itemType = itemElement.dataset.type;
                editItem(itemId, itemType);
            });
        });
        document.querySelectorAll('.remove-item').forEach(button => {
            button.addEventListener('click', (e) => {
                const itemElement = e.target.closest('.inventory-item');
                const itemId = itemElement.dataset.id;
                const itemType = itemElement.dataset.type;
                removeItem(itemId, itemType);
            });
        });
        document.querySelectorAll('.toggle-equipped').forEach(checkbox => {
            checkbox.addEventListener('change', (e) => {
                const itemElement = e.target.closest('.inventory-item');
                toggleEquipItem(itemElement);
                e.preventDefault();
            });
        });
    }
    
    function toggleEquipItem(itemElement) {
        const itemId = itemElement.dataset.id;
        const itemType = itemElement.dataset.type;
        const isCurrentlyEquipped = itemElement.classList.contains('item-equipped');
        console.log(`Toggling equip for item ${itemId} (${itemType}): currently ${isCurrentlyEquipped ? 'equipped' : 'unequipped'}`);
        if (!isCurrentlyEquipped) {
            let availableSlots = [];
            if (equipmentStatus) {
                switch (itemType) {
                    case 'weapons':
                        const isTwoHanded = itemElement.querySelector('.item-properties')?.textContent.toLowerCase().includes('two-handed');
                        if (isTwoHanded) {
                            if (equipmentStatus.available_slots.includes('main_hand') &&
                                equipmentStatus.available_slots.includes('off_hand')) {
                                availableSlots = ['main_hand'];
                            }
                        } else {
                            availableSlots = equipmentStatus.available_slots.filter(
                                slot => slot === 'main_hand' || slot === 'off_hand'
                            );
                        }
                        break;
                    case 'shields':
                        if (equipmentStatus.available_slots.includes('off_hand')) {
                            availableSlots = ['off_hand'];
                        }
                        break;
                    case 'armor':
                        if (equipmentStatus.available_slots.includes('body')) {
                            availableSlots = ['body'];
                        }
                        break;
                    case 'rings':
                        availableSlots = equipmentStatus.available_slots.filter(
                            slot => slot === 'ring_left' || slot === 'ring_right'
                        );
                        break;
                    case 'equipment':
                        // Allow equipment to be equipped in various slots based on what's available
                        availableSlots = equipmentStatus.available_slots;
                        break;
                    default:
                        console.log(`No client-side slot mapping for ${itemType}, letting server decide`);
                }
            }
            if (availableSlots.length === 0) {
                updateEquipStatus(itemId, true);
            } else if (availableSlots.length === 1) {
                updateEquipStatus(itemId, true, availableSlots[0]);
            } else {
                showSlotSelector(availableSlots, (selectedSlot) => {
                    updateEquipStatus(itemId, true, selectedSlot);
                });
            }
        } else {
            updateEquipStatus(itemId, false);
        }
    }
    
    function updateEquipStatus(itemId, equip, slot = null) {
        const payload = {
            is_equipped: equip
        };
        if (equip && slot) {
            payload.slot = slot;
        } else if (!equip) {
            payload.slot = "";
        }
        console.log(`Updating item ${itemId}: equip=${equip}, slot=${payload.slot || 'none'}`);
        // Call the API to update the item
        fetch(`/api/inventories/${inventoryId}/items/${itemId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })
        .then(response => {
            console.log(`API response status: ${response.status}`);
            if (!response.ok) {
                // If error, try to parse error details
                return response.json().then(data => {
                    console.error("Error response:", data);
                    throw new Error(data.message || 'Failed to update item');
                }).catch(err => {
                    // If can't parse JSON, use status text
                    if (err instanceof SyntaxError) {
                        throw new Error(`Server error: ${response.status} ${response.statusText}`);
                    }
                    throw err;
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Item updated successfully:", data);
            fetchInventory(); // Refresh the inventory display
        })
        .catch(error => {
            console.error('Error updating item:', error);
            alert(`Failed to update item: ${error.message}`);
            fetchInventory(); // Still refresh to ensure UI matches server state
        });
    }
    
    function showSlotSelector(availableSlots, callback) {
        const slotDisplayNames = {
            'head': 'Head',
            'body': 'Body (Armor)',
            'main_hand': 'Main Hand',
            'off_hand': 'Off Hand',
            'ring_left': 'Left Ring Finger',
            'ring_right': 'Right Ring Finger',
            'neck': 'Neck',
            'back': 'Back',
            'belt': 'Belt',
            'feet': 'Feet',
            'hands': 'Hands'
        };
        const modal = document.createElement('div');
        modal.classList.add('modal');
        modal.style.display = 'flex';
        modal.innerHTML = `
            <div class="modal-content" style="max-width: 400px;">
                <h2 class="modal-title">Select Equipment Slot</h2>
                <div class="slot-options">
                    ${availableSlots.map(slot => `
                        <button class="btn btn-primary slot-btn" data-slot="${slot}" style="margin-bottom: 8px; width: 100%; text-align: left;">
                            ${slotDisplayNames[slot] || slot}
                        </button>
                    `).join('')}
                </div>
                <div class="modal-actions">
                    <button class="btn btn-secondary" id="cancelSlotBtn">Cancel</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
        const slotButtons = modal.querySelectorAll('.slot-btn');
        slotButtons.forEach(button => {
            button.addEventListener('click', () => {
                const selectedSlot = button.dataset.slot;
                document.body.removeChild(modal);
                callback(selectedSlot);
            });
        });
        document.getElementById('cancelSlotBtn').addEventListener('click', () => {
            document.body.removeChild(modal);
        });
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                document.body.removeChild(modal);
            }
        });
    }
    
    async function populateItemTypes() {
        const itemTypeSelect = document.getElementById('itemType');
        const itemIdSelect = document.getElementById('itemId');
        itemIdSelect.innerHTML = '<option value="" disabled selected>Select Item Type First</option>';
        itemTypeSelect.addEventListener('change', async () => {
            const itemType = itemTypeSelect.value;
            itemIdSelect.innerHTML = '<option value="" disabled selected>Loading items...</option>';
            try {
                const endpointMap = {
                    'weapon': '/api/weapons',
                    'armor': '/api/armors',
                    'shield': '/api/shields',
                    'potion': '/api/potions',
                    'magic_item': '/api/magic-items',
                    'ring': '/api/rings',
                    'ammo': '/api/ammo',
                    'spell_scroll': '/api/spell-scrolls',
                    'container': '/api/containers',
                    'equipment': '/api/equipment'
                };
                const endpoint = endpointMap[itemType];
                if (!endpoint) {
                    throw new Error('Invalid item type');
                }
                const response = await fetch(endpoint);
                if (!response.ok) {
                    throw new Error('Failed to fetch items');
                }
                const items = await response.json();
                itemIdSelect.innerHTML = '<option value="" disabled selected>Select an Item</option>';
                items.forEach(item => {
                    const option = document.createElement('option');
                    option.value = item.id;
                    option.textContent = item.name;
                    itemIdSelect.appendChild(option);
                });
            } catch (error) {
                console.error('Error:', error);
                itemIdSelect.innerHTML = '<option value="" disabled selected>Error loading items</option>';
            }
        });
    }
    
    async function editItem(itemId, itemType) {
        try {
            const response = await fetch(`/api/inventories/${inventoryId}/items/${itemId}`);
            if (!response.ok) {
                throw new Error('Failed to fetch item details');
            }
            const item = await response.json();
            const newQuantity = prompt('Enter new quantity:', item.quantity || 1);
            if (newQuantity === null) return;
            const updateResponse = await fetch(`/api/inventories/${inventoryId}/items/${itemId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    quantity: parseInt(newQuantity)
                })
            });
            if (!updateResponse.ok) {
                throw new Error('Failed to update item');
            }
            fetchInventory();
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to edit item: ' + error.message);
        }
    }
    
    async function removeItem(itemId, itemType) {
        const confirmDelete = confirm('Are you sure you want to remove this item from your inventory?');
        if (!confirmDelete) return;
        try {
            const response = await fetch(`/api/inventories/${inventoryId}/items/${itemId}`, {
                method: 'DELETE'
            });
            if (!response.ok) {
                throw new Error('Failed to remove item');
            }
            fetchInventory();
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to remove item: ' + error.message);
        }
    }
    
    const inventoryForm = document.getElementById('inventoryForm');
    if (inventoryForm) {
        inventoryForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const itemType = document.getElementById('itemType').value;
            const itemId = document.getElementById('itemId').value;
            const quantity = document.getElementById('quantity').value;
            if (!itemType || !itemId || !quantity) {
                alert('Please fill in all fields');
                return;
            }
            try {
                const response = await fetch(`/api/inventories/${inventoryId}/items`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        item_type: itemType,
                        item_id: parseInt(itemId),
                        quantity: parseInt(quantity)
                    })
                });
                if (!response.ok) {
                    throw new Error('Failed to add item to inventory');
                }
                document.getElementById('inventoryModal').style.display = 'none';
                fetchInventory();
            } catch (error) {
                console.error('Error:', error);
                alert('Failed to add item: ' + error.message);
            }
        });
    }
    
    const cancelInventoryBtn = document.getElementById('cancelInventoryBtn');
    if (cancelInventoryBtn) {
        cancelInventoryBtn.addEventListener('click', () => {
            document.getElementById('inventoryModal').style.display = 'none';
        });
    }
    
    window.addEventListener('click', (e) => {
        const inventoryModal = document.getElementById('inventoryModal');
        if (e.target === inventoryModal) {
            inventoryModal.style.display = 'none';
        }
    });
});

document.head.insertAdjacentHTML('beforeend', `
<style>
.equipment-display {
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 1.5rem;
}
.equipment-slots {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 0.75rem;
}
.equipment-slot {
    background-color: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
    padding: 0.75rem;
}
.slot-name {
    font-size: 0.8rem;
    color: #999;
    margin-bottom: 0.25rem;
}
.slot-item {
    font-weight: bold;
    color: var(--primary-color);
}
.item-slot {
    font-size: 0.8rem;
    color: #74b9ff;
    margin-top: 0.5rem;
}
.item-notes {
    font-size: 0.8rem;
    color: #999;
    margin-top: 0.5rem;
    font-style: italic;
}
.slot-options {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
}
</style>
`);