# Cross-Language gRPC Person Service | Python & Go

**One Protocol Buffers contract. Two backend implementations. Four CRUD operations.**

A hands-on backend engineering project exploring how independently implemented **Python and Go** services communicate through **gRPC over HTTP/2**. Both languages use generated, strongly defined API interfaces from a shared `.proto` contract. Each server implements the same Person CRUD operations using in-memory storage, making the networking and service boundaries easy to understand without introducing database complexity.

**Tech stack:** Python · Go · gRPC · Protocol Buffers (`proto3`) · HTTP/2 · In-memory storage

## 1. Why This Project?

In a real microservices environment, teams often build services in different languages. They need a reliable way to agree on request/response formats without sharing internal application code.

This project demonstrates **contract-first API development**: define `PersonService` once in `proto/person.proto`, generate language-specific code, and implement compatible clients and servers independently in Python and Go.

**Core engineering topics:** service contracts, binary serialization, remote procedure calls, typed code generation, language interoperability, API error semantics, and concurrency-safe state management.

## 2. Architecture & Communication

```text
                   proto/person.proto
                    (shared contract)
                      /         \
                     v           v
            Python generated   Go generated
               stubs/types     stubs/types
                 |                 |
         Python client/server  Go client/server
                 \                 /
                  \   gRPC/HTTP2  /
                   +-------------+

       Client → generated stub → network → server handler
                                            |
                                    in-memory records
                                            |
       Client ← protobuf response ← server handler
```

The same RPC contract supports four client/server combinations:

| Client | Server | Purpose |
|---|---|---|
| Python | Python | Baseline CRUD workflow |
| Go | Go | Equivalent Go implementation |
| Python | Go | Cross-language interoperability |
| Go | Python | Reverse cross-language interoperability |

**Local transport:** `localhost:50051` with plaintext gRPC for development. Production deployment would require suitable TLS credentials, access control, and durable storage.

## 3. API Contract & CRUD Operations

The shared `proto/person.proto` defines `PersonService` and the `Person` model (`id`, `name`, `email`, `phone`). Each RPC is **unary**, meaning one request yields one response.

| RPC | Request | Response | Behavior |
|---|---|---|---|
| `CreatePerson` | Name, email, phone | `Person` | Generates an ID and stores a record |
| `GetPerson` | Person ID | `Person` | Retrieves an existing record |
| `UpdatePerson` | ID and replacement fields | `Person` | Replaces an existing record |
| `DeletePerson` | Person ID | `DeletePersonResponse` | Removes a record and returns confirmation |

Example service definition:

```protobuf
syntax = "proto3";
package person;
option go_package = "example.com/grpc-person-crud/go/generated;personpb";

service PersonService {
  rpc CreatePerson(CreatePersonRequest) returns (Person);
  rpc GetPerson(GetPersonRequest) returns (Person);
  rpc UpdatePerson(UpdatePersonRequest) returns (Person);
  rpc DeletePerson(DeletePersonRequest) returns (DeletePersonResponse);
}
```

The Python compiler produces `person_pb2.py` and `person_pb2_grpc.py`; the Go compiler plugins produce `person.pb.go` and `person_grpc.pb.go`. These generated types, client stubs, and server interfaces keep both implementations aligned with the same wire contract.

## 4. Repository Layout

```text
grpc-person-crud/
├── proto/
│   └── person.proto           # Language-neutral service contract
├── python/
│   ├── server.py              # Python CRUD handlers
│   ├── client.py              # Interactive Python client
│   ├── person_pb2.py          # Generated protobuf types
│   ├── person_pb2_grpc.py     # Generated stub/servicer
│   └── requirements.txt
├── go/
│   ├── server/main.go         # Go CRUD handlers
│   ├── client/main.go         # Interactive Go client
│   ├── generated/
│   │   ├── person.pb.go       # Generated protobuf types
│   │   └── person_grpc.pb.go  # Generated gRPC interfaces
│   ├── go.mod
│   └── go.sum
└── README.md
```

> Generated files are compiler outputs, not hand-written API implementations. The `go_package` import path must match the module path in `go/go.mod`.

## 5. Getting Started

**Prerequisites:** Python 3, Go, `protoc`, and the Go plugins `protoc-gen-go` and `protoc-gen-go-grpc` available on your `PATH`.

**Run the Python version** (from the repository root):

```bash
cd python
python3 -m venv .venv
source .venv/bin/activate             # macOS/Linux
python -m pip install -r requirements.txt
python -m grpc_tools.protoc -I../proto --python_out=. --grpc_python_out=. ../proto/person.proto
python server.py                      # Terminal 1
```

Open **Terminal 2** from the `python/` directory with the virtual environment activated:

```bash
python client.py
```

**Run the Go version** (from the repository root; stop the Python server first):

```bash
protoc -I proto \
  --go_out=go/generated --go_opt=paths=source_relative \
  --go-grpc_out=go/generated --go-grpc_opt=paths=source_relative \
  proto/person.proto
cd go
go mod tidy
go run ./server                      # Terminal 1
```

In **Terminal 2**, from `go/`:

```bash
go run ./client
```

**Important:** Both servers use port `50051`, so only start one server at a time. Both clients can connect to either server using the same contract. For cross-language testing, start the Go server and Python client, or the Python server and Go client.

## 6. Example Workflow & Error Handling

Use either interactive client to exercise a complete lifecycle:

```text
CreatePerson(name="John", email="john@example.com", phone="1234567890")
    → Person(id=1, name="John", ...)

GetPerson(id=1)
    → Person(id=1, name="John", ...)

UpdatePerson(id=1, name="John Smith", ...)
    → Person(id=1, name="John Smith", ...)

DeletePerson(id=1)
    → "Person 1 deleted successfully"

GetPerson(id=1)
    → NOT_FOUND
```

**Error semantics:** The server returns `NOT_FOUND` for unknown IDs and `INVALID_ARGUMENT` for missing required fields such as name or email. Clients handle these as gRPC status errors rather than treating them as successful business responses.

**Validation note:** The sample implementation uses basic non-empty checks, not comprehensive email or phone-number validation. Update is a full-record replacement, not a partial field patch.

## 7. Engineering Decisions & Trade-offs

| Decision | Why it matters |
|---|---|
| **Shared `proto3` definition** | Maintains one language-neutral API across Python and Go |
| **Generated stubs and server interfaces** | Avoids manually writing serialization and low-level HTTP request handling |
| **Unary gRPC methods** | Fits straightforward request/response CRUD operations |
| **Python dictionary / Go map** | Keeps persistence out of scope so RPC mechanics remain visible |
| **Go `sync.Mutex`** | Protects the shared map and ID counter from concurrent RPC access |
| **gRPC status codes** | Makes missing records and invalid requests explicit across languages |
| **Per-request Go contexts** | Lets the client set a deadline and cancel slow RPCs |

**Limitations by design:** Data disappears when a server restarts; Python and Go servers hold **separate** in-memory stores; there is no database synchronization, authentication, or production transport security. This is a learning implementation, not a production-ready customer directory.

## 8. Validation & Interview Discussion

**Recommended verification matrix** (run these tests locally; this README does not claim automated test results):

- Complete Create → Read → Update → Delete → Read (`NOT_FOUND`) in both language implementations.
- Confirm `NOT_FOUND` for an unknown ID and `INVALID_ARGUMENT` for missing required fields.
- Test **Python client → Go server** and **Go client → Python server** without modifying the contract.
- Restart a server to verify records are not persisted.
- Try multiple Go requests to understand why the in-memory map is protected by a mutex.

**Interview talking points:**

- **Why gRPC?** A contract-first RPC interface with generated clients and compact Protocol Buffer messages, suitable for service-to-service communication.
- **What enables interoperability?** Both languages implement the same service and message schema on the wire, rather than sharing source code.
- **What does the stub do?** Exposes remote RPCs through language-native methods while the gRPC runtime handles serialization and transport.
- **What would change for production?** Persistent storage, TLS/authentication, stronger validation, structured logging, automated tests, deadlines/retries, and graceful shutdown.

**Project summary:** Implemented a contract-first Person CRUD API in Python and Go using gRPC and Protocol Buffers, with generated client/server interfaces, unary RPCs, in-memory storage, and explicit error handling. The project is designed to demonstrate cross-language service interoperability through a shared `.proto` contract.
