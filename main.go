//go:generate goversioninfo

package main

import (
	"fmt"
	"github.com/Delta456/box-cli-maker/v2"
	"github.com/gookit/color"
	"os"
	"time"
)

func main() {

	// Start the main application loop
	for {

		Box := box.New(box.Config{Px: 4, Py: 2, Type: "Single", Color: "Green", TitlePos: "Top"})
		Box.Print("Discord Tag Scraper by top",
			"1. Input Username \n"+
				"2. Input Discord Token \n"+
				"0. Exit")

		// Get the user input
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			return
		}

		// Switch on the input
		switch input {
		case "1": // Input Username

			// Get the token
			tokenBytes, _ := os.ReadFile("token.txt")
			token := string(tokenBytes)
			if token == "" {
				color.Error.Println("Please input your discord token")
				continue
			}

			// Check if the token is valid
			statusCode, _ := CheckToken(token)
			if statusCode != 200 {
				color.Error.Println("Invalid Discord Token")
				continue
			}

			fmt.Print("Input Username: ")
			var username string
			_, err = fmt.Scanln(&username)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				continue
			}

			// Create the output file
			filename := username + "-" + time.Now().Format("2006-01-02_15-04-05") + ".txt"
			file, _ := os.Create(filename)

			// Check the discriminators
			for i := 0001; i < 9999; i++ {
				time.Sleep(time.Millisecond * 1000)
				b, err := CheckDiscriminator(token, username, i)
				if err != nil {
					fmt.Println("Error: " + err.Error())
					continue
				}
				if b {
					color.Success.Printf("[%s#%d] Discriminator %d is valid \n", username, i, i)
					// Output the tag to a file and use the current timestamp as the filename

					if err != nil {
						fmt.Println("Error: " + err.Error())
						continue
					}
					defer file.Close()
					// Write the tag to the file and every token gets a new line
					file.WriteString(fmt.Sprintf("%s#%d\n", username, i))

				} else {
					color.Red.Printf("[%s#%d] Discriminator %d is invalid \n", username, i, i)
				}
			}

		case "2": // Input Discord Token
			fmt.Print("Input Discord Token: ")
			var discordToken string
			_, err := fmt.Scanln(&discordToken)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				continue
			}

			// Check if the token is valid
			statusCode, _ := CheckToken(discordToken)
			if statusCode != 200 {
				color.Error.Println("Invalid Discord Token")
				continue
			}

			err = os.WriteFile("token.txt", []byte(discordToken), 0644)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				continue
			}
		case "0": // Exit
			os.Exit(0)

		default:
			color.Error.Println("Invalid input")
			continue
		}

	}

}
