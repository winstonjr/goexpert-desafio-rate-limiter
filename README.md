# goexpert-desafio-rate-limiter

1. Entrar na pasta raiz do projeto
2. Compilar a imagem do docker
```shell
docker build -t st .
```
3. Executar a imagem do docker
```shell
docker run \
       -e RATE_LIMITER_RULES='{"abc123": {"maxRequests": 10, "limitInSeconds": 1, "blockInSeconds": 1}, "[::1]:59291": {"maxRequests": 10, "limitInSeconds": 2, "blockInSeconds": 2}}' \
       -p 8080:8080 \
       rl
```
4. Executar o teste de carga
```shell
ab -t 13 -c 1 -n 100000000 -H 'API_KEY:abc123' http://localhost:8080/
```

Outra possibilidade é utilizar o docker compose (esta versão utiliza redis):

```shell
docker compose up --detach
```

O teste de carga foi feito da mesma forma que o anterior, usando apache benchmark
```shell
ab -t 13 -c 1 -n 100000000 -H 'API_KEY:abc123' http://localhost:8080/
```
