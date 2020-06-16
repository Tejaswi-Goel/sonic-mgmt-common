package custom_validation

import (
    "github.com/go-redis/redis/v7"
    util "github.com/Azure/sonic-mgmt-common/cvl/internal/util"
    "fmt"
    "net"
    "strings"
)

func checkAttrAlignment(vc *CustValidationCtxt) (string, error) {
    keys := strings.Split(vc.CurCfg.Key, "|")
    if len(keys) < 2 || keys[0] != "STATIC_ROUTE" {
        return "", fmt.Errorf("Invalid key format: %s", vc.CurCfg.Key)
    }
    prefix := keys[len(keys) - 1]
    if vc.SessCache.Data == nil {
        vc.SessCache.Data = make(map[string]bool)
    }
    nhCheckMap, ok := vc.SessCache.Data.(map[string]bool)
    if !ok {
        return "", fmt.Errorf("Invalid data type in session cache")
    }
    if _, ok := nhCheckMap[prefix]; !ok {
        // verify attr alignment and get list item number
        cfgNhNum := 0
        for _, fv := range vc.CurCfg.Data {
            if cfgNhNum == 0 {
                cfgNhNum = len(strings.Split(fv, ","))
            } else if cfgNhNum != len(strings.Split(fv, ",")) {
                return "", fmt.Errorf("NH config attributes not aligned")
            }
        }
        // check DB data alignment
        attrs, err := vc.RClient.HGetAll(vc.CurCfg.Key).Result()
	    if err != nil && err != redis.Nil {
            return "", fmt.Errorf("Failed to read NH attribute from DB, key: %s", vc.CurCfg.Key)
        }
        dbFields := make(map[string]bool)
        var dbNhNum int
	    if err != redis.Nil && len(attrs) != 0 {
            for dfn, dfv := range attrs {
                itemNum := len(strings.Split(dfv, ","))
                if dbNhNum == 0 {
                    dbNhNum = itemNum
                } else if dbNhNum != itemNum {
                    util.CVL_LEVEL_LOG(util.INFO, "Original DB data is not aligned, bypass checking")
                    return prefix, nil
                }
                dbFields[dfn] = true
            }
        }
        if cfgNhNum != dbNhNum {
            // check if update covered all non-delete fields
            for _, cfgData := range vc.ReqData {
                if cfgData.VOp == OP_DELETE {
                    if cfgData.Data == nil || len(cfgData.Data) == 0 {
                        return prefix, nil
                    }
                    for fn := range cfgData.Data {
                        delete(dbFields, fn)
                    }
                }
            }
            for fn := range dbFields {
                if _, ok := vc.CurCfg.Data[fn]; !ok {
                    return "", fmt.Errorf("Field %s in DB need to be updated to keep alignment", fn);
                }
            }
        }
        nhCheckMap[prefix] = true
    }
    return prefix, nil
}

type checkNHAttrHdlr func(string, string, ...interface{}) error

func checkNexthopGateway(prefix string, gwIP string, args ...interface{}) error {
    pfxIpStr := strings.Split(prefix, "/")[0]
    pfxIp := net.ParseIP(pfxIpStr)
    if pfxIp == nil {
        return fmt.Errorf("Invalid static route IP prefix: %s", prefix)
    }
    pfxIpv4 := pfxIp.To4() != nil
    ip := net.ParseIP(gwIP)
    if ip == nil {
        return fmt.Errorf("Invalid gateway IP format %s", gwIP)
    }
    gwIpv4 := ip.To4() != nil
    if gwIpv4 != pfxIpv4 {
        return fmt.Errorf("Address family of NH gateway %s not same as prefix %s", gwIP, pfxIpStr)
    }
    return nil
}

func checkTableKeyExists(db *redis.Client, tableList []string, key string) bool {
    for _, table := range tableList {
        fullKey := fmt.Sprintf("%s|%s", table, key)
        attrs, err := db.HGetAll(fullKey).Result()
        if err == nil && attrs != nil && len(attrs) > 0 {
            return true
        }
    }
    return false
}

func checkNexthopIntfVrf(_, name string, args ...interface{}) error {
    if len(args) < 2 {
        return fmt.Errorf("No enough arguments given for NH interface or VRF check")
    }
    if len(name) == 0 {
        return nil
    }
    tableList, ok := args[0].([]string)
    if !ok {
        return fmt.Errorf("Invalid argument type of table list for NH interface or VRF check")
    }
    db, ok := args[1].(*redis.Client)
    if !ok {
        return fmt.Errorf("Invalid argument type of DB client for NH interface or VRF check")
    }
    if found := checkTableKeyExists(db, tableList, name); !found {
        return fmt.Errorf("Interface or VRF %s not found in config DB", name)
    }
    return nil
}

func validateNexthopAttrCmn(vc *CustValidationCtxt, hdlr checkNHAttrHdlr,
                            args ...interface{}) CVLErrorInfo {
    if vc.CurCfg.VOp == OP_DELETE || len(vc.YNodeVal) == 0 {
	    return CVLErrorInfo{ErrCode: CVL_SUCCESS}
    }
    util.CVL_LEVEL_LOG(util.INFO, "Nexthop attribute validation: Oper %d Key %s Node %s Value %s",
                       vc.CurCfg.VOp, vc.CurCfg.Key, vc.YNodeName, vc.YNodeVal)
    prefix, err := checkAttrAlignment(vc)
    if err != nil {
        errMsg := fmt.Sprintf("Failed checking alignment: %s", err)
        return CVLErrorInfo{ErrCode: CVL_ERROR, Keys:[]string{vc.CurCfg.Key}, Value: vc.YNodeVal, Field: vc.YNodeName,
                            Msg: errMsg}
    }
    attrs := strings.Split(vc.YNodeVal, ",")
    for _, attrVal := range attrs {
        err = hdlr(prefix, attrVal, args...)
        if err != nil {
            errMsg := fmt.Sprintf("Validate failed for attribute %s", attrVal)
            return CVLErrorInfo{ErrCode: CVL_ERROR, Keys:[]string{vc.CurCfg.Key}, Value: vc.YNodeVal, Field: vc.YNodeName,
                                Msg: errMsg}
        }
    }
    return CVLErrorInfo{ErrCode: CVL_SUCCESS}
}

<<<<<<< HEAD
// ValidateNexthopGateway check if every item in comma separated list is valid IP address
// Path : /sonic-static-route/STATIC_ROUTE/nexthop
||||||| merged common ancestors
//Path : /sonic-static-route/STATIC_ROUTE/nexthop
// Purpose: To check if every item in comma separated list is valid IP address
=======
// ValidateNexthopGateway checks if every item in comma separated list is valid IP address
>>>>>>> origin/broadcom_sonic_3.x_share
// Returns -  CVL Error object
// Path /sonic-static-route/STATIC_ROUTE/nexthop
func (t *CustomValidation) ValidateNexthopGateway(
	vc *CustValidationCtxt) CVLErrorInfo {
    return validateNexthopAttrCmn(vc, checkNexthopGateway)
}

<<<<<<< HEAD
// ValidateNexthopInterface check if every item in comma separated list is an active interface name
// Path : /sonic-static-route/STATIC_ROUTE/ifname
||||||| merged common ancestors
//Path : /sonic-static-route/STATIC_ROUTE/ifname
// Purpose: To check if every item in comma separated list is an active interface name
=======
// ValidateNexthopInterface checks if every item in comma separated list is an active interface name
>>>>>>> origin/broadcom_sonic_3.x_share
// Returns -  CVL Error object
// Path /sonic-static-route/STATIC_ROUTE/ifname
func (t *CustomValidation) ValidateNexthopInterface(
	vc *CustValidationCtxt) CVLErrorInfo {
    var tableList = []string{"PORT", "PORTCHANNEL", "VLAN", "LOOPBACK_INTERFACE"}
    return validateNexthopAttrCmn(vc, checkNexthopIntfVrf, tableList, vc.RClient)
}

<<<<<<< HEAD
// ValidateNexthopVrf check if every item in comma separated list is an active VRF name
// Path : /sonic-static-route/STATIC_ROUTE/nexthop-vrf
||||||| merged common ancestors
//Path : /sonic-static-route/STATIC_ROUTE/nexthop-vrf
// Purpose: To check if every item in comma separated list is an active VRF name
=======
// ValidateNexthopVrf checks if every item in comma separated list is an active VRF name
>>>>>>> origin/broadcom_sonic_3.x_share
// Returns -  CVL Error object
// Path /sonic-static-route/STATIC_ROUTE/nexthop-vrf
func (t *CustomValidation) ValidateNexthopVrf(
	vc *CustValidationCtxt) CVLErrorInfo {
    var tableList = []string{"VRF"}
    return validateNexthopAttrCmn(vc, checkNexthopIntfVrf, tableList, vc.RClient)
}
