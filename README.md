# mdbx sql

MDBX-SQL is a prototype/idea project with goal of enabling SQL capabilities with the high-performance Memory-Mapped Database (MDBX). 
By leveraging the strengths of both SQL and MDBX, 
this project offers developers a versatile and powerful tool for data management and querying. 
MDBX-SQL is specifically optimized for near zero-memory allocation and data analytical workloads, 
ensuring exceptional performance and efficiency in handling complex data operations.

^ Well this is what I want to build, now lets see if this is possible and how to achieve it.

## Playground results...

Used json parser. This is for sure gonna go out... It's ABSOLUTELY forbidden to use json but heck, tested it now...

```
go clean -testcache && go test ./... -v -cover
	github.com/0x19/mdbx-sql/parser		coverage: 0.0% of statements
=== RUN   TestParserAndDatabase
2024/07/07 22:47:22 SQL Parsing completed in 115.713µs
2024/07/07 22:47:22 AST: &{Columns:[name age] TableName:users Condition:active}
2024/07/07 22:47:22 Insert operation completed in 212.396µs
2024/07/07 22:47:22 Get operation completed in 171.055µs
2024/07/07 22:47:22 Retrieved User: {ID:1 Name:John Doe Age:30}
2024/07/07 22:47:22 Update operation completed in 2.658472ms
2024/07/07 22:47:22 Get operation (post-update) completed in 8.02µs
2024/07/07 22:47:22 Updated User: {ID:1 Name:John Doe Age:31}
2024/07/07 22:47:22 Delete operation completed in 3.389402ms
2024/07/07 22:47:22 Get operation (post-delete) completed in 4.98µs
--- PASS: TestParserAndDatabase (0.02s)
=== RUN   TestPlayground
AST: &{Columns:[name age] TableName:users Condition:active} in 42.901µs 
--- PASS: TestPlayground (0.00s)
PASS
coverage: 85.5% of statements
ok  	github.com/0x19/mdbx-sql	0.029s	coverage: 85.5% of statements
```

Using capnp

```
=== RUN   TestParserAndDatabase
2024/07/08 08:22:09 SQL Parsing completed in 68.572µs
2024/07/08 08:22:09 AST: &{Columns:[name age] TableName:users Condition:active}
2024/07/08 08:22:09 Insert operation completed in 4.991858ms
2024/07/08 08:22:09 Get operation completed in 117.313µs
2024/07/08 08:22:09 Retrieved User: &{ID:1 Name:John Doe Age:30}
2024/07/08 08:22:09 Update operation completed in 3.316886ms
2024/07/08 08:22:09 Get operation (post-update) completed in 18.631µs
2024/07/08 08:22:09 Updated User: &{ID:1 Name:John Doe Age:31}
2024/07/08 08:22:09 Delete operation completed in 3.320577ms
2024/07/08 08:22:09 Get operation (post-delete) completed in 4.79µs
--- PASS: TestParserAndDatabase (0.02s)
=== RUN   TestPlayground
AST: &{Columns:[name age] TableName:users Condition:active} in 2.07µs 
--- PASS: TestPlayground (0.00s)
=== RUN   TestSchemaCreateAndGetTable
--- PASS: TestSchemaCreateAndGetTable (0.00s)
PASS
coverage: 76.7% of statements
ok  	github.com/0x19/mdbx-sql	0.028s	coverage: 76.7% of statements
```

Well I've screwed something up as json cannot be faster by any means then capnp....