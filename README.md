Nazo Playground API
===================

## How to use

![how-to-use](https://user-images.githubusercontent.com/29615019/153715441-a23dcc2a-8435-4c05-b458-a3493a80d628.gif)

1. Run the Nazo API server.

    ```
    $ go run cmd/nazo/main.go -c examples/simple-nazo
    ```

2. Get first nazo.

    ```
    $ curl http://localhost:8080/api/v1/nazo/
    1 + 1 = ?
    ```

3. Post your answer.

    ```
    $ curl -X POST -d '{"answer": "2"}' http://localhost:8080/api/v1/nazo/answer
    {"result":true,"next":"05b1f9cf-8543-44e5-b54b-e367d8b93ab3"}
    ```

4. Get next nazo.

    ```
    $ curl http://localhost:8080/api/v1/nazo/05b1f9cf-8543-44e5-b54b-e367d8b93ab3
    tone - one = ?
    ```

5. Post your answer.

    ```
    $ curl -X POST -d '{"answer": "999kg"}' http://localhost:8080/api/v1/nazo/05b1f9cf-8543-44e5-b54b-e367d8b93ab3/answer
    {"result":false}
    ```

6. Edit `examples/simple-nazo/conf.yaml` as you like and restart the Nazo API server.

## API Document

https://hiochiai.github.io/nazo-playground/
