# IITK Coin

Please look at the request format for each endpoint in their respective functions.

In order to access /secretpage please provide an authorization token in the **header**.

Deliberately left .env out of .gitignore for the purposes of checking.

For my reference:
- [ ] Handle errors in signup.go

- [X] send a json object in the response

- [X] GenJWT, SecretKey

- [X] Handle errors: panic(err) -> print(something went wrong...), print === Fprintf/Errorf...

- [ ] isAdmin -> /secretPage

- [ ] http.StatusBadRequest, http.StatusUnauthorized

- [ ] http.HandleFunc => http.POST/http.Handle().Methods("POST")

- [ ] User struct {
  ID       uint64 `json:"id"`
  Username string `json:"username"`
  Password string `json:"password"`
}
