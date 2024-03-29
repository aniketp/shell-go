package main

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"

	"github.com/eiannone/keyboard"
)

func main() {
	/* Handle SIGINT properly */
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		/* Silently exit the program */
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		/* Retrieve the primary info */
		curruser, currhost, currdir := getUserHostandDir()
		fmt.Print(curruser, "@", currhost, ":", currdir, "$ ")

		/* Get the first key for further development */
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			log.Fatal(err)
		} else if key == keyboard.KeyArrowUp {
			os.Exit(1)
		}
		fmt.Print(string(char))

		/* Read the keyboard input */
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		/* Append the first key with rest of the command */
		list := []string{string(char), input}
		var str bytes.Buffer
		for _, value := range list {
			str.WriteString(value)
		}

		/* Handle the execution of the input */
		err = execInput(str.String())
		if err != nil {
			fmt.Println(err)
		}

	}
}

func execInput(input string) error {
	/* Remove the newline character */
	input = strings.TrimSuffix(input, "\n")

	/* Split the input to retrieve the arguments */
	args := strings.Split(input, " ")

	/* Prepare the command to execute */
	cmd := exec.Command(args[0], args[1:]...)

	/* Check for any built-in commands */
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return errors.New("Path required")
		} else {
			err := os.Chdir(args[1])
			if err != nil {
				return err
			}
			return nil
		}
	case "exit":
		os.Exit(0)
	}

	/* Set the correct output device */
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	/* Execute the command and save its output */
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func getUserHostandDir() (string, string, string) {
	/* Get the logged in user */
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	/* Get the current working directory */
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	/* Get the Machine hostname */
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	return user.Username, name, dir
}
