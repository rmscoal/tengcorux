package attribute

import (
	"reflect"
	"testing"
)

func TestAttribute_KeyValuePair(t *testing.T) {
	got := KeyValuePair("some_key", "some_val")
	want := KeyValue{Key: Key("some_key"), Value: any("some_value")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}

	got = KeyValuePair("some_other_key", 1)
	want = KeyValue{Key: Key("some_other_key"), Value: any(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPUrl(t *testing.T) {
	got := HTTPUrl("some_url")
	want := KeyValue{Key: HTTPUrlKey, Value: any("some_url")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPUrlQuery(t *testing.T) {
	got := HTTPUrlQuery("{id: 1}")
	want := KeyValue{Key: HTTPUrlQueryKey, Value: any("{id: 1}")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPUrlPath(t *testing.T) {
	got := HTTPUrlPath("/path")
	want := KeyValue{Key: HTTPUrlPathKey, Value: any("/path")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPRequestBody(t *testing.T) {
	got := HTTPRequestBody("{\"loanId\": \"some_loan_id\"")
	want := KeyValue{Key: HTTPRequestBodyKey, Value: any("{\"loanId\": \"some_loan_id\"")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPRequestMethod(t *testing.T) {
	got := HTTPRequestMethod("POST")
	want := KeyValue{Key: HTTPRequestMethodKey, Value: any("POST")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPRequestID(t *testing.T) {
	got := HTTPRequestID("some_uuid")
	want := KeyValue{Key: HTTPRequestIDKey, Value: any("some_uuid")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPResponseStatus(t *testing.T) {
	got := HTTPResponseStatus(200)
	want := KeyValue{Key: HTTPResponseStatusKey, Value: any(200)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_HTTPResponseBody(t *testing.T) {
	got := HTTPResponseBody("{\"message\": \"Success\"}")
	want := KeyValue{Key: HTTPResponseBodyKey, Value: any("{\"message\": \"Success\"}")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_DBSystem(t *testing.T) {
	got := DBSystem("mysql")
	want := KeyValue{Key: DBSystemKey, Value: any("mysql")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_DBName(t *testing.T) {
	got := DBName("some_db_name")
	want := KeyValue{Key: DBNameKey, Value: any("some_db_name")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_DBInstanceID(t *testing.T) {
	got := DBInstanceID("some_value")
	want := KeyValue{Key: DBInstanceIDKey, Value: any("some_value")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_DBStatement(t *testing.T) {
	got := DBStatement("SELECT * FROM test WHERE id = ? AND status = ?")
	want := KeyValue{Key: DBStatementKey, Value: any("SELECT * FROM test WHERE id = ? AND status = ?")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_DBOperation(t *testing.T) {
	got := DBOperation("BEGIN")
	want := KeyValue{Key: DBOperationKey, Value: any("BEGIN")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_DBTable(t *testing.T) {
	got := DBTable("some_table")
	want := KeyValue{Key: DBTableKey, Value: any("some_table")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_MQSystem(t *testing.T) {
	got := MQSystem("Kafka")
	want := KeyValue{Key: MQSystemKey, Value: any("Kafka")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_MQInstanceID(t *testing.T) {
	got := MQInstanceID("localhost:6379")
	want := KeyValue{Key: MQInstanceIDKey, Value: any("localhost:6379")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_MQTopic(t *testing.T) {
	got := MQTopic("some_topic")
	want := KeyValue{Key: MQTopicKey, Value: any("some_topic")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_MQSubscriber(t *testing.T) {
	got := MQSubscriber("some_subscriber")
	want := KeyValue{Key: MQSubscriberKey, Value: any("some_subscriber")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_MQConsumerGroup(t *testing.T) {
	got := MQConsumerGroup("group:hello_world")
	want := KeyValue{Key: MQConsumerGroupKey, Value: any("group:hello_world")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_MQMessageBody(t *testing.T) {
	got := MQMessageBody("{\"content\":\"hello world\"}")
	want := KeyValue{Key: MQMessageBodyKey, Value: any("{\"content\":\"hello world\"}")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}

func TestAttribute_MQMessageID(t *testing.T) {
	got := MQMessageID(3432804932)
	want := KeyValue{Key: MQMessageIDKey, Value: any(3432804932)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", got, want)
	}
}
