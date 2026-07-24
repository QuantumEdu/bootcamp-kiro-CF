# Requirements Document

## Introduction

Deploy the POS AI-First Go application to AWS Lambda using a container image, with API Gateway fronting the chi router via algnhsa. The migration replaces only infrastructure adapters — domain and application layers remain unchanged. The system operates in dual-mode: SQLite for local development, PostgreSQL (RDS) for AWS production. AI queries transition from OpenRouter (local) to Amazon Bedrock (production). Infrastructure is defined as code using AWS SAM, with CI/CD through GitHub Actions.

## Glossary

- **Lambda_Handler**: The AWS Lambda function entry point that wraps the chi router using the algnhsa adapter to process API Gateway events
- **API_Gateway**: AWS API Gateway HTTP API that routes incoming HTTP requests to the Lambda_Handler
- **PostgreSQL_Adapter**: Infrastructure adapter implementing domain repository ports using PostgreSQL (pgx/v5 driver) instead of SQLite
- **Bedrock_Adapter**: Infrastructure adapter implementing the AIQueryService port using Amazon Bedrock Runtime API with Claude models
- **SAM_Template**: AWS Serverless Application Model template.yaml file defining all cloud infrastructure resources
- **Secrets_Manager**: AWS Secrets Manager service used to store and retrieve sensitive configuration values
- **Dual_Mode_Startup**: The application bootstrap logic that selects SQLite or PostgreSQL adapters based on the APP_ENV environment variable
- **Session_Store**: The alexedwards/scs session persistence layer, using SQLite3 store locally and PostgreSQL store on AWS
- **CI_CD_Pipeline**: GitHub Actions workflow that tests, builds the container image, and deploys via SAM
- **Health_Endpoint**: The GET /health route providing application liveness and dependency status information
- **Static_Assets_CDN**: S3 bucket with CloudFront distribution serving static files (CSS, JS, images)
- **Migration_Runner**: PostgreSQL DDL migration scripts equivalent to the existing SQLite migrations

## Requirements

### Requirement 1: Lambda Handler and API Gateway Integration

**User Story:** As a developer, I want the existing chi router to run inside AWS Lambda via algnhsa, so that the application serves all HTTP routes through API Gateway without code changes to handlers or templates.

#### Acceptance Criteria

1. WHEN the Lambda function receives an API Gateway HTTP event, THE Lambda_Handler SHALL route the event through the chi router via algnhsa and return the HTTP response to API Gateway.
2. THE Lambda_Handler SHALL serve all existing routes (pages, API endpoints, HTMX fragments, static file references) through a single Lambda function.
3. WHEN the APP_ENV environment variable equals "lambda", THE Lambda_Handler SHALL start the algnhsa adapter instead of the net/http ListenAndServe server.
4. WHEN the APP_ENV environment variable equals "local" or is unset, THE Lambda_Handler SHALL start the standard net/http server on the configured port.
5. THE API_Gateway SHALL be configured as an HTTP API (v2) with a catch-all route proxying all requests to the Lambda_Handler.
6. THE Lambda_Handler SHALL use a container image runtime with the Go binary, templates directory, and compiled assets packaged in the image.

### Requirement 2: PostgreSQL Repository Adapters

**User Story:** As a developer, I want PostgreSQL implementations of all seven repository ports, so that the application persists data in RDS when running on AWS.

#### Acceptance Criteria

1. THE PostgreSQL_Adapter SHALL implement the ProductRepository port with identical behavior to the SQLite implementation, using pgx/v5 for database operations.
2. THE PostgreSQL_Adapter SHALL implement the SaleRepository port, persisting sales and sale items within a database transaction.
3. THE PostgreSQL_Adapter SHALL implement the UserRepository port, including FindByID, FindByPINHash, FindAll, IncrementFailedAttempts, Lock, and ResetAttempts operations.
4. THE PostgreSQL_Adapter SHALL implement the ClientRepository port with Create and List operations.
5. THE PostgreSQL_Adapter SHALL implement the InventoryRepository port with Create and FindByProduct operations.
6. THE PostgreSQL_Adapter SHALL implement the ConfigRepository port with Get and Set (upsert) operations.
7. THE PostgreSQL_Adapter SHALL implement the MetricsRepository port with time-based aggregation queries using PostgreSQL date functions.
8. WHEN a PostgreSQL query fails, THE PostgreSQL_Adapter SHALL return a wrapped error with context identifying the operation and entity.
9. THE PostgreSQL_Adapter SHALL use parameterized queries for all database operations to prevent SQL injection.
10. THE PostgreSQL_Adapter SHALL accept a *pgxpool.Pool connection pool as its constructor dependency.

### Requirement 3: PostgreSQL Migrations

**User Story:** As a developer, I want PostgreSQL-equivalent DDL migrations for the existing SQLite schema, so that the RDS database has identical table structures with PostgreSQL-native types.

#### Acceptance Criteria

1. THE Migration_Runner SHALL create all tables defined in the SQLite migrations (usuarios, categorias, productos, clientes, ventas, venta_items, inventario_movimientos, configuracion, sessions) using PostgreSQL-native data types.
2. WHEN translating the SQLite schema, THE Migration_Runner SHALL convert INTEGER PRIMARY KEY AUTOINCREMENT to SERIAL PRIMARY KEY, REAL to NUMERIC(12,2), TEXT datetime to TIMESTAMPTZ with DEFAULT NOW(), and INTEGER boolean to BOOLEAN.
3. THE Migration_Runner SHALL create all indexes defined in the SQLite migrations with equivalent PostgreSQL CREATE INDEX statements.
4. THE Migration_Runner SHALL create a trigger function equivalent to trg_inventario_actualiza_stock using PostgreSQL PL/pgSQL.
5. THE Migration_Runner SHALL create the sessions table compatible with the alexedwards/scs pgxstore format (token TEXT PRIMARY KEY, data BYTEA, expiry TIMESTAMPTZ).
6. THE Migration_Runner SHALL be executable as standalone SQL files that can run via psql or an embedded migration tool at application startup.

### Requirement 4: Bedrock Adapter for NL-to-SQL

**User Story:** As a developer, I want a Bedrock adapter implementing the AIQueryService port, so that natural language queries use Amazon Bedrock instead of OpenRouter when running on AWS.

#### Acceptance Criteria

1. THE Bedrock_Adapter SHALL implement the AIQueryService port interface (GenerateSQL method with context, question input and sql, explanation, error output).
2. THE Bedrock_Adapter SHALL invoke the Amazon Bedrock Runtime InvokeModel API using the AWS SDK Go v2 bedrockruntime client.
3. THE Bedrock_Adapter SHALL use the same system prompt and response parsing logic as the existing OpenRouter adapter for NL-to-SQL translation.
4. WHEN the Bedrock API returns an error or times out, THE Bedrock_Adapter SHALL return a wrapped error without exposing internal AWS error details to the caller.
5. THE Bedrock_Adapter SHALL accept the model ID, AWS region, and inference parameters (temperature, max tokens) as constructor configuration.
6. THE Bedrock_Adapter SHALL use the Anthropic Messages API format (Claude models) when constructing the Bedrock request payload.
7. WHEN the generated SQL fails validation (non-SELECT, disallowed tables), THE Bedrock_Adapter SHALL return an error consistent with the existing NL-SQL validation rules.

### Requirement 5: Secrets Manager Integration

**User Story:** As a developer, I want the application to load sensitive configuration from AWS Secrets Manager when running on Lambda, so that credentials are not stored in environment variables or code.

#### Acceptance Criteria

1. WHEN APP_ENV equals "lambda", THE Secrets_Manager integration SHALL retrieve configuration values from AWS Secrets Manager at application startup.
2. THE Secrets_Manager integration SHALL retrieve the PostgreSQL connection string from a secret identified by the SECRET_DB_ARN environment variable.
3. THE Secrets_Manager integration SHALL retrieve the session encryption key from a secret identified by the SECRET_SESSION_ARN environment variable.
4. THE Secrets_Manager integration SHALL retrieve Bedrock configuration (model ID, region, parameters) from a secret identified by the SECRET_AI_ARN environment variable.
5. IF a Secrets Manager retrieval fails at startup, THEN THE Lambda_Handler SHALL log the error and terminate with a non-zero exit code.
6. WHEN APP_ENV equals "local" or is unset, THE application SHALL load configuration from environment variables and .env file as it does currently.
7. THE Secrets_Manager integration SHALL cache retrieved secret values in memory for the lifetime of the Lambda execution environment to avoid repeated API calls on warm starts.

### Requirement 6: SAM Template (Infrastructure as Code)

**User Story:** As a developer, I want a SAM template.yaml defining all AWS resources, so that the infrastructure is reproducible and deployable with a single sam deploy command.

#### Acceptance Criteria

1. THE SAM_Template SHALL define an AWS::Serverless::Function resource for the Lambda function using the container image package type.
2. THE SAM_Template SHALL define an AWS::Serverless::HttpApi resource as the API Gateway HTTP API with a catch-all route to the Lambda function.
3. THE SAM_Template SHALL define parameters for environment name, database secret ARN, session secret ARN, AI secret ARN, and VPC configuration.
4. THE SAM_Template SHALL configure the Lambda function with a memory size of 512MB, timeout of 30 seconds, and the us-east-1 region.
5. THE SAM_Template SHALL grant the Lambda function IAM permissions to read from Secrets Manager (specified secret ARNs) and invoke Bedrock models.
6. THE SAM_Template SHALL configure the Lambda function within a VPC (private subnets) to enable connectivity to RDS PostgreSQL.
7. THE SAM_Template SHALL define Outputs for the API Gateway endpoint URL and the Lambda function ARN.
8. THE SAM_Template SHALL include a Globals section setting the function runtime, architecture (arm64), and common environment variables.

### Requirement 7: CI/CD Pipeline

**User Story:** As a developer, I want a GitHub Actions workflow that tests, builds, and deploys the application to AWS on push to main, so that deployments are automated and reliable.

#### Acceptance Criteria

1. WHEN code is pushed to the main branch, THE CI_CD_Pipeline SHALL trigger the deployment workflow.
2. THE CI_CD_Pipeline SHALL run go test ./... with coverage and golangci-lint run as a prerequisite before building.
3. IF tests or linting fail, THEN THE CI_CD_Pipeline SHALL abort the workflow and report the failure.
4. THE CI_CD_Pipeline SHALL build the container image using the project Dockerfile and push it to Amazon ECR.
5. THE CI_CD_Pipeline SHALL execute sam deploy with the built image URI and configured parameters to deploy the stack.
6. THE CI_CD_Pipeline SHALL use OIDC-based authentication (aws-actions/configure-aws-credentials) with a dedicated IAM role for deployments.
7. THE CI_CD_Pipeline SHALL run a post-deploy health check by calling the API Gateway Health_Endpoint and verifying a 200 response.
8. IF the post-deploy health check fails, THEN THE CI_CD_Pipeline SHALL report the failure status without automatic rollback.

### Requirement 8: Dual-Mode Startup

**User Story:** As a developer, I want the application to select adapters (SQLite vs PostgreSQL, OpenRouter vs Bedrock) based on the APP_ENV variable, so that the same codebase runs locally and on AWS without code changes.

#### Acceptance Criteria

1. WHEN APP_ENV equals "local" or is unset, THE Dual_Mode_Startup SHALL initialize SQLite database connection, SQLite repository adapters, SQLite session store, and the OpenRouter AI adapter.
2. WHEN APP_ENV equals "lambda", THE Dual_Mode_Startup SHALL initialize PostgreSQL connection pool, PostgreSQL repository adapters, pgx session store, and the Bedrock AI adapter.
3. THE Dual_Mode_Startup SHALL use the same dependency injection pattern as the current main.go, passing adapters to use cases and handlers via constructors.
4. WHEN APP_ENV equals "lambda", THE Dual_Mode_Startup SHALL initialize the PostgreSQL connection pool with a maximum of 5 connections and a 5-second connect timeout.
5. THE Dual_Mode_Startup SHALL validate that all required configuration values are present before initializing adapters, and terminate with a descriptive error if any are missing.
6. THE Dual_Mode_Startup SHALL run PostgreSQL migrations automatically on first Lambda cold start when APP_ENV equals "lambda".

### Requirement 9: S3 and CloudFront for Static Assets

**User Story:** As a developer, I want static files served from S3 via CloudFront, so that Lambda does not handle static asset requests and users get fast CDN delivery.

#### Acceptance Criteria

1. THE Static_Assets_CDN SHALL serve all files from the static/ directory (CSS, JS, images) through a CloudFront distribution backed by an S3 origin.
2. THE SAM_Template SHALL define an S3 bucket resource for static assets with appropriate bucket policy allowing CloudFront access.
3. THE CI_CD_Pipeline SHALL sync the static/ directory to the S3 bucket during deployment using aws s3 sync.
4. WHEN running on Lambda, THE Lambda_Handler SHALL exclude the /static/* route from chi routing, as CloudFront handles those requests directly.
5. WHEN running locally, THE Lambda_Handler SHALL continue serving static files from the filesystem via the existing http.FileServer handler.
6. THE Static_Assets_CDN SHALL configure CloudFront with a cache policy of 24 hours for immutable assets and 1 hour for other static files.

### Requirement 10: Health Check and Monitoring

**User Story:** As a developer, I want a health endpoint that reports application and dependency status, so that monitoring systems can detect failures and deployments can verify success.

#### Acceptance Criteria

1. THE Health_Endpoint SHALL respond to GET /health with a JSON payload containing status, database connectivity, AI service availability, and the current APP_ENV value.
2. WHEN the database connection is healthy, THE Health_Endpoint SHALL include "database": "ok" in the response.
3. WHEN the database connection fails a ping, THE Health_Endpoint SHALL include "database": "error" in the response and return HTTP 503.
4. THE Health_Endpoint SHALL respond within 3 seconds, using a 2-second timeout for dependency checks.
5. WHEN running on Lambda, THE Health_Endpoint SHALL report the Lambda function memory allocation and cold-start status in the response metadata.
6. THE Lambda_Handler SHALL emit structured JSON logs (with timestamp, level, request_id fields) to CloudWatch Logs for all requests and errors.
7. IF an unhandled panic occurs, THEN THE Lambda_Handler SHALL recover, log the panic with stack trace, and return HTTP 500 without crashing the Lambda execution environment.
