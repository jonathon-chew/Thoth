#!/usr/bin/env bash

set -euo pipefail

# ----------------------------
# Styling for clarity
# ----------------------------
GREEN="\033[1;32m"
RED="\033[1;31m"
YELLOW="\033[1;33m"
CYAN="\033[1;36m"
RESET="\033[0m"

# ----------------------------
# Step 1: Go vet
# ----------------------------
echo -e "${CYAN}🔍 Running go vet...${RESET}"
if go vet ./...; then
  echo -e "${GREEN}✅ go vet passed!${RESET}"
else
  echo -e "${RED}❌ go vet found issues!${RESET}"
  exit 1
fi

# ----------------------------
# Step 2: Build all packages
# ----------------------------
echo -e "${CYAN}🛠 Building all packages...${RESET}"
if go build ./...; then
  echo -e "${GREEN}✅ Build succeeded!${RESET}"
else
  echo -e "${RED}❌ Build failed!${RESET}"
  exit 1
fi

# ----------------------------
# Step 3: Run all tests
# ----------------------------
echo -e "${CYAN}🧪 Running tests...${RESET}"

if go test -v ./...; then
  echo -e "${GREEN}✅ All tests passed!${RESET}"
else
  echo -e "${RED}❌ Some tests failed!${RESET}"
  exit 1
fi

# ----------------------------
# Step 4: Incriment the tag version
# ----------------------------
echo -e "${CYAN}🧪 Updating git tags...${RESET}"

if ./Thoth -i; then 
  echo -e "${GREEN}✅ Successfully updated the tags!${RESET}"
else
  echo -e "${RED}❌ Failed to update the tags successfully !${RESET}"
  exit 1
fi

echo -e "${GREEN}🎉 CI pipeline completed successfully!${RESET}"