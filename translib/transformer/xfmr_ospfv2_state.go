package transformer

import (
    "errors"
    "strings"
    "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
    "fmt"
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "github.com/Azure/sonic-mgmt-common/translib/tlerr"
    "github.com/openconfig/ygot/ygot"
    "encoding/json"
    log "github.com/golang/glog"
)


func init () {

    XlateFuncBind("DbToYang_ospfv2_global_state_xfmr", DbToYang_ospfv2_global_state_xfmr)
    XlateFuncBind("DbToYang_ospfv2_global_timers_spf_state_xfmr", DbToYang_ospfv2_global_timers_spf_state_xfmr)
    XlateFuncBind("DbToYang_ospfv2_global_timers_lsa_generation_state_xfmr", DbToYang_ospfv2_global_timers_lsa_generation_state_xfmr)
    XlateFuncBind("DbToYang_ospfv2_areas_area_state_xfmr", DbToYang_ospfv2_areas_area_state_xfmr)
    XlateFuncBind("DbToYang_ospfv2_neighbors_state_xfmr", DbToYang_ospfv2_neighbors_state_xfmr)
    XlateFuncBind("DbToYang_ospfv2_vlink_state_xfmr", DbToYang_ospfv2_vlink_state_xfmr)
    XlateFuncBind("DbToYang_ospfv2_stub_state_xfmr", DbToYang_ospfv2_stub_state_xfmr)
    XlateFuncBind("DbToYang_ospfv2_route_table_xfmr", DbToYang_ospfv2_route_table_xfmr)
    XlateFuncBind("ospfv2_router_area_tbl_xfmr", ospfv2_router_area_tbl_xfmr)

    XlateFuncBind("rpc_clear_ospfv2", rpc_clear_ospfv2)
}


func ospfv2_fill_only_global_state (output_state map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2) error {
    var err error
    var ospfv2Gbl_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global
    var ospfv2GblState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global_State
    var numint64 int64
    var ospfv2Zero bool = false
    var ospfv2One bool = true
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Global State"

    log.Infof("ospfv2_fill_only_global_state - start")

    ospfv2Gbl_obj = ospfv2_obj.Global
    if ospfv2Gbl_obj == nil {
        log.Errorf("%s failed !! Error: OSPFv2-Global container missing", cmn_log)
        return  oper_err
    }

    ospfv2GblState_obj = ospfv2Gbl_obj.State
    if ospfv2GblState_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2-Global-State container missing", cmn_log)
        return oper_err
    }

    if _routerId,ok := output_state["routerId"].(string); ok {
        ospfv2GblState_obj.RouterId = &_routerId
    }

    if _rfc1583Compatibility,ok := output_state["rfc1583Compatibility"].(bool); ok {
        ospfv2GblState_obj.OspfRfc1583Compatible = &_rfc1583Compatibility
    }

    if _opaqueCapable,ok := output_state["opaqueCapable"].(bool); ok {
        ospfv2GblState_obj.OpaqueLsaCapability = &_opaqueCapable
    }

    if value,ok := output_state["holdtimeMultplier"] ; ok {
        _holdtime_multiplier := uint32(value.(float64))
        ospfv2GblState_obj.HoldTimeMultiplier = &_holdtime_multiplier
    }

    if value,ok := output_state["spfLastExecutedMsecs"]; ok {
        _spfLastExecutedMsecs  := uint64(value.(float64))
        ospfv2GblState_obj.LastSpfExecutionTime = &_spfLastExecutedMsecs
    }

    if value,ok := output_state["spfLastDurationUsecs"] ; ok {
        _spfLastDurationUsecs   := uint32(value.(float64))
        ospfv2GblState_obj.LastSpfDuration = &_spfLastDurationUsecs
    }

    if value,ok := output_state["writeMultiplier"] ; ok {
        _write_multiplier := uint8(value.(float64))
        ospfv2GblState_obj.WriteMultiplier = &_write_multiplier
    }
    if value,ok := output_state["lsaExternalCounter"] ; ok {
        _lsaExternalCounter := uint32(value.(float64))
        ospfv2GblState_obj.ExternalLsaCount = &_lsaExternalCounter
    }
    if value,ok := output_state["lsaAsopaqueCounter"] ; ok {
        _lsaAsopaqueCounter := uint32(value.(float64))
        ospfv2GblState_obj.OpaqueLsaCount = &_lsaAsopaqueCounter
    }
    if value,ok := output_state["lsaExternalChecksum"]; ok {
        numint64 = int64(value.(float64))
        numstr := fmt.Sprintf("0x%08x", numint64)
        ospfv2GblState_obj.ExternalLsaChecksum = &numstr
    }
    if value,ok := output_state["lsaAsOpaqueChecksum"]; ok {
        numint64 = int64(value.(float64))
        numstr := fmt.Sprintf("0x%08x", numint64)
        ospfv2GblState_obj.OpaqueLsaChecksum = &numstr
    }
    if value,ok := output_state["attachedAreaCounter"] ; ok {
        _attachedAreaCounter  := uint32(value.(float64))
        ospfv2GblState_obj.AreaCount = &_attachedAreaCounter
    }
    ospfv2GblState_obj.AbrType = ocbinds.OpenconfigOspfv2Ext_OSPF_ABR_TYPE_UNSET
    if _abrtype,ok := output_state["abrType"].(string); ok {
        if _abrtype == "Alternative Cisco" {
            ospfv2GblState_obj.AbrType = ocbinds.OpenconfigOspfv2Ext_OSPF_ABR_TYPE_CISCO
        }
        if _abrtype == "Alternative IBM" {
            ospfv2GblState_obj.AbrType = ocbinds.OpenconfigOspfv2Ext_OSPF_ABR_TYPE_IBM
        }
        if _abrtype == "Alternative Shortcut" {
            ospfv2GblState_obj.AbrType = ocbinds.OpenconfigOspfv2Ext_OSPF_ABR_TYPE_SHORTCUT
        }
        if _abrtype == "Standard (RFC2328)" {
            ospfv2GblState_obj.AbrType = ocbinds.OpenconfigOspfv2Ext_OSPF_ABR_TYPE_STANDARD
        }
    }

    if _stubAdvertisement, ok := output_state["stubAdvertisement"].(bool); ok {
        if (!_stubAdvertisement) {
            ospfv2GblState_obj.StubAdvertisement = &ospfv2Zero
        } else {
            ospfv2GblState_obj.StubAdvertisement = &ospfv2One
        }
    }
    
    if value,ok := output_state["preShutdownEnabledSecs"] ; ok {
        _preShutdownEnabledSecs := uint32(value.(float64))
        ospfv2GblState_obj.PreShutdownEnabledSecs = &_preShutdownEnabledSecs
    }
    if value,ok := output_state["postStartEnabledSecs"] ; ok {
        _postStartEnabledSecs := uint32(value.(float64))
        ospfv2GblState_obj.PostStartEnabledSecs = &_postStartEnabledSecs
    }
    return err
}


func ospfv2_fill_global_timers_spf_state (output_state map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2) error {
    var err error
    var ospfv2Gbl_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global
    var ospfv2Timers_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global_Timers
    var ospfv2GblTimersSpfState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global_Timers_Spf_State 
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Global State"

    log.Infof("ospfv2_fill_global_timers_spf_state - start")

    ospfv2Gbl_obj = ospfv2_obj.Global
    if ospfv2Gbl_obj == nil {
        log.Errorf("%s failed !! Error: OSPFv2-Global container missing", cmn_log)
        return  oper_err
    }

    if nil == ospfv2Gbl_obj.Timers {
        log.Info("OSPF global Timers is nil")
        ospfv2Timers_obj = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global_Timers)
        if nil == ospfv2Timers_obj {
            log.Errorf("%s failed !! Error: Failed to create timers Tree under global", cmn_log)
            return oper_err
        }
        ygot.BuildEmptyTree(ospfv2Timers_obj)
        ospfv2Gbl_obj.Timers = ospfv2Timers_obj
    }

    if nil == ospfv2Gbl_obj.Timers.Spf {
        log.Info("OSPF global Timers Spf is nil")
        ygot.BuildEmptyTree(ospfv2Gbl_obj.Timers.Spf)
    }

    ospfv2GblTimersSpfState_obj = ospfv2Gbl_obj.Timers.Spf.State
    if ospfv2GblTimersSpfState_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2-Global-State container missing", cmn_log)
        return  oper_err
    }

    if value,ok := output_state["spfScheduleDelayMsecs"]; ok {
        _throttle_delay := uint32(value.(float64))
        ospfv2GblTimersSpfState_obj.ThrottleDelay = &_throttle_delay
    }

    if value,ok := output_state["holdtimeMinMsecs"] ; ok {
        _holdtime_minMsec := uint32(value.(float64))
        ospfv2GblTimersSpfState_obj.InitialDelay = &_holdtime_minMsec
    }
    
    if value,ok := output_state["holdtimeMaxMsecs"] ; ok {
        _holdtime_maxMsec := uint32(value.(float64))
        ospfv2GblTimersSpfState_obj.MaximumDelay = &_holdtime_maxMsec
    }

    var _spfTimerDueInMsecs uint32 = 0
    ospfv2GblTimersSpfState_obj.SpfTimerDue = &_spfTimerDueInMsecs
    if value,ok := output_state["spfTimerDueInMsecs"] ; ok {
        _spfTimerDueInMsecs = uint32(value.(float64))
        ospfv2GblTimersSpfState_obj.SpfTimerDue = &_spfTimerDueInMsecs
    }

    return err
}
func ospfv2_fill_route_table (ospf_info map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2) error {
    var err error
    var ospfv2RouteTables_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables
    var ospfv2RouteTable_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable
    var ospfv2RouteTableListState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable_State
    var ospfv2RouteTableState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable_State_RouteTableState
    var prefixStr string
    var ospfv2Route *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable_State_RouteTableState_Route
    var ospfv2RouteState *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable_State_RouteTableState_Route_State
    var ospfv2Nexthop *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable_State_RouteTableState_Route_State_NextHopsList_NextHops
    var ospfv2NexthopState *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable_State_RouteTableState_Route_State_NextHopsList_NextHops_State
    var nexthop_ip, nexthop_ifname string
    var ospfv2Zero bool = false
    var ospfv2One bool = true
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF Route Table"

    log.Infof("ospfv2_fill_route_table - start")

    ospfv2RouteTables_obj = ospfv2_obj.RouteTables
    if ospfv2RouteTables_obj == nil {
        log.Errorf("%s failed !! Error: OSPFv2 Route Tables container missing", cmn_log)
        return  oper_err
    }
    if nil == ospfv2RouteTables_obj.RouteTable {
        log.Info("Creating route table for router LSA")
        _, err = ospfv2RouteTables_obj.NewRouteTable(ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TABLE_ROUTER_ROUTE_TABLE)
        if nil != err { 
            log.Errorf("%s failed !! Error: Creating route table for router LSA failed", cmn_log)
            return  oper_err
        }
        log.Info("Creating route table for Network LSA")
        _, err = ospfv2RouteTables_obj.NewRouteTable(ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TABLE_NETWORK_ROUTE_TABLE)
        if nil != err { 
            log.Errorf("%s failed !! Error: Creating route table for Network LSA failed", cmn_log)
            return  oper_err
        }
        log.Info("Creating route table for external LSA")
        _, err = ospfv2RouteTables_obj.NewRouteTable(ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TABLE_EXTERNAL_ROUTE_TABLE)
        if nil != err { 
            log.Errorf("%s failed !! Error: Creating route table for external LSA failed", cmn_log)
            return  oper_err
        }
    }
    for key,value := range ospf_info {
        if (key == "vrfId" || key == "vrfName") {
            log.Infof("Skipping key with name %s, as it is not useful", key)
            continue;
        }
        route_info := value.(map[string]interface{})
        prefixStr = fmt.Sprintf("%v", key)
        log.Infof("Prefix string %s", prefixStr)
        switch(route_info["routeType"]) {
            case "R " : 
                ospfv2RouteTable_obj = ospfv2RouteTables_obj.RouteTable[ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TABLE_ROUTER_ROUTE_TABLE]
            case "N" :
                ospfv2RouteTable_obj = ospfv2RouteTables_obj.RouteTable[ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TABLE_NETWORK_ROUTE_TABLE]
            case "N E2" :
                ospfv2RouteTable_obj = ospfv2RouteTables_obj.RouteTable[ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TABLE_EXTERNAL_ROUTE_TABLE] 
            case "N E1" :
                ospfv2RouteTable_obj = ospfv2RouteTables_obj.RouteTable[ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TABLE_EXTERNAL_ROUTE_TABLE] 
            default:
                ospfv2RouteTable_obj = nil
        }
        if (nil == ospfv2RouteTable_obj) {
            log.Errorf("failed !! Error: RouteTable not found for routeType %s", route_info["routeType"])
            return oper_err
        }
        if nil == ospfv2RouteTable_obj.State {
            log.Info("Routetable state information is missing");
            ospfv2RouteTableListState_obj = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_RouteTables_RouteTable_State)
            if nil == ospfv2RouteTableListState_obj {
                log.Info("Failed to create State information for route Table list state")
                return oper_err
            }
            ygot.BuildEmptyTree(ospfv2RouteTableListState_obj)
            ospfv2RouteTable_obj.State = ospfv2RouteTableListState_obj
        }
        ospfv2RouteTableState_obj = ospfv2RouteTableListState_obj.RouteTableState
        if nil == ospfv2RouteTableState_obj.Route {
            ospfv2Route, err = ospfv2RouteTableState_obj.NewRoute(prefixStr) 
            ygot.BuildEmptyTree(ospfv2Route)
        } else {
            ospfv2Route = ospfv2RouteTableState_obj.Route[prefixStr]
            if nil == ospfv2Route {
                ospfv2Route, err = ospfv2RouteTableState_obj.NewRoute(prefixStr)
                ygot.BuildEmptyTree(ospfv2Route)
            }
        }
        if nil == ospfv2Route {
            log.Errorf(" failed !! Error,  prefix %s cannot be added in route table tree", prefixStr)
            return  oper_err
        }  
        ospfv2RouteState = ospfv2Route.State
        if value,ok := route_info["cost"] ; ok {
            _cost  := uint32(value.(float64))
            ospfv2RouteState.Cost = &_cost
        }
        
        if value,ok := route_info["type2_cost"] ; ok {
            _type2cost  := uint32(value.(float64))
            ospfv2RouteState.Type2Cost = &_type2cost
        }
        
        if value,ok := route_info["nexthops"] ; ok {
            nexthops := value.([]interface{})
            for _, value = range nexthops {
                nexthop := value.(map[string]interface{})
                if _intf_name, ok := nexthop["via"].(string); ok {
                    nexthop_ifname = fmt.Sprintf("%v",_intf_name)
                }
                if _ip, ok := nexthop["ip"].(string); ok {
                    nexthop_ip = fmt.Sprintf("%v",_ip)
                }
                if _direct_intf, ok := nexthop["directly attached to"].(string); ok {
                    nexthop_ifname = fmt.Sprintf("%v",_direct_intf)
                    nexthop_ip = "0.0.0.0"
                }
                ospfv2Nexthop, err = ospfv2RouteState.NextHopsList.NewNextHops(nexthop_ip, nexthop_ifname)
                if nil != ospfv2Nexthop {
                    ygot.BuildEmptyTree(ospfv2Nexthop)
                    ospfv2NexthopState = ospfv2Nexthop.State
                    if _area_id, ok := route_info["area"].(string); ok {
                        ospfv2NexthopState.AreaId = &_area_id
                    } else {
                        if area_id, ok := nexthop["area"].(string); ok {
                            ospfv2NexthopState.AreaId = &area_id
                        }      
                    }
                }
            }
        }
        if _ia, ok := route_info["IA"].(bool); ok {
            if !_ia {
                ospfv2RouteState.InterArea = &ospfv2Zero
            } else {
                ospfv2RouteState.InterArea = &ospfv2One
            }
        }
        if _routertype, ok := route_info["routerType"].(string); ok {
            if _routertype == "abr" {
                ospfv2RouteState.RouterType = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTER_TYPE_ABR
            }
            if _routertype == "asbr" {
                ospfv2RouteState.RouterType = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTER_TYPE_ASBR
            }
            if _routertype == "abr asbr" {
                ospfv2RouteState.RouterType = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTER_TYPE_ABRASBR
            }
        }
        switch(route_info["routeType"]) {
            case "R " :
                ospfv2RouteState.Type = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TYPE_ROUTER_ROUTE 
            case "N" : 
                ospfv2RouteState.Type = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TYPE_NETWORK_ROUTE 
            case "N E2" :
                ospfv2RouteState.Type = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TYPE_EXTERNAL_ROUTE 
                ospfv2RouteState.SubType = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_PATH_TYPE_EXTERNAL_ROUTE_TYPE_2
            case "N E1" :
                ospfv2RouteState.Type = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TYPE_EXTERNAL_ROUTE 
                ospfv2RouteState.SubType = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_PATH_TYPE_EXTERNAL_ROUTE_TYPE_1
            default:
                ospfv2RouteState.Type = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_TYPE_UNSET 
                ospfv2RouteState.SubType = ocbinds.OpenconfigOspfv2Ext_OSPF_ROUTE_PATH_TYPE_UNSET
        }
    }
    return err
}

func ospfv2_fill_global_timers_lsa_generation_state (output_state map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2) error {
    var err error
    var ospfv2Gbl_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global
    var ospfv2Timers_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global_Timers
    var ospfv2GblTimersLsaGenState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global_Timers_LsaGeneration_State
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Global State"

    log.Infof("ospfv2_fill_neighbors_state - start")

    ospfv2Gbl_obj = ospfv2_obj.Global
    if ospfv2Gbl_obj == nil {
        log.Errorf("%s failed !! Error: OSPFv2-Global container missing", cmn_log)
        return  oper_err
    }

    if nil == ospfv2Gbl_obj.Timers {
        log.Info("OSPF global Timers is nil")
        ospfv2Timers_obj = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Global_Timers)
        if nil == ospfv2Timers_obj {
            log.Errorf("%s failed !! Error: Failed to create timers Tree under global", cmn_log)
            return oper_err
        }
        ygot.BuildEmptyTree(ospfv2Timers_obj)
        ospfv2Gbl_obj.Timers = ospfv2Timers_obj
    }

    if nil == ospfv2Gbl_obj.Timers.LsaGeneration {
        log.Info("OSPF global Timers LsaGeneration is nil")
        ygot.BuildEmptyTree(ospfv2Gbl_obj.Timers.LsaGeneration)
    }

    ospfv2GblTimersLsaGenState_obj = ospfv2Gbl_obj.Timers.LsaGeneration.State
    if ospfv2GblTimersLsaGenState_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2-Global-Timers Lsa generation State container missing", cmn_log)
        return  oper_err
    }

    if value,ok := output_state["lsaMinIntervalMsecs"] ; ok {
        _lsaMinIntervalMsecs := uint32(value.(float64))
        ospfv2GblTimersLsaGenState_obj.LsaMinIntervalTimer = &_lsaMinIntervalMsecs
    }
    if value,ok := output_state["lsaMinArrivalMsecs"] ; ok {
        _lsaMinArrivalMsecs  := uint32(value.(float64))
        ospfv2GblTimersLsaGenState_obj.LsaMinArrivalTimer = &_lsaMinArrivalMsecs
    }
    if value,ok := output_state["refreshTimerMsecs"] ; ok {
        _refreshTimerMsecs     := uint32(value.(float64))
        ospfv2GblTimersLsaGenState_obj.RefreshTimer = &_refreshTimerMsecs
    }
    
    return err
}
func ospfv2_find_area_by_key(ospfv2Areas_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas, 
areaNameStr string) (*ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area, error) {
    var err error
    var ospfv2AreaKey1 *string
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    log.Infof("Entered ospfv2_find_area_by_key %s", areaNameStr)
    if ((nil == ospfv2Areas_obj) || (nil == ospfv2Areas_obj.Area)) {
        return nil, err
    }
    for _, ospfv2Area_obj = range ospfv2Areas_obj.Area {
        ospfv2AreaKey1 = ospfv2Area_obj.Identifier  
        log.Info("Key are ", ospfv2AreaKey1, areaNameStr)
        if(*ospfv2AreaKey1 == areaNameStr) {
            log.Info("Match found")
            return ospfv2Area_obj, nil
        }
    }
    return nil, err
} 
func ospfv2_create_new_area(ospfv2Areas_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas, 
areaNameStr string) (*ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area, error) {
    var err  error
    var ok   bool
    oper_err := errors.New("Operational error")
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    log.Infof("Entered ospfv2_create_new_area %s", areaNameStr)
    if (nil == ospfv2Areas_obj) {
        return nil, oper_err
    }    
    if  ospfv2Area_obj, ok = ospfv2Areas_obj.Area[areaNameStr]; !ok {
        ospfv2Area_obj, err = ospfv2Areas_obj.NewArea(areaNameStr)
        if (err != nil) {
            log.Info("Failed to create a new area")
            return  nil, err
        }
        ygot.BuildEmptyTree(ospfv2Area_obj)
    }
        
    ospfv2Area_obj.Config.Identifier = &areaNameStr
    return ospfv2Area_obj, err
} 
func ospfv2_get_or_create_area (output_state map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2, area_id interface{}, vrfName interface{}) (*ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area, map[string]interface{}, error) {
    var err error
    var ospfv2Areas_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    oper_err := errors.New("Operational error")
    area_not_found_err := errors.New("Area not found")
    cmn_log := "function: ospfv2_get_or_create_area"
    var areaNameStr string

    ospfv2Areas_obj = ospfv2_obj.Areas
    if ospfv2Areas_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2 areas list missing", cmn_log)
        return  nil, nil, oper_err
    }

    if value, ok := output_state["areas"]; ok {
        areas_map := value.(map[string]interface {})
        for key, area := range areas_map {
            area_info := area.(map[string]interface{})
            if (key != area_id) {
                log.Infof("Skip filling area state information for area %s", key)
                continue;
            }
            areaNameStr = fmt.Sprintf("%v",key)
            ospfv2Area_obj, err = ospfv2_find_area_by_key(ospfv2Areas_obj, areaNameStr)
            if nil == ospfv2Area_obj {
                log.Infof("Area object missing, add new area=%s", area_id)
                ospfv2Area_obj, err = ospfv2_create_new_area(ospfv2Areas_obj, areaNameStr)
                if (err != nil) {
                    log.Info("Failed to create a new area")
                    return  nil, nil, oper_err
                }
            }
            return ospfv2Area_obj, area_info, err
        }
    }
    return nil, nil, area_not_found_err
}

func ospfv2_fill_area_state (output_state map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2, area_id interface{}, vrfName interface{}) error {
    var err error
    var ospfv2Areas_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    var ospfv2AreaInfo_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_State
    var ospfv2Zero bool = false
    var ospfv2One bool = true
    var numint64 int64
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Areas-Area State"
    var areaNameStr string

    log.Infof("ospfv2_fill_area_state - start")

    ospfv2Areas_obj = ospfv2_obj.Areas
    if ospfv2Areas_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2 areas list missing", cmn_log)
        return  oper_err
    }

    if value, ok := output_state["areas"]; ok {
        areas_map := value.(map[string]interface {})
        for key, area := range areas_map {
            area_info := area.(map[string]interface{})
            if (key != area_id) {
                log.Infof("Skip filling area state information for area %s", key)
                continue;
            }
            areaNameStr = fmt.Sprintf("%v",key)
            ospfv2Area_obj, err = ospfv2_find_area_by_key(ospfv2Areas_obj, areaNameStr)
            if nil == ospfv2Area_obj {
                log.Infof("Area object missing, add new area=%s", area_id)
                ospfv2Area_obj, err = ospfv2_create_new_area(ospfv2Areas_obj, areaNameStr)
                if (err != nil) {
                    log.Info("Failed to create a new area")
                    return  oper_err
                }
            }
            ospfv2AreaInfo_obj = ospfv2Area_obj.State
            if ospfv2AreaInfo_obj == nil {
                log.Infof("Area State missing, add new area state for area %s", area_id)
                ospfv2AreaInfo_obj = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_State)
                ygot.BuildEmptyTree (ospfv2AreaInfo_obj)
                if ospfv2AreaInfo_obj == nil {
                    log.Errorf("%s failed !! Error: Area information missing", cmn_log)
                    return  oper_err
                }
                ospfv2Area_obj.State = ospfv2AreaInfo_obj
            }
            
            if _authtype,ok := area_info["authentication"].(string); ok {
                if _authtype == "authenticationNone" {
                    authType := "no"
                    ospfv2AreaInfo_obj.AuthenticationType = &authType 
                }
                if _authtype == "authenticationSimplePassword" {
                    authType := "simple password"
                    ospfv2AreaInfo_obj.AuthenticationType = &authType
                }
                if _authtype == "authenticationMessageDigest" {
                    authType := "message digest"
                    ospfv2AreaInfo_obj.AuthenticationType = &authType
                }
            }
            
            if value,ok := area_info["areaIfTotalCounter"] ; ok {
                _areaIfTotalCounter  := uint32(value.(float64))
                ospfv2AreaInfo_obj.InterfaceCount = &_areaIfTotalCounter
            }
            if value,ok := area_info["areaIfActiveCounter"] ; ok {
                _areaIfActiveCounter  := uint32(value.(float64))
                ospfv2AreaInfo_obj.ActiveInterfaceCount = &_areaIfActiveCounter
            }
            if value,ok := area_info["nbrFullAdjacentCounter"] ; ok {
                _nbrFullAdjacentCounter  := uint32(value.(float64))
                ospfv2AreaInfo_obj.AdjacencyCount = &_nbrFullAdjacentCounter
            }
            
            if value,ok := area_info["spfExecutedCounter"] ; ok {
                _spfExecutedCounter := uint32(value.(float64))
                ospfv2AreaInfo_obj.SpfExecutionCount = &_spfExecutedCounter
            }
            
            if value,ok := area_info["lsaNumber"] ; ok {
                _lsaNumber  := uint32(value.(float64))
                ospfv2AreaInfo_obj.LsaCount = &_lsaNumber
            }
            if value,ok := area_info["lsaRouterNumber"]; ok {
                _lsaRouterNumber := uint32(value.(float64))
                ospfv2AreaInfo_obj.RouterLsaCount = &_lsaRouterNumber
            }
            if value,ok := area_info["lsaRouterChecksum"]; ok {
                numint64 = int64(value.(float64))
                numstr := fmt.Sprintf("0x%08x", numint64)
                ospfv2AreaInfo_obj.RouterLsaChecksum = &numstr
            }
            if value,ok := area_info["lsaNetworkNumber"]; ok {
                _lsaNetworkNumber := uint32(value.(float64))
                ospfv2AreaInfo_obj.NetworkLsaCount = &_lsaNetworkNumber
            }
            if value,ok := area_info["lsaNetworkChecksum"]; ok {
                numint64 = int64(value.(float64))
                numstr := fmt.Sprintf("0x%08x", numint64)
                ospfv2AreaInfo_obj.NetworkLsaChecksum = &numstr
            }
            if value,ok := area_info["lsaSummaryNumber"]; ok {
                _lsaSummaryNumber := uint32(value.(float64))
                ospfv2AreaInfo_obj.SummaryLsaCount = &_lsaSummaryNumber
            }
            if value,ok := area_info["lsaSummaryChecksum"]; ok {
                numint64 = int64(value.(float64))
                numstr := fmt.Sprintf("0x%08x", numint64)
                ospfv2AreaInfo_obj.SummaryLsaChecksum = &numstr
            }
            if value,ok := area_info["lsaAsbrNumber"]; ok {
                _lsaAsbrNumber := uint32(value.(float64))
                ospfv2AreaInfo_obj.AsbrSummaryLsaCount = &_lsaAsbrNumber
            }
            if value,ok := area_info["lsaAsbrChecksum"]; ok {
                numint64 = int64(value.(float64))
                numstr := fmt.Sprintf("0x%08x", numint64)
                ospfv2AreaInfo_obj.AsbrSummaryLsaChecksum = &numstr
            }
            if value,ok := area_info["lsaNssaNumber"]; ok {
                _lsaNssaNumber := uint32(value.(float64))
                ospfv2AreaInfo_obj.NssaLsaCount = &_lsaNssaNumber
            }
            if value,ok := area_info["lsaNssaChecksum"]; ok {
                numint64 = int64(value.(float64))
                numstr := fmt.Sprintf("0x%08x", numint64)
                ospfv2AreaInfo_obj.NssaLsaChecksum = &numstr
            }
            if value,ok := area_info["lsaOpaqueLinkNumber"]; ok {
                _lsaOpaqueLinkNumber := uint32(value.(float64))
                ospfv2AreaInfo_obj.OpaqueLinkLsaCount = &_lsaOpaqueLinkNumber
            }
            if value,ok := area_info["lsaOpaqueLinkChecksum"]; ok {
                numint64 = int64(value.(float64))
                numstr := fmt.Sprintf("0x%08x", numint64)
                ospfv2AreaInfo_obj.OpaqueLinkLsaChecksum = &numstr
            }
            if value,ok := area_info["lsaOpaqueAreaNumber"]; ok {
                _lsaOpaqueAreaNumber := uint32(value.(float64))
                ospfv2AreaInfo_obj.OpaqueAreaLsaCount = &_lsaOpaqueAreaNumber
            }
            if value,ok := area_info["lsaOpaqueAreaChecksum"]; ok {
                numint64 = int64(value.(float64))
                numstr := fmt.Sprintf("0x%08x", numint64)
                ospfv2AreaInfo_obj.OpaqueAreaLsaChecksum = &numstr
            }
            ospfv2AreaInfo_obj.Shortcut = ocbinds.OpenconfigOspfv2Ext_OSPF_CONFIG_TYPE_UNSET
            if (areaNameStr != "0.0.0.0") {
                ospfv2AreaInfo_obj.Shortcut = ocbinds.OpenconfigOspfv2Ext_OSPF_CONFIG_TYPE_DEFAULT
            }
            if _shortcut,ok := area_info["shortcuttingMode"].(string); ok {
                if _shortcut == "Enabled" {
                    ospfv2AreaInfo_obj.Shortcut = ocbinds.OpenconfigOspfv2Ext_OSPF_CONFIG_TYPE_ENABLE
                }
                if _shortcut == "Disabled" {
                    ospfv2AreaInfo_obj.Shortcut = ocbinds.OpenconfigOspfv2Ext_OSPF_CONFIG_TYPE_DISABLE
                }
                if _shortcut == "Default" {
                    ospfv2AreaInfo_obj.Shortcut = ocbinds.OpenconfigOspfv2Ext_OSPF_CONFIG_TYPE_DEFAULT
                }
            }
            if value,ok := area_info["virtualAdjacenciesPassingCounter"]; ok {
                _virtualAdjacenciesPassingCounter := uint32(value.(float64))
                ospfv2AreaInfo_obj.VirtualLinkAdjacencyCount = &_virtualAdjacenciesPassingCounter
            }
            /*if _stubEnable, ok := area_info["stubEnable"].(bool); ok {
                if _stubEnable ==  true {
                    ospfv2_fill_area_stub_state(ospfv2Area_obj, area_info)    
                }
            }*/
            if _originStubMaxDistRouterLsa,ok := area_info["originStubMaxDistRouterLsa"].(bool); ok {
                if !_originStubMaxDistRouterLsa {
                    ospfv2AreaInfo_obj.OriginStubMaxDistRouterLsa = &ospfv2Zero
                } else {
                    ospfv2AreaInfo_obj.OriginStubMaxDistRouterLsa = &ospfv2One
                }
            }

            if _indefiniteActiveAdmin,ok := area_info["indefiniteActiveAdmin"].(bool); ok {
                if !_indefiniteActiveAdmin {
                    ospfv2AreaInfo_obj.IndefiniteActiveAdmin = &ospfv2Zero
                } else {
                    ospfv2AreaInfo_obj.IndefiniteActiveAdmin = &ospfv2One
                }
            }
        }
    }    
    return err
}

func ospfv2_fill_area_stub_state(ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area,
        area_info map[string]interface{}) error {
    var err error
    var ospfv2Zero bool = false
    var ospfv2One bool = true
    oper_err := errors.New("Operational error")
    cmn_log := "GET: Stub  State for area "
    var stubState *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Stub_State

    log.Infof("ospfv2_fill_area_stub_state - start")

    if (nil == ospfv2Area_obj.Stub) {
        log.Info("Stub information not present in area")
        ygot.BuildEmptyTree (ospfv2Area_obj.Stub)
    }
    stubState = ospfv2Area_obj.Stub.State
    if nil == stubState {
        log.Infof("state under area stub is  missing, add stub state for area ")
        stubState = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Stub_State)
        if stubState == nil {
            log.Errorf("%s failed !! Error: Failed to create Stub State Tree under area", cmn_log)
            return  oper_err
        }
        ygot.BuildEmptyTree (stubState)
        ospfv2Area_obj.Stub.State = stubState
    }
    if _stubEnable, ok := area_info["stubEnable"].(bool); ok {
        if !_stubEnable {
            stubState.Enable = &ospfv2Zero
        } else {
            stubState.Enable = &ospfv2One
        }
    }
    if _stubNoSummary, ok := area_info["stubNoSummary"].(bool); ok {
        if !_stubNoSummary {
            stubState.NoSummary = &ospfv2Zero
        } else {
            stubState.NoSummary = &ospfv2One
        }
    }
    return err
}

func ospfv2_fill_neighbors_state (output_state map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2, area_id string, intf_name string, vrfName interface{}) error {
    var err error
    var ospfv2Areas_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    var ospfv2Interfaces_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces
    var ospfv2Interface_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface
    var ospfv2Neighbors_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_NeighborsList
    var ospfv2NeighborKey ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_NeighborsList_Neighbor_Key
    var ospfv2Neighbor_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_NeighborsList_Neighbor
    var ospfv2NeighborState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_NeighborsList_Neighbor_State
    var ospfv2NeighborAreaKey ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_NeighborsList_Neighbor_State_AreaId_Union_String
    var ospfv2Zero bool = false
    var ospfv2One bool = true
    var areaNameStr string
    var nbr_area_id string

    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Neighbors State"

    log.Infof("ospfv2_fill_neighbors_state - start area %s intf %s vrf %v", area_id, intf_name, vrfName)

    ospfv2Areas_obj = ospfv2_obj.Areas
    if ospfv2Areas_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2 areas list missing", cmn_log)
        return  oper_err
    }
    areaNameStr = fmt.Sprintf("%v",area_id)
    ospfv2Area_obj, err = ospfv2_find_area_by_key(ospfv2Areas_obj, areaNameStr)
    if nil == ospfv2Area_obj {
        log.Infof("Area object missing, add new area=%s", area_id)
        ospfv2Area_obj, err = ospfv2_create_new_area(ospfv2Areas_obj, areaNameStr)
        if (err != nil) {
            log.Info("Failed to create a new area")
            return  oper_err
        }
    }
    ospfv2Interfaces_obj = ospfv2Area_obj.Interfaces
    if ospfv2Interfaces_obj == nil {
        log.Infof("Interfaces Tree under area is  missing, add new Interfaces tree for area %s", area_id)
        ospfv2Interfaces_obj = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces)
        if ospfv2Interfaces_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Interfaces Tree under area", cmn_log)
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2Interfaces_obj)
        ospfv2Area_obj.Interfaces = ospfv2Interfaces_obj
    }
    ospfv2Interface_obj = ospfv2Interfaces_obj.Interface[intf_name]
    if ospfv2Interface_obj == nil {
        log.Infof("Interface object missing under Interfaces Tree, add new Interface=%s", intf_name)
        ospfv2Interface_obj, err = ospfv2Interfaces_obj.NewInterface(intf_name)
        if (err != nil) {
            log.Info("Failed to create a new interface under Interfaces tree")
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2Interface_obj)
    }
    ospfv2Neighbors_obj = ospfv2Interface_obj.NeighborsList
    if ospfv2Neighbors_obj == nil {
        log.Infof("NeighborList Tree under Interface is  missing, add new NeighborList tree for area %s, interface %s", area_id, intf_name)
        ospfv2Neighbors_obj = 
            new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_NeighborsList)
        if ospfv2Neighbors_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Neighbors Tree under Interface", cmn_log)
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2Neighbors_obj)
        ospfv2Interface_obj.NeighborsList = ospfv2Neighbors_obj
    }

    if value, ok := output_state["neighbors"]; ok {
        neighbors_map := value.(map[string]interface {})
        for nbr_rtr_id, value := range neighbors_map {
            nbr_list := value.([]interface{})
            log.Info(nbr_rtr_id)
            for _, nbr := range nbr_list {
                nbr_info := nbr.(map[string]interface {})
                if _area_id,ok := nbr_info["areaId"].(string); ok {
                    result := strings.Split(_area_id, " ") 
                    nbr_area_id = result[0]
                    if nbr_area_id != area_id {
                        log.Infof("Neighbor area-id %s does not match %s ,skipping this neighbor", _area_id, area_id)
                        continue;
                    }
                }
                if _ntv_intf_name,ok := nbr_info["ifaceName"].(string); ok {
                    _intf_name, _ := ospfGetUIIntfName(_ntv_intf_name)
                    if _intf_name != intf_name {
                        log.Infof("Neighbor interface Name %s does not match %s ,skipping this neighbor", _intf_name, intf_name)
                        continue;
                    }
                }
                if _ifaceAddress,ok := nbr_info["ifaceAddress"].(string); ok {
                    //Prepare a new neighbor node
                    ospfv2NeighborKey.NeighborId = nbr_rtr_id
                    ospfv2NeighborKey.NeighborAddress = _ifaceAddress
                    ospfv2Neighbor_obj = ospfv2Neighbors_obj.Neighbor[ospfv2NeighborKey]
                    if (nil == ospfv2Neighbor_obj) {
                        log.Infof("Neighbor object missing, create a new neighbor under area%s, interface %s", area_id, intf_name)
                        ospfv2Neighbor_obj, err = ospfv2Neighbors_obj.NewNeighbor(nbr_rtr_id, _ifaceAddress)
                        if (err != nil) {
                            log.Info("Failed to create a new neighbor")
                            return  oper_err
                        }
                    } 
                    ygot.BuildEmptyTree(ospfv2Neighbor_obj)
                }
                ospfv2NeighborState_obj = ospfv2Neighbor_obj.State
                if nil == ospfv2NeighborState_obj {
                    log.Info("State information not present for OSPF neighbor")
                    ygot.BuildEmptyTree(ospfv2NeighborState_obj)
                    ospfv2Neighbor_obj.State = ospfv2NeighborState_obj
                }    
                
                ospfv2NeighborState_obj.InterfaceName = &intf_name

                if _area_id,ok := nbr_info["areaId"].(string); ok {
                    ospfv2NeighborAreaKey.String = _area_id
                    ospfv2NeighborState_obj.AreaId = &ospfv2NeighborAreaKey
                }

                if _ipAddress, ok := nbr_info["ifaceLocalAddress"].(string); ok {
                    ospfv2NeighborState_obj.InterfaceAddress = &_ipAddress
                }
                
                if value,ok := nbr_info["nbrPriority"] ; ok {
                    _nbrPriority  := uint8(value.(float64))
                    ospfv2NeighborState_obj.Priority = &_nbrPriority
                }
                
                ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_INIT 
                if _nbr_state,ok := nbr_info["nbrState"].(string); ok {
                    switch (_nbr_state) {
                        case  "Full" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_FULL 
                        case "2-Way" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_TWO_WAY 
                        case "ExStart" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_EXSTART
                        case "Down" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_DOWN
                        case "Attempt" : 
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_ATTEMPT
                        case "Init" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_INIT
                        case "Exchange" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_EXCHANGE
                        case "Loading" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_LOADING
                        default:
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_INIT
                    }
                }
                if value,ok := nbr_info["stateChangeCounter"] ; ok {
                    _stateChangeCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.StateChanges = &_stateChangeCounter
                }

                if value,ok := nbr_info["lastPrgrsvChangeMsec"] ; ok {
                    _lastPrgrsvChangeMsec  := uint64(value.(float64))
                    ospfv2NeighborState_obj.LastEstablishedTime = &_lastPrgrsvChangeMsec
                }
                
                if _routerDesignatedId, ok := nbr_info["routerDesignatedId"].(string); ok {
                    ospfv2NeighborState_obj.DesignatedRouter = &_routerDesignatedId
                }

                if _routerDesignatedBackupId, ok := nbr_info["routerDesignatedBackupId"].(string); ok {
                    ospfv2NeighborState_obj.BackupDesignatedRouter = &_routerDesignatedBackupId
                }

                if value,ok := nbr_info["optionsCounter"] ; ok {
                    _optionsCounter  := uint8(value.(float64))
                    ospfv2NeighborState_obj.OptionValue = &_optionsCounter
                }
                
                if _OptionalCapabilities, ok := nbr_info["optionsList"].(string); ok {
                    ospfv2NeighborState_obj.OptionalCapabilities = &_OptionalCapabilities
                }

                if value,ok := nbr_info["routerDeadIntervalTimerDueMsec"] ; ok {
                    _DeadTime  := uint64(value.(float64))
                    ospfv2NeighborState_obj.DeadTime = &_DeadTime
                }

                if value,ok := nbr_info["databaseSummaryListCounter"] ; ok {
                    _databaseSummaryListCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.DatabaseSummaryQueueLength = &_databaseSummaryListCounter
                }

                if value,ok := nbr_info["linkStateRetransmissionListCounter"] ; ok {
                    _linkStateRetransmissionListCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.RetransmitSummaryQueueLength = &_linkStateRetransmissionListCounter
                }

                if value,ok := nbr_info["linkStateRequestListCounter"] ; ok {
                    _linkStateRequestListCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.LinkStateRequestQueueLength = &_linkStateRequestListCounter
                }

                if _threadInactivityTimer, ok := nbr_info["threadInactivityTimer"].(string); ok {
                    if(_threadInactivityTimer == "on") {
                        ospfv2NeighborState_obj.ThreadInactivityTimer = &ospfv2One;
                    } else {
                        ospfv2NeighborState_obj.ThreadInactivityTimer = &ospfv2Zero;
                    }
                }

                if _threadLinkStateRequestRetransmission, ok := nbr_info["threadLinkStateRequestRetransmission"].(string); ok {
                    if(_threadLinkStateRequestRetransmission == "on") {
                        ospfv2NeighborState_obj.ThreadLsRequestRetransmission = &ospfv2One;
                    } else {
                        ospfv2NeighborState_obj.ThreadLsRequestRetransmission = &ospfv2Zero;
                    }
                }
                
                if _threadLinkStateUpdateRetransmission, ok := nbr_info["threadLinkStateUpdateRetransmission"].(string); ok {
                    if(_threadLinkStateUpdateRetransmission == "on") {
                        ospfv2NeighborState_obj.ThreadLsUpdateRetransmission = &ospfv2One;
                    } else {
                        ospfv2NeighborState_obj.ThreadLsUpdateRetransmission = &ospfv2Zero; 
                    }
                }
                if _bfdmap,ok := nbr_info["peerBfdInfo"] ; ok {
                    ospfv2NeighborState_obj.BfdState = &ospfv2One
                    bfdmap := _bfdmap.(map[string]interface{})
                    if _status, ok := bfdmap["status"].(string); ok {
                          ospfv2NeighborState_obj.BfdStatus = &_status
                    }
                    if _BfdPeerType, ok := bfdmap["type"].(string); ok {
                          ospfv2NeighborState_obj.BfdPeerType = &_BfdPeerType
                    }
                    if _lastUpdate, ok := bfdmap["lastUpdate"].(string); ok {
                          ospfv2NeighborState_obj.BfdPeerLastUpdateTime = &_lastUpdate
                    }
                }    
            }
        }
    }    
    return err
}
func ospfv2_fill_interface_vlink_state(intf_info map[string]interface{}, 
                ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area,
                intf_name string, vrfName interface{}) error {
    var err error
    var ok bool
    var _vlinkPeer string
    var ospfv2Zero bool = false
    var numint64 int64
    var ospfv2One bool = true
    var ospfv2IntfState string = "Down"
    var area_id = ospfv2Area_obj.Identifier
    oper_err := errors.New("Operational error")
    cmn_log := "GET: VLINK State for area "
    var ospfv2Vlinks_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks
    var ospfv2Vlink_obj  *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink
    var ospfv2VlinkState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State
    var ospfv2AreaId ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_AreaId_Union_String
    var interfaceName string = intf_name
    log.Infof("%s", cmn_log)

    log.Infof("ospfv2_fill_interface_vlink_state - start intf_name %s vrfName %v", intf_name, vrfName)

    if (nil == ospfv2Area_obj.VirtualLinks) {
        log.Info("Virtual Links information not present in area")
        return oper_err
    }
    ospfv2Vlinks_obj = ospfv2Area_obj.VirtualLinks

    if _vlinkPeer, ok = intf_info["vlinkRemoteRouterId"].(string); ok {
        ospfv2Vlink_obj = ospfv2Vlinks_obj.VirtualLink[_vlinkPeer]
        if nil == ospfv2Vlink_obj {
            log.Infof("Vlink interface missing for %s, peer %s, add new vlink", *area_id, _vlinkPeer)
            ospfv2Vlink_obj, err = ospfv2Vlinks_obj.NewVirtualLink(_vlinkPeer)
            if (err != nil) {
                log.Info("Failed to create a new vlink under vlink tree")
                return  oper_err
            }
            ygot.BuildEmptyTree (ospfv2Vlink_obj)
        }
    }
    if nil == ospfv2Vlink_obj {
        log.Infof("key parameter remote router id, needed for creating vlink interface not present, returning")
        return oper_err
    }
    ospfv2VlinkState_obj = ospfv2Vlink_obj.State
    if nil == ospfv2VlinkState_obj {
        log.Infof("Vlink interface State missing for %s, peer %s, returning", area_id, _vlinkPeer)
        return  oper_err
    }
    ygot.BuildEmptyTree (ospfv2VlinkState_obj)

    log.Infof("vlink intf_name %s ptr %x", intf_name, &intf_name)
    ospfv2VlinkState_obj.Name = &interfaceName
    if _intf_state,ok := intf_info["ifUp"].(bool); ok {
        if !_intf_state { 
            ospfv2VlinkState_obj.OperationalState = &ospfv2IntfState
        } else {
            ospfv2IntfState = "Up"
            ospfv2VlinkState_obj.OperationalState = &ospfv2IntfState
        }
    }

    if value,ok := intf_info["ifIndex"] ; ok {
        _ifIndex  := uint32(value.(float64))
        ospfv2VlinkState_obj.Index = &_ifIndex
    }

    if value,ok := intf_info["mtuBytes"] ; ok {
        _mtuBytes  := uint32(value.(float64))
        ospfv2VlinkState_obj.Mtu = &_mtuBytes
    }

    if value,ok := intf_info["bandwidthMbit"] ; ok {
        _bandwidthMbit  := uint32(value.(float64))
        ospfv2VlinkState_obj.Bandwidth = &_bandwidthMbit
    }

    if _ifFlags, ok := intf_info["ifFlags"].(string); ok {
        ospfv2VlinkState_obj.IfFlags = &_ifFlags
    }

    if _peerAddr, ok := intf_info["vlinkPeer"].(string); ok {
        ospfv2VlinkState_obj.PeerAddress = &_peerAddr
    }

    if _ospfEnabled,ok := intf_info["ospfEnabled"].(bool); ok {
        if !_ospfEnabled { 
            ospfv2VlinkState_obj.OspfEnable = &ospfv2Zero
        } else {
            ospfv2VlinkState_obj.OspfEnable = &ospfv2One
        }
    }

    if _ipAddress, ok := intf_info["ipAddress"].(string); ok {
        ospfv2VlinkState_obj.Address = &_ipAddress
    }

    if value,ok := intf_info["ipAddressPrefixlen"] ; ok {
        _ipAddressPrefixlen  := uint8(value.(float64))
        ospfv2VlinkState_obj.AddressLen = &_ipAddressPrefixlen
    }

    if _ospfIfType, ok := intf_info["ospfIfType"].(string); ok {
        ospfv2VlinkState_obj.OspfInterfaceType = &_ospfIfType
    }

    if _localIfUsed, ok := intf_info["localIfUsed"].(string); ok {
        ospfv2VlinkState_obj.BroadcastAddress = &_localIfUsed
    }
    
    if _areaStr, ok := intf_info["area"].(string); ok {
        ospfv2AreaId.String = _areaStr
        ospfv2VlinkState_obj.AreaId = &ospfv2AreaId
    }

    if _routerId, ok := intf_info["routerId"].(string); ok {
        ospfv2VlinkState_obj.RouterId = &_routerId
    }

    if _networkType, ok := intf_info["networkType"].(string); ok {
        if _networkType == "BROADCAST" {
            ospfv2VlinkState_obj.NetworkType = ocbinds.OpenconfigOspfTypes_OSPF_NETWORK_TYPE_BROADCAST_NETWORK
        }
        if _networkType == "VIRTUALLINK" {
            ospfv2VlinkState_obj.NetworkType = ocbinds.OpenconfigOspfTypes_OSPF_NETWORK_TYPE_VIRTUALLINK_NETWORK
        }
    }
    if value,ok := intf_info["cost"] ; ok {
        _cost  := uint32(value.(float64))
        ospfv2VlinkState_obj.Cost = &_cost
    }
    if value,ok := intf_info["transmitDelaySecs"] ; ok {
        _transmitDelaySecs  := uint32(value.(float64))
        ospfv2VlinkState_obj.TransmitDelay = &_transmitDelaySecs
    }

    if _AdjacencyState, ok := intf_info["state"].(string); ok {
        ospfv2VlinkState_obj.State = &_AdjacencyState
    }
    if value,ok := intf_info["priority"] ; ok {
        _priority  := uint8(value.(float64))
        ospfv2VlinkState_obj.Priority = &_priority
    }
    if _bdrId, ok := intf_info["bdrId"].(string); ok {
        ospfv2VlinkState_obj.BackupDesignatedRouter = &_bdrId
    }
    if _bdrAddress, ok := intf_info["bdrAddress"].(string); ok {
        ospfv2VlinkState_obj.BackupDesignatedRouterAddress = &_bdrAddress
    }
    if value,ok := intf_info["networkLsaSequence"] ; ok {
        numint64 = int64(value.(float64))
        numstr := fmt.Sprintf("0x%08x", numint64)
        ospfv2VlinkState_obj.NetworkLsaSequenceNumber = &numstr
    }
    if _mcastMemberOspfAllRouters,ok := intf_info["mcastMemberOspfAllRouters"].(bool); ok {
        if !_mcastMemberOspfAllRouters { 
            ospfv2VlinkState_obj.MemberOfOspfAllRouters = &ospfv2Zero
        } else {
            ospfv2VlinkState_obj.MemberOfOspfAllRouters = &ospfv2One
        }
    }
    if _mtuMismatchDetect,ok := intf_info["mtuMismatchDetect"].(bool); ok {
        if !_mtuMismatchDetect { 
            ospfv2VlinkState_obj.MtuIgnore = &ospfv2Zero
        } else {
            ospfv2VlinkState_obj.MtuIgnore = &ospfv2One
        }
    }
    if _mcastMemberOspfDesignatedRouters,ok := intf_info["mcastMemberOspfDesignatedRouters"].(bool); ok {
        if !_mcastMemberOspfDesignatedRouters {
            ospfv2VlinkState_obj.MemberOfOspfDesignatedRouters = &ospfv2Zero
        } else {
            ospfv2VlinkState_obj.MemberOfOspfDesignatedRouters = &ospfv2One
        }
    }
    if value,ok := intf_info["nbrCount"] ; ok {
        _nbrCount  := uint32(value.(float64))
        ospfv2VlinkState_obj.NeighborCount = &_nbrCount
    }
    if value,ok := intf_info["nbrAdjacentCount"] ; ok {
        _nbrAdjacentCount  := uint32(value.(float64))
        ospfv2VlinkState_obj.AdjacencyCount = &_nbrAdjacentCount
    }    
    if value,ok := intf_info["timerMsecs"] ; ok {
        _timerMsecs  := uint32(value.(float64))
        ospfv2VlinkState_obj.HelloInterval = &_timerMsecs
    } 
    if value,ok := intf_info["timerDeadSecs"] ; ok {
        _timerDeadSecs  := uint32(value.(float64))
        ospfv2VlinkState_obj.DeadInterval = &_timerDeadSecs
    } 
    if value,ok := intf_info["timerWaitSecs"] ; ok {
        _timerWaitSecs  := uint32(value.(float64))
        ospfv2VlinkState_obj.WaitTime = &_timerWaitSecs
    } 
    if value,ok := intf_info["timerRetransmitSecs"] ; ok {
        _timerRetransmitSecs  := uint32(value.(float64))
        ospfv2VlinkState_obj.RetransmissionInterval = &_timerRetransmitSecs
    } 
    if value,ok := intf_info["timerHelloInMsecs"] ; ok {
        _timerHelloInMsecs  := uint32(value.(float64))
        ospfv2VlinkState_obj.HelloDue = &_timerHelloInMsecs
    } 
    
    return err
}
func ospfv2_fill_interface_state (intf_info map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2, area_id string, intf_name string, vrfName interface{},
        output_interfaces_traffic map[string]interface{}) error {
    var err error
    var ospfv2Areas_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    var ospfv2Interfaces_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces
    var ospfv2Interface_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface
    var ospfv2InterfaceState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_State
    var ospfv2Zero bool = false
    var ospfv2One bool = true
    var numint64 int64
    var ospfv2IntfState string = "Down"
    var areaNameStr string
    var ospfv2AreaId ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_State_AreaId_Union_String
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Interface State"

    log.Infof("ospfv2_fill_interface_state - start area %s intf_name %s vrfName %v", area_id, intf_name, vrfName)

    ospfv2Areas_obj = ospfv2_obj.Areas
    if ospfv2Areas_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2 areas list missing", cmn_log)
        return  oper_err
    }
    areaNameStr = fmt.Sprintf("%v",area_id)
    ospfv2Area_obj, err = ospfv2_find_area_by_key(ospfv2Areas_obj, areaNameStr)
    if nil == ospfv2Area_obj {
        log.Infof("Area object missing, add new area=%s", area_id)
        ospfv2Area_obj, err = ospfv2_create_new_area(ospfv2Areas_obj, areaNameStr)
        if (err != nil) {
            log.Info("Failed to create a new area")
            return  oper_err
        }
    }
    ospfv2Interfaces_obj = ospfv2Area_obj.Interfaces
    if ospfv2Interfaces_obj == nil {
        log.Infof("Interfaces Tree under area is  missing, add new Interfaces tree for area %s", area_id)
        ospfv2Interfaces_obj = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces)
        if ospfv2Interfaces_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Interfaces Tree under area", cmn_log)
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2Interfaces_obj)
        ospfv2Area_obj.Interfaces = ospfv2Interfaces_obj
    }
    ospfv2Interface_obj = ospfv2Interfaces_obj.Interface[intf_name]
    if ospfv2Interface_obj == nil {
        log.Infof("Interface object missing under Interfaces Tree, add new Interface=%s", intf_name)
        ospfv2Interface_obj, err = ospfv2Interfaces_obj.NewInterface(intf_name)
        if (err != nil) {
            log.Info("Failed to create a new interface under Interfaces tree")
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2Interface_obj)
    }
    ospfv2InterfaceState_obj = ospfv2Interface_obj.State
    if ospfv2InterfaceState_obj == nil {
        log.Infof("State under Interface is  missing, add new Interface State for area %s, interface %s", area_id, intf_name)
        ospfv2InterfaceState_obj = 
            new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_State)
        if ospfv2InterfaceState_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create State under Interface", cmn_log)
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2InterfaceState_obj)
        ospfv2Interface_obj.State = ospfv2InterfaceState_obj
    }
    ospfv2InterfaceState_obj.Id = &intf_name
    
    if _intf_state,ok := intf_info["ifUp"].(bool); ok {
        if !_intf_state { 
            ospfv2InterfaceState_obj.OperationalState = &ospfv2IntfState
        } else {
            ospfv2IntfState = "Up"
            ospfv2InterfaceState_obj.OperationalState = &ospfv2IntfState
        }
    }
    if value,ok := intf_info["ifIndex"] ; ok {
        _ifIndex  := uint32(value.(float64))
        ospfv2InterfaceState_obj.Index = &_ifIndex
    }

    if value,ok := intf_info["mtuBytes"] ; ok {
        _mtuBytes  := uint32(value.(float64))
        ospfv2InterfaceState_obj.Mtu = &_mtuBytes
    }

    if value,ok := intf_info["bandwidthMbit"] ; ok {
        _bandwidthMbit  := uint32(value.(float64))
        ospfv2InterfaceState_obj.Bandwidth = &_bandwidthMbit
    }

    if _ifFlags, ok := intf_info["ifFlags"].(string); ok {
        ospfv2InterfaceState_obj.IfFlags = &_ifFlags
    }

    if _ospfEnabled,ok := intf_info["ospfEnabled"].(bool); ok {
        if !_ospfEnabled { 
            ospfv2InterfaceState_obj.OspfEnable = &ospfv2Zero
        } else {
            ospfv2InterfaceState_obj.OspfEnable = &ospfv2One
        }
    }

    if _passiveEnabled,ok := intf_info["timerPassiveIface"].(bool); ok {
        if !_passiveEnabled { 
            ospfv2InterfaceState_obj.Passive = &ospfv2Zero
        } else {
            ospfv2InterfaceState_obj.Passive = &ospfv2One
        }
    }

    if _ipAddress, ok := intf_info["ipAddress"].(string); ok {
        ospfv2InterfaceState_obj.Address = &_ipAddress
    }

    if value,ok := intf_info["ipAddressPrefixlen"] ; ok {
        _ipAddressPrefixlen  := uint8(value.(float64))
        ospfv2InterfaceState_obj.AddressLen = &_ipAddressPrefixlen
    }

    if _ospfIfType, ok := intf_info["ospfIfType"].(string); ok {
        ospfv2InterfaceState_obj.OspfInterfaceType = &_ospfIfType
    }

    if _localIfUsed, ok := intf_info["localIfUsed"].(string); ok {
        ospfv2InterfaceState_obj.BroadcastAddress = &_localIfUsed
    }
    
    if _areaStr, ok := intf_info["area"].(string); ok {
        ospfv2AreaId.String = _areaStr
        ospfv2InterfaceState_obj.AreaId = &ospfv2AreaId
    }

    if _routerId, ok := intf_info["routerId"].(string); ok {
        ospfv2InterfaceState_obj.RouterId = &_routerId
    }

    ospfv2InterfaceState_obj.NetworkType = ocbinds.OpenconfigOspfTypes_OSPF_NETWORK_TYPE_BROADCAST_NETWORK
    if _networkType, ok := intf_info["networkType"].(string); ok {
        if _networkType == "BROADCAST" {
            ospfv2InterfaceState_obj.NetworkType = ocbinds.OpenconfigOspfTypes_OSPF_NETWORK_TYPE_BROADCAST_NETWORK
        }
        if _networkType == "POINTOPOINT" {
            ospfv2InterfaceState_obj.NetworkType = ocbinds.OpenconfigOspfTypes_OSPF_NETWORK_TYPE_POINT_TO_POINT_NETWORK
        }
        if _networkType == "NBMA" {
            ospfv2InterfaceState_obj.NetworkType = ocbinds.OpenconfigOspfTypes_OSPF_NETWORK_TYPE_NON_BROADCAST_NETWORK
        }
        if _networkType == "VIRTUALLINK" {
            ospfv2InterfaceState_obj.NetworkType = ocbinds.OpenconfigOspfTypes_OSPF_NETWORK_TYPE_VIRTUALLINK_NETWORK
        }
    }
    if value,ok := intf_info["cost"] ; ok {
        _cost  := uint32(value.(float64))
        ospfv2InterfaceState_obj.Cost = &_cost
    }
    if value,ok := intf_info["transmitDelaySecs"] ; ok {
        _transmitDelaySecs  := uint32(value.(float64))
        ospfv2InterfaceState_obj.TransmitDelay = &_transmitDelaySecs
    }

    if _AdjacencyState, ok := intf_info["state"].(string); ok {
        ospfv2InterfaceState_obj.AdjacencyState = &_AdjacencyState
    }
    if value,ok := intf_info["priority"] ; ok {
        _priority  := uint8(value.(float64))
        ospfv2InterfaceState_obj.Priority = &_priority
    }
    if _bdrId, ok := intf_info["bdrId"].(string); ok {
        ospfv2InterfaceState_obj.BackupDesignatedRouterId = &_bdrId
    }
    if _bdrAddress, ok := intf_info["bdrAddress"].(string); ok {
        ospfv2InterfaceState_obj.BackupDesignatedRouterAddress = &_bdrAddress
    }
    if value,ok := intf_info["networkLsaSequence"] ; ok {
        numint64 = int64(value.(float64))
        numstr := fmt.Sprintf("0x%08x", numint64)
        ospfv2InterfaceState_obj.NetworkLsaSequenceNumber = &numstr
    }
    if _mcastMemberOspfAllRouters,ok := intf_info["mcastMemberOspfAllRouters"].(bool); ok {
        if !_mcastMemberOspfAllRouters { 
            ospfv2InterfaceState_obj.MemberOfOspfAllRouters = &ospfv2Zero
        } else {
            ospfv2InterfaceState_obj.MemberOfOspfAllRouters = &ospfv2One
        }
    }
    if _mtuMismatchDetect,ok := intf_info["mtuMismatchDetect"].(bool); ok {
        if !_mtuMismatchDetect { 
            ospfv2InterfaceState_obj.MtuIgnore = &ospfv2Zero
        } else {
            ospfv2InterfaceState_obj.MtuIgnore = &ospfv2One
        }
    }
    if _mcastMemberOspfDesignatedRouters,ok := intf_info["mcastMemberOspfDesignatedRouters"].(bool); ok {
        if !_mcastMemberOspfDesignatedRouters {
            ospfv2InterfaceState_obj.MemberOfOspfDesignatedRouters = &ospfv2Zero
        } else {
            ospfv2InterfaceState_obj.MemberOfOspfDesignatedRouters = &ospfv2One
        }
    }
    if _ifUnnumbered,ok := intf_info["ifUnnumbered"].(bool); ok {
        if !_ifUnnumbered {
            ospfv2InterfaceState_obj.Unnumbered = &ospfv2Zero
        } else {
            ospfv2InterfaceState_obj.Unnumbered = &ospfv2One
        }
    }
    if value,ok := intf_info["nbrCount"] ; ok {
        _nbrCount  := uint32(value.(float64))
        ospfv2InterfaceState_obj.NeighborCount = &_nbrCount
    }
    if value,ok := intf_info["nbrAdjacentCount"] ; ok {
        _nbrAdjacentCount  := uint32(value.(float64))
        ospfv2InterfaceState_obj.AdjacencyCount = &_nbrAdjacentCount
    }    
    if _,ok := intf_info["peerBfdInfo"] ; ok {
        ospfv2InterfaceState_obj.BfdState = &ospfv2One
    }    
    return err
}


func ospfv2_fill_interface_message_stats (output_state map[string]interface{}, 
        ospfv2Interface_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface, intf_name string) error {
    var err error
    var ospfv2IntfStats_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_MessageStatistics
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Interface Message Statistics"

    log.Infof("ospfv2_fill_interface_message_stats - start ")

    ospfv2IntfStats_obj = ospfv2Interface_obj.MessageStatistics
    if ospfv2IntfStats_obj == nil {
        log.Infof("message statistics under Interface is  missing, add new Interface msg statistics for interface %s", intf_name)
        ospfv2IntfStats_obj = 
            new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_MessageStatistics)
        if ospfv2IntfStats_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Message Statistics under Interface", cmn_log)
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2IntfStats_obj)
        ospfv2Interface_obj.MessageStatistics = ospfv2IntfStats_obj
    }
    
    for _,value := range output_state {
        interfaces_info := value.(map[string]interface{})
        for ntv_key, value := range interfaces_info {
            key, _ := ospfGetUIIntfName(ntv_key)
            if(key != intf_name) {
                log.Infof("skipping interface %s as stats needed for interface %s ", key, intf_name)
                continue
            }
            intf_info := value.(map[string]interface{})
            if value,ok := intf_info["helloIn"] ; ok {
                _helloIn  := uint32(value.(float64))
                ospfv2IntfStats_obj.HelloReceive = &_helloIn
            }
            if value,ok := intf_info["helloOut"] ; ok {
                _helloOut  := uint32(value.(float64))
                ospfv2IntfStats_obj.HelloTransmit = &_helloOut
            }
            if value,ok := intf_info["dbDescIn"] ; ok {
                _dbDescIn  := uint32(value.(float64))
                ospfv2IntfStats_obj.DbDescriptionReceive = &_dbDescIn
            }
            if value,ok := intf_info["dbDescOut"] ; ok {
                _dbDescOut  := uint32(value.(float64))
                ospfv2IntfStats_obj.DbDescriptionTransmit = &_dbDescOut
            }
            if value,ok := intf_info["lsReqIn"] ; ok {
                _lsReqIn  := uint32(value.(float64))
                ospfv2IntfStats_obj.LsRequestReceive = &_lsReqIn
            }
            if value,ok := intf_info["lsReqOut"] ; ok {
                _lsReqOut  := uint32(value.(float64))
                ospfv2IntfStats_obj.LsRequestTransmit = &_lsReqOut
            }
            if value,ok := intf_info["lsUpdIn"] ; ok {
                _lsUpdIn  := uint32(value.(float64))
                ospfv2IntfStats_obj.LsUpdateReceive = &_lsUpdIn
            }
            if value,ok := intf_info["lsUpdOut"] ; ok {
                _lsUpdOut  := uint32(value.(float64))
                ospfv2IntfStats_obj.LsUpdateTransmit = &_lsUpdOut
            }
            if value,ok := intf_info["lsAckIn"] ; ok {
                _lsAckIn  := uint32(value.(float64))
                ospfv2IntfStats_obj.LsAcknowledgeReceive = &_lsAckIn
            }
            if value,ok := intf_info["lsAckOut"] ; ok {
                _lsAckOut  := uint32(value.(float64))
                ospfv2IntfStats_obj.LsAcknowledgeTransmit = &_lsAckOut
            }
        } 
    }
    
    return err
}

func ospfv2_fill_interface_vlink_state_traffic (intf_info map[string]interface{}, output_state map[string]interface{}, 
        ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area, intf_name string) (*ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink, error) {
    var err error
    var ok bool
    var _vlinkPeer string
    var area_id = ospfv2Area_obj.Identifier
    oper_err := errors.New("Operational error")
    cmn_log := "GET: VLINK State Traffic for area "
    var ospfv2Vlinks_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks
    var ospfv2Vlink_obj  *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink
    var ospfv2VlinkState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State
    var ospfv2VlinkStats_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_MessageStatistics
    log.Infof("%s", cmn_log)

    log.Infof("ospfv2_fill_interface_vlink_state_traffic - start")

    if (nil == ospfv2Area_obj.VirtualLinks) {
        log.Info("Virtual Links information not present in area")
        return nil, oper_err
    }
    ospfv2Vlinks_obj = ospfv2Area_obj.VirtualLinks

    if _vlinkPeer, ok = intf_info["vlinkRemoteRouterId"].(string); ok {
        ospfv2Vlink_obj = ospfv2Vlinks_obj.VirtualLink[_vlinkPeer]
        if nil == ospfv2Vlink_obj {
            log.Infof("Vlink interface missing for %s, peer %s, add new vlink", area_id, _vlinkPeer)
            ospfv2Vlink_obj, err = ospfv2Vlinks_obj.NewVirtualLink(_vlinkPeer)
            if (err != nil) {
                log.Info("Failed to create a new vlink under vlink tree")
                return  nil, oper_err
            }
            ygot.BuildEmptyTree (ospfv2Vlink_obj)
            ospfv2Vlink_obj.RemoteRouterId = &_vlinkPeer
        }
    }
    if nil == ospfv2Vlink_obj {
        log.Infof("key parameter remote router id, needed for creating vlink interface not present, returning")
        return nil, oper_err
    }
    ospfv2VlinkState_obj = ospfv2Vlink_obj.State
    if nil == ospfv2VlinkState_obj {
        log.Infof("Vlink interface State missing for %s, peer %s, returning", area_id, _vlinkPeer)
        return  nil, oper_err
    }
    ospfv2VlinkState_obj.Name = &intf_name
    
    ospfv2VlinkStats_obj = ospfv2VlinkState_obj.MessageStatistics
    if nil == ospfv2VlinkStats_obj {
        log.Infof("Statistics under Vlink State missing for %s, peer %s, returning", area_id, _vlinkPeer)
        return  nil, oper_err
    }
    for _,value := range output_state {
        interfaces_info := value.(map[string]interface{})
        for ntv_key, value := range interfaces_info {
            key, _ := ospfGetUIIntfName(ntv_key)
            if(key != intf_name) {
                log.Infof("skipping interface %s as stats needed for interface %s ", key, intf_name)
                continue
            }
            intf_info := value.(map[string]interface{})
            if value,ok := intf_info["helloIn"] ; ok {
                _helloIn  := uint32(value.(float64))
                ospfv2VlinkStats_obj.HelloReceive = &_helloIn
            }
            if value,ok := intf_info["helloOut"] ; ok {
                _helloOut  := uint32(value.(float64))
                ospfv2VlinkStats_obj.HelloTransmit = &_helloOut
            }
            if value,ok := intf_info["dbDescIn"] ; ok {
                _dbDescIn  := uint32(value.(float64))
                ospfv2VlinkStats_obj.DbDescriptionReceive = &_dbDescIn
            }
            if value,ok := intf_info["dbDescOut"] ; ok {
                _dbDescOut  := uint32(value.(float64))
                ospfv2VlinkStats_obj.DbDescriptionTransmit = &_dbDescOut
            }
            if value,ok := intf_info["lsReqIn"] ; ok {
                _lsReqIn  := uint32(value.(float64))
                ospfv2VlinkStats_obj.LsRequestReceive = &_lsReqIn
            }
            if value,ok := intf_info["lsReqOut"] ; ok {
                _lsReqOut  := uint32(value.(float64))
                ospfv2VlinkStats_obj.LsRequestTransmit = &_lsReqOut
            }
            if value,ok := intf_info["lsUpdIn"] ; ok {
                _lsUpdIn  := uint32(value.(float64))
                ospfv2VlinkStats_obj.LsUpdateReceive = &_lsUpdIn
            }
            if value,ok := intf_info["lsUpdOut"] ; ok {
                _lsUpdOut  := uint32(value.(float64))
                ospfv2VlinkStats_obj.LsUpdateTransmit = &_lsUpdOut
            }
            if value,ok := intf_info["lsAckIn"] ; ok {
                _lsAckIn  := uint32(value.(float64))
                ospfv2VlinkStats_obj.LsAcknowledgeReceive = &_lsAckIn
            }
            if value,ok := intf_info["lsAckOut"] ; ok {
                _lsAckOut  := uint32(value.(float64))
                ospfv2VlinkStats_obj.LsAcknowledgeTransmit = &_lsAckOut
            }
        } 
    }
    
    return ospfv2Vlink_obj, err
}
func ospfv2_fill_interface_timers_state (intf_info map[string]interface{}, 
        ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2, area_id string, intf_name string, vrfName interface{}) (*ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface, error) {
    var err error
    var ospfv2Areas_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    var ospfv2Interfaces_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces
    var ospfv2Interface_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface
    var ospfv2InterfaceTimers_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_Timers
    var ospfv2InterfaceTimersState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_Timers_State
    var areaNameStr string

    log.Infof("ospfv2_fill_interface_timers_state - start area %s intf %s vrf %v", area_id, intf_name, vrfName)

    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Interface Timers State"

    ospfv2Areas_obj = ospfv2_obj.Areas
    if ospfv2Areas_obj == nil {
        log.Errorf("%s failed !! Error: Ospfv2 areas list missing", cmn_log)
        return  nil, oper_err
    }
    areaNameStr = fmt.Sprintf("%v",area_id)
    ospfv2Area_obj, err = ospfv2_find_area_by_key(ospfv2Areas_obj, areaNameStr)
    if nil == ospfv2Area_obj {
        log.Infof("Area object missing, add new area=%s", area_id)
        ospfv2Area_obj, err = ospfv2_create_new_area(ospfv2Areas_obj, areaNameStr)
        if (err != nil) {
            log.Info("Failed to create a new area")
            return  nil, oper_err
        }
    }
    ospfv2Interfaces_obj = ospfv2Area_obj.Interfaces
    if ospfv2Interfaces_obj == nil {
        log.Infof("Interfaces Tree under area is  missing, add new Interfaces tree for area %s", area_id)
        ospfv2Interfaces_obj = new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces)
        if ospfv2Interfaces_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Interfaces Tree under area", cmn_log)
            return  nil, oper_err
        }
        ygot.BuildEmptyTree (ospfv2Interfaces_obj)
        ospfv2Area_obj.Interfaces = ospfv2Interfaces_obj
    }
    ospfv2Interface_obj = ospfv2Interfaces_obj.Interface[intf_name]
    if ospfv2Interface_obj == nil {
        log.Infof("Interface object missing under Interfaces Tree, add new Interface=%s", intf_name)
        ospfv2Interface_obj, err = ospfv2Interfaces_obj.NewInterface(intf_name)
        if (err != nil) {
            log.Info("Failed to create a new interface under Interfaces tree")
            return  nil, oper_err
        }
        ygot.BuildEmptyTree (ospfv2Interface_obj)
    }
    ospfv2InterfaceTimers_obj = ospfv2Interface_obj.Timers
    if ospfv2InterfaceTimers_obj == nil {
        log.Infof("Timers under Interface is  missing, add new Interface Timers for area %s, interface %s", area_id, intf_name)
        ospfv2InterfaceTimers_obj = 
            new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_Timers)
        if ospfv2InterfaceTimers_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Timers under Interface", cmn_log)
            return  ospfv2Interface_obj, oper_err
        }
        ygot.BuildEmptyTree (ospfv2InterfaceTimers_obj)
        ospfv2Interface_obj.Timers = ospfv2InterfaceTimers_obj
    }
    ospfv2InterfaceTimersState_obj = ospfv2InterfaceTimers_obj.State
    if ospfv2InterfaceTimersState_obj == nil {
        log.Infof("Timers State under Interface is  missing, add new Interface Timers State for area %s, interface %s", area_id, intf_name)
        ospfv2InterfaceTimersState_obj = 
            new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface_Timers_State)
        if ospfv2InterfaceTimersState_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Timers State under Interface", cmn_log)
            return  ospfv2Interface_obj, oper_err
        }
        ygot.BuildEmptyTree (ospfv2InterfaceTimersState_obj)
        ospfv2InterfaceTimers_obj.State = ospfv2InterfaceTimersState_obj
    }
    if value,ok := intf_info["timerMsecs"] ; ok {
        _timerMsecs  := uint32(value.(float64))
        ospfv2InterfaceTimersState_obj.HelloInterval = &_timerMsecs
    } 
    if value,ok := intf_info["timerDeadSecs"] ; ok {
        _timerDeadSecs  := uint32(value.(float64))
        ospfv2InterfaceTimersState_obj.DeadInterval = &_timerDeadSecs
    } 
    if value,ok := intf_info["timerWaitSecs"] ; ok {
        _timerWaitSecs  := uint32(value.(float64))
        ospfv2InterfaceTimersState_obj.WaitTime = &_timerWaitSecs
    } 
    if value,ok := intf_info["timerRetransmitSecs"] ; ok {
        _timerRetransmitSecs  := uint32(value.(float64))
        ospfv2InterfaceTimersState_obj.RetransmissionInterval = &_timerRetransmitSecs
    } 
    if value,ok := intf_info["timerHelloInMsecs"] ; ok {
        _timerHelloInMsecs  := uint32(value.(float64))
        ospfv2InterfaceTimersState_obj.HelloDue = &_timerHelloInMsecs
    } 
    return ospfv2Interface_obj, err
}
var DbToYang_ospfv2_global_timers_spf_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Global Timer Spf  State"
    var vtysh_cmd string

    log.Info("DbToYang_ospfv2_global_timers_spf_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }
    log.Info(vrfName)

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)

    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_state)
    log.Info(vrfName)
    
    for key,value := range output_state {
        ospf_info := value.(map[string]interface{})
        log.Info(key)
        log.Info(ospf_info)
        err = ospfv2_fill_global_timers_spf_state (ospf_info, ospfv2_obj)
    }
    
    return  err;
}
var DbToYang_ospfv2_global_timers_lsa_generation_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Global Timer LSA generation  State"
    var vtysh_cmd string

    log.Info("DbToYang_ospfv2_global_timers_lsa_generation_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)

    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    for key,value := range output_state {
        ospf_info := value.(map[string]interface{})
        log.Info(key)
        log.Info(ospf_info)
        err = ospfv2_fill_global_timers_lsa_generation_state (ospf_info, ospfv2_obj)
    }
    
    return  err;
}


var DbToYang_ospfv2_route_table_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Route Table"
    var vtysh_cmd string

    log.Info("DbToYang_ospfv2_route_table_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)

    targetUriPath, _ := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " route json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_state)
    log.Info(vrfName)
    ospf_info := output_state[vrfName].(map[string]interface{})
    err = ospfv2_fill_route_table (ospf_info, ospfv2_obj)
    return  err;
}

var DbToYang_ospfv2_global_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Global State"
    var vtysh_cmd string

    log.Info("DbToYang_ospfv2_global_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }
    log.Info(vrfName)

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)

    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    for key,value := range output_state {
        ospf_info := value.(map[string]interface{})
        log.Info(key)
        log.Info(ospf_info)
        err = ospfv2_fill_only_global_state(ospf_info, ospfv2_obj)
    }
    
    return  err;
}
var DbToYang_ospfv2_areas_area_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF- Areas Area State"
    var vtysh_cmd string

    log.Info("DbToYang_ospfv2_areas_area_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }
    log.Info("vrfName=", vrfName)

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)
    area_id :=pathInfo.Var("identifier#2")
    if(len(area_id) == 0) {
        log.Info("Area Id is not specified, key is missing")
        log.Errorf ("%s failed !! Error", cmn_log);
        return  oper_err
    } else {
        area_id = getAreaDotted(area_id)
        log.Infof("Area Id %s", area_id)
    }
    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_state)
    
    for key,value := range output_state {
        ospf_info := value.(map[string]interface{})
        log.Info(key)
        log.Info(ospf_info)
        err = ospfv2_fill_area_state (ospf_info, ospfv2_obj, area_id, vrfName)
    }
    
    return  err;
}

var DbToYang_ospfv2_vlink_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error in  DbToYang_ospfv2_vlink_state_xfmr")
    cmn_log := "GET: xfmr for OSPF- Areas Area Virtual Link State"
    var vtysh_cmd string
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    var intf_name string
    var ntv_intf_name string
    var temp interface{}
    var intf_area_id string
    var ospfv2Vlink_obj  *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink

    log.Info("DbToYang_ospfv2_vlink_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }
    log.Info("vrfName=", vrfName)

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)
    area_id :=pathInfo.Var("identifier#2")
    if(len(area_id) == 0) {
        log.Info("Area Id is not specified, key is missing")
        log.Errorf ("%s failed !! Error", cmn_log);
        return  oper_err
    } else {
        area_id = getAreaDotted(area_id)
        log.Infof("Area Id %s", area_id)
    }
    remote_rtr_id :=pathInfo.Var("remote-router-id")
    if(len(remote_rtr_id) == 0) {
        log.Info("Remote Rtr Id is not specified, key is missing")
        log.Errorf ("%s failed !! Error", cmn_log);
        return  oper_err
    } else {
        log.Infof("remote rtr Id %s", remote_rtr_id)
    }

    targetUriPath, _ := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_state)
    
    ospf_info := output_state[vrfName].(map[string]interface{})
    ospfv2Area_obj, _, err = ospfv2_get_or_create_area (ospf_info, ospfv2_obj, area_id, vrfName)

    if nil == ospfv2Area_obj {
        log.Errorf("Failed to create a new area:%s, err=%s", area_id, err)
        return oper_err
    }

    vtysh_cmd = "show ip ospf vrf " + vrfName + " interface json"
    output_interfaces, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf interfaces:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_interfaces)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " interface traffic json"
    output_interfaces_traffic, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf interfaces traffic:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_interfaces_traffic)

    vtysh_cmd = "show ip ospf vrf " + vrfName + " neighbor detail json"
    output_nbrs_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf neighbor detail:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_nbrs_state)
    for _,value := range output_interfaces { 
        interfaces_info := value.(map[string]interface{})
        interface_map := interfaces_info["interfaces"].(map[string]interface{})
        for ntv_intf_name, temp = range interface_map {
            intf_name, _ = ospfGetUIIntfName(ntv_intf_name)
            if !strings.Contains(intf_name, "VLINK") {
                log.Info("Skip non vlink interface ", intf_name)
                continue
            }
            intf_info := temp.(map[string]interface{})
            intf_area_id = ""
            if intf_area_str,ok := intf_info["vlinkTransitArea"].(string); ok {
                result := strings.Split(intf_area_str, " ") 
                intf_area_id = result[0]
                if (intf_area_id != area_id) {
                    log.Infof("Skipping Interface %s belonging to area %s, as given area %s", intf_name, intf_area_id, area_id)
                    continue
                }
            }
            if (intf_area_id == "") {
                log.Infof("vlinkTransitArea attribute not present in vlink state, skip interface %s", intf_name)
                continue
            }
            if _vlinkPeer, ok := intf_info["vlinkRemoteRouterId"].(string); ok {
                if _vlinkPeer != remote_rtr_id {
                    log.Infof("Skipping Vlink interface (%s) as we need to fill vlink (%s)", _vlinkPeer, remote_rtr_id)
                    continue
                }
            }
            
            ospfv2_fill_interface_vlink_state(intf_info, ospfv2Area_obj, intf_name, vrfName)
            ospfv2Vlink_obj, err = ospfv2_fill_interface_vlink_state_traffic(intf_info, output_interfaces_traffic, ospfv2Area_obj, intf_name)
            if nil != ospfv2Vlink_obj {
                neighbors_info := output_nbrs_state[vrfName].(map[string]interface{})
                ospfv2_fill_vlink_neighbors_state(neighbors_info, ospfv2Vlink_obj, area_id, remote_rtr_id, intf_name)
            }
        }
    }
    
    return  err;
}

func  ospfv2_fill_vlink_neighbors_state (output_state map[string]interface{},  ospfv2Vlink_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink, area_id string, remote_rtr_id string, intf_name string) error {
    var err error
    var ospfv2Neighbors_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_NeighborsList
    var ospfv2NeighborKey ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_NeighborsList_Neighbor_Key
    var ospfv2Neighbor_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_NeighborsList_Neighbor
    var ospfv2NeighborState_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_NeighborsList_Neighbor_State
    var ospfv2NeighborAreaKey ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_NeighborsList_Neighbor_State_AreaId_Union_String
    var ospfv2Zero bool = false
    var ospfv2One bool = true
    var nbr_area_id string

    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Vlink Neighbors State"

    log.Infof("ospfv2_fill_vlink_neighbors_state - area %s intf %s rtrid %s", area_id, intf_name, remote_rtr_id)

    ospfv2Neighbors_obj = ospfv2Vlink_obj.State.NeighborsList
    if ospfv2Neighbors_obj == nil {
        log.Infof("NeighborList Tree under Vlink Interface is  missing, add new NeighborList tree for area %s, interface %s", area_id, intf_name)
        ospfv2Neighbors_obj = 
            new(ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_VirtualLinks_VirtualLink_State_NeighborsList)
        if ospfv2Neighbors_obj == nil {
            log.Errorf("%s failed !! Error: Failed to create Neighbors Tree under Vlink Interface", cmn_log)
            return  oper_err
        }
        ygot.BuildEmptyTree (ospfv2Neighbors_obj)
        ospfv2Vlink_obj.State.NeighborsList = ospfv2Neighbors_obj
    }

    if value, ok := output_state["neighbors"]; ok {
        neighbors_map := value.(map[string]interface {})
        for nbr_rtr_id, value := range neighbors_map {
            nbr_list := value.([]interface{})
            log.Info(nbr_rtr_id)
            for _, nbr := range nbr_list {
                nbr_info := nbr.(map[string]interface {})
                nbr_area_id = ""
                if _area_id,ok := nbr_info["vlinkTransitArea"].(string); ok {
                    result := strings.Split(_area_id, " ") 
                    nbr_area_id = result[0]
                    if nbr_area_id != area_id {
                        log.Infof("Neighbor area-id %s does not match %s ,skipping this neighbor", _area_id, area_id)
                        continue;
                    }
                }
                if nbr_area_id == "" {
                    log.Infof("Neighbor area-id not found in vlink ,skipping this neighbor")
                    continue
                }
                if _ntv_intf_name,ok := nbr_info["ifaceName"].(string); ok {
                    _intf_name, _ := ospfGetUIIntfName(_ntv_intf_name) 
                    if _intf_name != intf_name {
                        log.Infof("Neighbor interface Name %s does not match %s ,skipping this neighbor", _intf_name, intf_name)
                        continue;
                    }
                }
                if _ifaceAddress,ok := nbr_info["ifaceAddress"].(string); ok {
                    //Prepare a new neighbor node
                    ospfv2NeighborKey.NeighborId = nbr_rtr_id
                    ospfv2NeighborKey.NeighborAddress = _ifaceAddress
                    ospfv2Neighbor_obj = ospfv2Neighbors_obj.Neighbor[ospfv2NeighborKey]
                    if (nil == ospfv2Neighbor_obj) {
                        log.Infof("Neighbor object missing, create a new neighbor under area%s, vlink interface %s", area_id, intf_name)
                        ospfv2Neighbor_obj, err = ospfv2Neighbors_obj.NewNeighbor(nbr_rtr_id, _ifaceAddress)
                        if (err != nil) {
                            log.Info("Failed to create a new neighbor")
                            return  oper_err
                        }
                    } 
                    ygot.BuildEmptyTree(ospfv2Neighbor_obj)
                }
                ospfv2NeighborState_obj = ospfv2Neighbor_obj.State
                if nil == ospfv2NeighborState_obj {
                    log.Info("State information not present for Vlink OSPF neighbor")
                    ygot.BuildEmptyTree(ospfv2NeighborState_obj)
                    ospfv2Neighbor_obj.State = ospfv2NeighborState_obj
                }    
                
                ospfv2NeighborState_obj.InterfaceName = &intf_name
                if _area_id,ok := nbr_info["areaId"].(string); ok {
                    ospfv2NeighborAreaKey.String = _area_id
                    ospfv2NeighborState_obj.AreaId = &ospfv2NeighborAreaKey
                }

                if _ipAddress, ok := nbr_info["ifaceLocalAddress"].(string); ok {
                    ospfv2NeighborState_obj.InterfaceAddress = &_ipAddress
                }
                
                if value,ok := nbr_info["nbrPriority"] ; ok {
                    _nbrPriority  := uint8(value.(float64))
                    ospfv2NeighborState_obj.Priority = &_nbrPriority
                }
                
                ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_INIT 
                if _nbr_state,ok := nbr_info["nbrState"].(string); ok {
                    switch (_nbr_state) {
                        case  "Full" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_FULL 
                        case "2-Way" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_TWO_WAY 
                        case "ExStart" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_EXSTART
                        case "Down" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_DOWN
                        case "Attempt" : 
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_ATTEMPT
                        case "Init" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_INIT
                        case "Exchange" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_EXCHANGE
                        case "Loading" :
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_LOADING
                        default:
                            ospfv2NeighborState_obj.AdjacencyState = ocbinds.OpenconfigOspfTypes_OSPF_NEIGHBOR_STATE_INIT
                    }
                }
                if value,ok := nbr_info["stateChangeCounter"] ; ok {
                    _stateChangeCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.StateChanges = &_stateChangeCounter
                }

                if value,ok := nbr_info["lastPrgrsvChangeMsec"] ; ok {
                    _lastPrgrsvChangeMsec  := uint64(value.(float64))
                    ospfv2NeighborState_obj.LastEstablishedTime = &_lastPrgrsvChangeMsec
                }
                
                if _routerDesignatedId, ok := nbr_info["routerDesignatedId"].(string); ok {
                    ospfv2NeighborState_obj.DesignatedRouter = &_routerDesignatedId
                }

                if _routerDesignatedBackupId, ok := nbr_info["routerDesignatedBackupId"].(string); ok {
                    ospfv2NeighborState_obj.BackupDesignatedRouter = &_routerDesignatedBackupId
                }

                if value,ok := nbr_info["optionsCounter"] ; ok {
                    _optionsCounter  := uint8(value.(float64))
                    ospfv2NeighborState_obj.OptionValue = &_optionsCounter
                }
                
                if _OptionalCapabilities, ok := nbr_info["optionsList"].(string); ok {
                    ospfv2NeighborState_obj.OptionalCapabilities = &_OptionalCapabilities
                }

                if value,ok := nbr_info["routerDeadIntervalTimerDueMsec"] ; ok {
                    _DeadTime  := uint64(value.(float64))
                    ospfv2NeighborState_obj.DeadTime = &_DeadTime
                }

                if value,ok := nbr_info["databaseSummaryListCounter"] ; ok {
                    _databaseSummaryListCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.DatabaseSummaryQueueLength = &_databaseSummaryListCounter
                }

                if value,ok := nbr_info["linkStateRetransmissionListCounter"] ; ok {
                    _linkStateRetransmissionListCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.RetransmitSummaryQueueLength = &_linkStateRetransmissionListCounter
                }

                if value,ok := nbr_info["linkStateRequestListCounter"] ; ok {
                    _linkStateRequestListCounter  := uint32(value.(float64))
                    ospfv2NeighborState_obj.LinkStateRequestQueueLength = &_linkStateRequestListCounter
                }

                if _threadInactivityTimer, ok := nbr_info["threadInactivityTimer"].(string); ok {
                    if(_threadInactivityTimer == "on") {
                        ospfv2NeighborState_obj.ThreadInactivityTimer = &ospfv2One;
                    } else {
                        ospfv2NeighborState_obj.ThreadInactivityTimer = &ospfv2Zero;
                    }
                }

                if _threadLinkStateRequestRetransmission, ok := nbr_info["threadLinkStateRequestRetransmission"].(string); ok {
                    if(_threadLinkStateRequestRetransmission == "on") {
                        ospfv2NeighborState_obj.ThreadLsRequestRetransmission = &ospfv2One;
                    } else {
                        ospfv2NeighborState_obj.ThreadLsRequestRetransmission = &ospfv2Zero;
                    }
                }
                
                if _threadLinkStateUpdateRetransmission, ok := nbr_info["threadLinkStateUpdateRetransmission"].(string); ok {
                    if(_threadLinkStateUpdateRetransmission == "on") {
                        ospfv2NeighborState_obj.ThreadLsUpdateRetransmission = &ospfv2One;
                    } else {
                        ospfv2NeighborState_obj.ThreadLsUpdateRetransmission = &ospfv2Zero; 
                    }
                }
                if _bfdmap,ok := nbr_info["peerBfdInfo"] ; ok {
                    ospfv2NeighborState_obj.BfdState = &ospfv2One
                    bfdmap := _bfdmap.(map[string]interface{})
                    if _status, ok := bfdmap["status"].(string); ok {
                          ospfv2NeighborState_obj.BfdStatus = &_status
                    }
                    if _BfdPeerType, ok := bfdmap["type"].(string); ok {
                          ospfv2NeighborState_obj.BfdPeerType = &_BfdPeerType
                    }
                    if _lastUpdate, ok := bfdmap["lastUpdate"].(string); ok {
                          ospfv2NeighborState_obj.BfdPeerLastUpdateTime = &_lastUpdate
                    }
                }    
            }
        }
    }    
    return err
}
var DbToYang_ospfv2_stub_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error in  DbToYang_ospfv2_stub_state_xfmr")
    cmn_log := "GET: xfmr for OSPF- Areas Area Stub State"
    var vtysh_cmd string
    var ospfv2Area_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area
    var area_info map[string]interface{}

    log.Info("DbToYang_ospfv2_stub_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }
    log.Info("vrfName=", vrfName)

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)
    area_id :=pathInfo.Var("identifier#2")
    if(len(area_id) == 0) {
        log.Info("Area Id is not specified, key is missing")
        log.Errorf ("%s failed !! Error", cmn_log);
        return  oper_err
    } else {
        area_id = getAreaDotted(area_id)
        log.Infof("Area Id %s", area_id)
    }
    targetUriPath, _ := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_state)
    
    ospf_info := output_state[vrfName].(map[string]interface{})
    ospfv2Area_obj, area_info, err = ospfv2_get_or_create_area (ospf_info, ospfv2_obj, area_id, vrfName)

    if nil == ospfv2Area_obj {
        log.Errorf("Failed to create a new area:%s, err=%s", area_id, err)
        return oper_err
    }

    if _stubEnable, ok := area_info["stubEnable"].(bool); ok {
        if _stubEnable {
            ospfv2_fill_area_stub_state(ospfv2Area_obj, area_info)    
        }
    }
    
    return  err;
}
var DbToYang_ospfv2_neighbors_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF- Neighbor State"
    var vtysh_cmd string
    var area_id, intf_name, ntv_intf_name string
    var temp interface{}
    var intf_area_id string
    var ospfv2Interface_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2_Areas_Area_Interfaces_Interface

    log.Info("DbToYang_ospfv2_neighbors_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }
    log.Info("vrfName=", vrfName)

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)
    area_id =pathInfo.Var("identifier#2")
    if(len(area_id) == 0) {
        log.Info("Area Id is not specified, key is missing")
        log.Errorf ("%s failed !! Error", cmn_log);
        return  oper_err
    } else {
        area_id = getAreaDotted(area_id)
        log.Infof("Area Id %s", area_id)
    }
    
    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " neighbor detail json"
    output_nbrs_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf neighbor detail:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_nbrs_state)
    log.Info(vrfName)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " interface json"
    output_interfaces, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf interfaces:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_interfaces)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " interface traffic json"
    output_interfaces_traffic, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf interfaces traffic:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_interfaces_traffic)
    for _,value := range output_interfaces { 
        interfaces_info := value.(map[string]interface{})
        interface_map := interfaces_info["interfaces"].(map[string]interface{})
        for ntv_intf_name, temp = range interface_map {
            log.Info("interface is ", ntv_intf_name)
            intf_name, _ = ospfGetUIIntfName(ntv_intf_name)
            
            intf_info := temp.(map[string]interface{})
            if intf_area_str,ok := intf_info["area"].(string); ok {
                result := strings.Split(intf_area_str, " ") 
                intf_area_id = result[0]
                if (intf_area_id != area_id) {
                    log.Infof("Skipping Interface %s belonging to area %s, as given area %s", intf_name, intf_area_id, area_id)
                    continue
                }
            }
            if !strings.Contains(intf_name, "VLINK") {
                for _,value := range output_nbrs_state {
                    neighbors_info := value.(map[string]interface{})
                    err = ospfv2_fill_neighbors_state (neighbors_info, ospfv2_obj, area_id, intf_name, vrfName)
					if err != nil {
						log.Info("Failed to fill neighbor state information")
					}
                }
                ospfv2_fill_interface_state(intf_info, ospfv2_obj, area_id, intf_name, vrfName, output_interfaces_traffic)
                ospfv2Interface_obj, err =  ospfv2_fill_interface_timers_state(intf_info, ospfv2_obj, area_id, intf_name, vrfName)
                if (nil != ospfv2Interface_obj) {
                    ospfv2_fill_interface_message_stats(output_interfaces_traffic, ospfv2Interface_obj, intf_name)
                }
            }
        }
    }
    return  err;
}
/*
var DbToYang_ospfv2_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    var cmd_err error
    oper_err := errors.New("Operational error")
    cmn_log := "GET: xfmr for OSPF-Global State"
    var vtysh_cmd string

    log.Info("DbToYang_ospfv2_state_xfmr ***", inParams.uri)
    var ospfv2_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Ospfv2
    ospfv2_obj, vrfName, err := getOspfv2Root (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return  oper_err
    }
    log.Info(vrfName)

    // get the values from the backend
    pathInfo := NewPathInfo(inParams.uri)

    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    log.Info(targetUriPath)
    vtysh_cmd = "show ip ospf vrf " + vrfName + " json"
    output_state, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
      log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
      return  cmd_err
    }
    
    log.Info(output_state)
    log.Info(vrfName)
    
    for key,value := range output_state {
        ospf_info := value.(map[string]interface{})
        log.Info(key)
        log.Info(ospf_info)
        //err = ospfv2_fill_global_state (ospf_info, ospfv2_obj)
        //err = ospfv2_fill_areas_state (ospf_info, ospfv2_obj)
    }
    
    return  err;
}
*/
var ospfv2_router_area_tbl_xfmr TableXfmrFunc = func (inParams XfmrParams)  ([]string, error) {
    var tblList []string
    var err error
    var vrf,key,areakey string

    log.Info("ospfv2_router_area_tbl_xfmr: ", inParams.uri)
    pathInfo := NewPathInfo(inParams.uri)

    vrf = pathInfo.Var("name")
    ospfId      := pathInfo.Var("identifier")
    protoName  := pathInfo.Var("name#2")
    pArea_Id   := pathInfo.Var("identifier#2")

    if len(pathInfo.Vars) <  3 {
        err = errors.New("Invalid Key length");
        log.Info("Invalid Key length", len(pathInfo.Vars))
        return tblList, err
    }

    if len(vrf) == 0 {
        err = errors.New("vrf name is missing");
        log.Info("VRF Name is Missing")
        return tblList, err
    }
    if !strings.Contains(ospfId,"OSPF") {
        err = errors.New("OSPF ID is missing");
        log.Info("OSPF ID is missing")
        return tblList, err
    }
    if len(protoName) == 0 {
        err = errors.New("Protocol Name is missing");
        log.Info("Protocol Name is Missing")
        return tblList, err
    }

    if (inParams.oper != GET) {
        tblList = append(tblList, "OSPFV2_ROUTER_AREA")
        return tblList, nil
    }

    tblList = append(tblList, "OSPFV2_ROUTER_AREA")
    if len(pArea_Id) != 0 {
        /* To see when it can be used later, curently not needed
        key = vrf + "|" + pArea_Id
        log.Info("ospfv2_router_area_tbl_xfmr: key - ", key)
        if (inParams.dbDataMap != nil) {
            if _, ok := (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"]; !ok {
                (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"] = make(map[string]db.Value)
            }

            areaCfgTblTs := &db.TableSpec{Name: "OSPFV2_ROUTER_AREA"}
            areaEntryKey := db.Key{Comp: []string{vrf, pArea_Id}}

            if _, ok := (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"][key]; !ok {
                var entryValue db.Value
                if entryValue, err = inParams.d.GetEntry(areaCfgTblTs, areaEntryKey) ; err == nil {
                    (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"][key] = entryValue
                }
            }
        }*/
    } else {
        if(inParams.dbDataMap != nil) {
            cmd := "show ip ospf vrf" + " " + vrf + " " + "json"
            output_state, cmd_err := exec_vtysh_cmd (cmd)
            if cmd_err != nil {
                log.Errorf("Failed to fetch ospf global state:, err=%s", cmd_err)
                return  tblList, cmd_err
            }
            for _,value := range output_state {
                ospf_info := value.(map[string]interface{})
                if value, ok := ospf_info["areas"]; ok {
                    areas_map := value.(map[string]interface {})
                    if(len(areas_map) == 0) {
                        log.Errorf("Does not contain any area")
                        err = errors.New("Does not contain any area");
                        return tblList, err
                    }
                    if _, ok := (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"]; !ok {
                        (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"] = make(map[string]db.Value)
                    }

                    for key = range areas_map {
                        log.Info(key)
                        areakey = vrf + "|" + key
                        log.Info("ospfv2_router_area_tbl_xfmr: OSPF Area key - ", areakey)
                        if _, ok := (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"][areakey]; !ok {
                            (*inParams.dbDataMap)[db.ConfigDB]["OSPFV2_ROUTER_AREA"][areakey] = db.Value{Field: make(map[string]string)}
                        }
                    }
                }
            }
        }

    }
    return tblList, nil
}

var rpc_clear_ospfv2 RpcCallpoint = func(body []byte, dbs [db.MaxDB]*db.DB) ([]byte, error) {

    var err error
    var status string
    var mapData map[string]interface{}

    err = json.Unmarshal(body, &mapData)
    if err != nil {
        log.Info("Failed to unmarshall given input data")
        return nil, err
    }

    var result struct {
        Output struct {
              Status string `json:"response"`
        } `json:"sonic-ospfv2-clear:output"`
    }

    input := mapData["sonic-ospfv2-clear:input"]
    mapData = input.(map[string]interface{})

    log.Info("rpc_clear_ospfv2: mapData ", mapData)

    vrfName := "default" 
    intfName := ""
    intfAll := false 

    if value, ok := mapData["vrf-name"].(string) ; ok {
        if (value != "") {
           vrfName = value
        }
    }

    if value, ok := mapData["interface-all"].(bool) ; ok {
        if value {
            intfAll = true
        }
    }

    if value, ok := mapData["interface"].(string) ; ok {
        if value != "" {
            value, err = ospfGetNativeIntfName(value) 
            if (err != nil) {
                return nil, tlerr.New("Invalid OSPF interface name.")
            }
            intfName = value 
        }
    }

    cmdStr := ""
    if (intfAll) {
        cmdStr = "clear ip ospf vrf " + vrfName + " interface"
    } else if (intfName != "") {
        cmdStr = "clear ip ospf vrf " + vrfName + " interface " + intfName
    }

    log.Infof("rpc_clear_ospfv2: vrf-%s intf-%s all-%v.", vrfName, intfName, intfAll)

    if cmdStr != "" {
       exec_vtysh_cmd(cmdStr)
       status = "Success"
    } else {
       log.Error("rpc_clear_ospfv2: Invalid input received mapData ", mapData)
       status = "Failed"
    }

    log.Infof("rpc_clear_ospfv2: %s", status)
    result.Output.Status = status
    return json.Marshal(&result)
}


