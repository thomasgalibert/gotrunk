# gotrunk

Petit utilitaire pour configurer facilement Berryer (et autres plugins futurs) sur les postes des clients non-techniques.

L'exécutable ouvre une fenêtre terminal, pose quelques questions en langage clair, et écrit le fichier de configuration au bon endroit du système.

## Profil livré : `berryer`

Demande les identifiants PISTE (Client ID, Client Secret, environnement) et écrit :

- `~/.config/berryer/credentials.json` par défaut (Mac, Linux, Windows)
- ou le chemin pointé par la variable d'environnement `BERRYER_CREDENTIALS_FILE` si elle est définie.

Le fichier produit est exactement au format attendu par le loader (`packages/core/src/config.ts`) :

```json
{
  "PISTE_CLIENT_ID": "xxx",
  "PISTE_CLIENT_SECRET": "yyy",
  "PISTE_ENV": "production"
}
```

Le mode est `merge` : les valeurs déjà présentes dans le fichier sont pré-remplies dans le formulaire et les clés non touchées sont préservées.

## Pour les développeurs

### Builder

```bash
make build           # binaire local
make build-all       # mac-intel + mac-apple + windows
make package-mac     # crée des .command (double-cliquables sur Mac)
make package-windows # crée BerryerSetup.exe
```

### Ajouter un nouveau profil

1. Créer `profile/monplugin.go` exportant une variable `MonPlugin profile.Profile`.
2. L'ajouter à `registry` dans `profile/registry.go`.
3. Rebuild.

### Tester localement

```bash
make run                     # menu de profils (un seul = lancement direct)
./dist/gotrunk --profile=berryer
```

## Pour les clients

### macOS

1. Télécharger `BerryerSetup-mac-apple.zip` (Mac M1/M2/M3) ou `BerryerSetup-mac-intel.zip` (Mac plus ancien).
2. Double-cliquer sur le `.zip` dans le Finder — un fichier `.command` est extrait à côté.
3. Double-cliquer sur le `.command` extrait.
4. **Si macOS affiche un avertissement de sécurité** ("application non identifiée") : faire **clic droit** sur le `.command` → **Ouvrir** → confirmer **Ouvrir** dans la boîte de dialogue. C'est nécessaire une seule fois.
5. Suivre les instructions à l'écran.

> Le `.zip` est nécessaire pour préserver les permissions d'exécution : un téléchargement direct d'un binaire macOS via le navigateur supprime le bit `+x`, ce qui force l'utilisateur à passer par le Terminal pour le réactiver.

### Windows

1. Télécharger `BerryerSetup.exe`.
2. Double-cliquer dessus.
3. **Si Windows SmartScreen affiche un avertissement** : cliquer sur **Informations complémentaires** puis **Exécuter quand même**.
4. Suivre les instructions à l'écran.
