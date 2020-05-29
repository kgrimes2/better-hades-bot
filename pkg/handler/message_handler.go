package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bwmarrin/discordgo"
	"github.com/kgrimes2/better-hades-bot/pkg/model"
)

// MessageCreateHandler is called whenever a message is sent in the channel
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if the message is "!airhorn"
	if strings.HasPrefix(m.Content, ".") {
		if m.Content == ".ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
		if strings.HasPrefix(m.Content, ".in") {
			queueSlice := strings.Split(m.Content, " ")
			if len(queueSlice) != 2 {
				s.ChannelMessageSend(m.ChannelID, "Invalid syntax. Try: .in <queue number>")
			} else {
				sess := session.Must(session.NewSessionWithOptions(session.Options{
					SharedConfigState: session.SharedConfigEnable,
				}))
				svc := dynamodb.New(sess)

				// Retrieve the queue, if it exists
				result, err := svc.GetItem(&dynamodb.GetItemInput{
					TableName: aws.String("better-hades-bot-queues"),
					Key: map[string]*dynamodb.AttributeValue{
						"queue": {
							S: aws.String(fmt.Sprintf("rs%s", queueSlice[1])),
						},
					},
				})
				if err != nil {
					fmt.Printf("Issue retrieving RS%s queue\n", queueSlice[1])
					fmt.Println(err.Error())
					return
				}

				existingQueue := model.Queue{}
				err = dynamodbattribute.UnmarshalMap(result.Item, &existingQueue)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				for _, member := range existingQueue.Members {
					if member == m.Author.Mention() {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you're already in for RS%s!", m.Author.Mention(), queueSlice[1]))
						return
					}
				}

				updatedQueue := model.Queue{}
				updatedQueue.Queue = fmt.Sprintf("rs%s", queueSlice[1])
				updatedQueue.Members = append(existingQueue.Members, m.Author.Mention())
				fmt.Println(updatedQueue)

				uqAttr, err := dynamodbattribute.MarshalMap(updatedQueue)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				updatedQueueInput := &dynamodb.PutItemInput{
					Item:      uqAttr,
					TableName: aws.String("better-hades-bot-queues"),
				}
				_, err = svc.PutItem(updatedQueueInput)
				if err != nil {
					log.Println(err.Error())
					return
				}

				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you're in for RS%s", m.Author.Mention(), queueSlice[1]))
				s.ChannelMessageSend(m.ChannelID, strings.Join(updatedQueue.Members, ","))

				if len(updatedQueue.Members) == 4 {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Time to run that RS%s: %s", queueSlice[1], strings.Join(updatedQueue.Members, ",")))

					// Need to delete the item
					deleteInput := &dynamodb.DeleteItemInput{
						Key: map[string]*dynamodb.AttributeValue{
							"queue": {
								S: aws.String(fmt.Sprintf("rs%s", queueSlice[1])),
							},
						},
						TableName: aws.String("better-hades-bot-queues"),
					}
					_, err := svc.DeleteItem(deleteInput)
					if err != nil {
						log.Println(err.Error())
						return
					}
				}
			}
		}
		if strings.HasPrefix(m.Content, ".out") {
			queueSlice := strings.Split(m.Content, " ")
			if len(queueSlice) != 2 {
				s.ChannelMessageSend(m.ChannelID, "Invalid syntax. Try: .out <queue number>")
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you're out of the RS%s queue", m.Author.Mention(), queueSlice[1]))
			}
		}
	}
}
