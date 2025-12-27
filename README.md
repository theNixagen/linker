# Linker

**Linker** é uma API para um aplicativo clone do **Linktree / Beacons**, permitindo que usuários centralizem múltiplos links em uma única página.

O projeto começou como um MVP simples para praticar **Go** e aprender **Vue.js**, já que esse tipo de aplicação tem uma complexidade inicial baixa.  
Após a conclusão da primeira API em Go, o projeto evoluiu para um objetivo maior:

> Criar a **mesma API** implementada em **diferentes linguagens**, usando todas as linguagens com as quais eu já desenvolvi pelo menos uma API.

A ideia é comparar abordagens, padrões, desempenho e experiência de desenvolvimento entre linguagens e frameworks distintos, mantendo o mesmo domínio de negócio.
A medida que eu terminar os outros projetos vou colocando os links aqui.

---

## 🧱 Arquitetura

A aplicação utiliza uma arquitetura baseada em serviços externos, todos gerenciados via **Docker Compose**:

- **PostgreSQL** — banco de dados relacional
- **Redis** — cache e possíveis filas
- **MinIO** — storage compatível com S3 (ex: imagens de perfil)
- **API** — implementada em múltiplas linguagens (Go, etc.)

---

## 🚀 Como rodar a aplicação

### 1. Pré-requisitos

Certifique-se de ter instalado:

- Docker
- Docker Compose
- Go
- Git

---

### 2. Subir os serviços de infraestrutura

Inicie os serviços de banco de dados, cache e storage:

```bash
docker compose up -d
```

### 3. Variáveis de ambiente

Copie o conteudo do .env.example para um .env com este comando
```bash
cp .env.example .env
```

### 4. Dependencias

Instale as dependencias com o seguinte comando
```bash
go mod tidy
```

### 5. build e run

Builde a aplicação usando o comando abaixo
```bash
go build -o linker ./cmd/main.go
```

Execute o binario com o comando
```
./linker
```

