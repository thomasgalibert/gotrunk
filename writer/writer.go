package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"

	"github.com/thomas/gotrunk/profile"
)

func Load(path string, format profile.Format) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}
		return nil, err
	}
	out := map[string]any{}
	if len(data) == 0 {
		return out, nil
	}
	switch format {
	case profile.FormatJSON:
		err = json.Unmarshal(data, &out)
	case profile.FormatTOML:
		err = toml.Unmarshal(data, &out)
	case profile.FormatYAML:
		err = yaml.Unmarshal(data, &out)
	default:
		return nil, fmt.Errorf("format inconnu: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", path, err)
	}
	return out, nil
}

func Write(path string, format profile.Format, mode profile.Mode, values map[string]string) error {
	merged := map[string]any{}
	if mode == profile.ModeMerge {
		existing, err := Load(path, format)
		if err != nil {
			return err
		}
		merged = existing
	}
	for k, v := range values {
		merged[k] = v
	}

	data, err := marshal(format, merged)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("création du dossier: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("écriture temporaire: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("renommage final: %w", err)
	}
	return nil
}

// Backup copie le fichier existant vers <path>.bak-<timestamp>.
// Retourne le chemin de la sauvegarde ou "" si le fichier n'existait pas.
func Backup(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	backupPath := fmt.Sprintf("%s.bak-%s", path, time.Now().Format("20060102-150405"))
	if err := os.WriteFile(backupPath, data, info.Mode().Perm()); err != nil {
		return "", fmt.Errorf("création de la sauvegarde: %w", err)
	}
	return backupPath, nil
}

func marshal(format profile.Format, m map[string]any) ([]byte, error) {
	switch format {
	case profile.FormatJSON:
		return json.MarshalIndent(m, "", "  ")
	case profile.FormatTOML:
		return toml.Marshal(m)
	case profile.FormatYAML:
		return yaml.Marshal(m)
	default:
		return nil, fmt.Errorf("format inconnu: %s", format)
	}
}
