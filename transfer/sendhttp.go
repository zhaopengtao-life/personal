package transfer

import (
    "bytes"
    "compress/gzip"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "math"
    "net"
    "net/http"
    "personal_work/model"
    "strconv"
    "strings"
    "time"

    log "github.com/sirupsen/logrus"
    nsema "github.com/toolkits/concurrent/semaphore"
    nlist "github.com/toolkits/container/list"
)

var (
    HttpClient      *http.Client
    LoseRegainQueue = nlist.NewSafeListLimited(LoseRegainQueueMaxSize)
    ToTransferQueue = nlist.NewSafeListLimited(ToTransferQueueMaxSize)
    PushMap         = make(map[string]int64, 0)
)

const (
    LoseRegainQueueMaxSize  = 50 * 200
    ToTransferQueueMaxSize  = 102400 * 20
    ToTransferBatchSendSize = 5000
)

//
// CreateHTTPClient
//  @Description: 建立单例链接
//  @param addrIP
//
func CreateHTTPClient() {
    // 使用单例创建client
    httpClient := &http.Client{
        Transport: &http.Transport{
            DialContext: (&net.Dialer{
                Timeout:   30 * time.Second, // 连接超时时间
                KeepAlive: 30 * time.Second, // 保持长连接的时间
                //LocalAddr: &net.TCPAddr{IP: net.ParseIP(addrIP)},
            }).DialContext,
            ForceAttemptHTTP2: true,
            //MaxIdleConns:          1000,                                  // 最大空闲连接数
            MaxIdleConnsPerHost:   2000,                                  // 对每个host的最大连接数量(MaxIdleConnsPerHost<=MaxIdleConns)
            IdleConnTimeout:       30 * time.Second,                      // 连接最大空闲时间，超过这个时间就会被关闭。
            TLSHandshakeTimeout:   10 * time.Second,                      // 限制TLS握手使用的时间
            ExpectContinueTimeout: 1 * time.Second,                       // 限制客户端在发送一个包含
            TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // 跳过https证书校验
        },
        Timeout: 6 * time.Second,
    }
    HttpClient = httpClient
}

//
// DataSendToTransfer
//  @Description: 数据过滤
//  @param metricValues
//  @param TaskKey
//
func DataSendToTransfer(metricValues []*model.MetricValue, TaskKey string) {
    log.Infof("TaskKey: %v, module: DataSendToTransfer, metricValues_size: %v", TaskKey, len(metricValues))
    newMetrics := make([]*model.MetricValue, 0)
    for _, v := range metricValues {
        if v.Metric == "" || v.Endpoint == "" {
            log.Errorf("TaskKey: %v, DataSendToTransfer Metric Or Endpoint Is Nil Data: %v", TaskKey, v)
            continue
        }
        if v.Timestamp == 0 {
            log.Errorf("TaskKey: %v, DataSendToTransfer Metrics Data: %v Timestamp is Nil ", v)
            continue
        }
        var value interface{}
        switch cv := v.Value.(type) {
        case float64:
            if math.IsNaN(v.Value.(float64)) {
                value = 0.01
                log.Errorf("TaskKey: %v, DataSendToTransfer Data: %v Value is IsNaN ", TaskKey, v)
            }
            data := fmt.Sprintf("%0.2f", cv)
            if strings.Contains(data, "+Inf") {
                value = 0.01
                log.Errorf("ProcessData Error Float64 data: %v  v.Value:%v ", data, cv)
            } else {
                value = cv
            }
        case uint64:
            data := strconv.FormatUint(cv, 10)
            values, err := strconv.ParseFloat(data, 64)
            if err != nil {
                value = 0.01
                log.Errorf("ProcessData Error Uint64 data: %v  v.Value:%v ", data, cv)
            } else {
                value = values
            }
        default:
            value = cv
        }
        v.Value = value
        newMetrics = append(newMetrics, v)
    }
    if len(newMetrics) != 0 {
        // http发送
        PushSendToTransferQueue(newMetrics)
    }
}

//
// PushSendToTransferQueue
//  @Description: 正常数据待发送队列
//  @param item
//
func PushSendToTransferQueue(data []*model.MetricValue) {
    items := make([]interface{}, 0)
    for _, v := range data {
        items = append(items, v)
    }
    isSuccess := ToTransferQueue.PushFrontBatch(items)
    if !isSuccess {
        log.Errorf("PushSendToTransferQueue Error QueueLength: %v, SendToTransferData: %v", ToTransferQueue.Len(), items)
    }
}

//
// ToTransfer
//  @Description: 数据发送：http，需初始化启动
//
func ToTransfer() {
    // 最大并发数
    sema := nsema.NewSemaphore(32)
    for {
        itemList := ToTransferQueue.PopBackBy(ToTransferBatchSendSize)
        if len(itemList) == 0 {
            time.Sleep(3 * time.Second)
            continue
        }
        //  同步Call + 有限并发 进行发送
        sema.Acquire()
        go func(itemList []interface{}) {
            defer sema.Release()
            items := make([]*model.MetricValue, 0)
            for _, v := range itemList {
                items = append(items, v.(*model.MetricValue))
            }
            SendToTransfer(items)
        }(itemList)
    }
}

// SendToTransfer
//  @Description: 数据发送，失败重试
//  @param metricValues
//  @param TaskKey
//
func SendToTransfer(metricValues []*model.MetricValue) {
    // 异常捕获
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("SendToTransfer Send Transfer Panic: %v", err)
            return
        }
    }()
    processData, err := ProcessData(metricValues)
    if err != nil {
        log.Errorf("SendToTransfer ProcessData Fail Error items: %v, Data: %v", err, metricValues)
        return
    }

    // https进行链接
    if HttpClient == nil {
        CreateHTTPClient()
    }
    url := "127.0.0.1" + "/api/new/push"
    request, err := http.NewRequest("POST", url, strings.NewReader(string(processData)))
    if err != nil {
        log.Errorf("SendToTransfer NewRequest Client  Url: %v  Fail Error: %v, Data: %v", url, err, metricValues)
        return
    }
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")
    request.Header.Set("Connection", "Keep-Alive")
    request.Header.Set("Content-Encoding", "gzip")
    request.Header.Set("Accept-Encoding", "gzip")
    resp, err := HttpClient.Do(request)
    if err != nil {
        log.Errorf("SendToTransfer HttpClient Url: %v Fail Error: %v, Data: %v", url, err, metricValues)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        log.Errorf("SendToTransfer MetricValue Send LoseResponse.Code: %v Error: %v", resp.StatusCode, err)
        return
    }
    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Errorf("SendToTransfer Response ioutil.ReadAll Error: %v", err)
        return
    }
}

//
// ProcessData
//  @Description: 压缩数据转换结构
//  @param metrics
//  @return []byte
//  @return error
//
func ProcessData(metrics []*model.MetricValue) (dataGzip []byte, err error) {
    metricArray := make([][]interface{}, 0)
    for _, v := range metrics {
        array := make([]interface{}, 9)
        array[0] = v.Endpoint
        array[1] = v.Metric
        array[2] = v.Value
        array[3] = fmt.Sprint(v.Step)
        array[4] = v.Type
        array[5] = v.Tags
        array[6] = fmt.Sprint(v.Timestamp)
        if v.Desc != "" {
            array[7] = v.Desc
        }
        if v.Index != 0 {
            array[8] = fmt.Sprint(v.Index)
        }
        metricArray = append(metricArray, array)
    }
    sendTransferValue := &model.TransferValue{
        MetricArray: metricArray,
    }
    data, err := json.Marshal(sendTransferValue)
    if err != nil {
        log.Errorf("ProcessData Json Marshal Error: %v Data: %v", err, sendTransferValue)
        return nil, err
    }
    // 进行数据压缩
    dataGzip, err = Encodes(data)
    if err != nil {
        log.Errorf("ProcessData Encodes Error: %v", err)
        return nil, err
    }
    return dataGzip, nil
}

//
// Encodes
//  @Description: gzip 压缩
//  @param input
//  @return []byte
//  @return error
//
func Encodes(input []byte) ([]byte, error) {
    // 创建一个新的 byte 输出流
    var buf bytes.Buffer
    // 创建一个新的 gzip 输出流
    gzipWriter := gzip.NewWriter(&buf)
    // 将 input byte 数组写入到此输出流中
    _, err := gzipWriter.Write(input)
    if err != nil {
        _ = gzipWriter.Close()
        return nil, err
    }
    if err := gzipWriter.Close(); err != nil {
        return nil, err
    }
    // 返回压缩后的 bytes 数组
    return buf.Bytes(), nil
}
