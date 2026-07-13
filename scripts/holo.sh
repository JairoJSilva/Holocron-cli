#!/bin/bash

holo() {
    if [ "$1" = "use" ] || [ "$1" = "clear" ]; then
        # Executa o binário Go e captura a saída e o código de retorno
        OUTPUT=$(./switch-cli-cloud "$@")
        EXIT_CODE=$?
        
        if [ $EXIT_CODE -eq 0 ]; then
            # Avalia a saída para exportar/limpar as variáveis de ambiente
            eval "$OUTPUT"
        else
            # Em caso de erro, apenas mostra a saída (sem eval)
            echo "$OUTPUT"
        fi
    else
        # Repassa para os outros comandos diretamente
        ./switch-cli-cloud "$@"
    fi
}
