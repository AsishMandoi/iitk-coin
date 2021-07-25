# IITK Coin

- ### Subpackages
  - My package is split into multiple sub-packages (i.e. I have made a few sub-directories - `global`, `handlers`, `server` and `database`).
  - <details>
      <summary>Tree Directory Structure</summary>
      
      ```
      iitk-coin
      ├── database
      │   └── init.go
      │   └── redis.go
      │   └── sqlite_txn_ops.go
      │   └── sqlite.go
      ├── global
      │   └── global_objs.go
      │   └── init.go
      ├── handlers
      │   ├── balance.go
      │   ├── confirm.go
      │   ├── logine.go
      │   ├── redeem.go
      │   ├── reward.go
      │   ├── secret_page.go
      │   └── signup.go
      │   └── transfer.go
      ├── server
      │   ├── jwt.go
      │   ├── otp.go
      │   └── respond.go
      ├── .env
      ├── .env.dev
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
  I have used this script - http://p.ip.fi/Kb_e to test the endpoints.

- ### Request-Response (examples)
  - ##### `/signup`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      POST /signup HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json

      {
        "rollno":   192197,
        "name":     "Anonymous3",
        "iitk_email": "devtest.asish@gmail.com",
        "password": "Str0NgP@$5w0rD",
        "batch": "Y19"
      }

      ```
      Response body:

      ```
      {
        "message": "Added user successfully",
        "error": null
      }
      ```
    </details>
  - ##### `/login`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      POST /login HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json

      {
        "rollno": 192197,
        "password": "Str0NgP@$5w0rD"
      }
      ```
      Response body:
      ```
      {
        "message": "Login successful; Token generated successfully",
        "error": null,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwNTUwMTUsInJvbGUiOiIiLCJyb2xsbm8iOjE5MjE5N30.kG52objNZ8sj1Ba1Ogs1JYG0W6xPGZ9sFelAofdo0qU"
      }
      ```
    </details>
  - ##### `/secretpage`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      GET /secret_page HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImV4cCI6MTYyNDkzMzIxNywicm9sZSI6IiIsInJvbGxubyI6MTkwMTk3fQ.86Iyllo03FGqxvpq1iQCl3Yqs1P3jq_mXlY4O-8F2wI
      ```
      Response body:

      ```
      {
        "message": "SUCCESS",
        "error": null,
        "data": 192197
      }
      ```
    </details>
  - ##### `/view_coins`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      GET /view_coins HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImV4cCI6MTYyNjM1NDk0OSwicm9sZSI6IiIsInJvbGxubyI6MTkwMTk3fQ.E52q8iJw1_m5mxwRZADcbNF6B5srbP0iM97f2tWg-ao
      ```
      Response body:

      ```
      {
        "message": "SUCCESS",
        "error": null,
        "coins": 100
      }
      ```
    </details>
  - ##### `/transfer_coins`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      POST /transfer_coins HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOCIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwNDQ4NzUsInJvbGUiOiJBZG1pbiIsInJvbGxubyI6MTgxMTk3fQ.moLUYSlffF3EPxTxI_6k5ePneLhGHzOnB5UmB9IbsQQ

      {
        "receiver": 192197,
        "amount": 100,
        "description": "testing for an eligible sender"
      }
      ```
      Response body:

      ```
      {
        "message": "Post your otp on http://localhost:8080/confirm_transfer to confirm your transaction",
        "error": null
      }
      ```
    </details>
  - ##### `/confirm_transfer`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      POST /confirm_transfer HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwMzU3MjEsInJvbGUiOiIiLCJyb2xsbm8iOjE5MTE5N30.f1vXV40Xb1kgEQQaLGYAymGPzwqBiKHpue7eHmHqZlQ

      {
        "otp": "612765",
        "resend": false
      }
      ```
      Response body:

      ```
      {
        "message": "Transaction Successful: User: #191197 transferred 98 coins to user: #192197",
        "error": null,
        "transaction_id": 1529
      }
      ```
    </details>
  - ##### `/reward_coins`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      POST /reward_coins HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwNDQ4MDcsInJvbGUiOiIiLCJyb2xsbm8iOjE5MTE5N30.HjwFS35GEVe4k0jz7mLwrJOyTM51hQZTyJmeJHvwTzo

      {
        "receiver": 190197,
        "amount": 100,
        "description": "Testing for admin"
      }
      ```
      Response body:

      ```
      {
        "message": "Reward Successful; User: #190197 was rewarded with 200 coins",
        "error": null,
        "transaction_id": 14
      }
      ```
    </details>
  - ##### `/redeem`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      POST /redeem HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwNDU1OTksInJvbGUiOiIiLCJyb2xsbm8iOjE5MTE5N30.4Fu80f4fWcdQwtxR1Ps4s5LPwqbD_dPeHucihz7Yi_A

      {
        "item_id": 91051,
        "amount": 100,
        "description": "Testing an eligible sender."
      }
      ```
      Response body:

      ```
      {
        "message": "Post your otp to http://localhost:8080/confirm_redeem_request to confirm your transaction",
        "error": null
      }
      ```
    </details>
  - ##### `/confirm_redeem_request`:
    <details>
      <summary>Click to view</summary>
      Request:
      
      ```http
      POST /confirm_redeem_request HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwNDU1OTksInJvbGUiOiIiLCJyb2xsbm8iOjE5MTE5N30.4Fu80f4fWcdQwtxR1Ps4s5LPwqbD_dPeHucihz7Yi_A

      {
        "otp": "273801",
        "resend": false
      }
      ```
      Response body:

      ```
      {
        "message": "Redeem request successful",
        "error": null,
        "request_id": 3
      }
      ```
    </details>
  - ##### `/redeem_requests`:
    <details>
      <summary>Click to view</summary>
      
      ```http
      GET /redeem_requests HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOCIsImV4cCI6MTYyNjQwOTM1Niwicm9sZSI6IkFkbWluIiwicm9sbG5vIjoxODExOTd9.aOwSdGSmEyaQYGhJNBAt449rcFi3fQ6JT0u6gu7Adtg
      ```
      Response body:

      ```
      {
        "message": null,
        "error": null,
        "data": [
          {
            "request_id": 2,
            "redeemer": 192197,
            "item_id": 91020,
            "amount": 30,
            "description": "Testing an eligible sender.",
            "requested_on": "2021-07-16T02:37:48Z"
          },
          {
            "request_id": 3,
            "redeemer": 192197,
            "item_id": 91021,
            "amount": 50,
            "description": "Testing an eligible sender.",
            "requested_on": "2021-07-16T03:52:25Z"
          }
        ]
      }
      ```
    </details>
  - ##### `/update_redeem_status`:
    <details>
      <summary>Click to view</summary>
      
      ```http
      POST /update_redeem_status HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOCIsImV4cCI6MTYyNjQxMTQ5MCwicm9sZSI6IkFkbWluIiwicm9sbG5vIjoxODExOTd9.mrNBbfpwp9GjNKb2G0OgNbKNX8kdoJbafidMFof3sd0

      {
        "request_id": 3,
        "user": 192197,
        "coins": 50,
        "status": "Accept",
        "description": "Testing an Admin"
      }
      ```
      Response body:

      ```
      {
        "message": "Redeem updated successfully",
        "error": null,
        "transaction_id": 3
      }
      ```
    </details>
  - ##### `/redeem_status`:
    <details>
      <summary>Click to view</summary>
      
      ```http
      GET /redeem_requests HTTP/1.1
      HOST: localhost:8080
      Content-Type: application/json
      Accept: application/json
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOCIsImV4cCI6MTYyNjQwOTM1Niwicm9sZSI6IkFkbWluIiwicm9sbG5vIjoxODExOTd9.aOwSdGSmEyaQYGhJNBAt449rcFi3fQ6JT0u6gu7Adtg
      ```
      Response body:

      ```
      {
        "message": null,
        "error": null,
        "data": [
          {
            "id": 1,
            "item_id": 91019,
            "amount": 10,
            "description": "Testing an Admin",
            "status": "Accepted",
            "requested_on": "2021-07-16T00:51:47Z",
            "responded_on": "2021-07-16T03:53:28Z"
          },
          {
            "id": 2,
            "item_id": 91020,
            "amount": 30,
            "description": "Testing an eligible sender.",
            "status": "Pending",
            "requested_on": "2021-07-16T02:37:48Z",
            "responded_on": "0001-01-01T00:00:00Z"
          },
          {
            "id": 3,
            "item_id": 91021,
            "amount": 50,
            "description": "Testing an Admin",
            "status": "Accepted",
            "requested_on": "2021-07-16T03:52:25Z",
            "responded_on": "2021-07-16T04:29:16Z"
          },
          {
            "id": 4,
            "item_id": 91021,
            "amount": 50,
            "description": "Testing an eligible sender.",
            "status": "Pending",
            "requested_on": "2021-07-16T04:26:43Z",
            "responded_on": "0001-01-01T00:00:00Z"
          }
        ]
      }
      ```
    </details>

- ### Common Response Method
  Although the endpoints have slightly different formats for their response object, all of them are handled using a `type-switch` in a common `server.Respond()` function which responds to requests for all the endpoints. This method has been used a lot of times in various files. It has greatly reduced the bulkiness of codes in individual files.

- ### HTTP Status Codes
  A suitable http status code is assigned to every response.

- ### Database
  I have used two database management systems in this application, `SQLite` and `Redis`. The `init()` function of the `database` package automatically initializes the databases. The initialization errors are handled before making any other database operations.

- ### .env
  - The `.env` file contains the `Secret Key` to sign the JWT, the variable `Maximum Cap` for the coins and the variable `Minimum Events` which is a for users to be eligible for transactions. It is deliberately left out of `.gitignore` for the purposes of checking.
  - If an `.env` file is not found the defult values of these environment variables will be used throughout.
  - My intention was to make it convenient for one who is running the backend to be able to update these varibles in the `.env` file without having to search for them in the code. And I have made it so that, if these environment variables are updated here these values will be overwritten to the variables defined inside the code.
  - There is also a file named `.env.dev` that contains a the gmailid and password for a test account from which all OTPs are sent.

- ### Access Token
  Expiry Time is currently set to 30 minutes.

- ### Refresh Token (or similar)
  To be implemented soon.

- ### Cap for Maximum Coins
  Currently set to `10001` coins.

- ### Minimum Events
  Currently set to `6`.

- ### Redeem  
  The functionalities currently available are:<br />
  - Users can send redeem requests which will be in pending state by default. This is present in the `/redeem` endpoint. Once a valid request is made, an OTP is send to the user's emailid (that was collected during signup).
  - An Admin can see a list of all pending redeem requests on the `/redeem_requests` endpoint
  - Users can see the status of all their requests on the `/redeem_status` endpoint
  - An Admin can "Accept" or "Reject" a redeem request on the `/update_redeem_status`endpoint

- ### OTP
  - `OTP` based confirmation systems are implemented on the `/redeem` and the `/transfer_coins`. The respective OTPs will have to be POSTed on the endpoints `/confirm_redeem_request` and `/confirm_transfer`.
  - There is also a `Resend OTP` option available (only at the confirmation endpoints). If a user wants to get another OTP, they have to POST a request with `resend` value set to `true`.

- ### The Process of Confirmation
  1. The user sends a request (on one of the endpoints - `/redeem` or `/transfer_coins`)
  2. If the request is invalid the server responds with the messages and errors
  3. If the request is valid -
      - An OTP is generated.
      - The OTP along with the data that needs to be stored in the required tables is temporarily saved in the `Redis` database. Expiry time is set to 2 mins currently.
      - The OTP is sent to the user's emailid.
      - If the correct OTP is not entered, the process ends with an error message unless the user sets the `resend` option to be true.
        - If the resend option is true, one can enter the OTP again with a new POST request on the same endpoint
      - If the OTP is successfully entered and there is no error while storing the data in the required tables
        - Immediately, the data along with the OTP is deleted from the `Redis` database.
      - If no request is made within this expiry time of 2 mins (not even a resend), the main data to be stored is lost
  
  *One can potentially delay the process (of transfer/redeem) if they keep on resending the OTP before the current one expires. But, this can be done until the `JWT` token expires, after which the user has to login again.*

---
### For my reference:
- [x] look up popular directory structures
- [x] send a json object in the response for every endpoint
- [x] use MDN: HTTP status codes -> http.StatusBadRequest, http.StatusUnauthorized, ...
- [x] batch, txn depends on batch
- [X] make redeem_coins endpts.
- [X] implement OTP
- [X] Use Redis
- [ ] use other modes of transaction - `IMMEDIATE`, `EXCLUSIVE`
- [ ] use refresh token/similar
- [ ] check isAdmin from token and then authorize to /secretPage
- [ ] a new table `auth` for storing just rollnos and passwords?
- [ ] keep checking for unhandled errors

---

#### Some incredibly helpful resources of many

> A common approach for invalidating tokens when a user changes their password is to sign the token with a hash of their password. Thus if the password changes, any previous tokens automatically fail to verify. You can extend this to logout by including a last-logout-time in the user's record and using a combination of the last-logout-time and password hash to sign the token. This requires a DB lookup each time you need to verify the token signature, but presumably you're looking up the user anyway. – [Travis Terry](https://stackoverflow.com/questions/21978658/invalidating-json-web-tokens/23089839#comment45057142_23089839)

> Turn on the Write-Ahead Logging, Disable connections pool – [link1](https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance/35805826), [link2](https://sqlite.org/wal.html)

> Once Commit or Rollback is called on the transaction, that transaction's connection is returned to DB's idle connection pool. The pool size can be controlled with SetMaxIdleConns. – [link](https://golang.org/pkg/database/sql/#DB)

> As a general rule of thumb, if you can use structs to represent your JSON data, you should use them. The only good reason to use maps would be if it were not possible to use structs due to the uncertain nature of the keys or values in the data. – [Soham Kamani](https://www.sohamkamani.com/golang/parsing-json/#what-to-use-structs-vs-maps)

> Two concurrent executions can interleave such that your read values become stale.\
Solutions:<br />
  1. Do the read, write and validation checks in a single sql statement which is of write nature (so that it acquires lock).
  2. Use other modes of transaction - `IMMEDIATE`, `EXCLUSIVE`, (more specific errors can be handled)         – Bhuvan Singla

**In which line is the DB actually locked in the default (`DEFERRED`) mode?**

  An sqlite DB is locked after one of the write statements (e.g. `UPDATE`, `INSERT`, ...) are executed, irrespective of whether they are in a transaction. This is the default behaviour.
