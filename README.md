# goexpert-desafio-rate-limiter

## Testando usando o docker compose
1. No diretório raiz do projeto, executar:
```shell
docker compose up --detach
```
2. Se não possuir o apache benchmark, instalar o mesmo.
3. Testar carga usando apache benchmark, como no exemplo abaixo:
```shell
ab -t 13 -c 1 -n 100000000 -H 'API_KEY:abc123' http://localhost:8080/
```
