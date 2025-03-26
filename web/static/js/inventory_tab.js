document.addEventListener('DOMContentLoaded', function() {
    const inventoryButton = document.querySelector('.tab-button[data-tab="inventory"]');
    const pathSegments = window.location.pathname.split('/');
    // Verify we're on a character detail page
    if (!pathSegments.includes('view')) {
        console.error("Not on a character detail page");
        return;
    }
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];
    // Ensure characterId is numeric
    if (!characterId || isNaN(parseInt(characterId))) {
        console.error("Invalid character ID:", characterId);
        return;
    }
    let inventoryId = null;
    
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
            // Use the consistent inventory endpoint 
            console.log(`Fetching inventory data from /api/inventories/character/${characterId}`);
            const response = await fetch(`/api/inventories/character/${characterId}`);
            
            // Log response status and text for debugging
            console.log(`Response status: ${response.status}`);
            
            if (!response.ok) {
                // Try to get the error message from the response
                let errorMessage = 'Failed to fetch inventory';
                try {
                    const errorData = await response.json();
                    errorMessage = errorData.message || errorMessage;
                } catch (e) {
                    // If we can't parse the error as JSON, try to get the response text
                    const errorText = await response.text();
                    if (errorText) errorMessage = errorText;
                }
                throw new Error(errorMessage);
            }
            
            const responseData = await response.json();
            console.log("API Response Structure:", responseData); // Debug log
            
            // Safely extract data with fallbacks
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
                const details = item.item_details;
                if (!details) continue;
                
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
                processedInventory[category].push({
                    id: item.id,
                    item_id: item.item_id,
                    name: details.name,
                    description: details.description,
                    weight: details.weight,
                    cost: details.cost,
                    quantity: item.quantity,
                    is_equipped: item.is_equipped,
                    notes: item.notes,
                    damage: details.damage,
                    properties: details.properties,
                    ac: details.ac,
                    defense_modifier: details.defense_modifier,
                    capacity: details.capacity,
                    value: details.value
                });
            }
            
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
                    itemDetails += `<span class="item-damage">Damage: ${item.damage}</span>`;
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
                    itemDetails += `<span class="item-capacity">Capacity: ${item.capacity} lb</span>`;
                    break;
                case 'treasure':
                    itemDetails += `<span class="item-value">Value: ${item.value} gp</span>`;
                    break;
            }
    
            // Determine if this item type can be equipped
            const equippableTypes = ['weapons', 'armor', 'shields', 'rings'];
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
    }
    
    async function populateItemTypes() {
        const itemTypeSelect = document.getElementById('itemType');
        const itemIdSelect = document.getElementById('itemId');
        
        // Clear options
        itemIdSelect.innerHTML = '<option value="" disabled selected>Select Item Type First</option>';
        
        // Listen for change to populate items based on type
        itemTypeSelect.addEventListener('change', async () => {
            const itemType = itemTypeSelect.value;
            itemIdSelect.innerHTML = '<option value="" disabled selected>Loading items...</option>';
            
            try {
                // Map the item type to the corresponding API endpoint
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
    
    // Handle form submit for adding an item
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

});

function setupItemActionListeners() {
    // Existing event listeners for edit and remove buttons
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
    
    // New event listener for equipped toggle
    document.querySelectorAll('.toggle-equipped').forEach(checkbox => {
        checkbox.addEventListener('change', (e) => {
            const itemElement = e.target.closest('.inventory-item');
            const itemId = itemElement.dataset.id;
            const isEquipped = e.target.checked;
            
            // Toggle visual state immediately for better UX
            const label = e.target.nextElementSibling;
            label.textContent = isEquipped ? 'Equipped' : 'Equip';
            
            if (isEquipped) {
                itemElement.classList.add('item-equipped');
            } else {
                itemElement.classList.remove('item-equipped');
            }
            
            // Send update to server
            updateItemEquipped(itemId, isEquipped);
        });
    });
}

// Add this new function to handle the AJAX request
async function updateItemEquipped(itemId, isEquipped) {
    try {
        const response = await fetch(`/api/inventories/${inventoryId}/items/${itemId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                is_equipped: isEquipped
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to update item');
        }
        
        // You could refresh the inventory data here if needed
        // or handle any specific update logic
        
    } catch (error) {
        console.error('Error updating item:', error);
        alert('Failed to update equipment status: ' + error.message);
        // Revert the UI change on error
        fetchInventory();
    }
}