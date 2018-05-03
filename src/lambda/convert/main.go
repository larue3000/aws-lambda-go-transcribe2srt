package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	jobname string
)

func Handler(ctx context.Context, eventinfo interface{}) (interface{}, error) {

	//Marshal the eventinfo
	data, _ := json.Marshal(eventinfo)
	streventinfo := string(data)

	//Lets log the request for reference
	log.Printf("-----------------EVENT INFO-------------------")
	log.Printf("%+v", streventinfo)
	log.Printf("-----------------EVENT INFO-------------------")

	//Lets have a look at the request
	//Alex Request
	if (strings.Contains(streventinfo, "LaunchRequest")) || (strings.Contains(streventinfo, "IntentRequest")) || (strings.Contains(streventinfo, "SessionEndedRequest")) {
		var request AlexaRequest
		err := json.Unmarshal(data, &request)
		log.Printf("Alexa request received")
		if err != nil {
			log.Printf("Problem encountered unmarshalling request to Alex Response struct")
		}

		var response AlexaResponse
		var arrdialog []DialogDelegate
		var dialog DialogDelegate

		dialog.Type = "Dialog.Delegate"
		dialog.UpdatedIntent = nil
		//dialog.UpdatedIntent.Name = "String"
		//dialog.UpdatedIntent.ConfirmationStatus = "NONE"
		//dialog.UpdatedIntent.Slots.String.Name = "string"
		//dialog.UpdatedIntent.Slots.String.Value = "string"
		//dialog.UpdatedIntent.Slots.String.ConfirmationStatus = "NONE"

		arrdialog = append(arrdialog, dialog)

		response.Version = "1.0"
		response.Response.OutputSpeech = nil
		response.Response.Card = nil
		response.Response.Reprompt = nil
		response.Response.Directives = &arrdialog

		strResponse, _ := json.Marshal(response)
		log.Printf("-----------------LAMBDA RESPONSE-------------------")
		log.Printf("%+v", string(strResponse))
		log.Printf("-----------------LAMBDA REPSONSE-------------------")
		return response, nil
		//jobname = "01524-63806-28098-90622"
		//subtitles, converterror := Convert(jobname)

	}
	// API Gateway Request
	if (strings.Contains(streventinfo, "httpMethod")) && (strings.Contains(streventinfo, "headers")) {
		var request events.APIGatewayProxyRequest
		err := json.Unmarshal(data, &request)
		if len(request.Body) < 1 {
			log.Printf("No content in body")
			return "", err
		}
		jobname = request.Body
		subtitles, converterror := Convert(jobname)
		var response events.APIGatewayProxyResponse
		if converterror != nil {
			response.Body = "Server error"
			response.StatusCode = 400
		} else {
			response.Body = subtitles
			response.StatusCode = 200
		}

		return response, nil
		//Lex Request
	} else if (strings.Contains(streventinfo, "currentIntent")) && (strings.Contains(streventinfo, "userId")) {
		var lexrq events.LexEvent
		err := json.Unmarshal(data, &lexrq)
		if err != nil {

		}
		jobname = lexrq.InputTranscript
		subtitles, converterror := Convert(jobname)
		var response LexResponse

		if converterror != nil {
			response.DialogAction.Type = "Close"
			response.DialogAction.FulfillmentState = "Fulfilled"
			response.DialogAction.Message.ContentType = "PlainText"
			response.DialogAction.Message.Content = "Error!!!!!"
			return response, converterror
		} else {
			response.DialogAction.Type = "Close"
			response.DialogAction.FulfillmentState = "Fulfilled"
			response.DialogAction.Message.ContentType = "PlainText"
			response.DialogAction.Message.Content = subtitles
			return response, nil
		}

	} else {
		unsupported := errors.New("Unsupported service")
		return nil, unsupported
	}
}

func main() {
	lambda.Start(Handler)
}
