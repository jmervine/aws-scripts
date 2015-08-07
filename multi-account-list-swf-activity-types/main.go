package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmervine/aws-scripts/shared/account"
	"github.com/jmervine/aws-scripts/shared/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/swf"
)

var printer *utils.Output

// CONFIGURATION
// - Changable via flags:
var output *string
var creds = aws.String("../creds.txt")
var region = aws.String("us-east-1")
var uniq *bool

// MAIN
func main() {
	// flags, overrides for testing
	creds = flag.String("c", *creds, "creds file")
	region = flag.String("r", *region, "region")
	output = flag.String("o", "", "tee result to, defaults to <metric>.csv")
	uniq = flag.Bool("u", true, "print only unique activity type names")
	flag.Parse()

	if *output == "" {
		output = aws.String("ActivityTypeNames.csv")
	}

	printer = utils.NewOutput(*output, true)

	// static
	os.Setenv("AWS_REGION", *region)

	accts := account.ParseAccounts(creds)
	account.EachAccountAsync(accts, func(a *account.Account) {
		// export creds to env (aws is weird) b/c it's the easiest way
		os.Setenv("AWS_ACCESS_KEY_ID", a.Key)
		os.Setenv("AWS_SECRET_ACCESS_KEY", a.Secret)

		config := &aws.Config{
			Region: region,

			// fetch creds from envs
			Credentials: credentials.NewEnvCredentials(),
		}

		sw := swf.New(config)

		getDomain := func() *string {
			dom, err := sw.ListDomains(&swf.ListDomainsInput{
				RegistrationStatus: aws.String(swf.RegistrationStatusRegistered),
			})

			if err != nil {
				panic(err)
			}

			if len(dom.DomainInfos) > 1 {
				fmt.Printf("%+v\n", dom.DomainInfos)
				panic(fmt.Errorf("unexpected domain count"))
			}

			if len(dom.DomainInfos) == 0 {
				return nil
			}

			return dom.DomainInfos[0].Name
		}

		domain := getDomain()

		if domain == nil || *domain == "" {
			return
		}

		act, err := sw.ListActivityTypes(&swf.ListActivityTypesInput{
			Domain:             getDomain(),
			RegistrationStatus: aws.String(swf.RegistrationStatusRegistered),
		})

		if err != nil {
			panic(err)
		}

		if len(act.TypeInfos) == 0 {
			panic(fmt.Errorf("unexpected activity type count"))
		}

		if *uniq {
			var lines []string
			for _, t := range act.TypeInfos {
				lines = utils.AppendStringIfMissing(lines, *t.ActivityType.Name)
			}

			for _, l := range lines {
				printer.Puts(fmt.Sprintf("%s\n", l))
			}
		} else {
			for _, t := range act.TypeInfos {
				printer.Puts(fmt.Sprintf("%s,%s,%s\n", a.Name, *domain, *t.ActivityType.Name))
			}
		}

		return
	})
}
