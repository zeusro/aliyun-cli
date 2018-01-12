package command

import (
	"github.com/aliyun/aliyun-cli/cli"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/aliyun-cli/core"
	"github.com/aliyun/aliyun-cli/meta"
	"time"
	"strings"
)

var products = meta.LoadProfile()

func InitOpenApiCaller(cmd *cli.Command) {
	cmd.Run = func(c *cli.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("invaild arguments")
		}

		product := args[0]
		api := args[1]
		parameters := make(map[string]string)

		for _, v := range c.UnknownFlags().Items() {
			parameters[v.Name] = v.GetValue()
		}

		s, err := CallOpenApi(product, api, parameters)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	}
}

func CallOpenApi(product string, api string, parameters map[string]string) (string, error) {
	conf, err := core.LoadConfiguration()
	if err != nil {
		return "", err
	}
	cp := conf.GetCurrentProfile()

	client, err := cp.GetClient()

	if err != nil {
		fmt.Errorf("failed with new client %v", err)
	}

	product = strings.ToLower(product)
	productInfo, ok := products[product]
	if !ok {
		return "", fmt.Errorf("unknown product %s", product)
	}

	method, roa := parameters["roa"]

	request := requests.NewCommonRequest()
	request.RegionId = cp.RegionId
	request.Product = productInfo.Name
	request.ApiName = api
	request.Version = productInfo.MajorVersion

	if roa {
		request.Method = method
		request.Headers["Date"] = time.Now().Format(time.RFC1123Z)
		request.PathPattern = api
	}

	for k, v := range parameters {
		switch k {
		case "region":
			request.RegionId = v
		case "endpoint":
			request.Domain = v
		case "version":
			request.Version = v
		case "body":
			request.SetContent([]byte(v))
		case "accept":
			request.Headers["Accept"] = v
			if strings.Contains(v, "xml") {
				request.SetAcceptFormat("XML")
			} else if strings.Contains(v, "json") {
				request.SetAcceptFormat("JSON")
			}
		default:
			request.QueryParams[k] = v
		}
	}

	if request.Domain == "" {
		request.Domain, err = productInfo.GetEndpoint(request.RegionId, client)
		if err != nil {
			return "", err
		}
	}

	resp, err := client.ProcessCommonRequest(request)
	if err != nil {
		if !strings.Contains(err.Error(), "unmarshal") {
			// fmt.Printf("%v\n", err)
			return "", err
		}
	}

	return resp.GetHttpContentString(), nil
}
