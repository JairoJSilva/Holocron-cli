function holo {
    if ($args[0] -eq "use" -or $args[0] -eq "clear") {
        # Executa o comando Go e captura o output
        $output = switch-cli-cloud.exe $args
        if ($LASTEXITCODE -eq 0) {
            # O output de use e clear deve ser os comandos PowerShell para exportar variaveis
            Invoke-Expression $output
        } else {
            # Exibe a mensagem de erro que veio via stderr
            Write-Host $output -ForegroundColor Red
        }
    } else {
        # Passa diretamente para os outros comandos (list, current)
        switch-cli-cloud.exe $args
    }
}
