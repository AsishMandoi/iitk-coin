# IITK Coin

According to my directory structure, I have split my package into multiple sub-packages (i.e. I have made a few sub-directories - `functions`, `global`, `handlers`).

Please look at the request format for each endpoint in the **`global`** package, where I have defined some variables and struct types, to be used in other functions.

In the three route handler functions I have used a `struct` (of types as in the `global` package) assigned to the variable named `payload`, for the body of the response to be sent. The three of them have similar structure. I have later converted these `struct`s into json objects.

In order to access /secretpage please provide an authorization token in the **header**.

Deliberately left `.env` out of `.gitignore` for the purposes of checking.
  
---
### For my reference:
- [ ] Use refresh token
- [X] send a json object in the response
- [X] GenJWT, SecretKey
- [X] handle errors better: panic(err), Fprintf, Errorf... -> return error in header/body
- [ ] isAdmin -> /secretPage
- [ ] keep checking for unhandled errors
- [X] use MDN: HTTP status codes -> http.StatusBadRequest, http.StatusUnauthorized
- [ ] (opt) http.HandleFunc => http.POST/http.Handle().Methods("POST")
- [X] Stu struct {
  Rollno string \`json:"rollno"\`,
  Password string \`json:"password"\`
}
- [X] look up popular directory structures