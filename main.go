
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"strings"
	"bytes"
	"mime/multipart"
	"encoding/json"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"io"
	"path/filepath"
)
var print = fmt.Println

type Cripto struct{
	Numero_casas int32
	Token string
	Cifrado string
	Decifrado string
	Resumo_criptografico string
}

func get_json(url string) Cripto{

	resp, err := http.Get(url) //Retorna um ponteiro para um valor do tipo "Response"
	
	if err != nil{	
		log.Fatalln(err)
	}
	
	// Sinaliza que a última ação a ser feita no programa é o fechamento da resposta
	defer resp.Body.Close() 
	
	// Verifica se o código de status é 200, indicando assim o sucesso da solicitação
	if resp.StatusCode != 200 {
		log.Fatalln(resp.StatusCode)	
	}

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil{
		log.Fatalln(err)
	}
	
	var dadosJson Cripto

	err = json.Unmarshal(response, &dadosJson)

	if err != nil{
		log.Fatalln(err)
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

// Creates a new file upload http request with optional extra params
func post_json(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func main() {

	api_get := "https://api.codenation.dev/v1/challenge/dev-ps/generate-data?token='SEU TOKEN AQUI'"
	api_post := "https://api.codenation.dev/v1/challenge/dev-ps/submit-solution?token='SEU TOKEN AQUI'"
	final := get_json(api_get)
	final.Decifrado = decode_json(final.Cifrado, final.Numero_casas)
	final.Resumo_criptografico = resume_sha1(final.Decifrado)
	write_json(final)

	request, err := post_json(api_post, "answer", "answer.json")
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
    if err != nil {
			log.Fatalln(err)
		}
    resp.Body.Close()
	}
}
