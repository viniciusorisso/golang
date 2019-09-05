
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

	response, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil{
		log.Fatalln(err2)
	}
	
	var dadosJson Cripto

	err3 := json.Unmarshal(response, &dadosJson)

	if err3 != nil{
		log.Fatalln(err3)
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

func post_json(url string, final Cripto){

	requestBody, err := json.Marshal(final)

	if err != nil{
		log.Fatalln(err)
	}

	resp, err1 := http.Post(url,"answer/json", bytes.NewBuffer(requestBody))
	
	if err1 != nil{	
		log.Fatalln(err1)
	}
	
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil{	
		log.Fatalln(err2)
	}

	print(string(body))
}
// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
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

/*
func teste(url string){
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", "answer.json")
	if err != nil{
		log.Fatalln(err)
	}

	fh, err1 := os.Open("answer.json")
	if err1 != nil{
		log.Fatalln(err1)
	}

	 //iocopy
	 _, err = io.Copy(fileWriter, fh)
	 if err != nil {
		 panic(err)
	 }
 
	 bodyWriter.FormDataContentType()
	 bodyWriter.Close()
 
	 return err
}*/

func main() {

	api_get := "https://api.codenation.dev/v1/challenge/dev-ps/generate-data?token=2ba5541ee2a9cac769de829db6ca75e9c1facf08"
	api_post := "https://api.codenation.dev/v1/challenge/dev-ps/submit-solution?token=2ba5541ee2a9cac769de829db6ca75e9c1facf08"
	final := get_json(api_get)
	final.Decifrado = decode_json(final.Cifrado, final.Numero_casas)
	final.Resumo_criptografico = resume_sha1(final.Decifrado)
	write_json(final)
	//post_json(api_post, final)

	path, _ := os.Getwd()
	path += "/answer.json"

	request, err := newfileUploadRequest(api_post, "multipart/form-data", "answer.json")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
    if err != nil {
			log.Fatal(err)
		}
    resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
	
}
