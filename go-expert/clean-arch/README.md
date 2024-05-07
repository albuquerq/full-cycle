# Clean Architecture - Orders Service

Desafio do Full Cycle para o módulo de clean architecture. Esse projeto estende a base de [código criada no curso do Full Cycle](https://github.com/devfullcycle/goexpert/tree/main/20-CleanArch) adicionando o caso de uso de listagem de pedidos e faz pequenos ajustes necessários para o seu funcionamento.

## Ambiente

A configuração padrão é estabelecida por meio de variáveis de ambiente, essas estão definidas no arquivo `cmd/ordersystem/.env` (ignora-se práticas de segurança por trata-se de um projeto para fins de aprendizagem).


### Migrações

As migrações são criadas automaticamente ao inciar a aplicação.

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
