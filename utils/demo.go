package utils

//func insert()  {
//    // 原始sql批量插入
//    for _, v := range transPortStandard {
//        transPort := fmt.Sprintf("(%d,%#v,%.2f,%.2f,%.2f,%.2f,%v,%d)", transportID, v.ProvinceIds, v.FirstPiece, v.FirstFee, v.ContinuousFee, v.ContinuousPiece, v.TransPortDefault, v.GroupId)
//        transPortStandards = transPortStandards + transPort + ","
//    }
//}

//func GlobalOidDict(oid string) (map[string]*global.OidDict, error) {
//    var oidDictConfig *global.OidDictConfig
//    oidDictMap := make(map[string]*global.OidDict, 0)
//    file, err := os.Open(oid)
//    if err != nil {
//        panic(err)
//    }
//    bytes, err := ioutil.ReadAll(file)
//    if err != nil {
//        panic(err)
//    }
//    err = json.Unmarshal(bytes, &oidDictConfig) // 注意：要将和json对应的结构体指针传进来，而不是结构体对象
//    if err != nil {
//        log.Warn(" GlobalOidDict failed to write file ", err)
//        return nil, err
//    }
//    for k, v := range oidDictConfig.OidDict {
//        if _, ok := oidDictMap[v.SwitchName]; !ok {
//            oidDictMap[v.SwitchName] = oidDictConfig.OidDict[k]
//        }
//    }
//    return oidDictMap, nil
//}
