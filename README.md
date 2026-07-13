# switch-cli-cloud (flowctl)

Uma ferramenta de linha de comando feita em Go para alternar rapidamente entre contextos e perfis de clientes em provedores de nuvem (AWS, GCP, OCI, Azure).

## Estrutura do Projeto

* `configs/`: Contém os arquivos `.yaml` com as definições de perfil e contas por cliente.
* `cmd/switch-cli-cloud/main.go`: Ponto de entrada da CLI.
* `internal/`: Lógica de gerenciamento de estado e geração dos comandos para os provedores.
* `scripts/flowctl.sh`: Wrapper em Bash (Linux/macOS/WSL) que avalia as variáveis de ambiente.
* `scripts/flowctl.ps1`: Wrapper em PowerShell (Windows) que avalia as variáveis de ambiente.

## Como instalar e usar

Como este utilitário modifica variáveis de ambiente no seu shell pai, em vez de chamar o executável binário diretamente, você usará a função wrapper no PowerShell.

### 1. Compilar o Projeto
Certifique-se de ter o Go instalado no WSL. No terminal, execute:
```bash
go build -o switch-cli-cloud ./cmd/switch-cli-cloud
```

### 2. Configurar o Perfil do Shell
Adicione o conteúdo do script `scripts/flowctl.sh` no seu arquivo de Profile (ex: `~/.bashrc` ou `~/.zshrc`) para que o comando `flowctl` esteja sempre disponível, ou chame o script (lembrando de ajustar o caminho do executável `switch-cli-cloud` no script de acordo).

### 3. Comandos Disponíveis

* **Listar contextos:**
  ```bash
  flowctl list
  ```

* **Usar um contexto:**
  ```bash
  flowctl use mv-aws
  ```

* **Ver o contexto atual:**
  ```bash
  flowctl current
  ```

* **Limpar o contexto:**
  ```bash
  flowctl clear
  ```

## Adicionar/Editar Configurações
As configurações residem na pasta `configs/`. Basta preencher os arquivos YAML com seus dados verdadeiros usando caminhos Linux (ex: `/home/user/.gcp/key.json`).

