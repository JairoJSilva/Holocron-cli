package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config representa a configuração completa de um cliente
type Config struct {
	AWS   []AWSAccount   `yaml:"aws,omitempty"`
	GCP   []GCPAccount   `yaml:"gcp,omitempty"`
	OCI   []OCIAccount   `yaml:"oci,omitempty"`
	Azure []AzureAccount `yaml:"azure,omitempty"`
}

type AWSAccount struct {
	Name            string `yaml:"name"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Region          string `yaml:"region"`
}

type GCPAccount struct {
	Name            string `yaml:"name"`
	CredentialsFile string `yaml:"credentials_file"`
}

type OCIAccount struct {
	Name    string `yaml:"name"`
	Profile string `yaml:"profile"`
}

type AzureAccount struct {
	Name      string `yaml:"name"`
	ConfigDir string `yaml:"config_dir"`
}

// ContextEntry representa um contexto disponível para uso
type ContextEntry struct {
	Client  string
	Cloud   string
	Name    string
	FullID  string // client-cloud-name
	Display string // informação extra para exibição (ex: region)
}

// LoadConfig carrega as configurações de um cliente específico
func LoadConfig(client string) (*Config, error) {
	configPath := filepath.Join("configs", fmt.Sprintf("%s.yaml", client))

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("não foi possível ler o arquivo de configuração para %s: %w", client, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("erro ao parsear configuração yaml: %w", err)
	}

	return &cfg, nil
}

// LoadOrCreateConfig carrega o config de um cliente ou retorna um config vazio se o arquivo não existir
func LoadOrCreateConfig(client string) (*Config, error) {
	configPath := filepath.Join("configs", fmt.Sprintf("%s.yaml", client))

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("não foi possível ler o arquivo de configuração para %s: %w", client, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("erro ao parsear configuração yaml: %w", err)
	}

	return &cfg, nil
}

// SaveConfig salva as configurações de um cliente no arquivo YAML correspondente
func SaveConfig(client string, cfg *Config) error {
	configDir := "configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de configurações: %w", err)
	}

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yaml", client))

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração yaml: %w", err)
	}

	header := fmt.Sprintf("# %s.yaml\n", client)
	content := append([]byte(header), data...)

	return os.WriteFile(configPath, content, 0644)
}

// ListClients retorna os nomes dos clientes a partir dos arquivos .yaml em configs/
func ListClients() []string {
	matches, err := filepath.Glob(filepath.Join("configs", "*.yaml"))
	if err != nil {
		return nil
	}
	var clients []string
	for _, m := range matches {
		base := filepath.Base(m)
		name := strings.TrimSuffix(base, ".yaml")
		clients = append(clients, name)
	}
	return clients
}

// ListAllContexts retorna todos os contextos de todos os clientes
func ListAllContexts() ([]ContextEntry, error) {
	clients := ListClients()
	var entries []ContextEntry

	for _, client := range clients {
		cfg, err := LoadConfig(client)
		if err != nil {
			continue
		}

		for _, acc := range cfg.AWS {
			entries = append(entries, ContextEntry{
				Client:  client,
				Cloud:   "aws",
				Name:    acc.Name,
				FullID:  fmt.Sprintf("%s-aws-%s", client, acc.Name),
				Display: acc.Region,
			})
		}
		for _, acc := range cfg.GCP {
			entries = append(entries, ContextEntry{
				Client:  client,
				Cloud:   "gcp",
				Name:    acc.Name,
				FullID:  fmt.Sprintf("%s-gcp-%s", client, acc.Name),
				Display: acc.CredentialsFile,
			})
		}
		for _, acc := range cfg.OCI {
			entries = append(entries, ContextEntry{
				Client:  client,
				Cloud:   "oci",
				Name:    acc.Name,
				FullID:  fmt.Sprintf("%s-oci-%s", client, acc.Name),
				Display: acc.Profile,
			})
		}
		for _, acc := range cfg.Azure {
			entries = append(entries, ContextEntry{
				Client:  client,
				Cloud:   "azure",
				Name:    acc.Name,
				FullID:  fmt.Sprintf("%s-azure-%s", client, acc.Name),
				Display: acc.ConfigDir,
			})
		}
	}

	return entries, nil
}

// FindAccount busca uma conta específica pelo client, cloud e name
func FindAccount(client, cloud, name string) (*Config, interface{}, error) {
	cfg, err := LoadConfig(client)
	if err != nil {
		return nil, nil, err
	}

	switch cloud {
	case "aws":
		for _, acc := range cfg.AWS {
			if acc.Name == name {
				return cfg, acc, nil
			}
		}
	case "gcp":
		for _, acc := range cfg.GCP {
			if acc.Name == name {
				return cfg, acc, nil
			}
		}
	case "oci":
		for _, acc := range cfg.OCI {
			if acc.Name == name {
				return cfg, acc, nil
			}
		}
	case "azure":
		for _, acc := range cfg.Azure {
			if acc.Name == name {
				return cfg, acc, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("conta '%s' não encontrada para %s/%s", name, client, cloud)
}
