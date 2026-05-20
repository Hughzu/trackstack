# Challenge 08 - Status From Error

Objectif:
- mapper les erreurs vers des status HTTP
- supporter `errors.Is` sur des erreurs wrappees

Mapping attendu:
- invalid input -> 400
- unauthorized -> 401
- not found -> 404
- tout le reste -> 500

Tu modifies seulement `challenge.go`.

Commande:

```bash
go test ./challenges/go/08_status_from_error
```
