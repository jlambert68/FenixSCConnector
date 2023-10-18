package incomingPubSubMessages

import "FenixSCConnector/common_config"

// Create the PubSub-topic from Domain-Uuid
func generatePubSubTopicNameForExecutionStatusUpdates() (statusExecutionTopic string) {

	var pubSubTopicBase string
	pubSubTopicBase = common_config.TestInstructionExecutionPubSubTopicBase

	var testerGuiApplicationUuid string
	testerGuiApplicationUuid = common_config.ThisDomainsUuid

	// Get the first 8 characters from TesterGui-ApplicationUuid
	var shortedAppUuid string
	shortedAppUuid = testerGuiApplicationUuid[0:8]

	// Build PubSub-topic
	statusExecutionTopic = pubSubTopicBase + "-" + shortedAppUuid

	return statusExecutionTopic
}

// Creates a Topic-Subscription-Name
func generatePubSubTopicSubscriptionNameForExecutionStatusUpdates() (topicSubscriptionName string) {

	const topicSubscriptionPostfix string = "-sub"

	// Get Topic-name
	var topicID string
	topicID = generatePubSubTopicNameForExecutionStatusUpdates()

	// Create the Topic-Subscription-name
	topicSubscriptionName = topicID + topicSubscriptionPostfix

	return topicSubscriptionName
}
