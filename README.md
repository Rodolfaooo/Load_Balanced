# Load Balancer - Least Latency (Go)

## Objetivo

Este projeto implementa um Load Balancer utilizando a linguagem Go com o algoritmo **Least Latency (Menor Latência)**.

O objetivo é distribuir as requisições para o servidor que apresenta o menor tempo de resposta no momento da decisão, reduzindo a latência percebida pelo cliente e aumentando a eficiência da utilização dos recursos disponíveis.

---

## Arquitetura

O ambiente é composto por:

* 1 Load Balancer
* 3 Servidores Backend
* Docker Compose para orquestração dos containers

Fluxo de funcionamento:

Cliente → Load Balancer → Backend com menor latência

Os servidores possuem tempos de resposta diferentes para simular um ambiente heterogêneo.

| Servidor | Latência Simulada |
| -------- | ----------------- |
| Server 1 | 50 ms             |
| Server 2 | 2000 ms           |
| Server 3 | 300 ms            |

---

## Estrutura do Projeto

```text
Load_Balanced/
│
├── loadBalancer/
│   └── main.go
│
├── server/
│   └── main.go
│
├── dockerCompose.yml
│
├── DockerFile.lb
├── DockerFile.server
│
└── README.md
```

---

## Tecnologias Utilizadas

* Go
* Docker
* Docker Compose
* HTTP Reverse Proxy
* Goroutines
* Mutexes (controle de concorrência)

---

## Funcionamento do Algoritmo

O Load Balancer executa periodicamente verificações de saúde (Health Checks) em todos os servidores.

Durante cada Health Check são coletadas duas informações:

* Disponibilidade do servidor (Alive)
* Tempo de resposta (Latency)

A cada nova requisição recebida:

1. O balanceador verifica quais servidores estão disponíveis.
2. Compara as latências registradas.
3. Seleciona o servidor com a menor latência.
4. Encaminha a requisição para esse servidor.

Caso um servidor esteja indisponível, ele é removido temporariamente do conjunto de candidatos.

---

## Health Check

Os Health Checks são executados periodicamente utilizando Goroutines.

Objetivos:

* Detectar falhas dos servidores.
* Atualizar as métricas de latência.
* Evitar envio de tráfego para servidores indisponíveis.

---

## Como Executar

### Pré-requisitos

* Docker instalado
* Docker Compose instalado

### Clonar o projeto

```bash
git clone <URL_DO_REPOSITORIO>
cd Load_Balanced
```

### Subir o ambiente

```bash
docker compose up --build
```

Após a inicialização:

* Load Balancer: localhost:8080
* Backends: containers internos da rede Docker

---

## Como Testar

Realizar uma requisição:

```bash
curl http://localhost:8080
```

Executar múltiplas requisições:

```bash
for i in {1..20}
do
  curl http://localhost:8080
done
```

---

## Exemplo de Saída

```text
Resposta do SERVER1
```

Logs do Load Balancer:

```text
SERVER1 -> 50ms
SERVER2 -> 2000ms
SERVER3 -> 300ms

Selecionado: SERVER1
```

---

## Concorrência

O projeto utiliza recursos nativos da linguagem Go:

* Goroutines para execução paralela dos Health Checks.
* RWMutex para sincronização do acesso às informações dos servidores.
* HTTP Server da biblioteca padrão.

Esses mecanismos garantem que múltiplas requisições possam ser processadas simultaneamente sem condições de corrida (race conditions).

---

## Análise dos Resultados

Durante os testes foi configurada uma latência artificial nos servidores para simular diferenças de desempenho entre os nós.

Observou-se que o algoritmo Least Latency adaptou-se dinamicamente às condições da infraestrutura. O servidor com menor tempo de resposta recebeu a maior parte das requisições, enquanto os servidores mais lentos foram utilizados apenas quando apresentavam latências competitivas.

Diferentemente de algoritmos como Round Robin, que distribuem as requisições de forma igualitária independentemente do desempenho dos servidores, o Least Latency considera o estado atual da infraestrutura e busca minimizar o tempo de resposta para o cliente.

Além disso, quando um servidor se torna indisponível, o mecanismo de Health Check identifica a falha e remove automaticamente o nó da seleção, aumentando a resiliência do sistema.
