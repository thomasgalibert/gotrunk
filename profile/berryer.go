package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var Berryer = Profile{
	ID:          "berryer",
	DisplayName: "Berryer — identifiants PISTE",
	Description: "Configure les identifiants PISTE utilisés par Berryer (lecture par le loader si les variables d'environnement sont absentes).",
	Fields: []Field{
		{
			Key:      "PISTE_CLIENT_ID",
			Label:    "PISTE Client ID",
			Help:     "Identifiant client PISTE fourni par votre administrateur Berryer.",
			Required: true,
			Validate: minLength(4),
		},
		{
			Key:      "PISTE_CLIENT_SECRET",
			Label:    "PISTE Client Secret",
			Help:     "Secret client PISTE. Sera masqué à la saisie.",
			Required: true,
			Secret:   true,
			Validate: minLength(8),
		},
		{
			Key:      "PISTE_ENV",
			Label:    "Environnement PISTE",
			Help:     "production ou sandbox.",
			Required: true,
			Default:  "production",
			Validate: minLength(2),
		},
	},
	Output: OutputTarget{
		Format: FormatJSON,
		Path:   BerryerCredentialsPath,
		Mode:   ModeMerge,
	},
}

// BerryerCredentialsPath retourne le chemin du fichier de credentials Berryer.
// Honore la variable d'environnement BERRYER_CREDENTIALS_FILE si elle est définie,
// sinon ~/.config/berryer/credentials.json sur tous les OS.
func BerryerCredentialsPath() (string, error) {
	if override := strings.TrimSpace(os.Getenv("BERRYER_CREDENTIALS_FILE")); override != "" {
		return override, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("résolution du dossier utilisateur: %w", err)
	}
	return filepath.Join(home, ".config", "berryer", "credentials.json"), nil
}

func minLength(n int) func(string) error {
	return func(s string) error {
		if len(strings.TrimSpace(s)) < n {
			return fmt.Errorf("au moins %d caractères requis", n)
		}
		return nil
	}
}
