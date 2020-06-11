//////////////////////////////////////////////////////////////////////////
//
// Copyright 2020 Broadcom.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//////////////////////////////////////////////////////////////////////////

package transformer

import (
    "os"
    "errors"
    "strings"
    "strconv"
    "syscall"
    "net"
    "io/ioutil"
    "encoding/json"
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "github.com/Azure/sonic-mgmt-common/translib/tlerr"
    log "github.com/golang/glog"
    "github.com/openconfig/ygot/ygot"
    "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
)

//CounterObj - Sub structure required based on the counters file in the DUT
type CounterObj  struct {
    Value         string    `json: "value"`
    Description   string    `json:"description"`
    }

// JSONDhcpCounters - Counters structure for DHCP
type JSONDhcpCounters  struct {
    BootrequestSent        CounterObj  `json:"bootrequest-sent"`
    BootreplySent          CounterObj  `json:"bootreply-sent"`
    TotalDropped           CounterObj  `json:"total-dropped"`
    InvalidOpcode          CounterObj  `json:"invalid-opcode"`
    InvalidOptions         CounterObj  `json:"invalid-options"`
    BootrequestReceived    CounterObj  `json:"bootrequest-received"`
    DhcpDeclineReceived    CounterObj  `json:"dhcp-decline-received"`
    DhcpDiscoverReceived   CounterObj  `json:"dhcp-discover-received"`
    DhcpInformReceived     CounterObj  `json:"dhcp-inform-received"`
    DhcpRequestReceived    CounterObj  `json:"dhcp-request-received"`
    DhcpReleaseReceived    CounterObj  `json:"dhcp-release-received"`
    DhcpOfferSent          CounterObj  `json:"dhcp-offer-sent"`
    DhcpAckSent            CounterObj  `json:"dhcp-ack-sent"`
    DhcpNackSent           CounterObj  `json:"dhcp-nack-sent"`
}

// JSONDhcpv6Counters - Counters structure for DHCPv6
type JSONDhcpv6Counters  struct {
    TotalDropped                CounterObj  `json:"total-dropped"`
    InvalidOpcode               CounterObj  `json:"invalid-opcode"`
    InvalidOptions              CounterObj  `json:"invalid-options"`
    Dhcpv6SolicitReceived       CounterObj  `json:"dhcpv6-solicit-received"`
    Dhcpv6DeclineReceived       CounterObj  `json:"dhcpv6-decline-received"`
    Dhcpv6RequestReceived       CounterObj  `json:"dhcpv6-request-received"`
    Dhcpv6ReleaseReceived       CounterObj  `json:"dhcpv6-release-received"`
    Dhcpv6ConfirmReceived       CounterObj  `json:"dhcpv6-confirm-received"`
    Dhcpv6RebindReceived        CounterObj  `json:"dhcpv6-rebind-received"`
    Dhcpv6InfoRequestReceived   CounterObj  `json:"dhcpv6-Info-request-received"`
    Dhcpv6RelayReplyReceived    CounterObj  `json:"dhcpv6-relay-reply-received"`
    Dhcpv6AdvertiseSent         CounterObj  `json:"dhcpv6-advertise-sent"`
    Dhcpv6ReplySent             CounterObj  `json:"dhcpv6-reply-sent"`
    Dhcpv6ReconfigureSent       CounterObj  `json:"dhcpv6-reconfigure-sent"`
    Dhcpv6RelayForwSent         CounterObj  `json:"dhcpv6-relay-forw-sent"`
}

//Sub structure required to translate yang to dbobject field
var  relayAgentFields []string = []string{
        "dhcp_servers@",
        "dhcp_relay_src_intf",
        "dhcp_relay_max_hop_count",
        "dhcp_relay_link_select"}

var  relayAgentV6Fields []string = []string{
        "dhcpv6_servers@",
        "dhcpv6_relay_src_intf",
        "dhcpv6_relay_max_hop_count"}

//PATH_PREFIX - global
const PATH_PREFIX = "/mnt/tmp/"

func init () {
    XlateFuncBind("relay_agent_table_xfmr", relay_agent_table_xfmr)
    XlateFuncBind("YangToDb_relay_agent_intf_tbl_key_xfmr", YangToDb_relay_agent_intf_tbl_key_xfmr)
    XlateFuncBind("DbToYang_relay_agent_intf_tbl_key_xfmr", DbToYang_relay_agent_intf_tbl_key_xfmr)
    XlateFuncBind("YangToDb_relay_agent_xfmr", YangToDb_relay_agent_xfmr)
    XlateFuncBind("DbToYang_relay_agent_xfmr", DbToYang_relay_agent_xfmr)
    XlateFuncBind("YangToDb_relay_agent_id_field_xfmr", YangToDb_relay_agent_id_field_xfmr)
    XlateFuncBind("DbToYang_relay_agent_counters_xfmr", DbToYang_relay_agent_counters_xfmr)
    XlateFuncBind("DbToYang_relay_agent_v6_counters_xfmr", DbToYang_relay_agent_v6_counters_xfmr)
}

// relay_agent_table_xfmr - Transformer function to loop over multiple interfaces
var relay_agent_table_xfmr TableXfmrFunc = func (inParams XfmrParams) ([]string, error) {
    var tblList []string
    var err error

    log.Info("RATableXfmrFunc - Uri: ", inParams.uri);
    pathInfo := NewPathInfo(inParams.uri)

    targetUriPath, err := getYangPathFromUri(pathInfo.Path)

    ifName := pathInfo.Var("id");
    log.Info(ifName)

    if ifName == "" {
        log.Info("TableXfmrFunc - intf_table_xfmr Intf key is not present")

        if _, ok := dbIdToTblMap[inParams.curDb]; !ok {
            log.Info("TableXfmrFunc - intf_table_xfmr db id entry not present")
            return tblList, errors.New("Key not present")
        } else {
            return dbIdToTblMap[inParams.curDb], nil
        }
    }

    intfType, _, ierr := getIntfTypeByName(ifName)
    if intfType == IntfTypeUnset || ierr != nil {
        log.Info("TableXfmrFunc - Invalid interface type IntfTypeUnset");
        return tblList, errors.New("Invalid interface type IntfTypeUnset");
    }
                      
    log.Info(intfType)

    intTbl := IntfTypeTblMap[intfType]
    log.Info("TableXfmrFunc - targetUriPath : ", targetUriPath)
    log.Info(intTbl)


    if (intfType == IntfTypeEthernet) || intfType == IntfTypePortChannel {
            tblList = append(tblList, intTbl.cfgDb.intfTN)
    } else if intfType == IntfTypeVlan {
            tblList = append(tblList, intTbl.cfgDb.portTN)
    }
    log.Info(tblList)
    return tblList, err

}

// YangToDb_relay_agent_intf_tbl_key_xfmr -Function to read the interface name from the interface table (interface, vlan, portchannel
var YangToDb_relay_agent_intf_tbl_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
    var err error

    pathInfo := NewPathInfo(inParams.uri)
    ifName := pathInfo.Var("id")

    return ifName, err
}

// DbToYang_relay_agent_intf_tbl_key_xfmr - Function to fetch the helper address from the appropriate interface table
var DbToYang_relay_agent_intf_tbl_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
    var err error
    var tblList string
    res_map := make(map[string]interface{})
    log.Info("DbToYang_relay_agent_intf_tbl_key_xfmr: ", inParams.key)

    if (inParams.key != "") {
        var configDb, _ = db.NewDB(getDBOptions(db.ConfigDB))

        intfType, _, _ := getIntfTypeByName(inParams.key)

        intTbl := IntfTypeTblMap[intfType]

    if (intfType == IntfTypeEthernet) || intfType == IntfTypePortChannel {
            tblList = intTbl.cfgDb.intfTN
        } else if intfType == IntfTypeVlan {
            tblList = intTbl.cfgDb.portTN
        }

        entry, dbErr := configDb.GetEntry(&db.TableSpec{Name:tblList}, db.Key{Comp: []string{inParams.key}})
        configDb.DeleteDB()
        if dbErr != nil {
            log.Info("Failed to read mgmt port status from config DB, " + tblList + " " + inParams.key)
            return res_map, dbErr
        }

        if (strings.HasPrefix(inParams.uri, "/openconfig-relay-agent:relay-agent/dhcp/")) && (entry.Get("dhcp_servers@") != "")  {

        //Check if config exist in table for the interface
            res_map["id"] = inParams.key
        }

        if (strings.HasPrefix(inParams.uri, "/openconfig-relay-agent:relay-agent/dhcpv6/")) && (entry.Get("dhcpv6_servers@") != "")  {
        //Check if config exist in table for the interface
            res_map["id"] = inParams.key
        }
    }

    return res_map, err
}

// YangToDb_relay_agent_id_field_xfmr- Function to transform id coming from Yang to vlan-id in the vlan table, Ethernet and Portchannel don't need special handling
var YangToDb_relay_agent_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
    res_map := make(map[string]string)
    var err error

    pathInfo := NewPathInfo(inParams.uri)
    ifName := pathInfo.Var("id")

    if strings.HasPrefix(ifName, VLAN) == true {
        vlanId := ifName[len("Vlan"):len(ifName)]
        res_map["vlanid"] = vlanId
    }
    log.Info("YangToDb_relay_agent_id_field_xfmr: res_map:", res_map)
    return res_map, err
}

// Helper function to read the DHCP counters from the file mounted in /mnt/tmp folder
func getRelayCountersFromFile (fileName string, counterObj interface{}) error {
   
    if log.V(7) {
    //If verbose logging is enabled, log info 
    log.Infof("getRelayCountersFromFile Enter")
    }

    tmpFileName := PATH_PREFIX + fileName
    log.Info(tmpFileName) 

    jsonFile, err := os.Open(tmpFileName)
    if err != nil {
        log.Warningf("opening of dhcp counters json file failed")
        errStr := "Information not available"
        terr := tlerr.NotFoundError{Format: errStr}
        return terr
    }
    syscall.Flock(int(jsonFile.Fd()),syscall.LOCK_EX)
    log.Infof("syscall.Flock done")

    defer jsonFile.Close()
    defer log.Infof("jsonFile.Close called")
    defer syscall.Flock(int(jsonFile.Fd()), syscall.LOCK_UN);
    defer log.Infof("syscall.Flock unlock  called")

    byteValue, _ := ioutil.ReadAll(jsonFile)
    err = json.Unmarshal(byteValue, counterObj)
    if err != nil {
        log.Warningf("unmarshal of the json counters failed")
        errStr := "json.Unmarshal failed"
        terr := tlerr.InternalError{Format: errStr}
        return terr
    }
    return nil
}

// Helper function to get the root object
func getRelayAgentRoot(s *ygot.GoStruct) *ocbinds.OpenconfigRelayAgent_RelayAgent {
    deviceObj := (*s).(*ocbinds.Device)
    return deviceObj.RelayAgent
}

// Helper function to convert Value from string to uint64
func getCounterValue(valStr string) *uint64 {
    if val, err := strconv.ParseUint(valStr, 10, 64); err == nil {
        return &val
    }
   return nil
}

//DbToYang_relay_agent_counters_xfmr - sub tree transformer - that will read the appropriate file and populate the DHCP counters
var DbToYang_relay_agent_counters_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var raObj *ocbinds.OpenconfigRelayAgent_RelayAgent_Dhcp_Interfaces_Interface 
    var jsonRelayAgentCounter JSONDhcpCounters

    log.Info("In DbToYang_relay_agent_counters_xfmr")
    if log.V(7) {
      log.Info(inParams)
    }
    relayAgentObj := getRelayAgentRoot(inParams.ygRoot)
    log.Info(relayAgentObj)

    pathInfo := NewPathInfo(inParams.uri)
    ifName := pathInfo.Var("id")
    log.Info(ifName)
    
    if ifName == "" { 
       return err 
    }
    
    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    
    fileName := "dhcp-relay-ipv4-stats-"+ ifName + ".json"
 
    err = getRelayCountersFromFile(fileName, &jsonRelayAgentCounter)
    log.Info(jsonRelayAgentCounter)
    if err != nil {
        log.Infof("getRelayCountersFromFile failed")
        return err
    }

    if strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcp/interfaces/interface/state/counters") {
        if relayAgentObj != nil && relayAgentObj.Dhcp != nil {
            var ok bool = false
            ygot.BuildEmptyTree(relayAgentObj.Dhcp)
            ygot.BuildEmptyTree(relayAgentObj.Dhcp.Interfaces)
         if raObj, ok = relayAgentObj.Dhcp.Interfaces.Interface[ifName]; !ok {
                raObj, _ = relayAgentObj.Dhcp.Interfaces.NewInterface(ifName)
            }
            ygot.BuildEmptyTree(raObj)
        } else if relayAgentObj != nil {
            ygot.BuildEmptyTree(relayAgentObj)
            ygot.BuildEmptyTree(relayAgentObj.Dhcp)
            ygot.BuildEmptyTree(relayAgentObj.Dhcp.Interfaces)
            raObj, _ = relayAgentObj.Dhcp.Interfaces.NewInterface(ifName)
            ygot.BuildEmptyTree(raObj)
        }
        ygot.BuildEmptyTree(raObj.State)
        ygot.BuildEmptyTree(raObj.State.Counters)
        } else {
            err = errors.New("Invalid URI : " + targetUriPath)
        }

    counterObj := relayAgentObj.Dhcp.Interfaces.Interface[ifName].State.Counters

    counterObj.TotalDropped = getCounterValue(jsonRelayAgentCounter.TotalDropped.Value)
    
    counterObj.InvalidOpcode = getCounterValue(jsonRelayAgentCounter.InvalidOpcode.Value)

    counterObj.InvalidOptions = getCounterValue(jsonRelayAgentCounter.InvalidOptions.Value)

    counterObj.BootrequestReceived = getCounterValue(jsonRelayAgentCounter.BootrequestReceived.Value)

    counterObj.DhcpDeclineReceived = getCounterValue(jsonRelayAgentCounter.DhcpDeclineReceived.Value)

    counterObj.DhcpDiscoverReceived = getCounterValue(jsonRelayAgentCounter.DhcpDiscoverReceived.Value)

    counterObj.DhcpInformReceived = getCounterValue(jsonRelayAgentCounter.DhcpInformReceived.Value)

    counterObj.DhcpRequestReceived = getCounterValue(jsonRelayAgentCounter.DhcpRequestReceived.Value)

    counterObj.DhcpReleaseReceived = getCounterValue(jsonRelayAgentCounter.DhcpReleaseReceived.Value)

    counterObj.BootrequestSent = getCounterValue(jsonRelayAgentCounter.BootrequestSent.Value)

    counterObj.BootreplySent = getCounterValue(jsonRelayAgentCounter.BootreplySent.Value)

    counterObj.DhcpOfferSent = getCounterValue(jsonRelayAgentCounter.DhcpOfferSent.Value)

    counterObj.DhcpAckSent = getCounterValue(jsonRelayAgentCounter.DhcpAckSent.Value)

    counterObj.DhcpNackSent = getCounterValue(jsonRelayAgentCounter.DhcpNackSent.Value)

    return err
}


//DbToYang_relay_agent_v6_counters_xfmr - sub tree transformer - that will read the appropriate file and populate the DHCPv6 counters
var DbToYang_relay_agent_v6_counters_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var raObj *ocbinds.OpenconfigRelayAgent_RelayAgent_Dhcpv6_Interfaces_Interface 
    var jsonV6RelayAgentCounter JSONDhcpv6Counters

    log.Info("In DbToYang_relay_agent_v6_counters_xfmr")
    if log.V(7) {
       log.Info(inParams)
    }
    relayAgentObj := getRelayAgentRoot(inParams.ygRoot)
    log.Info(relayAgentObj)

    pathInfo := NewPathInfo(inParams.uri)
    ifName := pathInfo.Var("id")
    
    if ifName == "" { 
       return err 
    }

    targetUriPath, err := getYangPathFromUri(pathInfo.Path)

    fileName := "dhcp-relay-ipv6-stats-"+ ifName + ".json"
 
    err = getRelayCountersFromFile(fileName, &jsonV6RelayAgentCounter)
    log.Info(jsonV6RelayAgentCounter)
    if err != nil {
        log.Infof("getRelayCountersFromFile failed")
        return err
    }

    if strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcpv6/interfaces/interface/state/counters") {
        if relayAgentObj != nil && relayAgentObj.Dhcpv6 != nil {
            var ok bool = false
            ygot.BuildEmptyTree(relayAgentObj.Dhcpv6)
            ygot.BuildEmptyTree(relayAgentObj.Dhcpv6.Interfaces)
            if raObj, ok = relayAgentObj.Dhcpv6.Interfaces.Interface[ifName]; !ok {
                raObj, _ = relayAgentObj.Dhcpv6.Interfaces.NewInterface(ifName)
            }
            ygot.BuildEmptyTree(raObj)
         } else if relayAgentObj != nil {
            ygot.BuildEmptyTree(relayAgentObj)
            ygot.BuildEmptyTree(relayAgentObj.Dhcpv6)
            ygot.BuildEmptyTree(relayAgentObj.Dhcpv6.Interfaces)
            raObj, _ = relayAgentObj.Dhcpv6.Interfaces.NewInterface(ifName)
            ygot.BuildEmptyTree(raObj)
        }
        ygot.BuildEmptyTree(raObj.State)
        ygot.BuildEmptyTree(raObj.State.Counters)
        } else {
            err = errors.New("Invalid URI : " + targetUriPath)
    }

    counterObj := relayAgentObj.Dhcpv6.Interfaces.Interface[ifName].State.Counters

    counterObj.TotalDropped = getCounterValue(jsonV6RelayAgentCounter.TotalDropped.Value)    

    counterObj.InvalidOpcode = getCounterValue(jsonV6RelayAgentCounter.InvalidOpcode.Value)

    counterObj.InvalidOptions = getCounterValue(jsonV6RelayAgentCounter.InvalidOptions.Value)
          
    counterObj.Dhcpv6SolicitReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6SolicitReceived.Value)
              
    counterObj.Dhcpv6DeclineReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6DeclineReceived.Value)

    counterObj.Dhcpv6RequestReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6RequestReceived.Value)

    counterObj.Dhcpv6ReleaseReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6ReleaseReceived.Value)

    counterObj.Dhcpv6ConfirmReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6ConfirmReceived.Value)

    counterObj.Dhcpv6RebindReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6RebindReceived.Value)

    counterObj.Dhcpv6InfoRequestReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6InfoRequestReceived.Value)

    counterObj.Dhcpv6RelayReplyReceived = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6RelayReplyReceived.Value)

    counterObj.Dhcpv6AdverstiseSent = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6AdvertiseSent.Value)

    counterObj.Dhcpv6ReplySent = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6ReplySent.Value)

    counterObj.Dhcpv6ReconfigureSent = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6ReconfigureSent.Value)

    counterObj.Dhcpv6RelayForwSent = getCounterValue(jsonV6RelayAgentCounter.Dhcpv6RelayForwSent.Value)

    return err

}


// Helper function to get the tableName
func getRelayAgentIntfTblByType(ifName string) string {
    var tblList string
    intfType, _, ierr := getIntfTypeByName(ifName)
    if intfType == IntfTypeUnset || ierr != nil {
        log.Info("getRelayAgentIntfTblByType - Invalid interface type IntfTypeUnset");
        return tblList;
    }

    intTbl := IntfTypeTblMap[intfType]

    if (intfType == IntfTypeEthernet) || intfType == IntfTypePortChannel {
            tblList = intTbl.cfgDb.intfTN
    } else if intfType == IntfTypeVlan {
            tblList = intTbl.cfgDb.portTN
    } else if intfType == IntfTypeLoopback {
            tblList = "LOOPBACK_INTERFACE"
    }
    return tblList;
}

// Helper function to get the tableName
func getDhcpDataFromDb(ifName string, relayAgentObj *ocbinds.OpenconfigRelayAgent_RelayAgent, configDb *db.DB) {
   var tblList string

   tblList = getRelayAgentIntfTblByType(ifName)
   log.Info(tblList)

   entry, dbErr := configDb.GetEntry(&db.TableSpec{Name:tblList}, db.Key{Comp: []string{ifName}})
   if dbErr != nil {
     log.Warning("Failed to read mgmt port status from config DB, " + tblList + " " + ifName)
     return
   }

   //continue only if there is data to proceed 
   if entry.Has("dhcp_servers@") != true {
      return
   }

   //Now create and assign the values to the object
   log.Info(relayAgentObj)

   raObj, ok := relayAgentObj.Dhcp.Interfaces.Interface[ifName]
   if !ok {
      raObj, _ = relayAgentObj.Dhcp.Interfaces.NewInterface(ifName)
   }

   ygot.BuildEmptyTree(raObj)
   ygot.BuildEmptyTree(raObj.Config)

   //Helper Address
   helperAddress := entry.GetList("dhcp_servers")
   raObj.Config.HelperAddress = helperAddress
   
   //Augmented Params
   if entry.Has("dhcp_relay_link_select") == true {
      linkSelectVal := entry.Get("dhcp_relay_link_select")
      if linkSelectVal == "enable" {
         raObj.Config.LinkSelect = ocbinds.OpenconfigRelayAgent_RelayAgent_Dhcp_Interfaces_Interface_Config_LinkSelect_enable
      } else {
        raObj.Config.LinkSelect = ocbinds.OpenconfigRelayAgent_RelayAgent_Dhcp_Interfaces_Interface_Config_LinkSelect_disable
     }
   }
   if entry.Has("dhcp_relay_max_hop_count") == true {
     mhCount, err := entry.GetInt("dhcp_relay_max_hop_count")
     if (err != nil) {
        log.Error("Unable to read max hop count")
     } else {
     mhCount32 := uint32(mhCount)
     raObj.Config.MaxHopCount = &mhCount32
     }
   }
   if entry.Has("dhcp_relay_src_intf") == true {
     srcIntf := entry.Get("dhcp_relay_src_intf")
     raObj.Config.SrcIntf = &srcIntf
   }
}

// Helper function to get the tableName
func getDhcpv6DataFromDb(ifName string, relayAgentObj *ocbinds.OpenconfigRelayAgent_RelayAgent, configDb *db.DB) {
   var tblList string

   tblList = getRelayAgentIntfTblByType(ifName)
   log.Info(tblList)

   entry, dbErr := configDb.GetEntry(&db.TableSpec{Name:tblList}, db.Key{Comp: []string{ifName}})
   if dbErr != nil {
     log.Warning("Failed to read mgmt port status from config DB, " + tblList + " " + ifName)
     return
   }

   //continue only if there is data to proceed 
   if entry.Has("dhcpv6_servers@") != true {
      return
   }

   //Now create and assign the values to the object
   log.Info(relayAgentObj)

   raObj, ok := relayAgentObj.Dhcpv6.Interfaces.Interface[ifName]
   if !ok {
      raObj, _ = relayAgentObj.Dhcpv6.Interfaces.NewInterface(ifName)
   }

   ygot.BuildEmptyTree(raObj)
   ygot.BuildEmptyTree(raObj.Config)

   //Helper Address
   helperAddress := entry.GetList("dhcpv6_servers@")
   raObj.Config.HelperAddress = helperAddress

   //Augmented Params
   if entry.Has("dhcpv6_relay_max_hop_count") == true {
      mhCount, err := entry.GetInt("dhcpv6_relay_max_hop_count")
     if (err != nil) {
        log.Error("Unable to read max hop count")
      } else {
        mhCount32 := uint32(mhCount)
        raObj.Config.MaxHopCount = &mhCount32
      }
   }
   if entry.Has("dhcpv6_relay_src_intf") == true {
      srcIntf := entry.Get("dhcpv6_relay_src_intf")
      raObj.Config.SrcIntf = &srcIntf
   }

}

//Helper function to fetch relay info for a given interface
func getRelayAgentInfoForInterface (ifName string, inParams XfmrParams, relayAgentObj *ocbinds.OpenconfigRelayAgent_RelayAgent) error {
   var err error
   var configDb, _ = db.NewDB(getDBOptions(db.ConfigDB))

   if (strings.HasPrefix(inParams.requestUri, "/openconfig-relay-agent:relay-agent/dhcpv6")) {
      getDhcpv6DataFromDb(ifName, relayAgentObj, configDb)
   } else if (strings.HasPrefix(inParams.requestUri, "/openconfig-relay-agent:relay-agent/dhcp")) {
      getDhcpDataFromDb(ifName, relayAgentObj, configDb)
   } else if (strings.HasPrefix(inParams.requestUri, "/openconfig-relay-agent:relay-agent")) {
      //top most level so display v4 and v6 data
      getDhcpDataFromDb(ifName, relayAgentObj, configDb)
      getDhcpv6DataFromDb(ifName, relayAgentObj, configDb)
   }

   configDb.DeleteDB()
   return err
}

//Helper function to fetch relay info for a given interface
func getRelayAgentInfoForAllInterfaces (inParams XfmrParams) error {
   var err error
   relayAgentObj := getRelayAgentRoot(inParams.ygRoot)

   tables := [3]string{"INTERFACE", "VLAN", "PORTCHANNEL_INTERFACE"}
   for _, table := range tables {
       if _, ok := dbIdToTblMap[inParams.curDb]; !ok {
            log.Info("getRelayAgentInfoForAllInterfaces - intf_table_xfmr db id entry not present")
        continue
        }

        var intfKeys []db.Key
        intfKeys, err = inParams.d.GetKeys(&db.TableSpec{Name:table, CompCt:2} )
        for _, intfKey := range intfKeys {
           if(err != nil) {
                continue
           }
       ifName := intfKey.Comp[0]
           ygot.BuildEmptyTree(relayAgentObj)
           ygot.BuildEmptyTree(relayAgentObj.Dhcp)
           ygot.BuildEmptyTree(relayAgentObj.Dhcp.Interfaces)
           ygot.BuildEmptyTree(relayAgentObj.Dhcpv6)
           ygot.BuildEmptyTree(relayAgentObj.Dhcpv6.Interfaces)
       log.Info(ifName)
           err = getRelayAgentInfoForInterface(ifName, inParams, relayAgentObj)
        }
    }
     return err
}

// DbToYang_relay_agent_xfmr - Subtree transformer supports CREATE, UPDATE and DELETE operations, need to write a sub-tree as relay agent doesnt have its own table
var DbToYang_relay_agent_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error

    log.Info("DbToYang_relay_agent_xfmr: ", inParams.uri)
    pathInfo := NewPathInfo(inParams.uri)
    ifName := pathInfo.Var("id")
    targetUriPath := inParams.requestUri

    if strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcp/interfaces/interface") ||
       strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcpv6/interfaces/interface"){
       relayAgentObj := getRelayAgentRoot(inParams.ygRoot)
       err = getRelayAgentInfoForInterface(ifName, inParams, relayAgentObj)
    } else {
       err = getRelayAgentInfoForAllInterfaces(inParams)
    }
    log.Info(pathInfo.Template)
    return err

}

func getIntfIpInfo(dbCl *db.DB, tblName string, ifName string, ipv4 bool, ipv6 bool) bool {
   log.Info("Updating Interface IP Info from DB to Internal DS for Interface Name : ", ifName)

    log.Info(tblName, ifName)
    table := db.TableSpec{Name: tblName, CompCt: 2}
    pattern := db.Key{Comp: []string{ ifName, "*" }}
    keys, err := dbCl.GetKeysPattern(&table, pattern)
 
    if( err != nil) {
        return false
    }
    log.Info(keys)
    for _, key := range keys {
        if len(key.Comp) < 2 {
            continue
        }
        ipB, _, _ := net.ParseCIDR(key.Get(1))
        if ((validIPv4(ipB.String()) && (ipv4 == true)) ||
            (validIPv6(ipB.String()) && (ipv6 == true))) {
                return true           
        }
      }
      return false
}

//Helper function to modify relay info for a given interface
func replaceRelayAgentObjectAttributes (inParams XfmrParams)  error{
   var tblList string
   var err error

   log.Info("replaceRelayAgentObjectAttributes: ", inParams.uri)
   updateMap := make(map[db.DBNum]map[string]map[string]db.Value)

   updateMap[db.ConfigDB] = make(map[string]map[string]db.Value)

   relayAgentObj := getRelayAgentRoot(inParams.ygRoot)
   intfsObj := getIntfsRoot(inParams.ygRoot)
   log.Info("intfsObj:", intfsObj)
 
   if (relayAgentObj.Dhcp != nil  && relayAgentObj.Dhcp.Interfaces != nil && relayAgentObj.Dhcp.Interfaces.Interface != nil) {
       err = replaceDhcpObjectAttributes(inParams, relayAgentObj, updateMap, tblList)
   } 
   if (relayAgentObj.Dhcpv6 != nil  && relayAgentObj.Dhcpv6.Interfaces != nil && relayAgentObj.Dhcpv6.Interfaces.Interface != nil) {
       err = replaceDhcpV6ObjectAttributes(inParams, relayAgentObj, updateMap, tblList)
   }
   inParams.subOpDataMap[UPDATE] = &updateMap
   return err
 }

func replaceDhcpObjectAttributes (inParams XfmrParams, relayAgentObj *ocbinds.OpenconfigRelayAgent_RelayAgent, updateMap map[db.DBNum]map[string]map[string]db.Value, tblList string)  error {
   var helperAddress string
   var index uint8
   var err error

   log.Info("replaceDhcpObjectAttributes")


   for ifName := range relayAgentObj.Dhcp.Interfaces.Interface {

       if ifName == "" {
           errStr := "ifName is NULL"
           err = tlerr.InvalidArgsError{Format: errStr}
           return err
       }

       tblList = getRelayAgentIntfTblByType(ifName)
       log.Info(tblList)

       ipTbl := tblList
       if ( ipTbl == "VLAN") {
          ipTbl = "VLAN_INTERFACE"
       }

       if updateMap[db.ConfigDB][tblList] == nil {
         //allocate only for the first time
         updateMap[db.ConfigDB][tblList] = make(map[string]db.Value)
       }

       _, ok := updateMap[db.ConfigDB][tblList][ifName] 
       if !ok {
          updateMap[db.ConfigDB][tblList][ifName] = db.Value{Field: make(map[string]string)}
       }

       intfObj := relayAgentObj.Dhcp.Interfaces
       intf := intfObj.Interface[ifName]
       log.Info("intf:", intf)

       if (len(intf.Config.HelperAddress) != 0) {
        ipConf := false
        if (strings.HasPrefix(ifName, "Vlan")) {
          entry,_ := inParams.d.GetEntry(&db.TableSpec{Name:"SAG_GLOBAL"}, db.Key{Comp: []string{"IP"}})
          if entry.IsPopulated() {
             status := entry.Field["IPv4"]
             if (status == "enable") {
                sagIPv4Entry, _ := inParams.d.GetEntry(&db.TableSpec{Name:"SAG"}, db.Key{Comp: []string{ifName, "IPv4"}})
                log.Info("sagIPv4Entry:", sagIPv4Entry)
                if sagIPv4Entry.IsPopulated() {
                   if (sagIPv4Entry.Has("gwip@")) {
                      ipConf = true
                   }
                }
             }
          }
        }   
       
        if (!getIntfIpInfo(inParams.d, ipTbl, ifName, true, false)) {
          if !(strings.HasPrefix(ifName, "Vlan") && ipConf) {
             errStr := "Dhcp Relay config not allowed on a non routing interface."
             err :=  tlerr.InvalidArgsError{Format: errStr}
             return err
          }
        }
       }
      
	//helperAddress
	for index = 0; (index < uint8(len(intf.Config.HelperAddress))  && index < 4 && intf.Config.HelperAddress[index] != ""); index++ {
             if !validIPv4(intf.Config.HelperAddress[index]) {
                    errStr := "Invalid IPv4 address" + intf.Config.HelperAddress[index]
                    err = tlerr.InvalidArgsError{Format: errStr}
                    return err
                }

	     if (index ==0) {
		helperAddress = intf.Config.HelperAddress[index]
	     } else {
		helperAddress = helperAddress + "," + intf.Config.HelperAddress[index]
	     }
	}
	updateMap[db.ConfigDB][tblList][ifName].Field["dhcp_servers@"] =  helperAddress
     
	//link-select
	if (intf.Config.LinkSelect == ocbinds.OpenconfigRelayAgent_RelayAgent_Dhcp_Interfaces_Interface_Config_LinkSelect_enable) {
            entry, dbErr := inParams.d.GetEntry(&db.TableSpec{Name:tblList}, db.Key{Comp: []string{ifName}})
            if dbErr != nil {
            log.Info("Failed to read dhcp relay config info from configdb")
            return  dbErr
            }
            log.Info("entry:", entry)
            if (!entry.Has("dhcp_relay_src_intf")) {
               if (intf.Config.SrcIntf == nil) {
                errStr := "Src Intf needs to be configured before enabling link-select"
                err = tlerr.InvalidArgsError{Format: errStr}
                return err
               }
            }
	    updateMap[db.ConfigDB][tblList][ifName].Field["dhcp_relay_link_select"] = "enable"
	} else if (intf.Config.LinkSelect == ocbinds.OpenconfigRelayAgent_RelayAgent_Dhcp_Interfaces_Interface_Config_LinkSelect_disable) {
	    updateMap[db.ConfigDB][tblList][ifName].Field["dhcp_relay_link_select"] = "disable"
	}
     
	//max hop count
	if intf.Config.MaxHopCount!= nil && *intf.Config.MaxHopCount > uint32(0) && *intf.Config.MaxHopCount < uint32(17) {
	   mhCountInt := int(*intf.Config.MaxHopCount)
	   mhCountStr := strconv.Itoa(mhCountInt)
	   updateMap[db.ConfigDB][tblList][ifName].Field["dhcp_relay_max_hop_count"] = mhCountStr
	}
     
	//src intf
	if intf.Config.SrcIntf != nil {
          ipTbl := getRelayAgentIntfTblByType(*intf.Config.SrcIntf)
          if (ipTbl == "VLAN") {
             ipTbl = "VLAN_INTERFACE"
          }
          if (!getIntfIpInfo(inParams.d, ipTbl, *intf.Config.SrcIntf, true, false)) {
          errStr := "No Ip configured on the src intf"
          err :=  tlerr.InvalidArgsError{Format: errStr}
          return err
        }
	   updateMap[db.ConfigDB][tblList][ifName].Field["dhcp_relay_src_intf"] = *intf.Config.SrcIntf
	}
     }
     return err
 }

func replaceDhcpV6ObjectAttributes (inParams XfmrParams, relayAgentObj *ocbinds.OpenconfigRelayAgent_RelayAgent, updateMap map[db.DBNum]map[string]map[string]db.Value, tblList string)  error {
   var helperAddress string
   var index uint8
   var err error

   log.Info("replaceDhcpV6ObjectAttributes")

   for ifName := range relayAgentObj.Dhcpv6.Interfaces.Interface {

       if ifName == "" {
           errStr := "ifName is NULL"
           err = tlerr.InvalidArgsError{Format: errStr}
           return err
       }

       tblList = getRelayAgentIntfTblByType(ifName)
       log.Info(tblList)

       if updateMap[db.ConfigDB][tblList] == nil {
         //allocate only for the first time
         updateMap[db.ConfigDB][tblList] = make(map[string]db.Value)
       }

       _, ok := updateMap[db.ConfigDB][tblList][ifName] 
       if !ok {
          updateMap[db.ConfigDB][tblList][ifName] = db.Value{Field: make(map[string]string)}
       }

       intfObj := relayAgentObj.Dhcpv6.Interfaces
       intf := intfObj.Interface[ifName]

       log.Info("intf:", intf)

       if (len(intf.Config.HelperAddress) != 0) {
        ipConf := false
        if (strings.HasPrefix(ifName, "Vlan")) {
          entry,_ := inParams.d.GetEntry(&db.TableSpec{Name:"SAG_GLOBAL"}, db.Key{Comp: []string{"IP"}})
          log.Info("entry:", entry)
          if entry.IsPopulated() {
             status := entry.Field["IPv6"]
             if (status == "enable") {
                sagIPv6Entry, _ := inParams.d.GetEntry(&db.TableSpec{Name:"SAG"}, db.Key{Comp: []string{ifName, "IPv6"}})
                log.Info("sagIPv6Entry:", sagIPv6Entry)
                if sagIPv6Entry.IsPopulated() {
                   if (sagIPv6Entry.Has("gwip@")) {
                      ipConf = true
                   }
                }
             }
          }
        }
        ipTbl := tblList
        if (ipTbl == "VLAN") {
          ipTbl = "VLAN_INTERFACE"
        }
        if (!getIntfIpInfo(inParams.d, ipTbl, ifName, false, true)) {
          if !(strings.HasPrefix(ifName, "Vlan") && ipConf) {
             errStr := "IPv6 address not configured on this interface."
             err :=  tlerr.InvalidArgsError{Format: errStr}
             return err
          }
        }
       } 
       //helperAddressV6
       for index = 0; (index < uint8(len(intf.Config.HelperAddress))  && index < 4 && intf.Config.HelperAddress[index] != ""); index++ {
            if !validIPv6(intf.Config.HelperAddress[index]) {
                    errStr := "Invalid IPv6 address " + intf.Config.HelperAddress[index]
                    err = tlerr.InvalidArgsError{Format: errStr}
                    log.Info("Invalid ipv6 address")
                    return err
                }

            if (index ==0) {
       	helperAddress = intf.Config.HelperAddress[index]
            } else {
       	helperAddress = helperAddress + "," + intf.Config.HelperAddress[index]
            }
       }
       updateMap[db.ConfigDB][tblList][ifName].Field["dhcpv6_servers@"] = helperAddress

      //max hop count
       if intf.Config.MaxHopCount != nil && *intf.Config.MaxHopCount > uint32(0) && *intf.Config.MaxHopCount < uint32(17) {
	  mhCountInt := int(*intf.Config.MaxHopCount)
	  mhCountStr := strconv.Itoa(mhCountInt)
	  updateMap[db.ConfigDB][tblList][ifName].Field["dhcpv6_relay_max_hop_count"] =  mhCountStr
       } 
    
       //src intf
       if intf.Config.SrcIntf != nil {
           tblName := getRelayAgentIntfTblByType(*intf.Config.SrcIntf)
           if (tblName == "VLAN") {
              tblName = "VLAN_INTERFACE"
           }
           if (!getIntfIpInfo(inParams.d, tblName, *intf.Config.SrcIntf, false, true)) {
           errStr := "No IPv6 address configured on the src intf"
           err :=  tlerr.InvalidArgsError{Format: errStr}
           return err
           }

	   updateMap[db.ConfigDB][tblList][ifName].Field["dhcpv6_relay_src_intf"] = *intf.Config.SrcIntf
       }
    
    }
    return err
 }
    
//Function to delete config for an interface level
func deleteRelayAgentObjectAttributes(inParams XfmrParams, ifName string) error {
   var tblList string   
   var fieldStr [] string
   var configDb, _ = db.NewDB(getDBOptions(db.ConfigDB))
   var helperAddress string
   var index uint8
   var err error

   targetUriPath := inParams.requestUri

   log.Info("deleteRelayAgentObjectAttributes: ", inParams.uri)

   if ifName == "" {
       errStr := "deleteRelayAgentObjectAttributes - ifName is NULL"
       err = tlerr.InvalidArgsError{Format: errStr}
       return err
   }
 
   relayAgentObj := getRelayAgentRoot(inParams.ygRoot)
   tblList = getRelayAgentIntfTblByType(ifName)
   log.Info(tblList)

   entry, dbErr := configDb.GetEntry(&db.TableSpec{Name:tblList}, db.Key{Comp: []string{ifName}})
   configDb.DeleteDB()
   if dbErr != nil {
     errStr := "Failed to read mgmt port status from config DB, " + tblList + " " + ifName
     err = tlerr.InvalidArgsError{Format: errStr}
     return err
   }
   
   deleteMap := make(map[db.DBNum]map[string]map[string]db.Value)
   deleteMap[db.ConfigDB] = make(map[string]map[string]db.Value)
   deleteMap[db.ConfigDB][tblList] = make(map[string]db.Value)

   deleteMap[db.ConfigDB][tblList][ifName] = db.Value{Field: make(map[string]string)}

   //check for attribute level deletes followed by interface level followed by both v4 and v6 params delete
   if  strings.Contains(targetUriPath, "dhcpv6") && strings.Contains(targetUriPath, "helper-address"){
      if (relayAgentObj.Dhcpv6 != nil  && relayAgentObj.Dhcpv6.Interfaces != nil && relayAgentObj.Dhcpv6.Interfaces.Interface != nil) {
          //We have a specific address to delete - delete only that address
	  intf := relayAgentObj.Dhcpv6.Interfaces.Interface[ifName]
	  helperAddress = ""
   	  for index = 0; (index < uint8(len(intf.Config.HelperAddress))  && index < 4 && intf.Config.HelperAddress[index] != ""); index++ {
	     if (index == 0) {
		helperAddress = intf.Config.HelperAddress[index]
	     } else {
		helperAddress = helperAddress + "," + intf.Config.HelperAddress[index]
	     }
	  }
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcpv6_servers@"] = helperAddress
      } else {
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcpv6_servers@"] = ""
      }
   } else if  strings.Contains(targetUriPath, "dhcpv6") && strings.Contains(targetUriPath, "src-intf"){
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcpv6_relay_src_intf"] = ""
   } else if strings.Contains(targetUriPath, "dhcpv6") && strings.Contains(targetUriPath, "max-hop-count"){
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcpv6_relay_max_hop_count"] = ""
   } else if  strings.Contains(targetUriPath, "dhcp") && strings.Contains(targetUriPath, "helper-address"){
      if (relayAgentObj.Dhcp != nil  && relayAgentObj.Dhcp.Interfaces != nil && relayAgentObj.Dhcp.Interfaces.Interface != nil) {
          //We have a specific address to delete - delete only that address
	  intf := relayAgentObj.Dhcp.Interfaces.Interface[ifName]
   	  for index = 0; (index < uint8(len(intf.Config.HelperAddress))  && index < 4 && intf.Config.HelperAddress[index] != ""); index++ {
	     if (index ==0) {
		helperAddress = intf.Config.HelperAddress[index]
	     } else {
		helperAddress = helperAddress + "," + intf.Config.HelperAddress[index]
	     }
	  }
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcp_servers@"] = helperAddress
      } else {
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcp_servers@"] = ""
      }
   } else if  strings.Contains(targetUriPath, "dhcp") && strings.Contains(targetUriPath, "src-intf"){
       entry, _ := inParams.d.GetEntry(&db.TableSpec{Name:tblList}, db.Key{Comp: []string{ifName}})
       log.Info("entry:", entry)
       if entry.IsPopulated() {
          if (entry.Field["dhcp_relay_link_select"]) == "enable" {
             errStr := "Cannot remove src-intf when link-select is enabled"
             err = tlerr.InvalidArgsError{Format: errStr}
             return err
          }
       }
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcp_relay_src_intf"] = ""
   } else if strings.Contains(targetUriPath, "dhcp") && strings.Contains(targetUriPath, "max-hop-count"){
       deleteMap[db.ConfigDB][tblList][ifName].Field["dhcp_relay_max_hop_count"] = ""
   } else if  strings.Contains(targetUriPath, "dhcp") && strings.Contains(targetUriPath, "link-select"){
     deleteMap[db.ConfigDB][tblList][ifName].Field["dhcp_relay_link_select"] = ""
   } else if strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcp") {
     //delete interface level attributes for DHCP
     fieldStr = append (relayAgentFields)
   } else if strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcpv6") {   
     fieldStr = append (relayAgentV6Fields)
   } else {
     fieldStr = append(relayAgentFields, relayAgentV6Fields...)
   }

   for _, field := range (fieldStr) {
       if entry.Has(field) {
           deleteMap[db.ConfigDB][tblList][ifName].Field[field] = ""
       }
   }

    _, ok := deleteMap[db.ConfigDB][tblList][ifName] 
    if ok {
       log.Info(deleteMap)
       inParams.subOpDataMap[DELETE] = &deleteMap
    } else {
      log.Warning("Delete map was not populated")
    }
  return err
}

//Allocate memory only when required
func allocateMemoryMaps(ifName string, table string, deleteMap map[db.DBNum]map[string]map[string]db.Value){
   if deleteMap[db.ConfigDB][table] == nil {
      //do not create the map if there is no dhcp entry
      deleteMap[db.ConfigDB][table] = make(map[string]db.Value)
   }
   deleteMap[db.ConfigDB][table][ifName] = db.Value{Field: make(map[string]string)}
}

//Helper function to fetch relay info for a given interface
func deleteAllIntfsRelayAgentObjectAttributes(inParams XfmrParams) error {
   var err error
   var fieldStr [] string

   targetUriPath := inParams.requestUri

   deleteMap := make(map[db.DBNum]map[string]map[string]db.Value)
   deleteMap[db.ConfigDB] = make(map[string]map[string]db.Value)

   tables := [3]string{"INTERFACE", "VLAN", "PORTCHANNEL_INTERFACE"}
   for _, table := range tables {
       intfTble, err := inParams.d.GetTable(&db.TableSpec{Name:table})
       if err != nil {
          continue
       }

       intfKeys, err := intfTble.GetKeys()
       for _, intfName := range intfKeys {
           intfEntry, err := intfTble.GetEntry(intfName)
           if(err != nil) {
               continue
           }

       if intfEntry.Has("dhcp_servers@") || intfEntry.Has("dhcpv6_servers@") {
          //delete only if there is provisioning
          ifName := intfName.Comp[0]
          log.Info(intfName)

         //check for top level deletes
         if (targetUriPath == "/openconfig-relay-agent:relay-agent/dhcpv6") || 
	    (targetUriPath == "/openconfig-relay-agent:relay-agent/dhcpv6/interfaces") && 
	    intfEntry.Has("dhcpv6_servers@") {  
	    //need to allocate memory after checking the uri, otherwise when the request is for v6
	    //interfaces which have v4 will get deleted too
	    allocateMemoryMaps(ifName, table, deleteMap) 
            fieldStr = append (relayAgentV6Fields) 
         } else if (targetUriPath == "/openconfig-relay-agent:relay-agent/dhcp") ||
	    (targetUriPath == "/openconfig-relay-agent:relay-agent/dhcp/interfaces")  &&
	    intfEntry.Has("dhcp_servers@") {
	    allocateMemoryMaps(ifName, table, deleteMap) 
            fieldStr = append(relayAgentFields)
         } else if (targetUriPath == "/openconfig-relay-agent:relay-agent") {   
	    allocateMemoryMaps(ifName, table, deleteMap) 
            fieldStr = append(relayAgentFields, relayAgentV6Fields...)
         } else {
	   log.Error("Incorrect Uri")
	   err = errors.New("Invalid URI : " + targetUriPath)
	 }

         for _, field := range (fieldStr) {
            if intfEntry.Has(field) {
	       if field == "dhcp_relay_link_select" {
                  deleteMap[db.ConfigDB][table][ifName].Field[field] = "disable"
	       } else {
                 deleteMap[db.ConfigDB][table][ifName].Field[field] = ""
	       }
            }
         }
       }
    }
  }

  inParams.subOpDataMap[DELETE] = &deleteMap
  return err
}


var YangToDb_relay_agent_xfmr SubTreeXfmrYangToDb = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {
   var err error
   res_map := make(map[string]map[string]db.Value)
   targetUriPath := inParams.requestUri
   
   log.Info("YangToDb_relay_agent_xfmr: ", inParams.uri)

   switch inParams.oper {
        case CREATE:
            fallthrough
        case REPLACE:
	    fallthrough
        case UPDATE:
              err = replaceRelayAgentObjectAttributes(inParams)
        case DELETE:
            if  strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcp/interfaces/interface") ||
                strings.HasPrefix(targetUriPath, "/openconfig-relay-agent:relay-agent/dhcpv6/interfaces/interface"){
	        //Delete for specific interface or attribute
	        pathInfo := NewPathInfo(inParams.uri)
                ifName := pathInfo.Var("id")
                err = deleteRelayAgentObjectAttributes(inParams, ifName)
	    } else {
	       //Delete for all interfaces
               //inParams.subOpDataMap = make(map[int]*map[db.DBNum]map[string]map[string]db.Value)
               err = deleteAllIntfsRelayAgentObjectAttributes(inParams)
	    }
   } 
   
    return res_map, err
}

