#!/bin/bash

# Script para converter modelo DSL para JSON e vice-versa
# Requer OpenFGA CLI instalado

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
OPENFGA_DIR="$PROJECT_DIR/docker/openfga"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_usage() {
    echo "Usage: $0 [dsl-to-json|json-to-dsl|convert-all|validate] [input-file] [output-file]"
    echo ""
    echo "Commands:"
    echo "  dsl-to-json      Convert DSL file to JSON"
    echo "  json-to-dsl      Convert JSON file to DSL"
    echo "  convert-all      Convert all DSL files in openfga directory"
    echo "  validate         Validate model file"
    echo ""
    echo "Examples:"
    echo "  $0 dsl-to-json model.fga model.json"
    echo "  $0 json-to-dsl model.json model.fga"
    echo "  $0 convert-all"
    echo "  $0 validate model.json"
}

check_fga_cli() {
    if ! command -v fga &> /dev/null; then
        echo -e "${RED}OpenFGA CLI not found. Please install it first:${NC}"
        echo "curl -L https://github.com/openfga/cli/releases/latest/download/fga_linux_amd64.tar.gz | tar -xzf - -C /usr/local/bin"
        exit 1
    fi
}

dsl_to_json() {
    local input_file="$1"
    local output_file="$2"
    
    if [ ! -f "$input_file" ]; then
        echo -e "${RED}Input file not found: $input_file${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Converting DSL to JSON...${NC}"
    echo -e "${BLUE}Input:  $input_file${NC}"
    echo -e "${BLUE}Output: $output_file${NC}"
    
    # Usar OpenFGA CLI para converter
    fga model transform --from dsl --to json --file "$input_file" > "$output_file"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Successfully converted $input_file to $output_file${NC}"
    else
        echo -e "${RED}❌ Failed to convert $input_file${NC}"
        exit 1
    fi
}

json_to_dsl() {
    local input_file="$1"
    local output_file="$2"
    
    if [ ! -f "$input_file" ]; then
        echo -e "${RED}Input file not found: $input_file${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Converting JSON to DSL...${NC}"
    echo -e "${BLUE}Input:  $input_file${NC}"
    echo -e "${BLUE}Output: $output_file${NC}"
    
    # Usar OpenFGA CLI para converter
    fga model transform --from json --to dsl --file "$input_file" > "$output_file"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Successfully converted $input_file to $output_file${NC}"
    else
        echo -e "${RED}❌ Failed to convert $input_file${NC}"
        exit 1
    fi
}

convert_all() {
    echo -e "${YELLOW}Converting all DSL files in $OPENFGA_DIR...${NC}"
    
    if [ ! -d "$OPENFGA_DIR" ]; then
        echo -e "${RED}OpenFGA directory not found: $OPENFGA_DIR${NC}"
        exit 1
    fi
    
    cd "$OPENFGA_DIR"
    
    # Converter todos os arquivos .fga para .json
    for fga_file in *.fga; do
        if [ -f "$fga_file" ]; then
            json_file="${fga_file%.fga}.json"
            echo -e "${BLUE}Converting: $fga_file → $json_file${NC}"
            
            fga model transform --from dsl --to json --file "$fga_file" > "$json_file"
            
            if [ $? -eq 0 ]; then
                echo -e "${GREEN}✅ $fga_file converted successfully${NC}"
            else
                echo -e "${RED}❌ Failed to convert $fga_file${NC}"
            fi
        fi
    done
    
    echo -e "${GREEN}🎉 All conversions completed!${NC}"
}

validate_model() {
    local model_file="$1"
    
    if [ ! -f "$model_file" ]; then
        echo -e "${RED}Model file not found: $model_file${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}Validating model: $model_file${NC}"
    
    # Validar sintaxe do modelo
    fga model validate --file "$model_file"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Model is valid${NC}"
        return 0
    else
        echo -e "${RED}❌ Model validation failed${NC}"
        return 1
    fi
}

main() {
    check_fga_cli
    
    if [ $# -lt 1 ]; then
        print_usage
        exit 1
    fi
    
    local command="$1"
    
    case "$command" in
        dsl-to-json)
            if [ $# -ne 3 ]; then
                print_usage
                exit 1
            fi
            dsl_to_json "$2" "$3"
            ;;
        json-to-dsl)
            if [ $# -ne 3 ]; then
                print_usage
                exit 1
            fi
            json_to_dsl "$2" "$3"
            ;;
        convert-all)
            convert_all
            ;;
        validate)
            if [ $# -ne 2 ]; then
                print_usage
                exit 1
            fi
            validate_model "$2"
            ;;
        *)
            echo -e "${RED}Invalid command: $command${NC}"
            print_usage
            exit 1
            ;;
    esac
}

main "$@"
