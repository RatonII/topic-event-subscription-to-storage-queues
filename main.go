package main

import (
	"flag"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid"
	aauth "github.com/Azure/go-autorest/autorest/azure/auth"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var eventSubscriptions EventSubscriptions
	var sub string
	var resourceGr string
	var sfile	*string

	subsfile := flag.String("subsfile", "", "The yaml file that contains the eventgrid subscriptions")
	subscriptionid := flag.String("subscription-id", "", "The name of the subscription for the eventgrid and the storage account")
	resourceGroup := flag.String("resource-group", "", "The name of the resource group for the eventgrid and the storage account")
	flag.Parse()
	if *subscriptionid != "" {
		sub = *subscriptionid
	} else {
		log.Fatalln("Please provide a subscription for your azure account: --subscription")
	}
	if *resourceGroup != "" {
		resourceGr = *resourceGroup
	} else {
		log.Fatalln("Please provide a  resource group for your azure account: --resource-group")
	}
	if  *subsfile != "" {
		sfile = subsfile
	}else {
		log.Fatalln("Please provide a  file with the appropriate event subscriptions: --subsfile path/to/file")
	}
	subclient := eventgrid.NewEventSubscriptionsClient(sub)

	authorizer, err := aauth.NewAuthorizerFromCLI()
	if err != nil {
		panic(err)
	}
	subclient.Authorizer = authorizer
	eventSubConf := eventSubscriptions.getConf(subsfile)
	wg.Add(len(*eventSubConf))
	for _, e := range *eventSubConf {
		go CreateUpdateEventSubscription(subclient, e.EventSubscriptionName, e.StorageAccountName, e.TopicName, sub, resourceGr, e.StorageQueueName, &wg)
	}
	wg.Wait()
}
