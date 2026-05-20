# Go Mini Challenges

Objectif: reconstruire une memoire musculaire en Go sans transformer l'apprentissage en deuxieme emploi.

Les challenges se trouvent dans `apps/server/challenges/go/`.

## Regles du jeu

- Frequence cible: un challenge tous les 2 ou 3 jours.
- Duree cible: 20 a 45 minutes.
- Pas de marathon. Si un exercice deborde, je le coupe et je reviens plus tard.
- J'ecris tout a la main.
- L'IA peut proposer un exercice ou generer des tests, mais pas ecrire la solution a ma place.
- Je privilegie la standard library.
- Chaque exercice doit produire un petit resultat concret: une fonction, un test, un mini handler, un parseur, un mapper.

## Format d'une session

1. Lire l'enonce.
2. Ecrire d'abord un ou deux tests simples.
3. Coder la solution.
4. Corriger jusqu'a ce que ca passe.
5. Noter en 2 lignes ce qui m'a freine et ce que j'ai retenu.

## Semaine 1 - Refluidifier les bases

### Challenge 1 - Normalize email

Ecrire une fonction qui:
- trim les espaces
- met en minuscule
- retourne une erreur si la valeur est vide

Competences:
- strings
- erreurs
- tests table-driven

### Challenge 2 - Parse int optionnel

Ecrire une fonction qui convertit une string en `*int`.
- `""` retourne `nil`
- `"42"` retourne pointeur vers `42`
- `"abc"` retourne `nil`

Competences:
- pointeurs
- strconv
- helpers de parsing

### Challenge 3 - Validate command input

Creer une struct `CreateExpenseInput` et une fonction `Validate() error`.
Valider:
- `UserID` requis
- `Label` requis
- `Amount` > 0

Competences:
- struct
- methodes
- validation
- erreurs metier simples

### Challenge 4 - Group by category

A partir d'une slice d'entrees, retourner une map `category -> total`.

Competences:
- slices
- maps
- boucle
- accumulation

## Semaine 2 - Go utile pour TrackStack

### Challenge 5 - Parse date

Ecrire une fonction qui accepte:
- une date RFC3339
- ou une date `YYYY-MM-DD`
et retourne un `time.Time` UTC.

Competences:
- package `time`
- parsing
- normalisation

### Challenge 6 - Build datetime from date + time

Ecrire une fonction qui combine:
- une date optionnelle
- une heure optionnelle
en une string RFC3339.

Competences:
- time
- valeurs par defaut
- tests sur les cas limites

### Challenge 7 - HTTP JSON decode

Ecrire un mini handler HTTP qui:
- lit un body JSON
- refuse les champs inconnus
- repond `400` si payload invalide
- repond `200` avec une reponse JSON sinon

Competences:
- `net/http`
- `encoding/json`
- `httptest`

### Challenge 8 - Error mapping

Ecrire une fonction `StatusFromError(err error) int`.
Exemple:
- invalid input -> 400
- not found -> 404
- tout le reste -> 500

Competences:
- `errors.Is`
- conventions de transport

## Semaine 3 - Preparation au reverse TDD

### Challenge 9 - Fake repository + service

Creer un service simple qui depend d'une interface repository.
Ecrire un fake manuel pour les tests.

Competences:
- interfaces
- injection de dependances
- tests unitaires sans framework magique

### Challenge 10 - Create target

Ecrire un service qui:
- cherche une target utilisateur
- en cree une par defaut si absente
- retourne une erreur si `UserID` est vide

Competences:
- use case
- repository port
- defaults
- logique applicative

### Challenge 11 - Password hash round-trip

Ecrire un test qui verifie:
- qu'un mot de passe hashe peut etre verifie
- qu'un mauvais mot de passe echoue
- qu'un hash invalide echoue proprement

Competences:
- tests robustes
- lecture de code existant
- securite basique sans magie

### Challenge 12 - Mini dashboard

A partir d'entrees datees, calculer:
- total periode courante
- total periode precedente
- delta
- pourcentage si possible

Competences:
- time windows
- petits calculs metier
- struct de sortie

## Bibliotheque de snacks bonus

- filtrer une slice par predicat
- ecrire un middleware HTTP simple
- parser un query param avec defaut
- mapper DTO JSON -> commande metier
- ecrire un test table-driven sur des erreurs
- ecrire une fonction de normalisation d'IP
- ecrire une fonction de hash SHA256 d'une string

## Critere de reussite

Je ne cherche pas la performance ni l'elegance parfaite.
Je cherche:
- de la fluidite
- des reflexes Go
- moins d'hesitation devant un fichier vide
- plus de confiance avant de toucher TrackStack
