package skywalking

// SkyWalking can automatically categorized a type of span depending
// on the span's component layer.
//
// The complete lists available here:
// https://github.com/apache/skywalking/blob/master/oap-server/server-starter/src/main/resources/component-libraries.yml
type ComponentLibrary int32

const (
	Unknown      ComponentLibrary = 0
	GoRedis      ComponentLibrary = 7
	PostgreSQL   ComponentLibrary = 22
	GoKafka      ComponentLibrary = 27
	RabbitMQ     ComponentLibrary = 51
	GoHttpServer ComponentLibrary = 5004
	GoMysql      ComponentLibrary = 5012
)

func (c ComponentLibrary) AsInt32() int32 {
	return int32(c)
}
