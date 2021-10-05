package main

import (
	"errors"
	"fmt"
	. "github.com/jonhadfield/packer-azure-image-version"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

var version, versionOutput, tag, sha, buildDate string

const strMoreThanOneTrue = "more than one is true"
const strNoneTrue = "none are true"

func main() {
	if tag != "" && buildDate != "" {
		versionOutput = fmt.Sprintf("[%s-%s] %s UTC", tag, sha, buildDate)
	} else {
		versionOutput = version
	}

	// cwd, _ := os.Getwd()

	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Name = "packer azure image version"
	app.Version = versionOutput
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		{
			Name:  "Jon Hadfield",
			Email: "jon@lessknown.co.uk",
		},
	}
	app.HelpName = "-"
	app.Usage = "Packer Azure Image Version"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "subscription-id",
			Usage:    "specify the subscription id containing the policies",
			EnvVars:  []string{"AZURE_SUBSCRIPTION_ID"},
			Aliases:  []string{"s", "subscription"},
			Required: false,
		},
		&cli.BoolFlag{Name: "quiet", Usage: "suppress output"},
	}
	app.Commands = []*cli.Command{
		{
			// without flags, get will return all versions
			Name:  "get",
			Usage: "get image versions",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "latest", Aliases: []string{"l"}},
				&cli.BoolFlag{Name: "oldest", Aliases: []string{"o"}},
				&cli.BoolFlag{Name: "inc-major"},
				&cli.BoolFlag{Name: "inc-minor"},
				&cli.BoolFlag{Name: "inc-patch"},
			},
			Action: func(c *cli.Context) error {
				input := c.Args().Slice()
				switch len(input) {
				case 0:
					return fmt.Errorf("image definition id is required")
				case 1:
					if idi := ParseImageDefinitionID(input[0]); idi.ImageName == "" {
						_ = cli.ShowSubcommandHelp(c)

						return fmt.Errorf("invalid image definition id")
					}

					return GetImageVersions(GetImageVersionsInput{
						SubscriptionID:    c.String("subscription-id"),
						ImageDefinitionID: input[0],
						Latest:            c.Bool("latest"),
						Oldest:            c.Bool("oldest"),
						IncMajor:          c.Bool("inc-major"),
						IncMinor:          c.Bool("inc-minor"),
						IncPatch:          c.Bool("inc-patch"),
					})

				default:
					return fmt.Errorf("only one image definition id is expected")
				}
			},
		},
		{
			Name:  "set",
			Usage: "set packer image gallery destination",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "inc-major"},
				&cli.BoolFlag{Name: "inc-minor"},
				&cli.BoolFlag{Name: "inc-patch"},
				&cli.BoolFlag{Name: "unattended"},
				&cli.BoolFlag{Name: "cli-auth"},
			},
			Action: func(c *cli.Context) error {
				input := c.Args().Slice()
				if len(input) == 0 {
					return errors.New("at least one path is required")
				}

				if err := checkOneTrue(c.Bool("inc-major"),
					c.Bool("inc-minor"),
					c.Bool("inc-patch")); err != nil {

					if err.Error() == strNoneTrue {
						_ = cli.ShowSubcommandHelp(c)

						return fmt.Errorf("increment option required")
					}
					_ = cli.ShowSubcommandHelp(c)

					return fmt.Errorf("only one increment option can be specified")
				}

				return SetImageVersions(SetImageVersionInput{
					Paths:      input,
					IncMajor:   c.Bool("inc-major"),
					IncMinor:   c.Bool("inc-minor"),
					IncPatch:   c.Bool("inc-patch"),
					Unattended: c.Bool("unattended"),
					CLIAuth:    c.Bool("cli-auth"),
					Quiet:      c.Bool("quiet"),
				})
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println()
		fmt.Printf("error: %v\n\n", err)
		os.Exit(1)
	}
}

func checkOneTrue(i ...bool) error {
	var foundTrue bool
	for x := range i {
		if i[x] {
			if foundTrue {
				return fmt.Errorf(strMoreThanOneTrue)
			}

			foundTrue = true
		}
	}

	if !foundTrue {
		return fmt.Errorf(strNoneTrue)
	}

	return nil
}
