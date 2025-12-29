# Linker

**Linker** é uma API para um aplicativo clone do **Linktree / Beacons**, permitindo que usuários centralizem múltiplos links em uma única página.

O projeto começou como um MVP simples para praticar **Go** e **React**, já que esse tipo de aplicação tem uma complexidade inicial baixa.  

---

## 🧱 Arquitetura

A aplicação utiliza uma arquitetura baseada em serviços externos, todos gerenciados via **Docker Compose**:

- **PostgreSQL** — banco de dados relacional
- **Redis** — cache
- **MinIO** — storage compatível com S3 (ex: imagens de perfil)
- **API** — implementada em Go

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

