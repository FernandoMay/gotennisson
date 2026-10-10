# gotennisson — API de pedidos (Gin)

API REST minimalista en Go con [Gin](https://gin-gonic.com/). Sirve el catálogo de una tienda de tenis y calcula totales de pedidos.

## Estado del revival (2026-10)

- **🐛 Bug fixeado**: `POST /orders` usaba `c.BindJSON` **ignorando el error** — un payload roto devolvía `200` con un pedido vacío y total 0. Ahora responde `400` con el detalle del error.
- **🧪 12 aserciones de test** (`go test`, sin dependencias extra): tabla de casos del cálculo (vacío, simple, duplicados, inexistentes, mixtos) y los endpoints con `httptest` — incluyendo el JSON inválido que documenta el bug.
- **🏗️ Handler separable**: `newRouter()` en `handler.go`, `main()` de una línea. Los tests ejercitan el router in-process.
- **❤️ /health** para healthchecks.
- **📦 Dockerfile** multi-stage (binario estático ~12MB).

## Requisitos

- Go 1.21+

## Uso

```bash
go run .          # escucha en :8080
go test ./...     # tests
go build -o gotennisson .
```

## API

### `GET /health`

```json
{"status": "ok"}
```

### `GET /products`

```json
[
  {"id": 1, "name": "Zapatillas Nike Air Max", "price": 14999},
  {"id": 2, "name": "Raqueta de tenis Wilson", "price": 9999},
  {"id": 3, "name": "Pelotas de tenis Penn", "price": 499}
]
```

### `POST /orders`

Body — los productos se identifican **por nombre** (contrato original preservado):

```json
{"products": ["Zapatillas Nike Air Max", "Pelotas de tenis Penn"]}
```

Respuesta:

```json
{
  "message": "Order created successfully",
  "products": ["Zapatillas Nike Air Max", "Pelotas de tenis Penn"],
  "total": 15498
}
```

- Productos inexistentes se ignoran
- Duplicados suman por cada aparición
- JSON inválido → `400`
- Pedido vacío → `200` con `total: 0`

## Docker

```bash
docker build -t gotennisson .
docker run -p 8080:8080 gotennisson
```

## Tests

```bash
go test ./... -v
```

Cubren: `calculateOrder` (6 casos de tabla) y endpoints HTTP (`/products`, `/orders` válido, `/orders` inválido, lista vacía, `/health`).

## Licencia

MIT — ver [LICENSE](LICENSE).

## Autor

[FernandoMay](https://github.com/FernandoMay)
