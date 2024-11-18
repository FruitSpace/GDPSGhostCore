#!/bin/bash

RED='\033[0;31m'
BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m'
GRAY='\033[0;37m'

PATH=$PATH:~/go/bin
echo -e "${GREEN}Checking for ineffectual assignments...${GRAY}"
ineffassign ./...
echo -e "${GREEN}Building...${GRAY}"
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o HalogenGhostCore .
echo -e "${GREEN}Done.${NC}"