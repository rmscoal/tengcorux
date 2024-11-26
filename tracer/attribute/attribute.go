package attribute

type Key string

type KeyValue struct {
	Key   Key
	Value any
}

func KeyValuePair(key string, val any) KeyValue {
	return Key(key).Val(val)
}

func (k Key) Val(val any) KeyValue {
	return KeyValue{
		Key:   k,
		Value: val,
	}
}

const (
	// HTTPUrlKey is the Key conforming to the "url.full" semantics.
	HTTPUrlKey = Key("url.full")
	// HTTPUrlQueryKey is the Key conforming to the "url.query" semantics.
	HTTPUrlQueryKey = Key("url.query")
	// HTTPUrlPathKey is the Key conforming to the "url.query" semantics.
	HTTPUrlPathKey = Key("url.path")
)

func HTTPUrl(val any) KeyValue {
	return HTTPUrlKey.Val(val)
}

func HTTPUrlQuery(val any) KeyValue {
	return HTTPUrlQueryKey.Val(val)
}

func HTTPUrlPath(val any) KeyValue {
	return HTTPUrlPathKey.Val(val)
}

const (
	// HTTPRequestBodyKey is the Key conforming to the "http.request.body" semantics.
	HTTPRequestBodyKey = Key("http.request.body")
	// HTTPRequestMethodKey is the Key conforming to the "http.request.method" semantics.
	HTTPRequestMethodKey = Key("http.request.method")
	// HTTPRequestIDKey is the Key conforming to the "http.request.id" semantics.
	HTTPRequestIDKey = Key("http.request.id")
)

func HTTPRequestBody(val any) KeyValue {
	return HTTPRequestBodyKey.Val(val)
}

func HTTPRequestMethod(val any) KeyValue {
	return HTTPRequestMethodKey.Val(val)
}

func HTTPRequestID(val any) KeyValue {
	return HTTPRequestIDKey.Val(val)
}

const (
	// HTTPResponseStatusKey is the Key conforming to the "http.response.status" semantics.
	HTTPResponseStatusKey = Key("http.response.status")
	// HTTPResponseBodyKey is the Key conforming to the "http.response.body" semantics.
	HTTPResponseBodyKey = Key("http.response.body")
)

func HTTPResponseStatus(val any) KeyValue {
	return HTTPResponseStatusKey.Val(val)
}

func HTTPResponseBody(val any) KeyValue {
	return HTTPResponseBodyKey.Val(val)
}

const (
	// DBSystemKey is the Key conforming to the "db.system" semantics.
	DBSystemKey = Key("db.system")
	// DBNameKey is the Key conforming to the "db.name" semantics.
	DBNameKey = Key("db.name")
	// DBInstanceIDKey is the Key conforming to the "db.instance.id" semantics.
	DBInstanceIDKey = Key("db.instance.id")
	// DBStatementKey is the Key conforming to the "db.statement" semantics.
	DBStatementKey = Key("db.statement")
	// DBOperationKey is the Key conforming to the "db.operation" semantics.
	DBOperationKey = Key("db.operation")
	// DBTableKey is the Key conforming to the "db.table" semantics.
	DBTableKey = Key("db.table")
)

func DBSystem(val any) KeyValue {
	return DBSystemKey.Val(val)
}

func DBName(val any) KeyValue {
	return DBNameKey.Val(val)
}

func DBInstanceID(val any) KeyValue {
	return DBInstanceIDKey.Val(val)
}

func DBStatement(val any) KeyValue {
	return DBStatementKey.Val(val)
}

func DBOperation(val any) KeyValue {
	return DBOperationKey.Val(val)
}

func DBTable(val any) KeyValue {
	return DBTableKey.Val(val)
}

const (
	// MQSystemKey is the Key conforming to the "mq.system" semantics.
	MQSystemKey = Key("mq.system")
	// MQInstanceIDKey is the Key conforming to the "mq.instance.id" semantics can be
	// either used with the message queue server address or the instance ID.
	MQInstanceIDKey = Key("mq.instance.id")
	// MQTopicKey is the Key conforming to the "mq.topic" semantics.
	MQTopicKey = Key("mq.topic")
	// MQSubscriberKey is the Key conforming to the "mq.subscriber" semantics.
	MQSubscriberKey = Key("mq.subscriber")
	// MQConsumerGroupKey is the Key conforming to the "mq.consumer.group" semantics.
	MQConsumerGroupKey = Key("mq.consumer.group")
	// MQMessageBodyKey is the Key conforming to the "mq.message.body" semantics.
	MQMessageBodyKey = Key("mq.message.body")
	// MQMessageIDKey is the Key conforming to the "mq.message.id" semantics.
	MQMessageIDKey = Key("mq.message.id")
	// NOTE: Perhaps extends to each system like otel's semconv
)

func MQSystem(val any) KeyValue {
	return MQSystemKey.Val(val)
}

func MQInstanceID(val any) KeyValue {
	return MQInstanceIDKey.Val(val)
}

func MQTopic(val any) KeyValue {
	return MQTopicKey.Val(val)
}

func MQSubscriber(val any) KeyValue {
	return MQSubscriberKey.Val(val)
}

func MQConsumerGroup(val any) KeyValue {
	return MQConsumerGroupKey.Val(val)
}

func MQMessageBody(val any) KeyValue {
	return MQMessageBodyKey.Val(val)
}

func MQMessageID(val any) KeyValue {
	return MQMessageIDKey.Val(val)
}
