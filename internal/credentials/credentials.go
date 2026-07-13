package credentials

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"switch-cli-cloud/internal/config"
)

var reader *bufio.Reader

func prompt(label string) string {
	fmt.Print(label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// RunAddWizard executa o fluxo interativo para adicionar credenciais de um cliente
func RunAddWizard() error {
	reader = bufio.NewReader(os.Stdin)

	// 1. Perguntar o nome do cliente
	client := prompt("? Nome do cliente (ex: mv, maida, dentalis): ")
	if client == "" {
		return fmt.Errorf("nome do cliente não pode ser vazio")
	}
	client = strings.ToLower(client)

	// 2. Perguntar qual cloud provider
	fmt.Println()
	fmt.Println("? Qual cloud provider?")
	fmt.Println("  1) AWS")
	fmt.Println("  2) GCP")
	fmt.Println("  3) OCI")
	fmt.Println("  4) Azure")
	choice := prompt("Escolha [1-4]: ")

	// 3. Carregar config existente ou criar novo
	cfg, err := config.LoadOrCreateConfig(client)
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	fmt.Println()

	// 4. Perguntar o nome da conta
	accountName := prompt("? Nome da conta (ex: dev, prod, main): ")
	if accountName == "" {
		return fmt.Errorf("nome da conta não pode ser vazio")
	}
	accountName = strings.ToLower(accountName)

	fmt.Println()

	// 5. Solicitar os campos de acordo com a cloud selecionada
	switch choice {
	case "1":
		fmt.Println("--- Configuração AWS ---")
		accessKey := prompt("? AWS Access Key ID: ")
		secretKey := prompt("? AWS Secret Access Key: ")
		region := prompt("? AWS Region (ex: us-east-1): ")

		acc := config.AWSAccount{
			Name:            accountName,
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			Region:          region,
		}
		cfg.AWS = append(cfg.AWS, acc)

	case "2":
		fmt.Println("--- Configuração GCP ---")
		credFile := prompt("? Caminho do arquivo de credenciais JSON: ")

		acc := config.GCPAccount{
			Name:            accountName,
			CredentialsFile: credFile,
		}
		cfg.GCP = append(cfg.GCP, acc)

	case "3":
		fmt.Println("--- Configuração OCI ---")
		profile := prompt("? OCI CLI Profile: ")

		acc := config.OCIAccount{
			Name:    accountName,
			Profile: profile,
		}
		cfg.OCI = append(cfg.OCI, acc)

	case "4":
		fmt.Println("--- Configuração Azure ---")
		configDir := prompt("? Caminho do diretório de config do Azure: ")

		acc := config.AzureAccount{
			Name:      accountName,
			ConfigDir: configDir,
		}
		cfg.Azure = append(cfg.Azure, acc)

	default:
		return fmt.Errorf("opção inválida: %s. Use 1, 2, 3 ou 4", choice)
	}

	// 6. Salvar no arquivo YAML
	if err := config.SaveConfig(client, cfg); err != nil {
		return fmt.Errorf("erro ao salvar configuração: %w", err)
	}

	cloudNames := map[string]string{"1": "aws", "2": "gcp", "3": "oci", "4": "azure"}
	cloudName := cloudNames[choice]

	fmt.Printf("\n✅ Conta '%s-%s-%s' adicionada em configs/%s.yaml\n", client, cloudName, accountName, client)
	return nil
}
