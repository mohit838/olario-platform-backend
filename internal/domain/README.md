# Domain Layer

The domain layer contains business concepts only. It should not import HTTP,
database, Redis, MinIO, logging, or framework packages.

Current bounded areas:

- `tenant`: tenant ownership boundary
- `identity`: users and roles
- `catalog`: categories, products, and money representation
- `inventory`: stock movement history
- `customer`: customer records
- `ordering`: orders and order items

Repository interfaces and use cases will be added here or in an application
layer when we start building real business endpoints.
