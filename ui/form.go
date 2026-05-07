package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"

	"github.com/thomas/gotrunk/profile"
)

type ExistingAction int

const (
	ActionEdit ExistingAction = iota
	ActionStartOver
	ActionCancel
)

func PromptExistingFile(path string, modTime time.Time) (ExistingAction, error) {
	desc := fmt.Sprintf(
		"Un fichier de configuration existe déjà à :\n  %s\n\nDernière modification : %s\n\nQue souhaitez-vous faire ?",
		path,
		modTime.Format("02/01/2006 à 15:04"),
	)
	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Configuration existante détectée").
				Description(desc).
				Options(
					huh.NewOption("Modifier les identifiants existants", "edit"),
					huh.NewOption("Repartir de zéro (effacer les valeurs actuelles)", "start_over"),
					huh.NewOption("Annuler et quitter", "cancel"),
				).
				Value(&choice),
		),
	).WithTheme(huh.ThemeCharm())
	if err := form.Run(); err != nil {
		return ActionCancel, err
	}
	switch choice {
	case "edit":
		return ActionEdit, nil
	case "start_over":
		return ActionStartOver, nil
	default:
		return ActionCancel, nil
	}
}

func SelectProfile(profiles []profile.Profile) (profile.Profile, error) {
	if len(profiles) == 1 {
		return profiles[0], nil
	}
	options := make([]huh.Option[string], 0, len(profiles))
	for _, p := range profiles {
		options = append(options, huh.NewOption(p.DisplayName, p.ID))
	}
	var chosen string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Quel plugin souhaitez-vous configurer ?").
				Options(options...).
				Value(&chosen),
		),
	).WithTheme(huh.ThemeCharm())
	if err := form.Run(); err != nil {
		return profile.Profile{}, err
	}
	return profile.ByID(chosen)
}

func RunForm(p profile.Profile, defaults map[string]string) (map[string]string, error) {
	holders := make([]*string, len(p.Fields))
	for i, f := range p.Fields {
		v := f.Default
		if d, ok := defaults[f.Key]; ok && d != "" {
			v = d
		}
		holders[i] = &v
	}

	fields := make([]huh.Field, 0, len(p.Fields))
	for i := range p.Fields {
		f := p.Fields[i]
		input := huh.NewInput().
			Title(f.Label).
			Description(f.Help).
			Value(holders[i]).
			Validate(buildValidator(f))
		if f.Secret {
			input = input.EchoMode(huh.EchoModePassword)
		}
		fields = append(fields, input)
	}

	form := huh.NewForm(huh.NewGroup(fields...)).WithTheme(huh.ThemeCharm())
	if err := form.Run(); err != nil {
		return nil, err
	}

	values := make(map[string]string, len(p.Fields))
	for i, f := range p.Fields {
		values[f.Key] = *holders[i]
	}
	return values, nil
}

func ConfirmSave(p profile.Profile, values map[string]string, targetPath string) (bool, error) {
	summary := fmt.Sprintf("Enregistrer la configuration dans :\n  %s\n\n", targetPath)
	for _, f := range p.Fields {
		v := values[f.Key]
		if f.Secret {
			v = mask(v)
		}
		summary += fmt.Sprintf("  • %s : %s\n", f.Label, v)
	}

	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirmation").
				Description(summary).
				Affirmative("Enregistrer").
				Negative("Annuler").
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return false, err
	}
	return confirmed, nil
}

func PressEnterToClose(message string) {
	var ack bool
	_ = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Terminé").
				Description(message).
				Affirmative("Fermer").
				Negative("").
				Value(&ack),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

func buildValidator(f profile.Field) func(string) error {
	return func(s string) error {
		if f.Required && len(s) == 0 {
			return fmt.Errorf("ce champ est requis")
		}
		if f.Validate != nil {
			if err := f.Validate(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func mask(s string) string {
	if len(s) <= 4 {
		return "••••"
	}
	return s[:2] + "••••" + s[len(s)-2:]
}
