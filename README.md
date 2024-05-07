# Desafio Stress-Test

## Descrição 

Objetivo: Criar um sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.


O sistema deverá gerar um relatório com informações específicas após a execução dos testes.

Entrada de Parâmetros via CLI:

-url: URL do serviço a ser testado.  
-requests: Número total de requests.  
-concurrency: Número de chamadas simultâneas.  

Execução do Teste:

Realizar requests HTTP para a URL especificada.
Distribuir os requests de acordo com o nível de concorrência definido.
Garantir que o número total de requests seja cumprido.
Geração de Relatório:

Apresentar um relatório ao final dos testes contendo:
- Tempo total gasto na execução
- Quantidade total de requests realizados.
- Quantidade de requests com status HTTP 200.
- Distribuição de outros códigos de status HTTP (como 404, 500, etc.).


## Configuração Projeto

### Pré-requisitos (testado no linux)

- Make versão 4.3
- Go versão 1.22.2
- Docker versão 24.0.7
- Docker Compose versão v2.3.3

### Comandos

## Via make
`build-docker`: Cria a imagem Docker para o projeto. A imagem Docker é marcada com o nome do projeto e "latest".

`build-go`: Este comando compila o código-fonte Go em um binário. O binário é nomeado de acordo com o nome do projeto (com hífens substituídos por sublinhados).

`docker-run`: Este comando primeiro cria a imagem Docker (usando o comando `build-docker`) e, em seguida, executa a imagem Docker. O contêiner Docker é executado com rede do host. A URL, o número de solicitações e o nível de concorrência para o teste de estresse podem ser especificados com as variáveis `url`, `requests` e `concurrency`, respectivamente.

`go-run`: Este comando primeiro compila o código-fonte Go em um binário (usando o comando `build-go`) e, em seguida, executa o binário. A URL, o número de solicitações e o nível de concorrência para o teste de estresse podem ser especificados com as variáveis `url`, `requests` e `concurrency`, respectivamente.

Recomenda-se executar via make, pois qualquer modificação no código-fonte será refletido no binário ou imagem docker
```bash
make go-run url=http://example.com requests=1000 concurrency=10
```
Para executar o teste de estresse com a imagem Docker:
```bash
make docker-run url=http://example.com requests=1000 concurrency=10
```

Substitua http://example.com, 1000 e 10 pela URL desejada, número de solicitações e nível de concorrência, respectivamente.

### Exemplos
Alguns exemplos foram criados para simplificar o teste.
Apenas troque o `%` por um número de 1 à 9:
```bash
make run-url-%
```

## Diretamente
Caso utilize diretamente o bínário ou imagem docker passe os parâmetros por flag:

```bash
make build-docker
docker run --network=host stress-test:latest -url=http://example.com -requests=1000 -concurrency=10
```

```bash
make build-go 
./stress_test stress-test:latest -url=http://example.com -requests=1000 -concurrency=10
```