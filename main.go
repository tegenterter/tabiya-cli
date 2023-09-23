package main

import (
	"fmt"
	"github.com/notnil/chess"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"strings"
)

// Simplify any given position in Forsyth–Edwards Notation by stripping the half move and full move numbers
func simplify(fen string) string {
	rgx := regexp.MustCompile(`(?i)(^[rnbqk1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqk1-8]{1,8}\s[wb]{1}\s[kq-]{1,4}\s[a-h1-8-]{1,2})\s\d+\s\d+$`)
	if rgx.MatchString(fen) == false {
		panic("Invalid FEN")
	}
	return rgx.FindStringSubmatch(fen)[1]
}

func main() {
	var (
		db   string
		conf string
	)

	app := &cli.App{
		Name:  "tabiya",
		Usage: "Advanced position search utility for PGN chess databases",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "pgn",
				Aliases:     []string{"p"},
				Usage:       "Path to PGN database to be scanned",
				Required:    true,
				Destination: &db,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Path to file that contains list of positions in FEN notation",
				Required:    true,
				Destination: &conf,
			},
		},
		Action: func(*cli.Context) error {
			// Read the given PGN database file
			f, err := os.Open(db)
			if err != nil {
				return cli.Exit("Could not read PGN database", 1)
			}
			defer f.Close()

			// Read the given configuration file
			y, err := os.ReadFile(conf)
			if err != nil {
				return cli.Exit("Could not read positions file", 1)
			}

			var pos []string
			err = yaml.Unmarshal(y, &pos)
			if err != nil {
				return cli.Exit("Could not parse positions file", 1)
			}

			scanner := chess.NewScanner(f)

		outer:
			for scanner.Scan() {
				game := scanner.Next()

				for _, position := range game.Positions() {
					sfen := simplify(position.String())

					// Search for the positions irrespective of move number
					for _, p := range pos {
						if strings.HasPrefix(p, sfen) {
							// Output PGN of matched game
							fmt.Println(game, "\n")
							// Skip already matches games
							continue outer
						}
					}
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
