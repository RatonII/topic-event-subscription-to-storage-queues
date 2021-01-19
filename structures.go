package main

type EventSubscriptions []EventSubscription

type EventSubscription struct {
	EventSubscriptionName string  `yaml:"eventSubscriptionName"`
	TopicName             string  `yaml:"topicName"`
	StorageAccountName    string  `yaml:"storageAccountName"`
	StorageQueueName      *string `yaml:"storageQueueName"`
}
