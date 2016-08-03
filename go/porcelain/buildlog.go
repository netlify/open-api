package porcelain

import (
	"errors"
	"fmt"

	logContext "github.com/docker/distribution/context"
	io "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	"github.com/netlify/open-api/go/porcelain/context"
)

const buildRequestHandle = "buildlog:request"
const buildErrorHandle = "buildlog:error"
const buildCompleteHandle = "buildlog:complete"
const buildLogFmtHandle = "buildlog:%s"

// BuildMessage represents a message coming from the API
type BuildMessage struct {
	ErrorMessage string
	Completed    bool
	Message      string
	SiteID       string
	DeployID     string
	BuildID      string
	Data         map[string]interface{}
	Time         string
}

// StreamBuildLog will connect to the websocket endpoint in the API and then
// push incoming messages back to the client. In the event of an error (e.g. the
// log doesn't exist) the 'ErrorMessage' field will be set. On completion the
// 'Completed' flag will be set. On error or completion the channel will be closed
func (n *Netlify) StreamBuildLog(ctx context.Context, buildID string) (<-chan *BuildMessage, chan<- bool, error) {
	msgChan := make(chan *BuildMessage)
	shutdownChan := make(chan bool)

	if buildID == "" {
		close(msgChan)
		close(shutdownChan)
		return msgChan, shutdownChan, errors.New("Must provide a build ID")
	}

	endpoint := fmt.Sprintf("%s", n.streamingURL)
	l := logContext.GetLoggerWithFields(ctx, context.Fields{
		"build_id": buildID,
		"endpoint": endpoint,
	})
	l.Debug("Connecting to streaming endpoint")

	// Connect to endpoint
	conn, err := io.Dial(endpoint, transport.GetDefaultWebsocketTransport())
	if err != nil {
		close(msgChan)
		close(shutdownChan)
		return msgChan, shutdownChan, err
	}

	l.Debug("Connected")

	// start a shutdown listener
	go listenForShutdown(conn, shutdownChan, msgChan)

	// now bind up to the different topics
	buildLogHandle := fmt.Sprintf(buildLogFmtHandle, buildID)
	conn.On(buildLogHandle, func(c *io.Channel, m logMessage) {
		msgChan <- m.convert()
	})
	conn.On(buildCompleteHandle, func(c *io.Channel, m completeMessage) {
		msgChan <- m.convert()
		shutdownChan <- true
	})
	conn.On(buildErrorHandle, func(c *io.Channel, m errorMessage) {
		msgChan <- m.convert()
		shutdownChan <- true
	})

	l.Debug("Bound different handlers - sending request")
	err = conn.Emit(buildRequestHandle, buildID)
	if err != nil {
		shutdownChan <- true
		return msgChan, shutdownChan, err
	}
	l.Debug("Request sent")

	return msgChan, shutdownChan, nil
}

func listenForShutdown(conn *io.Client, shutdown chan bool, msgChan chan *BuildMessage) {
	select {
	case <-shutdown:
		conn.Close()
		close(shutdown)
		close(msgChan)
	}
}

type errorMessage struct {
	BuildID string `json:"build_id"`
	Reason  string `json:"reason"`
}

func (m errorMessage) convert() *BuildMessage {
	return &BuildMessage{
		ErrorMessage: m.Reason,
		Completed:    true,
		BuildID:      m.BuildID,
	}
}

type completeMessage string

func (m completeMessage) convert() *BuildMessage {
	return &BuildMessage{
		Completed: true,
		BuildID:   fmt.Sprintf("%s", m),
	}
}

type logMessage struct {
	BuildID  string                 `json:"build_id"`
	DeployID string                 `json:"deploy_id"`
	SiteID   string                 `json:"site_id"`
	Time     string                 `json:"time"`
	Msg      string                 `json:"msg"`
	Phase    string                 `json:"phase"`
	Data     map[string]interface{} `json:"data"`
}

func (m logMessage) convert() *BuildMessage {
	return &BuildMessage{
		Completed: false,
		Message:   m.Msg,
		DeployID:  m.DeployID,
		SiteID:    m.SiteID,
		BuildID:   m.BuildID,
		Data:      m.Data,
		Time:      m.Time,
	}
}
