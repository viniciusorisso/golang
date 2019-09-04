
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"strings"
	"bytes"
)
var print = fmt.Println

func jaison (url string) string{

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

	// Por fim escreve o conteúdo do feed de rss
	bodyBytes, err2 := ioutil.ReadAll(resp.Body) 
	
	if err2 != nil{
		log.Fatal(err2)
	}
	
	bodyString := string(bodyBytes) 
	
	escreve_json(bodyString)

	return bodyString
}

func escreve_json(json string) {

	mydata := []byte(json) // Texto em formato para ser escrito

	// the WriteFile method returns an error if unsuccessful - Retirado do site https://tutorialedge.net/golang/reading-writing-files-in-go/
	err := ioutil.WriteFile("answer.json", mydata, 0777)
	// handle this error
	if err != nil {
		// print it out
		fmt.Println(err)
	}
}

func encode_json(texto string, numero int32){

	saida := bytes.Buffer{}
	texto = strings.ToLower(texto)

	for i := range texto{
		if(texto[i] > 97 && texto[i] < 122){
			novo := int32(texto[i]) + numero
			if(novo > 122){
				aux := novo - 122
				print(aux)
				novo = 97 + aux - 1
				print(string(novo))
			}
			final := string(novo)
			saida.WriteString(final)
		}else{
			novo := int32(texto[i])
			saida.WriteString(string(novo))
		}
	}
	print(saida.String())
}

func main() {

	//api := "https://api.codenation.dev/v1/challenge/dev-ps/generate-data?token=2ba5541ee2a9cac769de829db6ca75e9c1facf08"
	//json := jaison(api)
	//fmt.Print(json)
	
	encode_json("az &**bc", 10)
}
