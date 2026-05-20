# TrackStack Rewrite Order

Objectif: choisir les premieres cibles de reverse TDD dans le bon ordre, du plus formateur au moins risque.

## Principe de selection

Je privilegie les zones:
- petites
- testables
- a faible couplage
- riches en apprentissage Go
- proches du vrai produit

Je repousse les zones:
- tres couplees a l'infra
- pleines d'effets de bord
- longues
- sensibles a l'auth ou au runtime complet

## Niveau 1 - Cibles parfaites pour demarrer

### 1. `apps/server/internal/platform/timeutil/parse.go`

Pourquoi:
- petit fichier
- parsing utile
- zero infra
- excellent terrain pour tests table-driven

J'y travaille pour apprendre:
- `time`
- validation
- erreurs
- cas limites

### 2. `apps/server/internal/contexts/heat/domain/time.go`

Pourquoi:
- minuscule
- pure fonction
- feedback rapide

J'y travaille pour apprendre:
- logique metier simple
- tests propres
- demarrage sans friction

### 3. `apps/server/internal/contexts/users/application/services/password.go`

Pourquoi:
- tres bon exercice de lecture
- fonctions petites
- melange standard library + robustesse

Attention:
- je ne reinvente pas la crypto
- je reecris pour comprendre, pas pour etre creatif

J'y travaille pour apprendre:
- fonctions utilitaires
- parsing de format
- tests de comportement
- comparaison sure

## Niveau 2 - Bons services applicatifs simples

### 4. `apps/server/internal/contexts/calories/application/services/target_service.go`

Pourquoi:
- use case clair
- validation simple
- create-or-default
- repository injectable

J'y travaille pour apprendre:
- interfaces
- fake repositories
- tests de service
- timestamps et creation conditionnelle

### 5. `apps/server/internal/contexts/expenses/application/services/settings_service.go`

Pourquoi:
- meme famille pedagogique que target service
- lecture facile
- pattern get-or-create tres concret

J'y travaille pour apprendre:
- orchestration de dependances
- valeurs par defaut
- enrichissement d'une vue

## Niveau 3 - HTTP simple avant le gros bordel

### 6. `apps/server/internal/contexts/calories/adapters/inbound/http/calories_handler.go`

Pourquoi:
- handler lisible
- beaucoup de petites fonctions testables
- parsing JSON, query params, mapping erreurs

Je ne reecris pas tout d'un coup.
Je commence par:
- `parseOptionalInt`
- `parseOptionalStringPtr`
- `decodeJSON`
- mapping d'erreurs
- puis seulement un handler public

J'y travaille pour apprendre:
- `httptest`
- JSON
- transport HTTP
- traduction request -> use case

## Niveau 4 - Calcul metier plus dense

### 7. `apps/server/internal/contexts/heat/application/services/refill_service.go`

Pourquoi:
- logique metier interessante
- calcul de saison
- pagination simple
- dates et snapshots

Attention:
- plus dense
- meilleur deuxieme ou troisieme vrai chantier, pas le premier

J'y travaille pour apprendre:
- calcul temporel
- composition de view model
- tests sur periodes

## Niveau 5 - A garder pour plus tard

### 8. `apps/server/internal/contexts/auth/application/services/auth_service.go`

Pourquoi attendre:
- plus de branches
- securite
- sessions
- rotation
- asynchronisme
- dependances plus nombreuses

Tres formateur, oui.
Bon point d'entree, non.

Je le garde pour quand je serai plus a l'aise avec:
- interfaces
- fakes
- tests de branches
- `context`
- temps et expiration

## Ordre concret recommande

1. `apps/server/internal/platform/timeutil/parse.go`
2. `apps/server/internal/contexts/heat/domain/time.go`
3. `apps/server/internal/contexts/users/application/services/password.go`
4. `apps/server/internal/contexts/calories/application/services/target_service.go`
5. `apps/server/internal/contexts/expenses/application/services/settings_service.go`
6. helpers internes de `apps/server/internal/contexts/calories/adapters/inbound/http/calories_handler.go`
7. un handler HTTP public dans `apps/server/internal/contexts/calories/adapters/inbound/http/calories_handler.go`
8. `apps/server/internal/contexts/heat/application/services/refill_service.go`
9. `apps/server/internal/contexts/auth/application/services/auth_service.go`

## Methode de travail par cible

Pour chaque cible:
1. lire le fichier
2. ecrire des tests de comportement
3. supprimer seulement le bloc vise
4. reecrire a la main
5. comparer avec l'existant si utile
6. noter les idiomes Go rencontres

## Feu vert / feu rouge

Feu vert:
- helpers purs
- services courts
- validation
- parsing
- mapping
- handlers simples

Feu rouge au debut:
- runtime monolith
- wiring complet
- auth complete
- Lambda adapter
- config globale
- DB reelle si un fake suffit

## Critere de progression

Je monte au niveau suivant quand:
- j'ecris les tests plus vite
- je bloque moins sur la syntaxe Go
- je comprends les dependances sans panique
- je peux expliquer ce que fait la cible en francais simple
