# gen_sql_executor

**Type-safe SQL execution wrapper over [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx) package.**
**Provides a thin layer for executing queries via [`session_manager.Executor`](../session_manager),**
**with built-in generic result mapping and unified `sql.ErrNoRows` handling via `context.Context`.**

Since it operates on the `session_manager.Executor` interface rather than a raw `*sqlx.Tx`,
the same query functions work both inside and outside an active transaction —
`Session.Executor()` transparently falls back to `*sqlx.DB` when no transaction was started.


## Usage examples

### Basic Exec
```go
package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

func DeleteOrder(ctx context.Context, e session_manager.Executor, orderID int) error {
	return gen_sql_executor.Exec(ctx, e, "DELETE FROM orders WHERE id = $1", orderID)
}

func main() {
	pgdb, err := sqlx.Connect("pgx", "pg_example:5432")
	if err != nil {
		log.Fatalln("could not connect to postgres database: ", err)
	}
	defer pgdb.Close()

	sessionManager := session_manager.NewManager(pgdb)
	sessionController := sessionManager.CreateController()

	if err := sessionController.Start(); err != nil {
		log.Fatalln("could not start transaction")
	}
	defer sessionController.Rollback()

	sess := sessionController.GetSession()

	if err := DeleteOrder(context.Background(), sess.Executor(), 42); err != nil {
		log.Fatalln("could not delete order: ", err)
	}

	if err := sessionController.Commit(); err != nil {
		log.Fatalln("could not commit transaction")
	}
}
```

### Get single row
```go
package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

type Order struct {
	ID     int    `db:"id"`
	Status string `db:"status"`
}

func GetOrder(ctx context.Context, e session_manager.Executor, orderID int) (Order, error) {
	return gen_sql_executor.Get[Order](ctx, e, "SELECT id, status FROM orders WHERE id = $1", orderID)
}

func main() {
	pgdb, err := sqlx.Connect("pgx", "pg_example:5432")
	if err != nil {
		log.Fatalln("could not connect to postgres database: ", err)
	}
	defer pgdb.Close()

	sessionManager := session_manager.NewManager(pgdb)
	sessionController := sessionManager.CreateController()

	if err := sessionController.Start(); err != nil {
		log.Fatalln("could not start transaction")
	}
	defer sessionController.Rollback()

	sess := sessionController.GetSession()

	// override default ErrNoData error for this query
	ctx := gen_sql_executor.SetCustomErrNoDataCtx(context.Background(), ErrOrderNotFound)

	order, err := GetOrder(ctx, sess.Executor(), 42)
	if err != nil {
		log.Fatalln("could not get order: ", err)
	}

	log.Printf("order: %+v", order)
}
```

### Select multiple rows
```go
package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

type Order struct {
	ID     int    `db:"id"`
	Status string `db:"status"`
}

func GetOrdersByStatus(ctx context.Context, e session_manager.Executor, status string) ([]Order, error) {
	return gen_sql_executor.Select[Order](ctx, e, "SELECT id, status FROM orders WHERE status = $1", status)
}

func main() {
	pgdb, err := sqlx.Connect("pgx", "pg_example:5432")
	if err != nil {
		log.Fatalln("could not connect to postgres database: ", err)
	}
	defer pgdb.Close()

	sessionManager := session_manager.NewManager(pgdb)
	sessionController := sessionManager.CreateController()

	if err := sessionController.Start(); err != nil {
		log.Fatalln("could not start transaction")
	}
	defer sessionController.Rollback()

	sess := sessionController.GetSession()

	orders, err := GetOrdersByStatus(context.Background(), sess.Executor(), "pending")
	if err != nil {
		log.Fatalln("could not get orders: ", err)
	}

	log.Printf("orders: %+v", orders)
}
```

### Insert with returning
```go
package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

type CreateOrderParams struct {
	Status string `db:"status"`
	Amount int    `db:"amount"`
}

func CreateOrder(ctx context.Context, e session_manager.Executor, params CreateOrderParams) (int, error) {
	return gen_sql_executor.ExecNamedReturningState[int](
		ctx, e,
		"INSERT INTO orders (status, amount) VALUES (:status, :amount) RETURNING id",
		params,
	)
}

func main() {
	pgdb, err := sqlx.Connect("pgx", "pg_example:5432")
	if err != nil {
		log.Fatalln("could not connect to postgres database: ", err)
	}
	defer pgdb.Close()

	sessionManager := session_manager.NewManager(pgdb)
	sessionController := sessionManager.CreateController()

	if err := sessionController.Start(); err != nil {
		log.Fatalln("could not start transaction")
	}
	defer sessionController.Rollback()

	sess := sessionController.GetSession()

	id, err := CreateOrder(context.Background(), sess.Executor(), CreateOrderParams{
		Status: "pending",
		Amount: 100,
	})
	if err != nil {
		log.Fatalln("could not create order: ", err)
	}

	log.Printf("created order id: %d", id)

	if err := sessionController.Commit(); err != nil {
		log.Fatalln("could not commit transaction")
	}
}
```

### With context
```go
package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

func DeleteOrder(ctx context.Context, orderID int) error {
	sess := session_manager.MustGetSession(ctx)
	return gen_sql_executor.Exec(ctx, sess.Executor(), "DELETE FROM orders WHERE id = $1", orderID)
}

func main() {
	pgdb, err := sqlx.Connect("pgx", "pg_example:5432")
	if err != nil {
		log.Fatalln("could not connect to postgres database: ", err)
	}
	defer pgdb.Close()

	sessionManager := session_manager.NewManager(pgdb)
	sessionController := sessionManager.CreateController()

	if err := sessionController.Start(); err != nil {
		log.Fatalln("could not start transaction")
	}
	defer sessionController.Rollback()

	ctx := session_manager.SetSession(context.Background(), sessionController.GetSession())

	if err := DeleteOrder(ctx, 42); err != nil {
		log.Fatalln("could not delete order: ", err)
	}

	if err := sessionController.Commit(); err != nil {
		log.Fatalln("could not commit transaction")
	}
}
```


## Error handling

By default, when a `Get`, `Select`, or `SelectNamed` query returns no rows, these functions return `gen_sql_executor.ErrNoData`.
Use `SetCustomErrNoDataCtx` to override this with a domain-specific error (e.g. `ErrOrderNotFound`) for a given call:

```go
ctx := gen_sql_executor.SetCustomErrNoDataCtx(ctx, ErrOrderNotFound)
order, err := gen_sql_executor.Get[Order](ctx, e, "SELECT ...", orderID)
// err is ErrOrderNotFound if no rows were found
```