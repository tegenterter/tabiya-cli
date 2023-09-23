## Tabiya

CLI tool that allows you to efficiently search large chess databases in [PGN format](https://en.wikipedia.org/wiki/Portable_Game_Notation) for games that match specific positions, and apply complex rating filters in addition. 

## Configuration

You'll need to set up one [YAML](https://en.wikipedia.org/wiki/YAML) configuration file that lists the positions to search for, adhering to the following syntax:

```yaml
positions:
  - fen: "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"
  - fen: "rnbqkb1r/1p2pppp/p2p1n2/8/3NP3/2N5/PPP2PPP/R1BQKB1R w KQkq - 0 6"
```

Positions are expressed using [Forsyth–Edwards Notation](https://en.wikipedia.org/wiki/Forsyth-Edwards_Notation). One convenient way to obtain these notations through [the analysis board on Lichess](https://lichess.org/analysis).

Note that move numbers are irrelevant to how positions are matched, thus supporting all possible move transpositions.

### Filters

In addition to searching for particular positions, games can also be filtered out based on the minimum specified rating of the players. For each position you can optionally specify any combination of the following rating filters:

```yaml
positions:
  - fen: "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"
    filter:
      rating:
        one: 2100
        average: 2200
        white: 2100
        black: 2300
```

You can choose to apply just one of the above filters. Two or more filters always work cumulatively.

Note that if no rating is specified in the database for any player, then the default rating of `1200` is applied for the purpose of applying rating based filters.

## Usage

Once you have a configuration set up, use it to search any PGN database (such as an issue of [The Week in Chess](https://theweekinchess.com/twic)) you have stored locally like so:

```bash
tabiya --pgn /path/to/games.pgn --config /path/to/config.yaml
```

Games that match the positions and possible rating filters listed in your configuration will be written as standard output, in the same PGN format as they appear in the source database.