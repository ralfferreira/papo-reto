# Papo Reto - Backend

Backend da plataforma de mensagens anônimas organizadas em grupos, desenvolvido em Go.

## Requisitos

- Go 1.24 ou superior
- PostgreSQL
- Redis

## Configuração

1. Clone o repositório:
```bash
git clone https://github.com/ralfferreira/papo-reto.git
cd papo-reto/backend
```

2. Configure as variáveis de ambiente no arquivo `.env` (já existe um arquivo de exemplo).

3. Instale as dependências:
```bash
go mod tidy
```

4. Crie o banco de dados PostgreSQL:
```bash
createdb papo_reto
```

## Executando a aplicação

### Usando Go diretamente

Para executar a aplicação em modo de desenvolvimento:

```bash
go run cmd/api/main.go
```

Ou construa e execute o binário:

```bash
go build -o papo-reto-api cmd/api/main.go
./papo-reto-api
```

O servidor estará disponível em `http://localhost:8080`.

### Usando Docker

O projeto está configurado para ser executado com Docker e Docker Compose, o que facilita a configuração do ambiente de desenvolvimento.

1. Certifique-se de ter o Docker e o Docker Compose instalados em sua máquina.

2. Execute o seguinte comando na raiz do projeto:

```bash
docker-compose up -d
```

Este comando irá:
- Construir a imagem Docker da aplicação
- Iniciar contêineres para a API, PostgreSQL e Redis
- Configurar a rede e os volumes necessários

Para verificar o status dos contêineres:

```bash
docker-compose ps
```

Para visualizar os logs da aplicação:

```bash
docker-compose logs -f api
```

Para parar todos os contêineres:

```bash
docker-compose down
```

O servidor estará disponível em `http://localhost:8080`.


## Estrutura do projeto

```
/backend
  /cmd
    /api          # Ponto de entrada da aplicação
  /internal
    /auth         # Autenticação e autorização
    /config       # Configurações da aplicação
    /handlers     # Handlers HTTP
    /middleware   # Middlewares
    /models       # Modelos de dados
    /repository   # Camada de acesso a dados
    /services     # Lógica de negócios
    /server       # Configuração do servidor HTTP
  /pkg            # Pacotes reutilizáveis
  /migrations     # Migrações de banco de dados
```

## Próximos passos

- Implementar WebSockets para comunicação em tempo real
- Adicionar testes automatizados
- Configurar CI/CD
- Implementar análise de sentimento para mensagens
