# IITK Coin
*A centralized pseudo-currency system to be used in IIT Kanpur*

- Containerized application accessible on [DockerHub](https://hub.docker.com/r/asishmandoi/iitk-coin)
- Source code accessible on [GitHub](https://github.com/AsishMandoi/iitk-coin)

## Highlights
- [X] Image size on DockerHub: `~10MB`, Total space required to run the application:`~50MB`
- [X] Suitable HTTP status codes assigned to all responses
- [X] OTP based endpoints for an added layer of security
- [X] Redis for temporary storage, fast retrieval of data
- [X] [Write-Ahead Log (WAL)](https://sqlite.org/wal.html) mode enabled in SQLite
- [X] A highly secure Access Token implementation

## Run the application locally
<!-- ### Using source code from GitHub
  **`Redis server` needs to be running locally**
  ``` bash
  export GOPATH=$(go env GOPATH)                                                          # Make sure `GOPATH` environment variable is set
  mkdir -p $GOPATH/src/github.com/AsishMandoi/ && cd $GOPATH/src/github.com/AsishMandoi/  # Make this directory and change the working directory as given

  git clone https://github.com/AsishMandoi/iitk-coin.git                                  # Clone this repository

  cd ./iitk-coin                                                                          # Change the working directory again as given

  export REDIS_CONTAINER_NAME=localhost                                                   # Required for redis

  ##### Set EMAIL_ID and PASSWORD in .env for OTPs to function #####
  # ...
  # EMAIL_ID=<enter_sender_email>
  # PASSWORD=<enter_sender_password>
  # ...

  go build -o iitk-coin-server && ./iitk-coin-server                                      # Build and run the executable binary
  ``` -->

### Using images from DockerHub *(recommended)*
  *Requires `docker` installed, no other installation required*
  ``` bash
  # Download and run the file `run-containers.sh`
  curl https://raw.githubusercontent.com/AsishMandoi/iitk-coin/main/scripts/run-containers.sh -o run-containers.sh -s && . run-containers.sh
  ```

### Using source code and docker-compose
  *Requires `docker-compose` installed, no other installation required*
  ``` bash
  # Download and run the file `run-from-source.sh`
  curl https://raw.githubusercontent.com/AsishMandoi/iitk-coin/main/scripts/run-from-source.sh -o run-from-source.sh -s && . run-from-source.sh
  ```

## Overview
- ### Subpackages
  The main package contains multiple sub-packages (i.e. I have made a few sub-directories - `global`, `handlers`, `server` and `database`).
  <details>
    <summary><b>Tree Directory Structure</b></summary>
    
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
    │   ├── auth.go
    │   ├── balance.go
    │   ├── confirm.go
    │   ├── redeem.go
    │   ├── reward.go
    │   ├── secret_page.go
    │   ├── signup.go
    │   └── transfer.go
    ├── server
    │   ├── jwt.go
    │   ├── otp.go
    │   └── respond.go
    ├── scripts
    │   ├── run-containers.sh
    │   └── run-from-source.sh
    ├── .env
    ├── .env.dev
    ├── .gitignore
    ├── .dockerignore
    ├── Dockerfile
    ├── docker-compose.yml
    ├── go.mod
    ├── go.sum
    ├── iitkusers.db
    ├── iitkusers.db-shm
    ├── iitkusers.db-wal
    ├── main.go
    └── README.md
    ```
  </details>

- ### Request-Response (examples)
  *[click to expand]*
    <details>
    <summary><b><code>Signup</code></b></summary>
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
    <details>
    <summary><b><code>Login</code></b></summary>
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
    <details>
    <summary><b><code>Reset Password (using old password)</code></b></summary>
    Request:
    
    ```http
    POST /reset_password HTTP/1.1
    HOST: localhost:8080
    Content-Type: application/json
    Accept: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2Mjc0NjUzMTMsInJvbGUiOiIiLCJyb2xsbm8iOjE5MDE5N30.HMdtutBN41UmKw9qTVE9RPSRCKfgDZK02FFyW8rFRgo

    {
      "send_otp": false,
      "old_password": "Str0NgP@$5w0rD",
      "new_password": "NewStr0NgP@$5w0rD"
    }
    ```
    Response body:
    ```
    {
      "message": "Password reset successful",
      "error": null
    }
    ```
    </details>
    <details>
    <summary><b><code>Reset Password (using OTP)</code></b></summary>
    Request(I):
    
    ```http
    POST /reset_password HTTP/1.1
    HOST: localhost:8080
    Content-Type: application/json
    Accept: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2Mjc0NjUzMTMsInJvbGUiOiIiLCJyb2xsbm8iOjE5MDE5N30.HMdtutBN41UmKw9qTVE9RPSRCKfgDZK02FFyW8rFRgo

    {
      "send_otp": true,
      "old_password": "",
      "new_password": "NewStr0NgP@$5w0rD"
    }
    ```
    Response(I) body:
    ```
    {
      "message": "Post your otp on http://localhost:8080/reset_password/confirm to confirm your transaction",
      "error": null
    }
    ```
    Request(II):
    
    ```http
    POST /reset_password/confirm HTTP/1.1
    HOST: localhost:8080
    Content-Type: application/json
    Accept: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2Mjc0NjUzMTMsInJvbGUiOiIiLCJyb2xsbm8iOjE5MDE5N30.HMdtutBN41UmKw9qTVE9RPSRCKfgDZK02FFyW8rFRgo

    {
      "otp": "554236",
      "resend": false
    }
    ```
    Response(II) body:
    ```
    {
      "message": "Password reset successful",
      "error": null
    }
    ```
    </details>
    <details>
    <summary><b><code>Access a Secret Page</code></b></summary>
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
    <details>
    <summary><b><code>Check Balance</code></b></summary>
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
    <details>
    <summary><b><code>Transfer Coins</code></b></summary>
    Request(I):
    
    ```http
    POST /transfer HTTP/1.1
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
    Response body(I):

    ```
    {
      "message": "Post your otp on http://localhost:8080/transfer/confirm to confirm your transaction",
      "error": null
    }
    ```
    Request(II):
    
    ```http
    POST /transfer/confirm HTTP/1.1
    HOST: localhost:8080
    Content-Type: application/json
    Accept: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwMzU3MjEsInJvbGUiOiIiLCJyb2xsbm8iOjE5MTE5N30.f1vXV40Xb1kgEQQaLGYAymGPzwqBiKHpue7eHmHqZlQ

    {
      "otp": "612765",
      "resend": false
    }
    ```
    Response(II) body:

    ```
    {
      "message": "Transaction Successful: User: #191197 transferred 98 coins to user: #192197",
      "error": null,
      "transaction_id": 1529
    }
    ```
    </details>
    <details>
    <summary><b><code>Reward Coins</code></b></summary>
    Request:
    
    ```http
    POST /reward HTTP/1.1
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
      "transaction_id": 1456
    }
    ```
    </details>
    <details>
    <summary><b><code>Make a Redeem request</code></b></summary>
    Request(I):
    
    ```http
    POST /redeems/new HTTP/1.1
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
    Response(I) body:

    ```
    {
      "message": "Post your otp to http://localhost:8080/redeems/new/confirm to confirm your transaction",
      "error": null
    }
    ```
    Request(II):
    
    ```http
    POST /redeems/new/confirm HTTP/1.1
    HOST: localhost:8080
    Content-Type: application/json
    Accept: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXRjaCI6IlkxOSIsImVtYWlsIjoiZGV2dGVzdC5hc2lzaEBnbWFpbC5jb20iLCJleHAiOjE2MjcwNDU1OTksInJvbGUiOiIiLCJyb2xsbm8iOjE5MTE5N30.4Fu80f4fWcdQwtxR1Ps4s5LPwqbD_dPeHucihz7Yi_A

    {
      "otp": "273801",
      "resend": false
    }
    ```
    Response(II) body:

    ```
    {
      "message": "Redeem request successful",
      "error": null,
      "request_id": 3441
    }
    ```
    </details>
    <details>
    <summary><b><code>See Pending Redeem Requests of all Users</code></b></summary>
    Request:

    ```http
    GET /redeems HTTP/1.1
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
    <details>
    <summary><b><code>Accept/Reject Redeem requests</code></b></summary>
    Request:

    ```http
    POST /redeems/update HTTP/1.1
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
    <details>
    <summary><b><code>Show status of previous Redeems of a user</code></b></summary>
    Request:

    ```http
    GET /redeems/status HTTP/1.1
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


- ### Database
  I have used two database management systems in this application, `SQLite` and `Redis`. The `init()` function of the `database` package automatically initializes the databases. The initialization errors are handled before making any other database operations.

- ### Write-Ahead Log
  - The `journal_mode` is set to `WAL` because of its [advantages](https://sqlite.org/wal.html#overview) over the default, `DELETE` mode in SQLite.
  - I personally tested in both modes and observed that the `WAL` mode works slightly faster (upto 10x faster) than the default mode while processing **parallelly requested** write operations into the database.
  - I also tested both the modes (again using parallel curl commands) intentionally keeping the DB locked for a certain time. In the default mode the concurrent requests are bound to be unsuccessful with an `database is locked` error. But, in `WAL` mode requests are handled sequentially and automatically once the db gets unlocked.

- ### Common Response Method
  Although the endpoints have slightly different formats for their response object, all of them are handled using a `type-switch` in a common `server.Respond()` function which responds to requests for all the endpoints. This method has been used a lot of times in various files. It has greatly reduced the bulkiness of codes in individual files.

- ### HTTP Status Codes
  A suitable http status code is assigned to every response.

- ### Access Token
  - The Access token has been made even more secure than before. The signature key of the access token changes every time the user logs in. In other words, **only the most recently generated access token is a valid token**.
  - How does this help?
    - Earlier, if the access token was somehow stolen, all the damage could have been done for the entire duration of the expiration period of the token, new tokens could have been generated but the damage won't be prevented.
    - Now, `any access token can be invalidated instantly` once the user relogs in, and one doesn't need to wait for the token to expire.
  - The signature key is a combination of a `Secret Key` stored in the `.env` file and an `uuid` that is generated along with the JWT.
  - Expiry time is currently set to `30 minutes`.

- ### .env
  - The `.env` file contains the following `enviroment variables`
    - `backend container name` and `redis container name`,
    - `redis password` (required to connect to the redis server),
    - `outgoing mail server` and its `port`,
    - `emailid` and `password` of the sender's account from which all OTPs are sent,
    - `secret key` required to sign the JWT,
    - `maximum cap` for the coins and the variable `minimum events` which is a criteria for users to be eligible for transactions,
    - `expiration time` for authorization tokens
  - If the `.env` file is not found the default values of these environment variables will be used throughout the application.
  - The admin can update these varibles in the `.env` file. The updated values will be overwritten to the default values of the variables defined in the source code.

  *The correct `EMAIL_ID` and `PASSWORD` needs to be set for the otp functionality to work.*
  *For running this application locally, the user will have the option to enter them*

- ### Cap for Maximum Coins
  Upper limit of the balance any user can hold. Currently set to `10001` coins.

- ### Minimum Events
  It is the minimum number of events to participate in to be eligible to make transactions. Currently set to `6`.

- ### Redeem
  - Users can send redeem requests which will be in pending state by default. This can be done on the `/redeems/new` endpoint. Once a valid request is made, an OTP is send to the user's emailid (that was collected during signup).
  - An Admin can see a list of all pending redeem requests made by all users on the `/redeems` endpoint.
  - Users can see the status of all their requests on the `/redeems/status` endpoint.
  - An Admin can "Accept" or "Reject" a redeem request on the `/redeems/update`endpoint.
  - Accepted redeem requests are stored in the DB as 1, Rejected ones as 0 and Pending ones as 2.

- ### OTP
  - `OTP` based confirmation systems are implemented on the `/redeems/new`, the `/transfer` and the `/reset_password` endpoints. The respective OTPs will have to be POSTed on the endpoints `/redeems/new/confirm`, `/transfer/confirm` and `/reset_password/confirm`.
  - There is also a `Resend OTP` option available (only at the confirmation endpoints). If a user wants to get another OTP, they have to POST a request with `resend` value set to `true`.

- ### The Process of Confirmation
  1. The user sends a request (on one of the endpoints - `/redeems/new` or `/transfer` or `reset/password`)
  2. If the request is invalid the server responds with suitable error messages
  3. If the request is valid -
      - An OTP is generated.
      - The OTP along with the rest of the data that needs to be stored/updated is temporarily saved in the `Redis` database. Expiry time is set to `2 mins` currently.
      - The OTP is sent to the user's emailid.
      - If the correct OTP is not entered/entered wrong, the process ends with an error message unless the user sets the `resend` option to be true.
        - If the resend option is true, one can enter the OTP again with a new POST request on the same endpoint
      - If the OTP is successfully entered and there is no error while storing/updating the necessary data
        - Immediately, the data along with the OTP is deleted from the `Redis` database.
      - If no request is made within this expiry time of 2 mins (not even a resend), the main data to be stored is lost
  
  *One can potentially delay the process (of transfer/redeem) if they keep on resending the OTP before the current one expires. But, this can be done until the `JWT` token expires, after which the user has to login again.*

- ### Testing
  I have used this script - http://p.ip.fi/A0uG to test the endpoints for multiple concurrent requests.

## More features to be added in the future
- [ ] Other modes of transaction - `IMMEDIATE`, `EXCLUSIVE` in SQLite
- [ ] Refresh token/similar for better user experience.
- [ ] ...

## Some incredibly helpful resources

> A common approach for invalidating tokens when a user changes their password is to sign the token with a hash of their password. Thus if the password changes, any previous tokens automatically fail to verify. You can extend this to logout by including a last-logout-time in the user's record and using a combination of the last-logout-time and password hash to sign the token. This requires a DB lookup each time you need to verify the token signature, but presumably you're looking up the user anyway.\
> [Travis Terry (stackoverflow)](https://stackoverflow.com/questions/21978658/invalidating-json-web-tokens/23089839#comment45057142_23089839)

> You can't change environment variables on a container (or any other process) after it's been created.\
> [David Maze (stackoverflow)](https://stackoverflow.com/a/65495853/15885436)

> Turn on the Write-Ahead Logging, Disable connections pool\
> [sqlite-concurrent-writing-performance (stackoverflow)](https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance/35805826), [Write-Ahead Logging (sqlite.org)](https://sqlite.org/wal.html)

> Once Commit or Rollback is called on the transaction, that transaction's connection is returned to DB's idle connection pool. The pool size can be controlled with SetMaxIdleConns.\
> [sql documentation (golang.org)](https://golang.org/pkg/database/sql/#DB)

> As a general rule of thumb, if you can use structs to represent your JSON data, you should use them. The only good reason to use maps would be if it were not possible to use structs due to the uncertain nature of the keys or values in the data.\
> [Soham Kamani](https://www.sohamkamani.com/golang/parsing-json/#what-to-use-structs-vs-maps)

> Some of the commands (in redis), especially with `the string data structure`, only make sense given specific type of data.\
> ...\
> `Hashes` are a good example of why calling Redis a key-value store isn’t quite accurate. You see, in a lot of ways, hashes are like strings. The important difference is that they provide an extra level of indirection: a field.\
> ...\
> The benefit would be the ability to pull and update/delete specific pieces of data, without having to get or write the entire value.\
> [The Little Redis Book - Karl Seguin](https://www.openmymind.net/redis.pdf)

> Two concurrent executions can interleave such that your read values become stale.\
  > Solutions:
  > 1. Do the read, write and validation checks in a single sql statement which is of write nature (so that it acquires lock), or
  > 2. Use other modes of transaction - `IMMEDIATE`, `EXCLUSIVE`, (more specific errors can be handled)
>
> [Bhuvan Singla](https://github.com/bhuvansingla)
