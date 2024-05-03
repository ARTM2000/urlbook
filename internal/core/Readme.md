This package is the core application logic and is divided further into sub-packages:
- `common`: It contains common utilities used throughout the application, such as "router.go" (for HTTP routing) and "logger.go" (for logging).
- `dto`: This package defines data transfer objects (DTOs) used for passing data between different layers.
- `entity`: It contains the domain entities, representing the core data structures used in the application.
- `model`: This package contains model structures representing specific HTTP request and response bodies.
- `port`: Here, you define __interfaces__ (ports) that represent the required functionalities of the application. For example, “repository” interfaces define methods for accessing data and “service” interfaces define methods for business logic.
- `server`: This package contains the HTTP server setup.
- `service`: This package contains the core application services that handle business logic