// web/static/js/thief_skills_tab.js
document.addEventListener('DOMContentLoaded', function() {
    const thiefSkillsTab = document.getElementById('thief-skills-tab');
    
    if (thiefSkillsTab) {
        thiefSkillsTab.addEventListener('click', function() {
            loadThiefSkills();
        });
    }
    
    // If thief skills tab is active on page load, load the skills
    if (thiefSkillsTab && thiefSkillsTab.classList.contains('active')) {
        loadThiefSkills();
    }
});

function loadThiefSkills() {
    const characterId = document.body.dataset.characterId;
    if (!characterId) {
        console.error('Character ID not found');
        return;
    }
    
    // Show loading indicator
    const thiefSkillsContent = document.querySelector('#thief-skills .card-body');
    if (thiefSkillsContent) {
        thiefSkillsContent.innerHTML = '<div class="text-center"><div class="spinner-border" role="status"><span class="visually-hidden">Loading...</span></div></div>';
    }
    
    // Fetch thief skills for this character
    fetch(`/api/characters/${characterId}/thief-skills`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load thief skills');
            }
            return response.json();
        })
        .then(skills => {
            displayThiefSkills(skills);
        })
        .catch(error => {
            console.error('Error loading thief skills:', error);
            const tableBody = document.querySelector('#thief-skills table tbody');
            if (tableBody) {
                tableBody.innerHTML = '<tr><td colspan="3" class="text-center text-danger">Error loading thief skills</td></tr>';
            }
        });
}

function displayThiefSkills(skills) {
    const tableBody = document.querySelector('#thief-skills table tbody');
    if (!tableBody) {
        console.error('Thief skills table body not found');
        return;
    }
    
    if (skills.length === 0) {
        tableBody.innerHTML = '<tr><td colspan="3" class="text-center">No thief skills available for this character</td></tr>';
        return;
    }
    
    // Clear existing rows
    tableBody.innerHTML = '';
    
    // Add a row for each skill
    skills.forEach(skill => {
        const row = document.createElement('tr');
        
        const nameCell = document.createElement('td');
        nameCell.textContent = skill.Name;
        
        const attributeCell = document.createElement('td');
        attributeCell.textContent = skill.Attribute;
        
        const chanceCell = document.createElement('td');
        chanceCell.textContent = skill.SuccessChance;
        // Add styling based on success chance
        if (skill.SuccessChance === 'N/A') {
            chanceCell.classList.add('text-muted');
        } else {
            const parts = skill.SuccessChance.split(':');
            if (parts.length === 2) {
                const successRate = parseInt(parts[0]) / parseInt(parts[1]);
                if (successRate >= 0.75) {
                    chanceCell.classList.add('text-success');
                } else if (successRate >= 0.5) {
                    chanceCell.classList.add('text-warning');
                } else {
                    chanceCell.classList.add('text-danger');
                }
            }
        }
        
        row.appendChild(nameCell);
        row.appendChild(attributeCell);
        row.appendChild(chanceCell);
        
        tableBody.appendChild(row);
    });
}