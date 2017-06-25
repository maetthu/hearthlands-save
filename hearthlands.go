package main

import (
	"flag"
	"fmt"
	"hearthlands"
	"os"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] savegame.hls\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	gold := flag.Int("gold", -1, "Set gold to specified amount.")
	population := flag.Int("population", -1, "Set population to specified amount.")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	filename := flag.Arg(0)

	_, err := os.Stat(filename)

	if len(filename) == 0 || os.IsNotExist(err) {
		flag.Usage()
		os.Exit(2)
	}

	game, err := hearthlands.OpenSaveFile(filename)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Printf("Gold: %[1]v (Hex: %[1]x)\n", game.Gold())
	fmt.Printf("Population: %[1]v (Hex: %[1]x)\n", game.Population())

	// just display current values if nothing set
	if *gold == -1 && *population == -1 {
		os.Exit(0)
	}

	if *gold != -1 {
		game.SetGold(*gold)
	}

	if *population != -1 {
		game.SetPopulation(*population)
	}

	if game.Save(filename) != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println("Values set to:")
	fmt.Printf("Gold: %[1]v (Hex: %[1]x)\n", game.Gold())
	fmt.Printf("Population: %[1]v (Hex: %[1]x)\n", game.Population())

	/*for i, b := range game.Entry {
		fmt.Printf("%02x", b)
		if i > 0 && (i-1)%16 == 0 {
			fmt.Println("")
		}
	}*/

}
