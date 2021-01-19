package main

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/mgmt/eventgrid"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"sync"
)

func CreateUpdateEventSubscription(client eventgrid.EventSubscriptionsClient,
	eventSubscriptionName, storageAccountName,eventgridName, topicName, subscription, resourceGroup string,
	storageQueueName *string, wg *sync.WaitGroup) {

	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/domains/%s/topics/%s", subscription, resourceGroup,eventgridName, topicName)
	storageAccountID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s", subscription, resourceGroup, storageAccountName)
	_, err := client.CreateOrUpdate(context.Background(), scope, eventSubscriptionName, eventgrid.EventSubscription{
		EventSubscriptionProperties: &eventgrid.EventSubscriptionProperties{
			Destination: eventgrid.StorageQueueEventSubscriptionDestination{
				StorageQueueEventSubscriptionDestinationProperties: &eventgrid.StorageQueueEventSubscriptionDestinationProperties{
					ResourceID: &storageAccountID,
					QueueName:  storageQueueName,
				},
				EndpointType: eventgrid.EndpointTypeStorageQueue,
			},
			EventDeliverySchema: eventgrid.EventGridSchema,
		},
	})
	defer wg.Done()
	if err != nil {
		panic(err)
	}
	fmt.Printf("EventSubscription that links topic %s with storage queue %s on the storage account %s was created or updated", topicName, storageQueueName, storageAccountName)
}

func (s *EventSubscriptions) getConf(eventSubscriptionsFile *string) *EventSubscriptions {

	yamlFile, err := ioutil.ReadFile(*eventSubscriptionsFile)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return s
}
