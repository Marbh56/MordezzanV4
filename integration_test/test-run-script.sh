#!/bin/bash

RESET="\033[0m"
BOLD="\033[1m"
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[34m"
MAGENTA="\033[35m"
CYAN="\033[36m"

# Check for --simple flag
SIMPLE_OUTPUT=false
if [[ "$1" == "--simple" ]]; then
  SIMPLE_OUTPUT=true
fi

print_header() {
  if [ "$SIMPLE_OUTPUT" = false ]; then
    echo -e "\n${BOLD}${BLUE}=================================================${RESET}"
    echo -e "${BOLD}${MAGENTA} $1 ${RESET}"
    echo -e "${BOLD}${BLUE}=================================================${RESET}\n"
  else
    echo -e "${BOLD}${MAGENTA}$1${RESET}"
  fi
}

print_section() {
  if [ "$SIMPLE_OUTPUT" = false ]; then
    echo -e "\n${BOLD}${CYAN}>>> $1${RESET}\n"
  fi
}

print_success() {
  echo -e "${GREEN}✓ $1${RESET}"
}

print_error() {
  echo -e "${RED}✗ $1${RESET}"
}

print_info() {
  if [ "$SIMPLE_OUTPUT" = false ]; then
    echo -e "${YELLOW}i ${RESET}$1"
  fi
}

# Only print this header if not in simple mode
if [ "$SIMPLE_OUTPUT" = false ]; then
  print_header "INTEGRATION TEST RUNNER"
  print_info "Starting integration tests using improved formatting..."
fi

if ! command -v go &> /dev/null; then
  print_error "Go is not installed or not in PATH"
  exit 1
fi

LOGDIR="test_logs"
FAILDIR="$LOGDIR/failures"
mkdir -p $LOGDIR
mkdir -p $FAILDIR

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
LOGFILE="$LOGDIR/integration_test_$TIMESTAMP.log"

run_test() {
  local test_name=$1
  local test_pattern=$2
  
  if [ "$SIMPLE_OUTPUT" = false ]; then
    print_section "Running $test_name integration tests"
  fi
  
  local temp_output=$(mktemp)
  
  # Determine the correct directory path based on script location
  local test_dir="./integration_test"
  # If script is being run from integration_test directory, adjust the path
  if [[ $(pwd) == */integration_test ]]; then
    test_dir="."
  fi
  
  if [ "$SIMPLE_OUTPUT" = true ]; then
    # In simple mode, only run the tests without printing the output
    go test $test_dir/... -v -run $test_pattern > "$temp_output" 2>&1
  else
    # In normal mode, display all test output
    go test $test_dir/... -v -run $test_pattern | tee "$temp_output" | tee -a "$LOGFILE"
  fi
  
  local test_status=${PIPESTATUS[0]}
  if [ $test_status -ne 0 ]; then
    local failure_file="$FAILDIR/${test_name}_FAILED_$TIMESTAMP.log"
    cp "$temp_output" "$failure_file"
    if [ "$SIMPLE_OUTPUT" = false ]; then
      print_error "$test_name integration tests failed. Output saved to $failure_file"
    fi
  fi
  
  rm "$temp_output"
  return $test_status
}

run_test "auth middleware" "TestAuthMiddlewareIntegration"
AUTH_TEST_STATUS=$?

run_test "user registration" "TestUserRegistrationIntegration"
USER_REG_TEST_STATUS=$?

run_test "user" "TestUserCRUDIntegration"
USER_TEST_STATUS=$?

run_test "character" "TestCharacterCRUDIntegration"
CHAR_TEST_STATUS=$?

run_test "spell" "TestSpellCRUDIntegration"
SPELL_TEST_STATUS=$?

run_test "armor" "TestArmorCRUDIntegration"
ARMOR_TEST_STATUS=$?

run_test "weapon" "TestWeaponCRUDIntegration"
WEAPON_TEST_STATUS=$?

run_test "equipment" "TestEquipmentCRUDIntegration"
EQUIPMENT_TEST_STATUS=$?

run_test "shield" "TestShieldCRUDIntegration"
SHIELD_TEST_STATUS=$?

run_test "potion" "TestPotionCRUDIntegration"
POTION_TEST_STATUS=$?

run_test "magic_item" "TestMagicItemCRUDIntegration"
MAGIC_ITEM_TEST_STATUS=$?

run_test "ring" "TestRingCRUDIntegration"
RING_TEST_STATUS=$?

run_test "ammo" "TestAmmoCRUDIntegration"
AMMO_TEST_STATUS=$?

run_test "spell_scroll" "TestSpellScrollCRUDIntegration"
SPELL_SCROLL_TEST_STATUS=$?

run_test "container" "TestContainerCRUDIntegration" 
CONTAINER_TEST_STATUS=$?

run_test "inventory" "TestInventoryCRUDIntegration"
INVENTORY_TEST_STATUS=$?

run_test "treasure" "TestTreasureCRUDIntegration"
TREASURE_TEST_STATUS=$?

if [ "$SIMPLE_OUTPUT" = false ]; then
  echo ""
fi

# Check if any tests failed and print the summary
if [ $AUTH_TEST_STATUS -eq 0 ] && \
   [ $USER_REG_TEST_STATUS -eq 0 ] && \
   [ $USER_TEST_STATUS -eq 0 ] && \
   [ $CHAR_TEST_STATUS -eq 0 ] && \
   [ $SPELL_TEST_STATUS -eq 0 ] && \
   [ $ARMOR_TEST_STATUS -eq 0 ] && \
   [ $WEAPON_TEST_STATUS -eq 0 ] && \
   [ $EQUIPMENT_TEST_STATUS -eq 0 ] && \
   [ $SHIELD_TEST_STATUS -eq 0 ] && \
   [ $POTION_TEST_STATUS -eq 0 ] && \
   [ $MAGIC_ITEM_TEST_STATUS -eq 0 ] && \
   [ $RING_TEST_STATUS -eq 0 ] && \
   [ $AMMO_TEST_STATUS -eq 0 ] && \
   [ $SPELL_SCROLL_TEST_STATUS -eq 0 ] && \
   [ $CONTAINER_TEST_STATUS -eq 0 ] && \
   [ $INVENTORY_TEST_STATUS -eq 0 ] && \
   [ $TREASURE_TEST_STATUS -eq 0 ]; then
  print_header "ALL TESTS PASSED SUCCESSFULLY!"
  
  if [ "$SIMPLE_OUTPUT" = false ]; then
    print_success "Auth Middleware integration tests: PASSED"
    print_success "User Registration integration tests: PASSED"
    print_success "User CRUD integration tests: PASSED"
    print_success "Character CRUD integration tests: PASSED"
    print_success "Spell CRUD integration tests: PASSED"
    print_success "Armor CRUD integration tests: PASSED"
    print_success "Weapon CRUD integration tests: PASSED"
    print_success "Equipment CRUD integration tests: PASSED"
    print_success "Shield CRUD integration tests: PASSED"
    print_success "Potion CRUD integration tests: PASSED"
    print_success "Magic Item CRUD integration tests: PASSED"
    print_success "Ring CRUD integration tests: PASSED"
    print_success "Ammo CRUD integration tests: PASSED"
    print_success "Spell Scroll CRUD integration tests: PASSED"
    print_success "Container CRUD integration tests: PASSED"
    print_success "Inventory CRUD integration tests: PASSED"
    print_success "Treasure CRUD integration tests: PASSED"
    print_info "Log file saved to: $LOGFILE"
  fi
  
  if [ -z "$(ls -A $FAILDIR)" ]; then
    rmdir $FAILDIR
  fi
  exit 0
else
  print_header "TEST FAILURES DETECTED"
  
  if [ $AUTH_TEST_STATUS -eq 0 ]; then
    print_success "Auth Middleware integration tests: PASSED"
  else
    print_error "Auth Middleware integration tests: FAILED"
  fi
  if [ $USER_REG_TEST_STATUS -eq 0 ]; then
    print_success "User Registration integration tests: PASSED"
  else
    print_error "User Registration integration tests: FAILED"
  fi
  if [ $USER_TEST_STATUS -eq 0 ]; then
    print_success "User CRUD integration tests: PASSED"
  else
    print_error "User CRUD integration tests: FAILED"
  fi
  if [ $CHAR_TEST_STATUS -eq 0 ]; then
    print_success "Character CRUD integration tests: PASSED"
  else
    print_error "Character CRUD integration tests: FAILED"
  fi
  if [ $SPELL_TEST_STATUS -eq 0 ]; then
    print_success "Spell CRUD integration tests: PASSED"
  else
    print_error "Spell CRUD integration tests: FAILED"
  fi
  if [ $ARMOR_TEST_STATUS -eq 0 ]; then
    print_success "Armor CRUD integration tests: PASSED"
  else
    print_error "Armor CRUD integration tests: FAILED"
  fi
  if [ $WEAPON_TEST_STATUS -eq 0 ]; then
    print_success "Weapon CRUD integration tests: PASSED"
  else
    print_error "Weapon CRUD integration tests: FAILED"
  fi
  if [ $EQUIPMENT_TEST_STATUS -eq 0 ]; then
    print_success "Equipment CRUD integration tests: PASSED"
  else
    print_error "Equipment CRUD integration tests: FAILED"
  fi
  if [ $SHIELD_TEST_STATUS -eq 0 ]; then
    print_success "Shield CRUD integration tests: PASSED"
  else
    print_error "Shield CRUD integration tests: FAILED"
  fi
  if [ $POTION_TEST_STATUS -eq 0 ]; then
    print_success "Potion CRUD integration tests: PASSED"
  else
    print_error "Potion CRUD integration tests: FAILED"
  fi
  if [ $MAGIC_ITEM_TEST_STATUS -eq 0 ]; then
    print_success "Magic Item CRUD integration tests: PASSED"
  else
    print_error "Magic Item CRUD integration tests: FAILED"
  fi
  if [ $RING_TEST_STATUS -eq 0 ]; then
    print_success "Ring CRUD integration tests: PASSED"
  else
    print_error "Ring CRUD integration tests: FAILED"
  fi
  if [ $AMMO_TEST_STATUS -eq 0 ]; then
    print_success "Ammo CRUD integration tests: PASSED"
  else
    print_error "Ammo CRUD integration tests: FAILED"
  fi
  if [ $SPELL_SCROLL_TEST_STATUS -eq 0 ]; then
    print_success "Spell Scroll CRUD integration tests: PASSED"
  else
    print_error "Spell Scroll integration tests: FAILED"
  fi
  if [ $CONTAINER_TEST_STATUS -eq 0 ]; then
    print_success "Container CRUD integration tests: PASSED"
  else
    print_error "Container CRUD integration tests: FAILED"
  fi
  if [ $INVENTORY_TEST_STATUS -eq 0 ]; then
    print_success "Inventory CRUD integration tests: PASSED"
  else
    print_error "Inventory CRUD integration tests: FAILED"
  fi
  if [ $TREASURE_TEST_STATUS -eq 0 ]; then
    print_success "Treasure CRUD integration tests: PASSED"
  else
    print_error "Treasure CRUD integration tests: FAILED"
  fi
  
  if [ "$SIMPLE_OUTPUT" = false ]; then
    print_info "Complete log file saved to: $LOGFILE"
    print_info "Failed test logs saved to: $FAILDIR/"
  fi
  
  exit 1
fi