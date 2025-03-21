async function editInventoryItem(itemId) {
    try {
        // Get the inventory ID
        const inventoryID = await getInventoryID();
        if (!inventoryID) {
            throw new Error('Could not find inventory');
        }

        // Fetch item details
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

        // Create or get the modal
        if (!document.getElementById('edit-item-modal')) {
            createEditItemModal();
        }

        // Fill the modal with item data
        document.getElementById('modal-title').textContent = 'Edit Inventory Item';
        
        // Set item type
        const itemTypeSelect = document.getElementById('edit-item-type');
        itemTypeSelect.value = itemData.item_type;
        itemTypeSelect.disabled = true; // Can't change item type when editing
        
        // Set item ID
        const itemIdInput = document.getElementById('edit-item-id');
        itemIdInput.value = itemData.item_id;
        
        // Set item name (read-only)
        const itemNameInput = document.getElementById('edit-item-name');
        if (itemData.item_details && itemData.item_details.name) {
            itemNameInput.value = itemData.item_details.name;
        } else {
            itemNameInput.value = 'Unknown Item';
        }
        
        // Set quantity
        document.getElementById('edit-item-quantity').value = itemData.quantity;
        
        // Set equipped status
        document.getElementById('edit-item-equipped').checked = itemData.is_equipped;
        
        // Set notes
        document.getElementById('edit-item-notes').value = itemData.notes || '';
        
        // Set the item ID in a data attribute for the save function
        document.getElementById('edit-item-form').setAttribute('data-item-id', itemId);
        
        // Show the modal
        document.getElementById('edit-item-modal').style.display = 'block';
    } catch (error) {
        console.error('Error preparing to edit item:', error);
        alert('Error: ' + error.message);
    }
}

// Function to create the edit item modal if it doesn't exist
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
    
    // Add modal to the document
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    
    // Add event listeners
    document.getElementById('close-edit-modal').addEventListener('click', closeEditModal);
    document.getElementById('cancel-edit').addEventListener('click', closeEditModal);
    document.getElementById('edit-item-form').addEventListener('submit', saveEditedItem);
    
    // Close on click outside
    window.addEventListener('click', function(event) {
        const modal = document.getElementById('edit-item-modal');
        if (event.target === modal) {
            closeEditModal();
        }
    });
}

// Function to close the edit modal
function closeEditModal() {
    const modal = document.getElementById('edit-item-modal');
    if (modal) {
        modal.style.display = 'none';
    }
}

// Function to save edited item
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
        
        // Validate quantity
        if (isNaN(quantity) || quantity < 1) {
            alert('Quantity must be at least 1');
            return;
        }
        
        // Prepare the request
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
        
        // Close the modal and refresh the inventory
        closeEditModal();
        fetchInventory();
        
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