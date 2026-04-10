# Nome do arquivo final executável
APP_NAME=validador_excel
MAIN_PATH=./src/main.go

# Diretório de saída
BIN_DIR=bin

.PHONY: all clean deps linux windows windows-docker

all: clean deps linux windows

# Limpa builds anteriores e recria a pasta bin
clean:
	@echo "Limpando diretórios de build..."
	rm -rf $(BIN_DIR)
	rm -rf fyne-cross
	mkdir -p $(BIN_DIR)

# Baixa e atualiza as dependências do Go
deps:
	@echo "Atualizando módulos do Go..."
	go mod tidy

# Build nativo para Linux
linux: deps
	@echo "Compilando para Linux..."
	go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME)_linux $(MAIN_PATH)
	@echo "✓ Build Linux concluído em: $(BIN_DIR)/$(APP_NAME)_linux"

# Build cruzado para Windows (Requer mingw-w64)
# O parâmetro -H=windowsgui esconde aquela tela preta de terminal do Windows
windows: deps
	@echo "Compilando para Windows (MinGW)..."
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui -s -w" -o $(BIN_DIR)/$(APP_NAME).exe $(MAIN_PATH)
	@echo "✓ Build Windows concluído em: $(BIN_DIR)/$(APP_NAME).exe"

# Build cruzado para Windows (Requer Docker + fyne-cross)
# Use este comando se não quiser instalar o MinGW no seu Linux
windows-docker: deps
	@echo "Compilando para Windows via fyne-cross..."
	fyne-cross windows -arch amd64 -output $(APP_NAME).exe ./src
	@echo "✓ Build Windows gerado na pasta fyne-cross/bin/windows-amd64/"