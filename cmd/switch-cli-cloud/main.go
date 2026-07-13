package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"switch-cli-cloud/internal/config"
	"switch-cli-cloud/internal/context"
	"switch-cli-cloud/internal/credentials"
	"switch-cli-cloud/internal/providers"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		handleList()
	case "use":
		if len(os.Args) < 3 {
			// Modo interativo: exibe menu numerado
			handleUseInteractive()
		} else {
			handleUse(os.Args[2])
		}
	case "current":
		handleCurrent()
	case "clear":
		handleClear()
	case "add":
		if err := credentials.RunAddWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("switch-cli-cloud - Gerenciador de Contextos de Nuvem")
	fmt.Println("Uso:")
	fmt.Println("  list      Lista todos os contextos disponíveis")
	fmt.Println("  use       Seleciona um contexto interativamente")
	fmt.Println("  use <ctx> Ativa um contexto específico (ex: dentalis-aws-dev)")
	fmt.Println("  current   Mostra o contexto atualmente ativado")
	fmt.Println("  clear     Limpa o contexto atual")
	fmt.Println("  add       Adiciona credenciais de um cliente (interativo)")
}

func handleList() {
	entries, err := config.ListAllContexts()
	if err != nil || len(entries) == 0 {
		fmt.Println("Nenhum contexto configurado. Use 'add' para cadastrar credenciais.")
		return
	}

	currentClient := ""
	for _, e := range entries {
		if e.Client != currentClient {
			if currentClient != "" {
				fmt.Println()
			}
			fmt.Printf("📂 %s\n", e.Client)
			currentClient = e.Client
		}
		display := ""
		if e.Display != "" {
			display = fmt.Sprintf("(%s)", e.Display)
		}
		fmt.Printf("   %-30s %-8s %s\n", e.FullID, strings.ToUpper(e.Cloud), display)
	}
}

func handleUseInteractive() {
	entries, err := config.ListAllContexts()
	if err != nil || len(entries) == 0 {
		fmt.Println("Nenhum contexto configurado. Use 'add' para cadastrar credenciais.")
		os.Exit(1)
	}

	fmt.Println("Selecione o contexto:")
	for i, e := range entries {
		display := ""
		if e.Display != "" {
			display = fmt.Sprintf(" - %s", e.Display)
		}
		fmt.Printf("  %d) %-30s (%s%s)\n", i+1, e.FullID, strings.ToUpper(e.Cloud), display)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Escolha [1-")
	fmt.Print(len(entries))
	fmt.Print("]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var choice int
	if _, err := fmt.Sscanf(input, "%d", &choice); err != nil || choice < 1 || choice > len(entries) {
		fmt.Fprintf(os.Stderr, "Erro: Opção inválida.\n")
		os.Exit(1)
	}

	selected := entries[choice-1]
	handleUse(selected.FullID)
}

// parseContextID separa "client-cloud-name" usando os clouds conhecidos como pivô
// Exemplo: "farmacia-digital-aws-prod" → client="farmacia-digital", cloud="aws", name="prod"
func parseContextID(target string) (client, cloud, name string, err error) {
	knownClouds := []string{"aws", "gcp", "oci", "azure"}

	for _, c := range knownClouds {
		sep := "-" + c + "-"
		idx := strings.Index(target, sep)
		if idx != -1 {
			client = target[:idx]
			cloud = c
			name = target[idx+len(sep):]
			return client, cloud, name, nil
		}
	}

	return "", "", "", fmt.Errorf("contexto inválido: '%s'. Use formato cliente-cloud-nome (ex: dentalis-aws-dev)", target)
}

func handleUse(target string) {
	client, cloud, name, err := parseContextID(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	_, account, err := config.FindAccount(client, cloud, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	var exportCmd string

	switch cloud {
	case "aws":
		acc := account.(config.AWSAccount)
		exportCmd = providers.GenerateAWSExport(acc.AccessKeyID, acc.SecretAccessKey, acc.Region)
	case "gcp":
		acc := account.(config.GCPAccount)
		exportCmd = providers.GenerateGCPExport(acc.CredentialsFile)
	case "oci":
		acc := account.(config.OCIAccount)
		exportCmd = providers.GenerateOCIExport(acc.Profile)
	case "azure":
		acc := account.(config.AzureAccount)
		exportCmd = providers.GenerateAzureExport(acc.ConfigDir)
	}

	// Salva o estado atual
	if err := context.SaveCurrent(client, strings.ToUpper(cloud), name); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar contexto: %v\n", err)
	}

	// Imprime o comando para o wrapper Bash executar via eval
	fmt.Println(exportCmd)
}

func handleCurrent() {
	client, cloud, account, err := context.GetCurrent()
	if err != nil || client == "" {
		fmt.Println("Nenhum contexto ativado no momento.")
		return
	}
	fmt.Printf("Cliente:       %s\n", client)
	fmt.Printf("Cloud:         %s\n", cloud)
	fmt.Printf("Conta/Perfil:  %s\n", account)
	fmt.Printf("Contexto:      %s-%s-%s\n", client, strings.ToLower(cloud), account)
}

func handleClear() {
	if err := context.ClearCurrent(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Erro ao limpar estado: %v\n", err)
		}
	}
	fmt.Println(providers.GenerateClearExport())
}
