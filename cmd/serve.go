package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kgrimes2/better-hades-bot/pkg/handler"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var (
	token string
)

func init() {
	serveCmd.Flags().StringVarP(&token, "token", "t", "", "discord bot token")
	serveCmd.MarkFlagRequired("token")

	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve up the Discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			fmt.Println("Error creating Discord session: ", err)
			return
		}

		// Add ReadyHandler
		dg.AddHandler(handler.ReadyHandler)

		// Add MessageCreateHandler
		dg.AddHandler(handler.MessageCreateHandler)

		// Open the websocket and begin listening.
		err = dg.Open()
		if err != nil {
			fmt.Println("Error opening Discord session: ", err)
		}

		// Wait here until CTRL-C or other term signal is received.
		fmt.Println("Better Hades Bot is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc

		// Cleanly close down the Discord session.
		dg.Close()
	},
}
