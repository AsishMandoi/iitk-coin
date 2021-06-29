# IITK Coin

- ### Subpackages
  - My package is split into multiple sub-packages (i.e. I have made a few sub-directories - `global`, `handlers`, `server` and `database`).
  - <details>
      <summary>Tree Directory Structure</summary>
      
      ```
      iitk-coin
      ├── database
      │   └── commonOps.go
      │   └── init.go
      │   └── txnOps.go
      ├── global
      │   └── globalObjects.go
      │   └── init.go
      ├── handlers
      │   ├── balance.go
      │   ├── loginpage.go
      │   ├── reward.go
      │   ├── secretpage.go
      │   └── signuppage.go
      │   └── transfer.go
      ├── server
      │   ├── jwt.go
      │   └── respond.go
      ├── .env
      ├── .gitignore
      ├── go.mod
      ├── go.sum
      ├── iitkusers.db
      ├── iitkusers.db-shm
      ├── iitkusers.db-wal
      ├── main.go
      └── README.md
      ```
    </details>

- ### Write-Ahead Log
  - The `journal_mode` is set to `WAL` because of its [advantages](https://sqlite.org/wal.html#overview) over the default, `DELETE` mode in SQLite.
  - I personally tested in both modes and observed that the `WAL` mode works slightly faster (upto 10x faster) than the default mode while processing **parallelly requested** write operations into the database.
  - I also tested both the modes (again using parallel curl commands) intentionally keeping the DB locked for a certain time. In the default mode the concurrent requests are bound to be unsuccessful with an `database is locked` error. But, in `WAL` mode requests are handled sequentially and automatically once the db gets unlocked.

- ### Testing
  I have used this script - http://p.ip.fi/8SxQ to test the endpoints.

- ### Request Format
  - ##### `/signup` page:
    <details>
      <summary>Click to view</summary>
      
      ```http
      POST /signup HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json

      {
        "rollno":   <Your_Rollno>,
        "name":     "<Your_Name>",
        "password": "<Your_Password>",
        "batch":    "<Your_Batch>"
      }
      ```
    </details>
  - ##### `/login` page:
    <details>
      <summary>Click to view</summary>
      
      ```http
      POST /login HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json

      {
        "rollno":   <Your_Rollno>,
        "password": "<Your_Password>"
      }
      ```
    </details>
  - ##### `/secretpage`:
    <details>
      <summary>Click to view</summary>
      
      ```http
      GET /secretpage HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer <Token>
      ```
    </details>
  - ##### `/view_coins`:
    <details>
      <summary>Click to view</summary>
      
      ```http
      GET /view_coins HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer <Token>
      ```
    </details>
  - ##### `/transfer_coins` page:
    <details>
      <summary>Click to view</summary>
      
      ```http
      POST /transfer_coins HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer <Token>

      {
        "receiver":    <Receiver_Rollno>,
        "amount":      <Amount>,
        "description": <Remarks>
      }
      ```
    </details>
  - ##### `/reward_coins` page:
    <details>
      <summary>Click to view</summary>
      
      ```http
      POST /reward_coins HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer <Token>

      {
        "receiver":    <Receiver_Rollno>,
        "amount":      <Amount>,
        "description": <Remarks>
      }
      ```
    </details>
  - ##### `/redeem_coins` page:
    <details>
      <summary>Click to view</summary>
      
      ```
      Not Implemented yet
      ```
    </details>
  Sample requests for each endpoint are specified at the beginning of each of the handler functions.

- ### Response Format
  For the response format for each endpoint please look at the `global` package, where some global objects (variables and struct types) are defined, to be used in other functions.

- ### Common Response Method
  Although the endpoints have slightly different formats for their response object, all of them are handled using a `type-switch` in a common `server.Respond()` function which responds to requests for all the endpoints. This method has been used a lot of times in various files. It has greatly reduced the bulkiness of codes in individual files.

- ### HTTP Status Codes
  A suitable http status code is assigned to every response.

- ### Database
  The `init()` function of the `database` package automatically initializes the database. The initialization errors are handled before making any other database operations. Currently `sqlite` is used as the database management system. If in the future I wish to switch to any other SQL based database management system, I will just have to change one line of code in the `database` package, and import the corresponding package required for it.

- ### .env
  - The `.env` file contains the `Secret Key` to sign the JWT, the variable `Maximum Cap` for the coins and the variable `Minimum Events` which is a for users to be eligible for transactions. It is deliberately left out of `.gitignore` for the purposes of checking.
  - If an `.env` file is not found the defult values of these environment variables will be used throughout.
  - My intention was to make it convenient for one who is running the backend to be able to update these varibles in the `.env` file without having to search for them in the code. And I have made it so that, if these environment variables are updated here these values will be overwritten to the variables defined inside the code.

- ### Access Token
  Expiry Time is currently set to 30 minutes.

- ### Refresh Token
  Not implemented yet.

- ### Cap for Maximum Coins
  Currently set to `1001` coins.

- ### Minimum Events
  Currently set to `6`.

---
### For my reference:
- [x] look up popular directory structures
- [x] send a json object in the response for every endpoint
- [x] use MDN: HTTP status codes -> http.StatusBadRequest, http.StatusUnauthorized, ...
- [x] batch, txn depends on batch
- [ ] use other modes of transaction - `IMMEDIATE`, `EXCLUSIVE`
- [ ] use refresh token
- [ ] make redeem_coins endpt.
- [ ] check isAdmin from token and then authorize to /secretPage
- [ ] a new table `auth` for storing just rollnos and passwords?
- [ ] keep checking for unhandled errors

A common approach for invalidating tokens when a user changes their password is to sign the token with a hash of their password. Thus if the password changes, any previous tokens automatically fail to verify. You can extend this to logout by including a last-logout-time in the user's record and using a combination of the last-logout-time and password hash to sign the token. This requires a DB lookup each time you need to verify the token signature, but presumably you're looking up the user anyway. – [Travis Terry](https://stackoverflow.com/questions/21978658/invalidating-json-web-tokens/23089839#comment45057142_23089839)

Turn on the Write-Ahead Logging, Disable connections pool --[link1](https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance/35805826), [link2](https://sqlite.org/wal.html)

Once Commit or Rollback is called on the transaction, that transaction's connection is returned to DB's idle connection pool. The pool size can be controlled with SetMaxIdleConns. --[link](https://golang.org/pkg/database/sql/#DB)

---

> Two concurrent executions can interleave such that your read values become stale.

Solutions:
1. Do the read, write and validation checks in a single sql statement which is of write nature (so that it acquires lock).
2. Use other modes of transaction - `IMMEDIATE`, `EXCLUSIVE`, (more specific errors can be handled)

**In which line is the DB actually locked in the default (`DEFERRED`) mode?**

  An sqlite DB is locked after one of the write statements (e.g. `UPDATE`, `INSERT`, ...) are executed, irrespective of whether they are in a transaction. This is the default behaviour.
