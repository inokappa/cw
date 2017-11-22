package main

import (
    "os"
    "flag"
    "fmt"
    "time"
    "io/ioutil"
    "sort"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "gopkg.in/yaml.v2"
)

const AppVersion = "0.0.1"

var (
    argProfile = flag.String("profile", "", "Profile 名を指定.")
    argRegion = flag.String("region", "ap-northeast-1", "Region 名を指定.")
    argEndpoint = flag.String("endpoint", "", "AWS API のエンドポイントを指定.")
    argConfig = flag.String("config", "", "YAML ファイルを指定.")
    argTarget = flag.String("target", "", "YAML ファイル内のターゲット名を指定.")
    argVersion = flag.Bool("version", false, "バージョンを出力.")
)

func awsCloudWatchClient(profile string, region string) *cloudwatch.CloudWatch {
    var config aws.Config
    if profile != "" {
        creds := credentials.NewSharedCredentials("", profile)
        config = aws.Config{Region: aws.String(region), Credentials: creds, Endpoint: aws.String(*argEndpoint)}
    } else {
        config = aws.Config{Region: aws.String(region), Endpoint: aws.String(*argEndpoint)}
    }
    sess := session.New(&config)
    cwClient := cloudwatch.New(sess)
    return cwClient
}

func readConfig(config_yaml string, target string) (config map[interface{}]interface{}) {
    yml, err := ioutil.ReadFile(config_yaml)
    if err != nil {
        fmt.Println(err)
    }

    m := make(map[interface{}]interface{})
    err = yaml.Unmarshal([]byte(yml), &m)
    if err != nil {
        fmt.Println(err)
    }

    config = make(map[interface{}]interface{})
    config["start_time"] = m[target].(map[interface {}]interface {})["start_time"].(int)
    config["metric_name"] = m[target].(map[interface {}]interface {})["metric_name"].(string)
    config["namespace"] = m[target].(map[interface {}]interface {})["namespace"].(string)
    config["period"] = m[target].(map[interface {}]interface {})["period"].(int)
    config["statistics"] = m[target].(map[interface {}]interface {})["statistics"].(string)
    config["dimensions"] = m[target].(map[interface {}]interface {})["dimensions"].([]interface {})
    config["unit"] = m[target].(map[interface {}]interface {})["unit"].(string)

    return config
}

func generateDimensions(dimensions []interface{}) []*cloudwatch.Dimension {
    var dims []*cloudwatch.Dimension
    for _, dimension := range dimensions {
        d, _ := dimension.(map[interface {}]interface {})
        dim := &cloudwatch.Dimension{
            Name: aws.String(d["name"].(string)),
            Value: aws.String(d["value"].(string)),
        }
        dims = append(dims, dim)
    }
    return dims
}

func main() {

    flag.Parse()
    if *argVersion {
        fmt.Println(AppVersion)
        os.Exit(0)
    }

    if *argConfig == "" {
        fmt.Println("-config オプションで YAML ファイルを指定して下さい.")
        os.Exit(1)
    }

    if *argTarget == "" {
        fmt.Println("-target オプションで YAML ファイルのターゲットを指定して下さい.")
        os.Exit(1)
    }

    config := readConfig(*argConfig, *argTarget)

    time_now := time.Now()
    duration, _ := config["start_time"]
    start_time := time_now.Add(time.Duration(duration.(int)) * time.Second)
    dimensions := generateDimensions(config["dimensions"].([]interface{}))
    period := int64(config["period"].(int))

    params := &cloudwatch.GetMetricStatisticsInput {
        EndTime: aws.Time(time_now),
        MetricName: aws.String(config["metric_name"].(string)),
        Namespace: aws.String(config["namespace"].(string)),
        Period: aws.Int64(period),
        StartTime: aws.Time(start_time),
        Statistics: []*string{
            aws.String(config["statistics"].(string)),
        },
        Dimensions: dimensions,
        Unit: aws.String(config["unit"].(string)),
    }

    cwClient := awsCloudWatchClient(*argProfile, *argRegion)
    res, err := cwClient.GetMetricStatistics(params)

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Timestamp でソートする
    sort.Slice(res.Datapoints, func(i, j int) bool {
        return (*res.Datapoints[i].Timestamp).Before(*res.Datapoints[j].Timestamp)
    })

    fmt.Println(res.Datapoints)

}
