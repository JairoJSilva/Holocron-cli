#!/bin/bash

# ==============================================================
# deploy.sh - Script de deploy com versionamento semântico
# Repositório: git@github.com:JairoJSilva/Holocron-cli.git
# ==============================================================

set -e

REPO_URL="git@github.com:JairoJSilva/Holocron-cli.git"
BRANCH="main"

# Cores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}   Holocron CLI - Deploy Script${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 1. Inicializar repositório Git se necessário
if [ ! -d ".git" ]; then
    echo -e "${YELLOW}⚙  Inicializando repositório Git...${NC}"
    git init
else
    echo -e "${GREEN}✔  Repositório Git já inicializado.${NC}"
fi

# 2. Configurar remote
CURRENT_REMOTE=$(git remote get-url origin 2>/dev/null || echo "")
if [ -z "$CURRENT_REMOTE" ]; then
    echo -e "${YELLOW}⚙  Adicionando remote origin...${NC}"
    git remote add origin "$REPO_URL"
elif [ "$CURRENT_REMOTE" != "$REPO_URL" ]; then
    echo -e "${YELLOW}⚙  Atualizando remote origin...${NC}"
    git remote set-url origin "$REPO_URL"
else
    echo -e "${GREEN}✔  Remote origin já configurado.${NC}"
fi

echo -e "${GREEN}   Remote: ${REPO_URL}${NC}"
echo ""

# 3. Adicionar todos os arquivos
echo -e "${YELLOW}⚙  Adicionando arquivos ao staging...${NC}"
git add -A

# 4. Mostrar o que será commitado
echo ""
echo -e "${CYAN}Arquivos no staging:${NC}"
git status --short
echo ""

# 5. Pedir a mensagem do commit
read -p "📝 Mensagem do commit: " COMMIT_MSG

if [ -z "$COMMIT_MSG" ]; then
    COMMIT_MSG="chore: atualização geral"
fi

# 6. Fazer o commit
echo ""
echo -e "${YELLOW}⚙  Criando commit...${NC}"
git commit -m "$COMMIT_MSG"

# 7. Garantir que a branch é main e fazer push
echo ""
git branch -M "$BRANCH"
echo -e "${YELLOW}⚙  Enviando para ${BRANCH}...${NC}"
git push -u origin "$BRANCH"

echo ""
echo -e "${GREEN}✅ Push para ${BRANCH} realizado com sucesso!${NC}"

# ==============================================================
# VERSIONAMENTO SEMÂNTICO
# ==============================================================

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}   Versionamento Semântico${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# Perguntar se quer criar tag
read -p "🏷  Deseja criar uma tag de versão? (s/n): " CREATE_TAG

if [ "$CREATE_TAG" != "s" ] && [ "$CREATE_TAG" != "S" ]; then
    echo ""
    echo -e "${GREEN}✅ Deploy finalizado sem criação de tag.${NC}"
    exit 0
fi

# Buscar última tag existente
LATEST_TAG=$(git tag -l "v*" --sort=-v:refname | head -n 1 2>/dev/null || echo "")

if [ -z "$LATEST_TAG" ]; then
    echo -e "${YELLOW}   Nenhuma tag encontrada. Primeira versão será criada.${NC}"
    LATEST_TAG="v0.0.0"
else
    echo -e "${GREEN}   Última tag: ${LATEST_TAG}${NC}"
fi

# Extrair major, minor, patch
VERSION="${LATEST_TAG#v}"
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

# Calcular sugestões
NEXT_PATCH="v${MAJOR}.${MINOR}.$((PATCH + 1))"
NEXT_MINOR="v${MAJOR}.$((MINOR + 1)).0"
NEXT_MAJOR="v$((MAJOR + 1)).0.0"

echo ""
echo -e "${CYAN}Qual tipo de versão deseja criar?${NC}"
echo -e "  1) ${GREEN}patch${NC}  → ${NEXT_PATCH}   (correções de bugs)"
echo -e "  2) ${YELLOW}minor${NC}  → ${NEXT_MINOR}   (novas funcionalidades)"
echo -e "  3) ${RED}major${NC}  → ${NEXT_MAJOR}   (breaking changes)"
echo -e "  4) ${CYAN}custom${NC} → Informar manualmente"
echo ""
read -p "Escolha [1-4]: " TAG_CHOICE

case "$TAG_CHOICE" in
    1) NEW_TAG="$NEXT_PATCH" ;;
    2) NEW_TAG="$NEXT_MINOR" ;;
    3) NEW_TAG="$NEXT_MAJOR" ;;
    4)
        read -p "Informe a versão (ex: v1.2.3): " NEW_TAG
        # Garantir que começa com 'v'
        if [[ ! "$NEW_TAG" =~ ^v ]]; then
            NEW_TAG="v${NEW_TAG}"
        fi
        # Validar formato
        if [[ ! "$NEW_TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo -e "${RED}❌ Formato inválido. Use o padrão vX.Y.Z (ex: v1.0.0)${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}❌ Opção inválida.${NC}"
        exit 1
        ;;
esac

# Pedir mensagem da tag
echo ""
read -p "📝 Mensagem da tag ${NEW_TAG} (Enter para usar a mensagem do commit): " TAG_MSG

if [ -z "$TAG_MSG" ]; then
    TAG_MSG="$COMMIT_MSG"
fi

# Criar e enviar a tag
echo ""
echo -e "${YELLOW}⚙  Criando tag ${NEW_TAG}...${NC}"
git tag -a "$NEW_TAG" -m "$TAG_MSG"

echo -e "${YELLOW}⚙  Enviando tag para o repositório...${NC}"
git push origin "$NEW_TAG"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ Deploy finalizado!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   Branch:  ${BRANCH}${NC}"
echo -e "${GREEN}   Tag:     ${NEW_TAG}${NC}"
echo -e "${GREEN}   Commit:  ${COMMIT_MSG}${NC}"
echo -e "${GREEN}========================================${NC}"
