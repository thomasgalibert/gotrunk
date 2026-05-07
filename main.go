package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thomas/gotrunk/profile"
	"github.com/thomas/gotrunk/ui"
	"github.com/thomas/gotrunk/writer"
)

func main() {
	profileID := flag.String("profile", "", "ID du profil à configurer (ex: berryer). Si omis : sélection directe s'il n'y a qu'un seul profil, sinon un menu s'affiche.")
	flag.Parse()

	if err := run(*profileID); err != nil {
		fmt.Fprintf(os.Stderr, "\nErreur : %v\n", err)
		ui.PressEnterToClose("Une erreur est survenue. Vous pouvez fermer cette fenêtre.")
		os.Exit(1)
	}
}

func run(profileID string) error {
	var p profile.Profile
	var err error
	if profileID != "" {
		p, err = profile.ByID(profileID)
	} else {
		p, err = ui.SelectProfile(profile.All())
	}
	if err != nil {
		return err
	}

	targetPath, err := p.ResolveOutputPath()
	if err != nil {
		return fmt.Errorf("résolution du chemin de sortie: %w", err)
	}

	info, statErr := os.Stat(targetPath)
	fileExists := statErr == nil

	mode := p.Output.Mode
	var defaults map[string]string

	if fileExists {
		action, err := ui.PromptExistingFile(targetPath, info.ModTime())
		if err != nil {
			return err
		}
		switch action {
		case ui.ActionCancel:
			ui.PressEnterToClose("Aucune modification effectuée.")
			return nil
		case ui.ActionEdit:
			existing, err := writer.Load(targetPath, p.Output.Format)
			if err != nil {
				return err
			}
			defaults = stringifyMap(existing)
		case ui.ActionStartOver:
			defaults = nil
			mode = profile.ModeOverwrite
		}
	}

	values, err := ui.RunForm(p, defaults)
	if err != nil {
		return err
	}

	confirmed, err := ui.ConfirmSave(p, values, targetPath)
	if err != nil {
		return err
	}
	if !confirmed {
		ui.PressEnterToClose("Configuration annulée. Aucun fichier n'a été modifié.")
		return nil
	}

	var backupPath string
	if fileExists {
		backupPath, err = writer.Backup(targetPath)
		if err != nil {
			return fmt.Errorf("sauvegarde du fichier existant: %w", err)
		}
	}

	if err := writer.Write(targetPath, p.Output.Format, mode, values); err != nil {
		return fmt.Errorf("écriture du fichier de configuration: %w", err)
	}

	msg := fmt.Sprintf("Configuration enregistrée dans :\n%s", targetPath)
	if backupPath != "" {
		msg += fmt.Sprintf("\n\nUne sauvegarde de l'ancien fichier a été créée :\n%s", backupPath)
	}
	ui.PressEnterToClose(msg)
	return nil
}

func stringifyMap(in map[string]any) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		switch s := v.(type) {
		case string:
			out[k] = s
		case fmt.Stringer:
			out[k] = s.String()
		default:
			out[k] = fmt.Sprintf("%v", v)
		}
	}
	return out
}
