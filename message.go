package twiliogo

import (
	"encoding/json"
	"net/url"
)

type Message struct {
	Sid                 string `json:"sid"`
	DateCreated         string `json:"date_created"`
	DateUpdated         string `json:"date_updated"`
	DateSent            string `json:"date_sent"`
	AccountSid          string `json:"account_sid"`
	From                string `json:"from"`
	MessagingServiceSid string `json:"messaging_service_sid"`
	To                  string `json:"to"`
	Body                string `json:"body"`
	NumSegments         string `json:"num_segments"`
	Status              string `json:"status"`
	Direction           string `json:"direction"`
	Price               string `json:"price"`
	PriceUnit           string `json:"price_unit"`
	ApiVersion          string `json:"api_version"`
	Uri                 string `json:"uri"`
}

func NewMessage(client Client, from string, to string, content ...Optional) (*Message, error) {
	return sendMessage(client, from, "", to, content...)
}

func NewMessageFromService(client Client, fromServiceSid string, to string, content ...Optional) (*Message, error) {
	return sendMessage(client, "", fromServiceSid, to, content...)
}

func sendMessage(client Client, fromNumber, fromMessagingServiceSid, to string, content ...Optional) (*Message, error) {
	var message *Message

	params := url.Values{}
	if fromNumber != "" {
		params.Set("From", fromNumber)
	} else if fromMessagingServiceSid != "" {
		params.Set("MessagingServiceSid", fromMessagingServiceSid)
	} else {
		return nil, Error{"Either From phone number or MessageServiceSid must be provided."}
	}

	params.Set("To", to)

	for _, optional := range content {
		param, value := optional.GetParam()

		if param != "Body" && param != "MediaUrl" && param != "StatusCallback" && param != "ApplicationSid" && param != "MessagingServiceSid" {
			return nil, Error{"Only allowed params are Body, MediaUrl, StatusCallback, ApplicationSid, MessagingServiceSid"}
		}

		params.Set(param, value)
	}

	if params.Get("Body") == "" && params.Get("MediaUrl") == "" {
		return nil, Error{"Must have at least a Body or MediaUrl"}
	}

	if len(params.Get("Body")) > 1560 {
		params.Set("Body", params.Get("Body")[:1560])
	}

	res, err := client.post(params, "/Messages.json")

	if err != nil {
		return message, err
	}

	message = new(Message)
	err = json.Unmarshal(res, message)

	return message, err
}

func GetMessage(client Client, sid string) (*Message, error) {
	var message *Message

	res, err := client.get(url.Values{}, "/Messages/"+sid+".json")

	if err != nil {
		return nil, err
	}

	message = new(Message)
	err = json.Unmarshal(res, message)

	return message, err
}
