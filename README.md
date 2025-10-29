# SMTS - Sign Me This Shit

Outil automatisé pour générer et signer les feuilles de présence des étudiants FIP 3A d'IMT Atlantique.

## Prérequis

- Google Chrome ou navigateur compatible Chromium installé (utilisé pour générer le PDF du planning en mode headless)
- Identifiants IMT Atlantique

## Installation

Pour installer **SMTS**, téléchargez le binaire pour votre plateforme depuis la page [Releases](https://github.com/mathismqn/smts/releases).

1. Allez sur la page [Releases](https://github.com/mathismqn/smts/releases)
2. Téléchargez le binaire approprié pour votre système (Windows, macOS ou Linux)
3. (Optionnel) Déplacez le binaire dans un répertoire inclus dans `$PATH` pour un accès facile

Exemple (Linux/macOS) :
```bash
# Déplacer le binaire téléchargé dans /usr/local/bin pour y accéder depuis n'importe où
sudo mv smts-linux-amd64 /usr/local/bin/smts
sudo chmod +x /usr/local/bin/smts
```

## Utilisation

### Préparation de votre signature

Vous devez avoir un fichier de signature au format PNG. Deux options :
- Nommez-la `signature.png` et placez-la dans le répertoire courant
- Ou utilisez un autre nom/emplacement (à préciser lors de l'utilisation)

Pour créer votre signature, vous pouvez utiliser : https://onlinesignature.com/fr/draw-a-signature-online

### Configuration

Avant la première utilisation, configurez vos identifiants IMT Atlantique :

```bash
smts setup
```

Cette commande vous demandera votre nom d'utilisateur et mot de passe IMT Atlantique. Ces informations sont nécessaires pour :
- Se connecter automatiquement à votre agenda PASS
- Récupérer votre planning de la semaine
- Détecter automatiquement votre campus

**Sécurité** : Vos identifiants sont stockés de manière sécurisée dans le trousseau de clés de votre système (Keychain sur macOS, Windows Credential Manager sur Windows, Secret Service sur Linux). Ils ne sont jamais envoyés ailleurs que sur les serveurs officiels d'IMT Atlantique.

### Génération de la feuille de présence

**Important** : À effectuer avant la fin de la semaine car l'outil génère la feuille pour la semaine en cours.

```bash
smts sign
```

Le fichier PDF généré est automatiquement nommé : `NOM Prénom – FIPA3[Campus] – S[Semaine].pdf`

Vous pouvez utiliser une signature personnalisée :
```bash
smts sign --signature /chemin/vers/signature.png
```

#### Options avancées

Vous pouvez personnaliser les informations de la feuille de présence avec les options suivantes :

```bash
# Spécifier un campus différent (Brest, Rennes, ou Nantes)
smts sign --campus Rennes

# Spécifier un nom personnalisé (prénom et nom requis ensemble)
smts sign --firstname Jean --lastname DUPONT

# Combiner plusieurs options
smts sign --campus Brest --firstname Marie --lastname MARTIN --signature ma-signature.png
```

**Options disponibles :**
- `--signature, -s` : chemin vers le fichier de signature (défaut: `signature.png`)
- `--campus` : campus (Brest, Rennes, ou Nantes) - auto-détecté si non fourni
- `--firstname` : prénom - auto-détecté si non fourni (requiert `--lastname`)
- `--lastname` : nom de famille - auto-détecté si non fourni (requiert `--firstname`)

**Note** : Par défaut, le campus et le nom sont automatiquement détectés depuis votre agenda PASS. Les options ci-dessus permettent de les remplacer si nécessaire.
