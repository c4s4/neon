# Code Review : Projet NeON

## Introduction

**NeON** est un outil de build écrit en Go, conçu pour automatiser des tâches de compilation et de déploiement via des fichiers de configuration YAML (`build.yml`). Le projet se distingue par sa capacité à gérer l'héritage entre fichiers de build, permettant ainsi de définir des configurations de base et de les spécialiser dans des projets dérivés. Il intègre des fonctionnalités avancées comme la gestion d'un "singleton" (pour éviter l'exécution simultanée de plusieurs instances d'un build) et un système d'expressions pour dynamiser les propriétés du build.

L'architecture globale est découpée en plusieurs packages : `build` (le cœur du moteur), `builtin` (les tâches natives), `task` (la définition des tâches) et `util` (les fonctions d'aide).

## Problèmes identifiés et propositions de solutions

### 1. Architecture et Design

#### Signature de fonction excessivement longue
Dans `neon/neon/main.go`, la fonction `ParseCommandLine` retourne **24 valeurs** distinctes :
```go
func ParseCommandLine() (string, bool, bool, string, bool, bool, string, bool,
	bool, string, bool, bool, bool, string, string, bool, bool, bool, string,
	bool, bool, string, bool, []string)
```
C'est un "anti-pattern" majeur en développement logiciel. Cela rend le code extrêmement difficile à lire, à maintenir et très fragile lors de l'ajout de nouvelles options.

*   **Solution :** Créer une structure `Options` ou `Config` qui regroupe tous ces paramètres et retourner un seul pointeur vers cette structure.

#### Gestion des erreurs et sortie brutale
La fonction `PrintError` dans `main.go` appelle systématiquement `os.Exit(1)`.
L'utilisation de `os.Exit` au sein de fonctions utilitaires rend le code impossible à tester unitairement (car le processus s'arrête brusquement) et empêche un nettoyage propre des ressources (les `defer` ne sont pas exécutés).

*   **Solution :** Faire remonter les erreurs (`return err`) jusqu'à la fonction `main`, et n'appeler `os.Exit` qu'à un seul endroit, tout à la fin du point d'entrée du programme.

### 2. Fiabilité et Robustesse

#### Fuite de Goroutine dans `ListenPort`
Dans `neon/neon/build/build.go`, la fonction `ListenPort` lance une goroutine qui boucle indéfiniment sur `listener.Accept()`.
```go
go func() {
    for {
        _, _ = listener.Accept()
        time.Sleep(100 * time.Millisecond)
    }
}()
```
Même quand le `listener` est fermé (via le `defer` dans `Run`), la boucle continue de tourner et de dormir, créant une fuite de ressources.

*   **Solution :** Vérifier l'erreur retournée par `listener.Accept()`. Si une erreur survient (ce qui arrive quand le listener est fermé), la boucle doit être interrompue (`break`).

#### Risque de boucle infinie (Récursion)
La méthode `GetParents()` dans `build.go` appelle `NewBuild` récursivement pour charger les fichiers parents définis dans le champ `extends`. Si un utilisateur crée une dépendance circulaire (le fichier A étend B, qui étend A), le programme plantera avec un `stack overflow`.

*   **Solution :** Implémenter un mécanisme de détection de cycles en gardant une trace des fichiers déjà chargés durant la phase de résolution des parents.

### 3. Maintenabilité et Modernisation

#### Dépendances obsolètes
Le projet utilise `gopkg.in/yaml.v2`. La version `v3` est désormais disponible et apporte des améliorations significatives en termes de performances, de gestion des erreurs et de support des spécifications YAML.

*   **Solution :** Migrer vers `gopkg.in/yaml.v3`.

#### Manque d'abstractions (Interfaces)
Le code s'appuie presque exclusivement sur des types concrets (`Build`, `Target`, `Context`). Cela rend le couplage très fort et complique la création de tests unitaires (mocks).

*   **Solution :** Introduire des interfaces pour les composants clés, notamment pour les opérations d'exécution de tâches, afin de pouvoir tester le moteur de build sans avoir à créer de vrais fichiers sur le disque.

### 4. Expérience Utilisateur (UX)

#### Rigidité de la configuration
Le chemin du fichier de configuration est codé en dur : `~/.neon/settings.yml`. Bien que ce soit standard pour un CLI, cela pose problème dans des environnements de CI/CD ou des containers où l'on souhaite injecter la configuration via des variables d'environnement.

*   **Solution :** Permettre la définition d'un chemin alternatif via une variable d'environnement (ex: `NEON_CONFIG_PATH`) avant de se rabattre sur le chemin par défaut.

## Correctifs

Les points suivants ont été corrigés :

- **Signature de `ParseCommandLine`** : La fonction retourne désormais une structure `Options` au lieu de 24 valeurs distinctes.
- **Gestion des erreurs** : La fonction `PrintError` ne provoque plus l'arrêt brutal du programme via `os.Exit(1)`. Les erreurs remontent désormais jusqu'à la fonction `main` qui gère la sortie.
- **Fuite de Goroutine dans `ListenPort`** : La boucle d'acceptation des connexions s'interrompt désormais correctement lorsque le listener est fermé, évitant ainsi la fuite de ressources.
- **Configuration flexible** : Le chemin du fichier de configuration peut être surchargé via la variable d'environnement `NEON_CONFIG_PATH`.
- **Détection de cycles** : Détection de cycles implémentée dans l’héritage des builds ; une erreur claire est renvoyée au lieu d’un débordement de pile.
- **Abstraction** : Introduire des interfaces pour les composants clés (`Build`, `Target`, `Context`) ne sera pas réalisé car inutile.
- **Migration Anko** : Mise à jour de Anko en version *0.1.12*.

## TODO

Les points suivants restent à traiter :

- **Migration YAML** : Migrer de `gopkg.in/yaml.v2` vers `gopkg.in/yaml.v3`.
