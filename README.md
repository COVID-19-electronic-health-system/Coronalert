# Basic SMS Notification Service

CoronaTracker's notifications service.

## Curent state
- Transitioning EC2 instances on @Carter Klein's personal AWS profile to a secure CoronaTracker cloud environment for launch

# Prerequisites
- Request access to the CoronaTracker database (approval dependent on prior commits to other CoronaTracker repositories) [here](https://docs.google.com/spreadsheets/d/1Y_l4oq_32q1IMhpLyFmFAF1xnKSWQz7TPSX27ndveuQ/edit#gid=0)
	- **NOTE:** When signing up for AWS, your username/password combination will be in the signup email you receive
- Create accounts for AWS, MongoDB, Twilio
	- **NOTE:** In AWS, ensure you're always located in N. Virginia (us-east-1)
- Ensure Go (we're on 1.13) is installed
- Set `MONGODB_URI` environment variable locally. This can be found in MongoDB Atlas -> Clusters -> Connect (in the box of the cluster your want) -> connect your application (see [here](https://studio3t.com/knowledge-base/articles/connect-to-mongodb-atlas/))
- Set `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN` environment variables locally. This can befound on the Twilio dashboard under "Account SID" and "Auth Token," respectively

# Creating a new lambda function
- Create new AWS Lambda (we would prefer using Go 1.x runtime). Use the "Simple microservices permissions" policy
- Set the `MONGODB_URI` environment variable (TODO: encrypt in transit)
- Open permissions -> <your-lambda-role> and attach the `coronalert-lambda-policy` to it
- Add this lambda to the CoronaTracker VPC, in the two private application subnets
- Create a representative test for the lambda for future use
- Attach the Coronalert API gateway as a trigger, with IAM security
- Remove the `ANY` type (see Actions button), and add requests specific to your lambda
- Deploy lambda (see Actions button) using the default deployment

**NOTE:** When building a Lambda function in Golang, follow these steps:

### OSX
1. Build your `main.go` using `GOOS=linux go build -o main`
2. Zip the executable `main` using `zip main.zip main`
3. Upload `main.zip` to the lambda

### Windows
1. Install `build-lambda-zip` using `go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip`
2. Build your `main.go` using `GOOS=linux go build -o main`
3. Zip the executable `main` using `~/go/bin/build-lambda-zip.exe main` (**~/go/bin/ could be different depending on your file system layout**)
4. Upload `main.zip` to the lambda

Golang is new to having native support by AWS Lambda, therefore you aren't yet able to edit lambdas inline.