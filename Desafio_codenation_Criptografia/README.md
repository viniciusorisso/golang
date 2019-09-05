# golang
# Criptografia de Júlio César

    Desafio proposto pelo site https://www.codenation.dev/ no programa Acelera dev.

# Resumo

Criar um algorítmo que faz um request de um arquivo json usando uma API recebendo:
    
    {
        "numero_casas": 10,
        "token":"token_do_usuario",
        "cifrado": "texto criptografado",
        "decifrado": "aqui vai o texto decifrado",
        "resumo_criptografico": "aqui vai o resumo"
    }
    
E preenche os campos de "decifrado" e "resumo_criptografico" (usando sha1), diante do "numero_casas" e o texto a ser cifrado.
Ao fim, envia por post o arquivo json atualizado e preenchido, obedecendo a seguinte observação:

    OBS: a API espera um arquivo sendo enviado como multipart/form-data, como se fosse enviado por um formulário HTML, com um campo do tipo file com o nome answer. Considere isso ao enviar o arquivo.
