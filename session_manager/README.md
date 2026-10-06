# session_manager

**Transaction lifecycle management wrapper over [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx) package.**
**Provides a thin layer for starting, committing, and rolling back transactions,**
**with support for inject an active transaction into `context.Context`.**


## Usage examples

### Basic
```go
package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

func UpdateData(ctx context.Context, e session_manager.Executor, data int) error {
	return gen_sql_executor.Exec(ctx, e, "UPDATE row SET row.data = $1", data)
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
		return
	}
	defer sessionController.Rollback()

	sess := sessionController.GetSession()

	if err := UpdateData(context.Background(), sess.Executor(), 10); err != nil {
		log.Fatalln("could not update data")
		return
	}

	if err := sessionController.Commit(); err != nil {
		log.Fatalln("could not commit transaction")
		return
	}
}
```

### With context
```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

func UpdateData(ctx context.Context, data int) error {
	sess := session_manager.MustGetSession(ctx)

	return gen_sql_executor.Exec(ctx, sess.Executor(), "UPDATE row SET row.data = $1", data)
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	pgdb, err := sqlx.Connect("pgx", "pg_example:5432")
	if err != nil {
		log.Fatalln("could not connect to postgres database: ", err)
	}
	defer pgdb.Close()

	sessionManager := session_manager.NewManager(pgdb)
	sessionController := sessionManager.CreateController()

	if err := sessionController.Start(); err != nil {
		log.Fatalln("could not start transaction")
		return
	}
	defer sessionController.Rollback()

	ctx = session_manager.SetSession(ctx, sessionController.GetSession())

	if err := UpdateData(ctx, 10); err != nil {
		log.Fatalln("could not update data")
		return
	}

	if err := sessionController.Commit(); err != nil {
		log.Fatalln("could not commit transaction")
		return
	}
}
```

### Without an active transaction
Since `Session.Executor()` falls back to the underlying `*sqlx.DB` when no transaction is active,
the same repository/query code works both inside and outside a transaction — no branching needed on the caller's side.
```go
package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab..uz/crm/sqlutils/gen_sql_executor"
	"gitlab..uz/crm/sqlutils/session_manager"
)

func UpdateData(ctx context.Context, e session_manager.Executor, data int) error {
	return gen_sql_executor.Exec(ctx, e, "UPDATE row SET row.data = $1", data)
}

func main() {
	pgdb, err := sqlx.Connect("pgx", "pg_example:5432")
	if err != nil {
		log.Fatalln("could not connect to postgres database: ", err)
	}
	defer pgdb.Close()

	sessionManager := session_manager.NewManager(pgdb)
	sessionController := sessionManager.CreateController()

	// no Start() call — GetSession() falls back to *sqlx.DB
	sess := sessionController.GetSession()

	if err := UpdateData(context.Background(), sess.Executor(), 10); err != nil {
		log.Fatalln("could not update data")
	}
}
```