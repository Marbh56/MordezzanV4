document.addEventListener('DOMContentLoaded', function() {
    const inventoryTab = document.querySelector('.tab-item[data-tab="inventory-tab"]');
    if (inventoryTab) {
        inventoryTab.addEventListener('click', function() {
            console.log('Inventory tab clicked');
            setupInventoryAddItem();
            fetchInventory();
        });
    }
    if (document.getElementById('inventory-tab') &&
        document.getElementById('inventory-tab').classList.contains('active')) {
        console.log('Inventory tab is active by default, initializing');
        setupInventoryAddItem();
        fetchInventory();
    }
});

// Encumbrance thresholds structure
class EncumbranceThresholds {
    constructor(baseEncumbered, baseHeavyEncumbered, maximumCapacity) {
        this.baseEncumbered = baseEncumbered;
        this.baseHeavyEncumbered = baseHeavyEncumbered;
        this.maximumCapacity = maximumCapacity;
    }
}

// Calculate encumbrance thresholds based on strength and constitution
function calculateEncumbranceThresholds(strength, constitution) {
    const baseThresholds = {
        baseEncumbered: 75,
        baseHeavyEncumbered: 150,
        maximumCapacity: 300 // Base maximum capacity
    };
    
    // Calculate strength modifier (in pounds)
    let strMod = 0;
    let maxMod = 0;
    
    if (strength <= 6) {
        strMod = -25;
        maxMod = -100;
    } else if (strength >= 7 && strength <= 8) {
        strMod = -15;
        maxMod = -50;
    } else if (strength >= 13 && strength <= 14) {
        strMod = 15;
        maxMod = 50;
    } else if (strength >= 15 && strength <= 16) {
        strMod = 25;
        maxMod = 100;
    } else if (strength === 17) {
        strMod = 35;
        maxMod = 150;
    } else if (strength === 18) {
        strMod = 50;
        maxMod = 200;
    }
    
    // Calculate constitution modifier (in pounds)
    let conMod = 0;
    let conMaxMod = 0;
    
    if (constitution <= 6) {
        conMod = -10;
        conMaxMod = -25;
    } else if (constitution >= 7 && constitution <= 8) {
        conMod = -5;
        conMaxMod = -15;
    } else if (constitution >= 13 && constitution <= 14) {
        conMod = 5;
        conMaxMod = 15;
    } else if (constitution >= 15 && constitution <= 16) {
        conMod = 10;
        conMaxMod = 25;
    } else if (constitution >= 17) {
        conMod = 15;
        conMaxMod = 35;
    }
    
    // Apply modifiers
    baseThresholds.baseEncumbered += strMod + conMod;
    baseThresholds.baseHeavyEncumbered += (strMod + conMod) * 2;
    baseThresholds.maximumCapacity += maxMod + conMaxMod;
    
    // Ensure minimum thresholds
    if (baseThresholds.baseEncumbered < 40) {
        baseThresholds.baseEncumbered = 40;
    }
    if (baseThresholds.baseHeavyEncumbered < 60) {
        baseThresholds.baseHeavyEncumbered = 60;
    }
    if (baseThresholds.maximumCapacity < 100) {
        baseThresholds.maximumCapacity = 100;
    }
    
    return new EncumbranceThresholds(
        baseThresholds.baseEncumbered,
        baseThresholds.baseHeavyEncumbered,
        baseThresholds.maximumCapacity
    );
}

// Determine encumbrance status from current weight and thresholds
function getEncumbranceStatus(currentWeight, thresholds) {
    if (currentWeight > thresholds.maximumCapacity) {
        return 'overloaded';
    } else if (currentWeight > thresholds.baseHeavyEncumbered) {
        return 'heavily-encumbered';
    } else if (currentWeight > thresholds.baseEncumbered) {
        return 'encumbered';
    } else {
        return 'unencumbered';
    }
}

// Get encumbrance modifiers based on status
function getEncumbranceModifiers(encumbranceStatus) {
    switch (encumbranceStatus) {
        case 'encumbered':
            return { movementPenalty: -10, acPenalty: -1 };
        case 'heavily-encumbered':
            return { movementPenalty: -20, acPenalty: -2 };
        case 'overloaded':
            return { movementPenalty: -30, acPenalty: -3 };
        default:
            return { movementPenalty: 0, acPenalty: 0 };
    }
}

// Update the inventory display with encumbrance information
function updateEncumbranceDisplay(currentWeight, thresholds) {
    const encumbranceStatus = getEncumbranceStatus(currentWeight, thresholds);
    const encumbranceElement = document.getElementById('encumbrance-status');
    const weightElement = document.getElementById('total-weight');
    const maxWeightElement = document.getElementById('max-weight');
    const encumberedThresholdElement = document.getElementById('encumbered-threshold');
    const heavyThresholdElement = document.getElementById('heavy-threshold');
    const maxThresholdElement = document.getElementById('max-threshold');
    
    if (weightElement) {
        weightElement.textContent = currentWeight.toFixed(2);
    }
    
    if (maxWeightElement) {
        maxWeightElement.textContent = thresholds.maximumCapacity.toFixed(2);
    }
    
    // Update threshold markers if they exist
    if (encumberedThresholdElement) {
        encumberedThresholdElement.textContent = thresholds.baseEncumbered;
        encumberedThresholdElement.style.left = `${(thresholds.baseEncumbered / thresholds.maximumCapacity) * 100}%`;
    }
    
    if (heavyThresholdElement) {
        heavyThresholdElement.textContent = thresholds.baseHeavyEncumbered;
        heavyThresholdElement.style.left = `${(thresholds.baseHeavyEncumbered / thresholds.maximumCapacity) * 100}%`;
    }
    
    if (maxThresholdElement) {
        maxThresholdElement.textContent = thresholds.maximumCapacity;
        maxThresholdElement.style.left = '100%';
    }
    
    if (encumbranceElement) {
        // Clear existing classes
        encumbranceElement.classList.remove('unencumbered', 'encumbered', 'heavily-encumbered', 'overloaded');
        
        // Add current status class
        encumbranceElement.classList.add(encumbranceStatus);
        
        // Update text
        let statusText = 'Unencumbered';
        if (encumbranceStatus === 'encumbered') {
            statusText = 'Encumbered (MV -10, AC -1)';
        } else if (encumbranceStatus === 'heavily-encumbered') {
            statusText = 'Heavily Encumbered (MV -20, AC -2)';
        } else if (encumbranceStatus === 'overloaded') {
            statusText = 'Overloaded (Cannot move)';
        }
        
        encumbranceElement.textContent = statusText;
    }
    
    // Update progress bar if it exists
    const progressBar = document.getElementById('encumbrance-bar');
    if (progressBar) {
        const percentage = Math.min(100, (currentWeight / thresholds.maximumCapacity) * 100);
        progressBar.style.width = `${percentage}%`;
        
        // Update progress bar color based on encumbrance
        progressBar.className = 'encumbrance-progress-bar';
        progressBar.classList.add(encumbranceStatus);
    }
    
    // Apply encumbrance effects to character
    applyEncumbranceEffects(encumbranceStatus);
    
    return encumbranceStatus;
}

// Apply encumbrance effects to character stats
function applyEncumbranceEffects(encumbranceStatus) {
    const modifiers = getEncumbranceModifiers(encumbranceStatus);
    
    // Find combat tab elements
    const movementRateElement = document.getElementById('movement-rate');
    const armorClassElement = document.getElementById('armor-class');
    
    // Store original values if not already stored
    if (!window.originalMovementRate && movementRateElement) {
        window.originalMovementRate = parseInt(movementRateElement.textContent) || 40;
    }
    
    if (!window.originalArmorClass && armorClassElement) {
        window.originalArmorClass = parseInt(armorClassElement.textContent) || 9;
    }
    
    // Apply modifiers if elements exist
    if (movementRateElement && window.originalMovementRate) {
        const newMovementRate = Math.max(0, window.originalMovementRate + modifiers.movementPenalty);
        movementRateElement.textContent = newMovementRate;
        
        // Add visual indicator of penalty
        if (modifiers.movementPenalty !== 0) {
            movementRateElement.setAttribute('data-original', window.originalMovementRate);
            movementRateElement.setAttribute('data-penalty', modifiers.movementPenalty);
            movementRateElement.classList.add('stat-penalized');
        } else {
            movementRateElement.removeAttribute('data-original');
            movementRateElement.removeAttribute('data-penalty');
            movementRateElement.classList.remove('stat-penalized');
        }
    }
    
    if (armorClassElement && window.originalArmorClass) {
        const newArmorClass = window.originalArmorClass + modifiers.acPenalty;
        armorClassElement.textContent = newArmorClass;
        
        // Add visual indicator of penalty
        if (modifiers.acPenalty !== 0) {
            armorClassElement.setAttribute('data-original', window.originalArmorClass);
            armorClassElement.setAttribute('data-penalty', modifiers.acPenalty);
            armorClassElement.classList.add('stat-penalized');
        } else {
            armorClassElement.removeAttribute('data-original');
            armorClassElement.removeAttribute('data-penalty');
            armorClassElement.classList.remove('stat-penalized');
        }
    }
}

// Calculate total weight from inventory items
function calculateTotalWeight(items) {
    if (!items || !Array.isArray(items)) {
        return 0;
    }
    
    return items.reduce((total, item) => {
        let itemWeight = 0;
        
        if (item.item_details && typeof item.item_details.weight !== 'undefined') {
            itemWeight = parseFloat(item.item_details.weight) || 0;
        }
        
        // Multiply by quantity
        return total + (itemWeight * item.quantity);
    }, 0);
}

// Update inventory with encumbrance information
function updateInventoryWithEncumbrance(inventoryData) {
    if (!window.characterData || !inventoryData || !inventoryData.items) {
        console.warn('Missing character or inventory data for encumbrance calculation');
        return;
    }
    
    const strength = window.characterData.strength || 10;
    const constitution = window.characterData.constitution || 10;
    
    // Calculate thresholds based on character attributes
    const thresholds = calculateEncumbranceThresholds(strength, constitution);
    
    // Calculate current weight
    const currentWeight = calculateTotalWeight(inventoryData.items);
    
    // Update the display
    updateEncumbranceDisplay(currentWeight, thresholds);
    
    // Save the calculated values for reference
    window.encumbranceData = {
        currentWeight,
        thresholds,
        status: getEncumbranceStatus(currentWeight, thresholds),
        modifiers: getEncumbranceModifiers(getEncumbranceStatus(currentWeight, thresholds))
    };
    
    return window.encumbranceData;
}

// Initialize encumbrance UI
function initializeEncumbrance() {
    // Create encumbrance UI elements if they don't exist
    if (!document.getElementById('encumbrance-container')) {
        const inventorySummary = document.querySelector('.inventory-summary');
        
        if (inventorySummary) {
            const encumbranceHTML = `
                <div id="encumbrance-container" class="encumbrance-container">
                    <div class="encumbrance-info">
                        <div id="encumbrance-status" class="encumbrance-status unencumbered">Unencumbered</div>
                        <div class="encumbrance-weight">
                            <span id="total-weight">0</span> / <span id="max-weight">0</span> lbs
                        </div>
                    </div>
                    <div class="encumbrance-progress">
                        <div id="encumbrance-bar" class="encumbrance-progress-bar unencumbered" style="width: 0%;"></div>
                    </div>
                    <div class="weight-thresholds">
                        <span class="threshold-mark" id="encumbered-threshold">0</span>
                        <span class="threshold-mark" id="heavy-threshold">0</span>
                        <span class="threshold-mark" id="max-threshold">0</span>
                    </div>
                </div>
            `;
            
            inventorySummary.innerHTML = encumbranceHTML;
        }
    }
    
    // Apply encumbrance calculation
    if (window.inventoryState && window.inventoryState.items) {
        updateInventoryWithEncumbrance({
            items: window.inventoryState.items,
            inventory: window.inventoryState.inventory
        });
    }
}

function setupInventoryAddItem() {
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
        await fetchInventory();
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
        
        // Initialize encumbrance system
        initializeEncumbrance();
        
        // Update the inventory with encumbrance data
        updateInventoryWithEncumbrance({
            inventory: inventory,
            items: enrichedItems
        });
        
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
            const equippedStatus = item.is_equipped ? 'Equipped' : 'Not equipped';
            row.innerHTML = `
                <td>${itemName}</td>
                <td>${formatItemType(item.item_type)}</td>
                <td>${item.quantity}</td>
                <td>${itemWeight}</td>
                <td>${equippedStatus}</td>
                <td>${item.notes || ''}</td>
                <td class="item-actions">
                    <button class="btn btn-item-edit" data-id="${item.id}">Edit</button>
                    <button class="btn btn-item-delete" data-id="${item.id}">Delete</button>
                </td>
            `;
            inventoryItems.appendChild(row);
        });
        addInventoryButtonListeners();
        inventoryTable.style.display = 'table';
        
        return data;
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
        return null;
    }
}

function addInventoryButtonListeners() {
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
}

function formatItemType(type) {
    if (!type) return 'Unknown';
    return type
        .split('_')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');
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
        await fetchInventory();
        
        // Update combat tab if it's active
        const combatTab = document.getElementById('combat-tab');
        if (combatTab && combatTab.classList.contains('active') && typeof loadCombatData === 'function') {
            loadCombatData();
        }
    } catch (error) {
        console.error('Error deleting inventory item:', error);
        alert('Failed to delete item: ' + error.message);
    }
}

function getCharacterIdFromURL() {
    const pathParts = window.location.pathname.split('/');
    return pathParts[pathParts.length - 1];
}

// Add CSS for encumbrance system
function addEncumbranceStyles() {
    if (document.getElementById('encumbrance-styles')) {
        return; // Styles already added
    }
    
    const styleElement = document.createElement('style');
    styleElement.id = 'encumbrance-styles';
    styleElement.textContent = `
        .encumbrance-container {
            margin-top: 15px;
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 10px;
            background-color: #f8f9fa;
        }
        
        .encumbrance-info {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;
        }
        
        .encumbrance-status {
            font-weight: bold;
            padding: 3px 8px;
            border-radius: 4px;
        }
        
        .encumbrance-status.unencumbered {
            background-color: #d4edda;
            color: #155724;
        }
        
        .encumbrance-status.encumbered {
            background-color: #fff3cd;
            color: #856404;
        }
        
        .encumbrance-status.heavily-encumbered {
            background-color: #f8d7da;
            color: #721c24;
        }
        
        .encumbrance-status.overloaded {
            background-color: #721c24;
            color: white;
        }
        
        .encumbrance-weight {
            font-size: 14px;
            color: #495057;
        }
        
        .encumbrance-progress {
            height: 8px;
            background-color: #e9ecef;
            border-radius: 4px;
            overflow: hidden;
        }
        
        .encumbrance-progress-bar {
            height: 100%;
            transition: width 0.3s ease;
        }
        
        .encumbrance-progress-bar.unencumbered {
            background-color: #28a745;
        }
        
        .encumbrance-progress-bar.encumbered {
            background-color: #ffc107;
        }
        
        .encumbrance-progress-bar.heavily-encumbered {
            background-color: #dc3545;
        }
        
        .encumbrance-progress-bar.overloaded {
            background-color: #6f1d1b;
        }
        
        .stat-penalized {
            color: #dc3545;
            position: relative;
        }
        
        .stat-penalized::after {
            content: attr(data-original) " (" attr(data-penalty) ")";
            position: absolute;
            font-size: 10px;
            bottom: -16px;
            left: 0;
            white-space: nowrap;
            color: #6c757d;
        }
        
        .weight-thresholds {
            display: flex;
            justify-content: space-between;
            margin-top: 4px;
            font-size: 11px;
            color: #6c757d;
        }
        
        .threshold-mark {
            position: relative;
        }
        
        .threshold-mark::before {
            content: '';
            position: absolute;
            top: -6px;
            left: 50%;
            transform: translateX(-50%);
            width: 1px;
            height: 4px;
            background-color: #6c757d;
        }
    `;
    
    document.head.appendChild(styleElement);
}

// Function to edit inventory item - called when edit button is clicked
async function editInventoryItem(itemId) {
    try {
        const inventoryID = await getInventoryID();
        if (!inventoryID) {
            throw new Error('Could not find inventory');
        }
        const token = localStorage.getItem('authToken');
        const response = await fetch(`/api/inventories/${inventoryID}/items/${itemId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!response.ok) {
            throw new Error('Failed to fetch item details');
        }
        const itemData = await response.json();
        console.log('Item to edit:', itemData);
        if (!document.getElementById('edit-item-modal')) {
            createEditItemModal();
        }
        document.getElementById('modal-title').textContent = 'Edit Inventory Item';
        const itemTypeSelect = document.getElementById('edit-item-type');
        itemTypeSelect.value = itemData.item_type;
        itemTypeSelect.disabled = true;
        const itemIdInput = document.getElementById('edit-item-id');
        itemIdInput.value = itemData.item_id;
        const itemNameInput = document.getElementById('edit-item-name');
        if (itemData.item_details && itemData.item_details.name) {
            itemNameInput.value = itemData.item_details.name;
        } else {
            itemNameInput.value = 'Unknown Item';
        }
        document.getElementById('edit-item-quantity').value = itemData.quantity;
        document.getElementById('edit-item-equipped').checked = itemData.is_equipped;
        document.getElementById('edit-item-notes').value = itemData.notes || '';
        // Set the item ID in a data attribute for the save function
        document.getElementById('edit-item-form').setAttribute('data-item-id', itemId);
        document.getElementById('edit-item-modal').style.display = 'block';
    } catch (error) {
        console.error('Error preparing to edit item:', error);
        alert('Error: ' + error.message);
    }
}

function createEditItemModal() {
    const modalHTML = `
    <div id="edit-item-modal" class="modal">
        <div class="modal-content">
            <span class="close" id="close-edit-modal">&times;</span>
            <h2 id="modal-title">Edit Inventory Item</h2>
            <form id="edit-item-form">
                <div class="form-group">
                    <label for="edit-item-type">Item Type:</label>
                    <input type="text" id="edit-item-type" readonly>
                </div>
                <div class="form-group">
                    <label for="edit-item-name">Item:</label>
                    <input type="text" id="edit-item-name" readonly>
                    <input type="hidden" id="edit-item-id">
                </div>
                <div class="form-group">
                    <label for="edit-item-quantity">Quantity:</label>
                    <input type="number" id="edit-item-quantity" min="1" required>
                </div>
                <div class="form-group">
                    <input type="checkbox" id="edit-item-equipped">
                    <label for="edit-item-equipped">Equipped</label>
                </div>
                <div class="form-group">
                    <label for="edit-item-notes">Notes:</label>
                    <textarea id="edit-item-notes" rows="3"></textarea>
                </div>
                <div class="form-actions">
                    <button type="submit" class="btn btn-primary">Save Changes</button>
                    <button type="button" id="cancel-edit" class="btn">Cancel</button>
                </div>
            </form>
        </div>
    </div>
    `;
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    document.getElementById('close-edit-modal').addEventListener('click', closeEditModal);
    document.getElementById('cancel-edit').addEventListener('click', closeEditModal);
    document.getElementById('edit-item-form').addEventListener('submit', saveEditedItem);
    window.addEventListener('click', function(event) {
        const modal = document.getElementById('edit-item-modal');
        if (event.target === modal) {
            closeEditModal();
        }
    });
}

function closeEditModal() {
    const modal = document.getElementById('edit-item-modal');
    if (modal) {
        modal.style.display = 'none';
    }
}

async function saveEditedItem(event) {
    event.preventDefault();
    try {
        const form = document.getElementById('edit-item-form');
        const itemId = form.getAttribute('data-item-id');
        const inventoryID = await getInventoryID();
        if (!inventoryID || !itemId) {
            throw new Error('Missing inventory or item ID');
        }
        const token = localStorage.getItem('authToken');
        const quantity = parseInt(document.getElementById('edit-item-quantity').value);
        const isEquipped = document.getElementById('edit-item-equipped').checked;
        const notes = document.getElementById('edit-item-notes').value;
        if (isNaN(quantity) || quantity < 1) {
            alert('Quantity must be at least 1');
            return;
        }
        const response = await fetch(`/api/inventories/${inventoryID}/items/${itemId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                quantity: quantity,
                is_equipped: isEquipped,
                notes: notes
            })
        });
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to update item: ${errorText}`);
        }
        closeEditModal();
        await fetchInventory();
        
        // Update combat tab if it's active
        const combatTab = document.getElementById('combat-tab');
        if (combatTab && combatTab.classList.contains('active') && typeof loadCombatData === 'function') {
            loadCombatData();
        }
    } catch (error) {
        console.error('Error saving edited item:', error);
        alert('Error saving changes: ' + error.message);
    }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    // Add the encumbrance styles
    addEncumbranceStyles();
    
    // Connect the inventory tab events
    const inventoryTab = document.querySelector('.tab-item[data-tab="inventory-tab"]');
    if (inventoryTab) {
        const originalClickHandler = inventoryTab.onclick;
        inventoryTab.onclick = function(event) {
            if (originalClickHandler) {
                originalClickHandler.call(this, event);
            }
            
            // Initialize encumbrance after a short delay
            setTimeout(initializeEncumbrance, 500);
        };
    }
    
    // If inventory tab is active by default, initialize
    if (document.getElementById('inventory-tab') && 
        document.getElementById('inventory-tab').classList.contains('active')) {
        setTimeout(initializeEncumbrance, 500);
    }
});