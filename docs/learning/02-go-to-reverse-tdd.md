# From Go Practice To Reverse TDD

Objectif: passer des mini-challenges a la reecriture progressive du projet sans me noyer dans la complexite.

## Philosophie

Le but n'est pas de "refactorer le repo".
Le but est de recuperer:
- de la comprehension
- de la confiance
- de la maitrise

Je ne reecris pas tout.
Je choisis de petits morceaux avec une valeur pedagogique forte et un risque faible.

## Etape 1 - Reapprendre a ecrire du Go sans friction

Duree cible: 2 a 4 semaines.

Je fais:
- 1 mini-challenge tous les 2 ou 3 jours
- surtout sur parsing, validation, erreurs, `time`, `http`, tests table-driven

Je passe a l'etape suivante quand je peux ecrire sans panique:
- une fonction avec `error`
- un test table-driven
- un petit handler HTTP
- un petit service avec interface repository

## Etape 2 - Lire TrackStack comme un atelier, pas comme une cathedrale

Avant de supprimer du code, je lis des zones ciblees:
- domaine simple
- service simple
- handler simple
- helper utilitaire

Je me limite a comprendre:
- les entrees
- les sorties
- les dependances
- les validations
- les erreurs

Je ne cherche pas a comprendre tout le runtime d'un coup.

## Etape 3 - Reverse TDD sur petits perimetres

Regle absolue:
- je commence par des fonctions/services isoles
- pas par le runtime monolithique
- pas par l'auth complete
- pas par Terraform
- pas par les couches pleines d'effets de bord

Workflow:
1. choisir une petite cible
2. ecrire ou faire generer des tests
3. lire les tests jusqu'a comprendre le comportement
4. supprimer seulement l'implementation visee
5. reecrire a la main
6. faire passer les tests
7. noter ce que j'ai appris

## Etape 4 - Remonter progressivement en difficulte

Ordre recommande:
1. helper pur
2. fonction de domaine
3. service applicatif simple
4. handler HTTP simple
5. service avec calcul temporel
6. auth et sessions plus tard

## Regles d'hygiene mentale

- Une seule cible active a la fois.
- Si je bloque plus de 30 minutes, je redescends d'un cran.
- Je ne prends pas un bloc "important" comme premier terrain de jeu.
- Je prefere finir petit que rever grand.
- Compiler, tester, corriger: le compilateur est le coach, pas l'ennemi.

## Usage de l'IA

L'IA sert a:
- proposer une cible d'apprentissage
- generer une premiere batterie de tests
- expliquer une erreur apres que j'ai essaye
- comparer ma solution avec l'existant

L'IA ne sert pas a:
- ecrire l'implementation a ma place
- me voler la boucle essai/erreur
- court-circuiter la comprehension

## Rythme realiste

Avec une vie de parent:
- mini-challenges courts en semaine
- reverse TDD leger seulement quand l'energie est bonne
- pas d'obligation de cadence heroique

Cadence recommandee:
- semaine A: 2 snacks Go
- semaine B: 1 snack Go + 1 reverse TDD
- semaine C: 1 snack Go + 1 reverse TDD
- puis ajustement selon fatigue reelle

## Definition de succes

Le succes n'est pas:
- reecrire le plus de code possible

Le succes est:
- comprendre plus de code qu'avant
- dependre moins de l'IA
- retrouver des reflexes Go
- pouvoir expliquer ce qu'un module fait sans reciter une doc
