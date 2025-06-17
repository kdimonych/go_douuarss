# The migration service

__Usage__

```bash
go run main.go [-migrations <Path to migration files>] up | down
```

Where:

- __-migrations__ - path to migration files
- __up__ - performs DB migration
- __down__ - performs roll back the DB migration

> __Reuired environment variable:__ `DATABASE_URL`
