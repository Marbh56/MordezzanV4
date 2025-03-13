#!/bin/bash

RESET="\033[0m"
BOLD="\033[1m"
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[34m"
MAGENTA="\033[35m"
CYAN="\033[36m"

print_header() {
  echo -e "\n${BOLD}${BLUE}=================================================${RESET}"
  echo -e "${BOLD}${MAGENTA} $1 ${RESET}"
  echo -e "${BOLD}${BLUE}=================================================${RESET}\n"
}

print_section() {
  echo -e "\n${BOLD}${CYAN}>>> $1${RESET}\n"
}

print_success() {
  echo -e "${GREEN}✓ $1${RESET}"
}

print_error() {
  echo -e "${RED}✗ $1${RESET}"
}

print_info() {
  echo -e "${YELLOW}i ${RESET}$1"
}

print_header "INTEGRATION TEST RUNNER"
print_info "Starting integration tests using improved formatting..."

if ! command -v go &> /dev/null; then
  print_error "Go is not installed or not in PATH"
  exit 1
fi

LOGDIR="test_logs"
mkdir -p $LOGDIR

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
LOGFILE="$LOGDIR/integration_test_$TIMESTAMP.log"

print_section "Running user integration tests"
go test ./integration_test/... -v -run TestUserCRUDIntegration | tee -a "$LOGFILE"
USER_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running character integration tests"
go test ./integration_test/... -v -run TestCharacterCRUDIntegration | tee -a "$LOGFILE"
CHAR_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running spell integration tests"
go test ./integration_test/... -v -run TestSpellCRUDIntegration | tee -a "$LOGFILE"
SPELL_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running armor integration tests"
go test ./integration_test/... -v -run TestArmorCRUDIntegration | tee -a "$LOGFILE"
ARMOR_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running weapon integration tests"
go test ./integration_test/... -v -run TestWeaponCRUDIntegration | tee -a "$LOGFILE"
WEAPON_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running equipment integration tests"
go test ./integration_test/... -v -run TestEquipmentCRUDIntegration | tee -a "$LOGFILE"
EQUIPMENT_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running shield integration tests"
go test ./integration_test/... -v -run TestShieldCRUDIntegration | tee -a "$LOGFILE"
SHIELD_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running potion integration tests"
go test ./integration_test/... -v -run TestPotionCRUDIntegration | tee -a "$LOGFILE"
POTION_TEST_STATUS=${PIPESTATUS[0]}

print_section "Running magic item integration tests"
go test ./integration_test/... -v -run TestMagicItemCRUDIntegration | tee -a "$LOGFILE"
MAGIC_ITEM_TEST_STATUS=${PIPESTATUS[0]}

echo ""
if [ $USER_TEST_STATUS -eq 0 ] &&
   [ $CHAR_TEST_STATUS -eq 0 ] &&
   [ $SPELL_TEST_STATUS -eq 0 ] &&
   [ $ARMOR_TEST_STATUS -eq 0 ] &&
   [ $WEAPON_TEST_STATUS -eq 0 ] &&
   [ $EQUIPMENT_TEST_STATUS -eq 0 ] &&
   [ $SHIELD_TEST_STATUS -eq 0 ] &&
   [ $POTION_TEST_STATUS -eq 0 ] &&
   [ $MAGIC_ITEM_TEST_STATUS -eq 0 ]; then
  print_header "ALL TESTS PASSED SUCCESSFULLY!"
  print_success "User CRUD integration tests: PASSED"
  print_success "Character CRUD integration tests: PASSED"
  print_success "Spell CRUD integration tests: PASSED"
  print_success "Armor CRUD integration tests: PASSED"
  print_success "Weapon CRUD integration tests: PASSED"
  print_success "Equipment CRUD integration tests: PASSED"
  print_success "Shield CRUD integration tests: PASSED"
  print_success "Potion CRUD integration tests: PASSED"
  print_success "Magic Item CRUD integration tests: PASSED"
  print_info "Log file saved to: $LOGFILE"
  exit 0
else
  print_header "TEST FAILURES DETECTED"
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
  print_info "Check the log file for details: $LOGFILE"
  exit 1
fi