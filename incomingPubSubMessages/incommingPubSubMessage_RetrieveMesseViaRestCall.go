package incomingPubSubMessages

import (
	"FenixSCConnector/common_config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const (
	googlePubsubURL = "https://pubsub.googleapis.com/v1/projects/%s/subscriptions/%s:pull"
)

type pullRequest struct {
	MaxMessages int `json:"maxMessages"`
}

type pullResponse struct {
	ReceivedMessages []struct {
		AckID   string `json:"ackId"`
		Message struct {
			Data []byte `json:"data"`
		} `json:"message"`
	} `json:"receivedMessages"`
}

func retrivePubSubMessages(subscriptionID string, oauth2Token string) (err error) {
	url := fmt.Sprintf(googlePubsubURL, common_config.GcpProject, subscriptionID)
	body := &pullRequest{
		MaxMessages: 10, // Number of messages you want to pull
	}

	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+oauth2Token)
	req.Header.Set("Content-Type", "application/json")

	var client *http.Client
	var resp *http.Response
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "5e557f2d-1340-4140-999a-74c144815adc",
			"err": err,
		}).Error("Error making request:")

		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		common_config.Logger.WithFields(logrus.Fields{
			"ID":                "5e557f2d-1340-4140-999a-74c144815adc",
			"resp.Status":       resp.Status,
			"resp.StatusCode":   resp.StatusCode,
			"string(bodyBytes)": string(bodyBytes),
		}).Error("Non http.StatsOK was received:")

		return errors.New(resp.Status)
	}

	var response pullResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "47f5c19b-ee09-40d5-8ab4-b4494985b205",
			"err": err,
		}).Error("Error decoding response:")

		return errors.New(fmt.Sprintf("Error decoding response: %s", err.Error()))

	}

	for _, message := range response.ReceivedMessages {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":      "28ac438d-81bd-4590-bc08-bb02cf2b98af",
			"message": string(message.Message.Data),
			"err":     err,
		}).Debug("Received message")

		// Trigger TestInstruction in parallel while processing next message
		go func() {
			err = triggerProcessTestInstructionExecution(message.Message.Data)
			if err == nil {

				// Acknowledge the message
				err = acknowledgeMessage(common_config.GcpProject, subscriptionID, message.AckID, oauth2Token)

				if err != nil {

					common_config.Logger.WithFields(logrus.Fields{
						"ID":            "439247e2-27d8-4451-9ad0-51c0d60e3dd8",
						"message.AckID": message.AckID,
						"err":           err,
					}).Error("Failed to acknowledge message")

				} else {

					common_config.Logger.WithFields(logrus.Fields{
						"ID":            "54aebdf9-888d-4e0a-911c-e6cf9165acba",
						"message.AckID": message.AckID,
					}).Debug("Success in Acknowledged message")

				}
			} else {

				common_config.Logger.WithFields(logrus.Fields{
					"ID": "657fd8b2-2d9b-4158-8dbb-9f12668b94b2",
				}).Error("Failed to Process TestInstructionExecution")

			}
		}()

	}

	return err
}

type ackRequest struct {
	AckIds []string `json:"ackIds"`
}

func acknowledgeMessage(projectID string, subscriptionID string, ackID string, oauth2Token string) error {
	url := fmt.Sprintf(googlePubsubURL, projectID, subscriptionID)

	body := &ackRequest{
		AckIds: []string{ackID},
	}

	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url+":acknowledge", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+oauth2Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}
