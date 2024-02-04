# MMOME

La société a racheté une startup qui développe le jeu MMOME (Massively Multiplayer Online Maitre Esprit).
Le projet présente des difficultés de maintenabilité, d'évolutivité, de performance et a des bugs connus.

En tant qu'experts Go, vous avez été appelé pour résoudre ces problèmes.

## Le jeu

Le jeu consiste à deviner une combinaison de 5 emojis (c'est un jeu moderne) choisis au hasard parmi cette liste
`🏴💀🔥🎉🚀🤡`.

À chaque essai, le jeu indique au joueur le nombre d'emojis bien placés et mal placés dans sa proposition.

Il n'y a pas de limite au nombre d'essais que peut faire un joueur.

Le gagnant est le joueur qui trouve la bonne combinaison le plus rapidement.

## Le serveur

Le jeu se joue au travers d'une API HTTP.

Voici les endpoints disponibles:

* `/new`: Crée une nouvelle partie si aucune partie n'est en cours.
* `/guess/{guess}`: Propose une combinaison et retourne le résultat sous la forme d'un JSON de la forme
  `{"correct": 2, "misplaced": 1}`.

L'API est décrite par la spécification [server/openapi.yaml](./server/openapi.yaml).

Lorsqu'un joueur trouve la bonne combinaison, la partie en cours s'arrête et il n'est plus possible de proposer une
nouvelle combinaison tant que l'endpoint `/new` n'a pas été appelé.

Le code actuel est dans le fichier [server/server.go](./server/server.go).

### Problèmes

Le code actuel a été écrit quasiment entièrement en acceptant les propositions de l'IA installée dans l'IDE de l'équipe
de développement.

Plusieurs problèmes ont été identifiés lors de l'audit du projet:

* Il mélange les responsabilités. L'auditeur nous suggère de séparer le code métier de la couche HTTP pour améliorer la
  maintenabilité et faciliter les tests.
* Le produit souhaite pouvoir faire évoluer le service à l'avenir. Avoir des règles différentes? Plusieurs parties en
  même temps? des json plus beaux pour un effet waouh? Le code n'est actuellement pas prêt pour ça.
* Le code n'est évidemment pas testé.
* Les utilisateurs se plaignent d'un bug, le compteur d'emojis mal placés ne fonctionne pas correctement, il compte
  plusieurs fois les emojis présent dans la solution.
    * Par exemple, si la combinaison à trouver est `🏴🏴🚀🏴🏴` et que le joueur
      propose `🚀🚀🔥🚀🚀`, le résultat devrait être `{"correct": 0, "misplaced": 1}` mais actuellement, il
      est `{"correct": 0, "misplaced": 4}`.

### Objectifs

* Réorganiser le code en séparant la logique métier de la couche HTTP, simplifier la maintenabilité et
  l'évolutivité.
* Écrire des tests qui vérifient au moins que la logique du jeu est correcte.
* Profiter de l'écriture des tests pour corriger le bug du compteur d'emojis mal placés.

## Le solutionneur automatique

L'utilisation de l'API est fastidieuse, et le jeu demande de la reflexion, ce qui rebute certain de nos utilisateurs.

Afin de le monétiser, la société propose donc un service premium qui fournit un outil permettant de trouver la
combinaison gagnante automatiquement.

Le code actuel est dans le fichier [solver/solver.go](./solver/solver.go).

### Problèmes

Le développeur n'ayant aucune idée de comment résoudre efficacement le problème, il a donc simplement choisi d'essayer
toutes les combinaisons possibles jusqu'à trouver la bonne.

C'est lent, mais ça fonctionne. Cependant, les utilisateurs australiens se plaignent de ne jamais gagner.

Aussi, l'outil d'analyse de qualité de code se plaint des 5 boucles for imbriquées, qui pourrait être bien plus si on
devait gérer des combinaisons de plus de 5 emojis. Encore une fois, la maintenabilité et l'évolutivité pourrait être
améliorés.

### Objectifs

* Réorganiser le code afin de faire la recherche brut force de la combinaison gagnante en asynchrone pour moins
  pénaliser les connexions lentes. Les performances devraient au moins être multipliées par 50!
* Améliorer la maintenabilité du code, et l'évolutivité de la commande quand c'est possible.
* v2: Améliorer l'algorithme pour qu'il trouve la combinaison gagnante en le moins d'essais possible.
