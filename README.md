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

Ou avec une signature personnalisée :
```bash
smts sign --signature /chemin/vers/signature.png
```

Le fichier PDF est automatiquement nommé : `NOM Prénom – FIPA3[Campus] – S[Semaine].pdf`
