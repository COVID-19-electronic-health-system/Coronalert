# Basic SMS Notification Service

Our notifications service.

## Curent state
- Subscribe functionality running on an EC2 instance on @Carter Klein's private AWS account
- Data is ephemeral, stored only while the instance is running

## Near-term state
- Running this Dockerized service on an EC2 instance behind an API Gateway, on a collaborative AWS account
- Data is persisted to a MongoDB cluster and encrypted at rest

## Future State:
- Run this service on a Lambda, or EC2 instances behind an autoscaling group
- iOS/Android push notifications
- Much more fleshed-out and nuanced notifications dependent on location, age, etc

# Install/run locally with Postman

### Prerequisites
1. Install Go
2. Set up your environment with `TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN` - for the time being, ask @Carter Klein

1. `git clone https://github.com/COVID-19-electronic-health-system/BasicNotificationService`
2. `cd BasicNotificationServce`
3. `go get -u ./...`
4. `go run main.go`
5. To ensure it's working, in Postman, send a POST request to `localhost:8080/api/subscribe`:
```json
{
	"number": "1234567890"
}
```
