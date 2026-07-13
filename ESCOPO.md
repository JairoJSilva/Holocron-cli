# Escopo do Projeto: Holocron CLI

## 1. Objetivo
O **Holocron CLI** tem como objetivo fornecer uma forma simples e rápida de alternar entre diferentes contextos e contas de provedores de nuvem (AWS, GCP, OCI, Azure) utilizando apenas variáveis de ambiente e arquivos de configuração YAML, sem a necessidade de um banco de dados.

## 2. Visão Geral
A ferramenta disponibiliza o comando base `holo`, que permite, com uma única instrução, carregar automaticamente as credenciais do perfil desejado no shell do usuário e limpar o contexto anterior.

### Comportamento por Provedor:
* **AWS**: Exporta as variáveis de ambiente necessárias (ex: `AWS_PROFILE`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, etc).
* **OCI**: Define a variável `OCI_CLI_PROFILE`.
* **Azure**: Define a variável `AZURE_CONFIG_DIR` apontando para o diretório de configuração do contexto.
* **GCP**: Configura a variável `GOOGLE_APPLICATION_CREDENTIALS` apontando para o arquivo JSON de credenciais.

## 3. Estrutura de Clientes e Contextos
O projeto foi dimensionado inicialmente para suportar 6 clientes, onde cada um pode ter ambientes em 4 provedores de nuvem, totalizando **24 contextos** disponíveis.

**Clientes Mapeados:**
1. MV
2. Maida
3. Dentalis
4. GIF
5. Consulfarma
6. Farmácia Digital

## 4. Arquivos de Configuração
Todas as definições de acesso ficam armazenadas no diretório `configs/`. Cada cliente possui seu próprio arquivo YAML (totalizando 6 arquivos). Por exemplo:
* `configs/mv.yaml` (contém até 4 contas de nuvem)
* `configs/maida.yaml` (contém até 4 contas de nuvem)
* `configs/dentalis.yaml` (contém até 4 contas de nuvem)

*Nota: Para a versão inicial, os arquivos YAML utilizam caminhos absolutos do sistema (ex: caminho para o `.json` do GCP ou o config da AWS). Futuramente, o projeto pode ser expandido para se integrar com cofres de senhas.*

## 5. Comandos Disponíveis

1. **Listar contextos:**
   ```bash
   holo list
   ```
   *Retorna a lista de todos os contextos configurados (ex: `mv-aws`, `dentalis-oci`, `maida-gcp`).*

2. **Ativar um contexto:**
   ```bash
   holo use <nome-do-contexto>
   ```
   *Exemplo: `holo use mv-aws`. Realiza a troca imediata preparando e injetando as variáveis no shell pai.*

3. **Ver o contexto atual:**
   ```bash
   holo current
   ```
   *Mostra um resumo instantâneo contendo cliente, nuvem, conta e perfil atualmente ativos.*

4. **Limpar contexto:**
   ```bash
   holo clear
   ```
   *Remove todas as variáveis de ambiente e credenciais atreladas às nuvens do shell atual, limpando o acesso.*

## 6. Estrutura Interna da Solução (Go)
O CLI, desenvolvido na linguagem Go, está organizado da seguinte maneira:

* `cmd/holo/main.go`: Ponto de entrada (entrypoint) do executável.
* `internal/config/`: Lógica responsável pela leitura e formatação dos arquivos YAML.
* `internal/context/`: Lógica central para gerenciar as trocas de contexto.
* `internal/providers/`: Implementações separadas para AWS, OCI, Azure e GCP.
* `configs/`: Armazena os YAMLs.
* `scripts/`: Contém os wrappers `holo.sh` (Bash) e `holo.ps1` (PowerShell). Como binários não podem exportar variáveis diretamente para o shell pai, os wrappers executam o binário e interpretam a sua saída (`eval` ou `Invoke-Expression`).

### Gerenciamento de Estado
A ferramenta mantém um registro persistente do contexto atualmente em uso no arquivo `~/.holo/current`. É através deste arquivo de estado que o comando `holo current` retorna os dados quase que instantaneamente.

## 7. Fluxo de Uso no Dia a Dia
Com essa ferramenta instalada, o desenvolvedor poderá transitar facilmente entre os 24 contextos sem precisar fazer logins manuais através do browser ou editar o `~/.aws/credentials` na mão. Um fluxo de trabalho real seria algo como:

1. `holo use mv-aws` *(Acessar a infraestrutura da AWS do cliente MV)*
2. `holo use farmacia-digital-oci` *(Mudar para a nuvem da OCI do cliente Farmácia Digital)*
3. `holo use maida-azure` *(Trocar para a Azure da Maida)*
4. `holo use dentalis-gcp` *(Trabalhar no ambiente GCP da Dentalis)*
5. `holo clear` *(Limpar qualquer acesso antes de finalizar o trabalho)*
