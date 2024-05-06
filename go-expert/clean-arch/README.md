# Clean Architecture - Orders Service

Desafio do Full Cycle para o módulo de clean architecture.

## Ambiente

A configuração padrão é estabelecida por meio de variáveis de ambiente, essas estão definidas no arquivo `cmd/ordersystem/.env` (ignora-se práticas de segurança por trata-se de um projeto para fins de aprendizagem).

> __Atenção!__ Não há migrações automáticas nesse projeto. Para criar a tabela de `orders`, siga as instruções de [criação da tabela orders](#criação-da-tabela-orders).


### Criação da Tabela `orders`

Para criar a tabela orders suba o ambiente via docker-compose, com o comando:

```shell
docker-compose up -d
```

Em seguida conecte-se ao container mysql e execute a ferramenta mysql para criação da tabela via prompt.

```shell
$ docker-compose exec mysql bash
$> mysql -proot -uroot orders
$mysql> CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id));
```

## Validando REST API

A API REST roda por padrão no endereço  base `http://localhost:8000`.

Para validar a criação de pedidos utilizando o plugin [rest-client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) do VSCode, abra o arquivo `api/create_order.http` e execute a requisição `POST http://localhost:8000/order`.

Para validar a listagem ainda nessa abordagem, abra o arquivo `api/list_order.http` e execute a requisição `GET http://localhost:8000/order`

## Validando gRPC

O servidor gRPC estará rodando por padrão no `localhost:50051`.

Para validar a listagem de `orders` primeiro precisamos ter algum pedido cadastrado no sistema. Para cadastrar utilizando a API gRPC, utilize a ferramenta [`evans`](https://github.com/ktr0731/evans?tab=readme-ov-file#installation), execute os seguintes comandos:

```shell
$ evans -r

> call CreateOrder
```

Após cadastro, liste os pedidos executando o comando:

```shell
$ evans -r

> call ListOrders
```

## Validando GraphQL

O serviço GraphQL roda por padrão no endereço [`http://localhost:8080`](http://localhost:8080).

Para listagem de `orders` precisamos cadastrar ao menos um pedido. Para cadastrar um pedido utilizando a API GraphQL, abra o navegador no endereço informado e execute o seguinte comando:

```graphql
mutation {
  createOrder (input: {
    id: "order-id",
    Price: 10.99,
    Tax: 0.99,
  }){
    id,
    Price,
    Tax,
    FinalPrice
  }
}
```
Com algum pedido já criado, execute a consulta de listagem, com o seguinte comando:

```graphql
query {
  listOrders {
    id,
    Price,
    FinalPrice,
    Tax
  }
}
```