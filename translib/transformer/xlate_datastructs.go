////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2020 Dell, Inc.                                                 //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//  http://www.apache.org/licenses/LICENSE-2.0                                //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

package transformer

import (
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "github.com/openconfig/ygot/ygot"
    "regexp"
)

var rgpIpv6, rgpMac, rgpIsMac, rgpKeyExtract, rgpSncKeyExtract *regexp.Regexp

//Map of map[uri][dbKey]
var keyXfmrCache map[string]string

type KeySpec struct {
        DbNum db.DBNum
        Ts    db.TableSpec
        Key   db.Key
        Child []KeySpec
        IgnoreParentKey bool
}

type NotificationType int
const (
    Sample NotificationType = iota
    OnChange
)

type XfmrTranslateSubscribeInfo struct {
    DbDataMap RedisDbMap
    MinInterval int
    NeedCache bool
    PType NotificationType
    OnChange bool
}

type xpathTblKeyExtractRet struct {
    xpath string
    tableName string
    dbKey string
    isVirtualTbl bool
}

type xlateFromDbParams struct {
	d *db.DB //current db
	dbs [db.MaxDB]*db.DB
	curDb db.DBNum
	ygRoot *ygot.GoStruct
	uri string
	requestUri string //original uri using which a curl/NBI request is made
	oper int
	dbDataMap *map[db.DBNum]map[string]map[string]db.Value
	// subOpDataMap map[int]*RedisDbMap // used to add an in-flight data with a sub-op
	// param interface{}
	txCache interface{}
	//  skipOrdTblChk *bool
	//  pCascadeDelTbl *[] string //used to populate list of tables needed cascade delete by subtree overloaded methods
	xpath string //curr uri xpath
	tbl string
	tblKey string
	resultMap map[string]interface{}
	validate bool
}

type xlateToParams struct {
        d *db.DB
        ygRoot *ygot.GoStruct
        oper int
        uri string
        requestUri string
        xpath string
        keyName string
        jsonData interface{}
        resultMap map[int]RedisDbMap
        result map[string]map[string]db.Value
        txCache interface{}
        tblXpathMap map[string]map[string]map[string]bool
        subOpDataMap map[int]*RedisDbMap
        pCascadeDelTbl *[]string
        xfmrErr *error
        name string
        value interface{}
        tableName string
        yangDefValMap map[string]map[string]db.Value
        yangAuxValMap map[string]map[string]db.Value
}

