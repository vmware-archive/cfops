package localpubsub

type pubsubInterface interface {
	Close() error
	Subscribe(channel ...interface{}) error
	PSubscribe(channel ...interface{}) error
	Unsubscribe(channel ...interface{}) error
	PUnsubscribe(channel ...interface{}) error
	Receive() interface{}
}
