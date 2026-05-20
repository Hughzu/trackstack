# Challenge 06 - Build Datetime

Objectif:
- combiner une date et une heure en RFC3339 UTC
- si la date est vide, utiliser la date UTC de `now`
- si l'heure est vide, utiliser l'heure UTC de `now` au format `HH:MM`
- erreur si la combinaison est invalide

Tu modifies seulement `challenge.go`.

Commande:

```bash
go test ./challenges/go/06_build_datetime
```
