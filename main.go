package main

import (
	"fmt"
	"github.com/notnil/chess"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DefaultRating = 1200 // Default rating assigned to unrated players for the purpose of applying rating filters

type Config struct {
	Positions []Position
}

type Position struct {
	FEN    string
	Filter Filter
}

type Rating struct {
	One     int
	White   int
	Black   int
	Average int
}

type Filter struct {
	Rating
}

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

			var c Config
			err = yaml.Unmarshal(y, &c)
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
					for _, p := range c.Positions {
						if strings.HasPrefix(p.FEN, sfen) {
							// Apply rating filters
							if p.Filter.Rating != (Rating{}) {
								white, black := DefaultRating, DefaultRating
								if game.GetTagPair("WhiteElo") != nil {
									white, _ = strconv.Atoi(game.GetTagPair("WhiteElo").Value)
								}
								if game.GetTagPair("BlackElo") != nil {
									black, _ = strconv.Atoi(game.GetTagPair("BlackElo").Value)
								}

								if p.Filter.Rating.Average > 0 && (white+black)/2 < p.Filter.Rating.Average {
									continue
								}
								if p.Filter.Rating.One > 0 && white < p.Filter.Rating.One && black < p.Filter.Rating.One {
									continue
								}
								if p.Filter.Rating.White > 0 && white < p.Filter.Rating.White {
									continue
								}
								if p.Filter.Rating.Black > 0 && black < p.Filter.Rating.Black {
									continue
								}
							}

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
