package main

import (
	"bufio"
	"errors"
	"log"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"
)

func main() {
	/* Handle SIGINT properly */
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Print("\n")
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		/* Retrieve the primary info */
		curruser, currdir := getUserandDir()
		fmt.Print(curruser, "@", currdir, "> ")

		/* Read the keyboard input */
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		/* Handle the execution of the input */
		err = execInput(input)
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

func getUserandDir() (string, string) {
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

	return user.Username, dir
}
