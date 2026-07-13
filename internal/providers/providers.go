package providers

import "fmt"

func GenerateAWSExport(accessKeyID, secretAccessKey, region string) string {
	return fmt.Sprintf(`export AWS_ACCESS_KEY_ID="%s"; export AWS_SECRET_ACCESS_KEY="%s"; export AWS_REGION="%s"`, accessKeyID, secretAccessKey, region)
}

func GenerateGCPExport(credentialsFile string) string {
	return fmt.Sprintf(`export GOOGLE_APPLICATION_CREDENTIALS="%s"`, credentialsFile)
}

func GenerateOCIExport(profile string) string {
	return fmt.Sprintf(`export OCI_CLI_PROFILE="%s"`, profile)
}

func GenerateAzureExport(configDir string) string {
	return fmt.Sprintf(`export AZURE_CONFIG_DIR="%s"`, configDir)
}

func GenerateClearExport() string {
	return `unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_REGION GOOGLE_APPLICATION_CREDENTIALS OCI_CLI_PROFILE AZURE_CONFIG_DIR`
}
