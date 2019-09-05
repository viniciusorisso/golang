
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"strings"
	"bytes"
	"encoding/json"
	"crypto/sha1"
	"encoding/hex"
)
var print = fmt.Println

type Cripto struct{
	Numero_casas int32
	Token string
	Cifrado string
	Decifrado string
	Resumo_criptografico string
}

func teste(url string) Cripto{

	resp, err := http.Get(url) //Retorna um ponteiro para um valor do tipo "Response"
	
	if err != nil{	
		log.Fatal(err)
	}
	
	// Sinaliza que a última ação a ser feita no programa é o fechamento da resposta
	defer resp.Body.Close() 
	
	// Verifica se o código de status é 200, indicando assim o sucesso da solicitação
	if resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)	
	}

	response, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil{
		log.Fatal(err2)
	}
	
	var dadosJson Cripto

	err3 := json.Unmarshal(response, &dadosJson)

	if err3 != nil{
		log.Fatal(err3)
	}

	return dadosJson
}

func encode_json(texto string, numero int32) string{

	saida := bytes.Buffer{}
	texto = strings.ToLower(texto)

	for i := range texto{
		if(texto[i] >= 97 && texto[i] <= 122){
			novo := int32(texto[i]) + numero
			if(novo > 122){
				aux := novo - 122
				novo = 96 + aux
			}
			final := string(novo)
			saida.WriteString(final)
		}else{
			novo := int32(texto[i])
			saida.WriteString(string(novo))
		}
	}
	print(saida.String())
	return saida.String()
}

func decode_json(texto string, numero int32) string{

	saida := bytes.Buffer{}
	texto = strings.ToLower(texto)

	for i := range texto{
		if(texto[i] >= 97 && texto[i] <= 122){
			novo := int32(texto[i]) - numero
			if(novo < 97){
				aux := 97 - novo
				novo = 123 - aux
			}
			final := string(novo)
			saida.WriteString(final)
		}else{
			novo := int32(texto[i])
			saida.WriteString(string(novo))
		}
	}
	return saida.String()
}

func write_json(final Cripto){

	file, _ := json.MarshalIndent(final, "", " ")
 
	_ = ioutil.WriteFile("answer.json", file, 0644)
}

func resume_sha1(texto string) string{

    h := sha1.New()
    h.Write([]byte(texto))
    resume := hex.EncodeToString(h.Sum(nil))
	
	return resume
}

func main() {

	api := "https://api.codenation.dev/v1/challenge/dev-ps/generate-data?token='COLOCAR O TOKEN AQUI'"
	final := teste(api)
	final.Decifrado = decode_json(final.Cifrado, final.Numero_casas)
	final.Resumo_criptografico = resume_sha1(final.Decifrado)
	write_json(final)

}
