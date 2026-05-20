# Go Challenge Suite

Ces challenges sont faits pour etre resolus en modifiant uniquement les fichiers `challenge.go`.

Regles simples:
- les tests sont la source de verite
- quand les tests passent, le challenge est termine
- le squelette Go est volontairement minimal
- tu peux reorganiser librement le fichier tant que les tests passent

Commandes utiles depuis `apps/server`:

```bash
go test ./challenges/go/...
```

```bash
go test ./challenges/go/01_normalize_email
```

Ordre recommande:
1. `01_normalize_email`
2. `02_optional_int`
3. `03_validate_expense_input`
4. `04_group_by_category`
5. `05_parse_date`
6. `06_build_datetime`
7. `07_json_handler`
8. `08_status_from_error`
9. `09_user_lookup_service`
10. `10_target_service`
11. `11_password_roundtrip`
12. `12_month_snapshot`

Par design, les tests echouent au debut. C'est le jeu.
