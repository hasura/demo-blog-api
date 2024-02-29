# Demo Blog API
This is simple API with in memory data created to be used during the demo of Hasura's Open API Connector.

### Warning
The Open API Connector is active under developement. It has known limitations and may breakdown in certain cases.

## Running the API
You can run this API using `go run .` command. Any missing dependencies can be downloaded using `go mod download`. This project uses [Swag](https://github.com/swaggo/swag) for generating Open API documentation for APIs.

## Using the Open API Connector
1. Clone the repository from GitHub: https://github.com/hasura/ndc-nodejs-lambda/tree/openapi-prototype
2. Checkout branch `openapi-prototype`
3. Run `npm run build && npm link`. This will build and install the package
4. You can generate the `Api.ts` and `functions.ts` using the following command:
```
npx yo hasura-ndc-nodejs-lambda --open-api $path-open-api-doc-file --base-url "http://localhost:9090"
```
5. Once those files have been generated, you can start the NodeJS Lambda Connector by running `npm run watch` in the directory where `Api.ts` and `functions.ts` were created.

You should have a connector running that you can use in your project's metadata.

# References
Open API Spec: https://swagger.io/specification/
NodeJS Lamda Connector: https://github.com/hasura/ndc-nodejs-lambda/tree/main
Supergraph Modeling/Hasura Metadata: https://hasura.io/docs/3.0/supergraph-modeling/introduction