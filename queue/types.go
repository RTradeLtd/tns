package queue

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Various variables used by our queue package

var (
	nilTime time.Time
	// DatabaseFileAddQueue is a queue used for simple file adds
	DatabaseFileAddQueue = "dfa-queue"
	// IpfsPinQueue is a queue used for ipfs pins
	IpfsPinQueue = "ipfs-pin-queue"
	// IpfsFileQueue is a queue used for advanced file adds
	IpfsFileQueue = "ipfs-file-queue"
	// IpfsClusterPinQueue is a queue used for ipfs cluster pins
	IpfsClusterPinQueue = "ipfs-cluster-add-queue"
	// EmailSendQueue is a queue used to handle sending email messages
	EmailSendQueue = "email-send-queue"
	// IpnsEntryQueue is a queue used to handle ipns entry creation
	IpnsEntryQueue = "ipns-entry-queue"
	// IpfsKeyCreationQueue is a queue used to handle ipfs key creation
	IpfsKeyCreationQueue = "ipfs-key-creation-queue"
	// PaymentCreationQueue is a queue used to handle payment processing
	PaymentCreationQueue = "payment-creation-queue"
	// PaymentConfirmationQueue is a queue used to handle payment confirmations
	PaymentConfirmationQueue = "payment-confirmation-queue"
	// DashPaymentConfirmationQueue is a queue used to handle confirming dash payments
	DashPaymentConfirmationQueue = "dash-payment-confirmation-queue"
	// MongoUpdateQueue is a queue used to trigger mongodb updates
	MongoUpdateQueue = "mongo-update-queue"
	// ZoneCreationQueue is a queue used to handle tns zone creations
	ZoneCreationQueue = "zone-creation-queue"
	// RecordCreationQueue is a queue used to handle tns record creation
	RecordCreationQueue = "record-creation-queue"
	// AdminEmail is the email used to notify RTrade about any critical errors
	AdminEmail = "temporal.reports@rtradetechnologies.com"
	// IpfsPinFailedContent is a to-be formatted message sent on IPFS pin failures
	IpfsPinFailedContent = "Pin failed for content hash %s on IPFS network %s, for reason %s"
	// IpfsPinFailedSubject is a subject for IPFS pin failed messages
	IpfsPinFailedSubject = "IPFS Pin Failed"
	// IpfsFileFailedContent is a to be formatted message sent on ipfs add failures
	IpfsFileFailedContent = "IPFS File Add Failed for object name %s on IPFS network %s"
	// IpfsFileFailedSubject is a subject for ipfs file add fails
	IpfsFileFailedSubject = "IPFS File Add Failed"
	// IpfsPrivateNetworkUnauthorizedSubject is a subject whenever someone tries to access a bad private network
	IpfsPrivateNetworkUnauthorizedSubject = "Unauthorized access to IPFS private network"
	// IpfsInitializationFailedSubject is a subject used when connecting to ipfs fails
	IpfsInitializationFailedSubject = "Connection to IPFS failed"
	// IpnsEntryFailedSubject is a subject sent upon IPNS failures
	IpnsEntryFailedSubject = "IPNS Entry Creation Failed"
	// IpnsEntryFailedContent is the content used when sending an email for IPNS entry creation failures
	IpnsEntryFailedContent = "IPNS Entry creation failed for content hash %s using key %s for reason %s"
	// PaymentConfirmationFailedSubject is a subject used when payment confirmations fail
	PaymentConfirmationFailedSubject = "Payment Confirmation Failed"
	// PaymentConfirmationFailedContent is a content used when a payment confirmation failure occurs
	PaymentConfirmationFailedContent = "Payment failed for content hash %s with error %s"
)

// Manager is a helper struct to interact with rabbitmq
type Manager struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	Queue        *amqp.Queue
	Logger       *log.Logger
	QueueName    string
	Service      string
	ExchangeName string
}

// Queue Messages - These are used to format messages to send through rabbitmq

// IPFSKeyCreation is a message used for processing key creation
// only supported for the public IPFS network at the moment
type IPFSKeyCreation struct {
	UserName    string  `json:"user_name"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Size        int     `json:"size"`
	NetworkName string  `json:"network_name"`
	CreditCost  float64 `json:"credit_cost"`
}

// IPFSPin is a struct used when sending pin request
type IPFSPin struct {
	CID              string  `json:"cid"`
	NetworkName      string  `json:"network_name"`
	UserName         string  `json:"user_name"`
	HoldTimeInMonths int64   `json:"hold_time_in_months"`
	CreditCost       float64 `json:"credit_cost"`
}

// IPFSFile is our message for the ipfs file queue
type IPFSFile struct {
	// MinioHostIP is the ip address of the minio host this object is stored on
	MinioHostIP      string  `json:"minio_host_ip"`
	FileName         string  `json:"file_name,omitempty"`
	FileSize         int64   `json:"file_size,omitempty"`
	BucketName       string  `json:"bucket_name"`
	ObjectName       string  `json:"object_name"`
	UserName         string  `json:"user_name"`
	NetworkName      string  `json:"network_name"`
	HoldTimeInMonths string  `json:"hold_time_in_months"`
	CreditCost       float64 `json:"credit_cost"`
	Encrypted        bool    `json:"encrypted"`
}

// IPFSClusterPin is a queue message used when sending a message to the cluster to pin content
type IPFSClusterPin struct {
	CID              string  `json:"cid"`
	NetworkName      string  `json:"network_name"`
	UserName         string  `json:"user_name"`
	HoldTimeInMonths int64   `json:"hold_time_in_months"`
	CreditCost       float64 `json:"credit_cost"`
}

// DatabaseFileAdd is a struct used when sending data to rabbitmq
type DatabaseFileAdd struct {
	Hash             string  `json:"hash"`
	HoldTimeInMonths int64   `json:"hold_time_in_months"`
	UserName         string  `json:"user_name"`
	NetworkName      string  `json:"network_name"`
	CreditCost       float64 `json:"credit_cost"`
}

// IPNSUpdate is our message for the ipns update queue
type IPNSUpdate struct {
	CID         string  `json:"content_hash"`
	IPNSHash    string  `json:"ipns_hash"`
	LifeTime    string  `json:"life_time"`
	TTL         string  `json:"ttl"`
	Key         string  `json:"key"`
	Resolve     bool    `json:"resolve"`
	UserName    string  `json:"user_name"`
	NetworkName string  `json:"network_name"`
	CreditCost  float64 `json:"credit_cost"`
}

// EmailSend is a helper struct used to contained formatted content ot send as an email
type EmailSend struct {
	Subject     string   `json:"subject"`
	Content     string   `json:"content"`
	ContentType string   `json:"content_type"`
	UserNames   []string `json:"user_names"`
	Emails      []string `json:"emails,omitempty"`
}

// IPNSEntry is used to hold relevant information needed to process IPNS entry creation requests
type IPNSEntry struct {
	CID         string        `json:"cid"`
	LifeTime    time.Duration `json:"life_time"`
	TTL         time.Duration `json:"ttl"`
	Resolve     bool          `json:"resolve"`
	Key         string        `json:"key"`
	UserName    string        `json:"user_name"`
	NetworkName string        `json:"network_name"`
	CreditCost  float64       `json:"credit_cost"`
}

// PaymentCreation is for the payment creation queue
type PaymentCreation struct {
	TxHash     string `json:"tx_hash"`
	Blockchain string `json:"blockchain"`
	UserName   string `json:"user_name"`
}

// DashPaymenConfirmation is a message used to signal processing of a dash payment
type DashPaymenConfirmation struct {
	UserName         string `json:"user_name"`
	PaymentForwardID string `json:"payment_forward_id"`
	PaymentNumber    int64  `json:"payment_number"`
}

// PaymentConfirmation is a message used to confirm a payment
type PaymentConfirmation struct {
	UserName      string `json:"user_name"`
	PaymentNumber int64  `json:"payment_number"`
}

// MongoUpdate is an update used to trigger
type MongoUpdate struct {
	DatabaseName   string            `json:"database_name"`
	CollectionName string            `json:"collection_name"`
	Fields         map[string]string `json:"fields"`
}

// ZoneCreation is used for creating tns zones
type ZoneCreation struct {
	Name           string `json:"name"`
	ManagerKeyName string `json:"manager_key_name"`
	ZoneKeyName    string `json:"zone_key_name"`
	UserName       string `json:"user_name"`
}

// RecordCreation is a messaged used when creating a record
type RecordCreation struct {
	ZoneName      string                 `json:"zone_name"`
	RecordName    string                 `json:"record_name"`
	RecordKeyName string                 `json:"record_key_name"`
	MetaData      map[string]interface{} `json:"meta_data"`
	UserName      string                 `json:"user_name"`
}
