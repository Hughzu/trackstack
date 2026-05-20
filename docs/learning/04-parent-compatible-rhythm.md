# Parent-Compatible Rhythm

Objectif: faire progresser Go et TrackStack avec une vraie vie de parent, pas avec un planning de moine imaginaire.

## Principe central

- Je n'optimise pas pour ma semaine ideale.
- J'optimise pour ma semaine fatiguee.
- Les sessions moches, courtes, mais reelles battent les grands blocs fantasmes qui sautent tout le temps.

## Machine par defaut

- Je commence sur le Mac, parce que c'est la machine que j'ouvrirai vraiment.
- Le PC fixe reste utile pour les sessions plus lourdes, les jeux, ou du vrai temps profond quand il existe.
- Je n'achete pas un laptop tout de suite.
- Si, apres quelques semaines, mon vrai workflow vit surtout sur laptop et que le partage du Mac devient un goulot d'etranglement, alors j'achete mon propre laptop.

## Cadence hebdomadaire cible

### Semaine

- 2 soirs madame
- 1 soir Go challenges
- 1 soir TrackStack
- 1 soir bonus: jeu, madame, TrackStack, Go, ou rien

### Week-end

- aucune dette morale de code
- jeu, temps de couple, repos, ou code seulement si une fenetre naturelle s'ouvre

## Regles de travail

- Une session de code = 40 a 60 minutes max.
- Les soirs fatigues servent a faire petit et clair, pas a penser un grand redesign.
- Le deep work n'est pas interdit, mais c'est un bonus, pas la base du systeme.
- TrackStack reste important parce que c'est la vitrine, mais il doit etre decoupe pour rentrer dans des sessions courtes.

## Repartition des types de sessions

### Soir Go

Utiliser ces sessions pour:
- un challenge Go
- un test table-driven
- une petite fonction
- un exercice de parsing, validation, HTTP, erreurs

### Soir TrackStack

Utiliser ces sessions pour:
- une seule tache de 30 a 45 minutes
- un helper
- un service court
- un test
- une mini reecriture reverse TDD
- une petite mise a jour de doc technique

## Regle d'or pour TrackStack

Une session ne doit jamais commencer par: "Bon, je fais quoi ?"

Chaque tache TrackStack doit:
- etre definie avant la session
- etre faisable en 30 a 45 minutes
- avoir une fin visible

Exemples de bonnes taches:
- ecrire les tests d'un helper
- reecrire une fonction de parsing
- couvrir un petit service
- faire passer un package de tests
- documenter un flux technique precis

Exemples de mauvaises taches:
- avancer sur l'auth
- reprendre TrackStack
- travailler l'architecture

Ce ne sont pas des taches. Ce sont des facons elegantes de perdre 40 minutes.

## Routine de fin de session

Avant de fermer le laptop, noter:
- la prochaine action exacte
- le fichier vise
- la commande de test a relancer

Format minimal:

```text
Prochaine action: ecrire les tests de ParseDate pour le cas RFC3339 invalide
Fichier: apps/server/internal/platform/timeutil/parse.go
Commande: go test ./internal/platform/timeutil
```

## Regle mentale

- Je n'attends pas de retrouver la vie d'avant avec des blocs de 2 heures.
- Je construis un systeme qui survit a la fatigue, au couple, au boulot, au sport, et a la vraie vie.
- Mon objectif n'est pas d'avoir l'air discipline.
- Mon objectif est de tenir dans le temps sans peter une durite.

## Verdict pratique

- Oui aux petites habitudes.
- Oui a 40 a 60 minutes un soir sur trois.
- Oui a TrackStack, mais decoupe en tickets de mercenaire.
- Oui aux jeux video aussi, sinon le systeme devient punitif et meurt.
- Non au fantasme du grand bloc parfait comme fondation unique.
