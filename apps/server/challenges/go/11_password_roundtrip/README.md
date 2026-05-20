# Challenge 11 - Password Roundtrip

Objectif:
- refuser un mot de passe vide
- hasher un mot de passe avec un format qui commence par `$scrypt$`
- verifier le bon mot de passe
- rejeter un mauvais mot de passe
- rejeter un hash invalide sans panic

Tu modifies seulement `challenge.go`.

Commande:

```bash
go test ./challenges/go/11_password_roundtrip
```
