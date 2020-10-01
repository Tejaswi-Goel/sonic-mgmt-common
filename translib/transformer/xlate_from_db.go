////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2019 Dell, Inc.                                                 //
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
    "fmt"
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "strings"
    "encoding/json"
    "strconv"
    "errors"
    "sync"
    "github.com/openconfig/goyang/pkg/yang"
    "github.com/Azure/sonic-mgmt-common/translib/tlerr"

    log "github.com/golang/glog"
)

type typeMapOfInterface map[string]interface{}

var mapCopyMutex = &sync.Mutex{}

func DbValToInt(dbFldVal string, base int, size int, isUint bool) (interface{}, error) {
	var res interface{}
	var err error
	if isUint {
		if res, err = strconv.ParseUint(dbFldVal, base, size); err != nil {
			log.Warningf("Non Yint%v type for yang leaf-list item %v", size, dbFldVal)
		}
	} else {
		if res, err = strconv.ParseInt(dbFldVal, base, size); err != nil {
			log.Warningf("Non Yint %v type for yang leaf-list item %v", size, dbFldVal)
		}
	}
	return res, err
}

func getLeafrefRefdYangType(yngTerminalNdDtType yang.TypeKind, fldXpath string) (yang.TypeKind) {
	if yngTerminalNdDtType == yang.Yleafref {
		var entry *yang.Entry
		var path string
		if _, ok := xDbSpecMap[fldXpath]; ok {
			path = xDbSpecMap[fldXpath].dbEntry.Type.Path
			entry = xDbSpecMap[fldXpath].dbEntry
		} else if _, ok := xYangSpecMap[fldXpath]; ok {
			path = xYangSpecMap[fldXpath].yangEntry.Type.Path
			entry = xYangSpecMap[fldXpath].yangEntry
		}
		path = stripAugmentedModuleNames(path)
		path = path[1:]
		xfmrLogInfoAll("Received path %v for FieldXpath %v", path, fldXpath)
		if strings.Contains(path, "..") {
			if entry != nil && len(path) > 0 {
				// Referenced path within same yang file
				xpath, _, err := XfmrRemoveXPATHPredicates(path)
				if  err != nil {
					log.Warningf("error in XfmrRemoveXPATHPredicates %v", path)
					return yngTerminalNdDtType
				}
				xpath = xpath[1:]
				pathList := strings.Split(xpath, "/")
				for _, x := range pathList {
					if x == ".." {
						entry = entry.Parent
					} else {
						if _,ok := entry.Dir[x]; ok {
							entry = entry.Dir[x]
						}
					}
				}
				if entry != nil && entry.Type != nil {
					yngTerminalNdDtType = entry.Type.Kind
					xfmrLogInfoAll("yangLeaf datatype %v", yngTerminalNdDtType)
					if yngTerminalNdDtType == yang.Yleafref {
						leafPath := getXpathFromYangEntry(entry)
						xfmrLogInfoAll("getLeafrefRefdYangType: xpath for leafref type:%v",leafPath)
						return getLeafrefRefdYangType(yngTerminalNdDtType, leafPath)
					}
				}
			}
		} else if len(path) > 0 {
			// Referenced path in a different yang file
			xpath, _, err := XfmrRemoveXPATHPredicates(path)
			if  err != nil {
				log.Warningf("error in XfmrRemoveXPATHPredicates %v", xpath)
				return yngTerminalNdDtType
			}
			// Form xpath based on sonic or non sonic yang path
			if strings.Contains(xpath, "sonic") {
				pathList := strings.Split(xpath, "/")
				xpath = pathList[SONIC_TABLE_INDEX]+ "/" + pathList[SONIC_FIELD_INDEX]
				if _, ok := xDbSpecMap[xpath]; ok {
					yngTerminalNdDtType = xDbSpecMap[xpath].dbEntry.Type.Kind
				}

			} else {
				xpath = replacePrefixWithModuleName(xpath)
				if _, ok := xYangSpecMap[xpath]; ok {
					yngTerminalNdDtType = xYangSpecMap[xpath].dbEntry.Type.Kind
				}
			}

		}
		xfmrLogInfoAll("yangLeaf datatype %v", yngTerminalNdDtType)
	}
	return yngTerminalNdDtType
}

func DbToYangType(yngTerminalNdDtType yang.TypeKind, fldXpath string, dbFldVal string) (interface{}, interface{}, error) {
	xfmrLogInfoAll("Received FieldXpath %v, yngTerminalNdDtType %v and Db field value %v to be converted to yang data-type.", fldXpath, yngTerminalNdDtType, dbFldVal)
	var res interface{}
	var resPtr interface{}
	var err error
	const INTBASE = 10

	if yngTerminalNdDtType == yang.Yleafref {
		yngTerminalNdDtType = getLeafrefRefdYangType(yngTerminalNdDtType, fldXpath)
	}

	switch yngTerminalNdDtType {
        case yang.Ynone:
                log.Warning("Yang node data-type is non base yang type")
		//TODO - enhance to handle non base data types depending on future use case
		err = errors.New("Yang node data-type is non base yang type")
        case yang.Yint8:
                res, err = DbValToInt(dbFldVal, INTBASE, 8, false)
		var resInt8 int8 = int8(res.(int64))
		resPtr = &resInt8
        case yang.Yint16:
                res, err = DbValToInt(dbFldVal, INTBASE, 16, false)
		var resInt16 int16 = int16(res.(int64))
		resPtr = &resInt16
        case yang.Yint32:
                res, err = DbValToInt(dbFldVal, INTBASE, 32, false)
		var resInt32 int32 = int32(res.(int64))
		resPtr = &resInt32
        case yang.Yuint8:
                res, err = DbValToInt(dbFldVal, INTBASE, 8, true)
		var resUint8 uint8 = uint8(res.(uint64))
		resPtr = &resUint8
        case yang.Yuint16:
                res, err = DbValToInt(dbFldVal, INTBASE, 16, true)
		var resUint16 uint16 = uint16(res.(uint64))
		resPtr = &resUint16
        case yang.Yuint32:
                res, err = DbValToInt(dbFldVal, INTBASE, 32, true)
		var resUint32 uint32 = uint32(res.(uint64))
		resPtr = &resUint32
        case yang.Ybool:
		if res, err = strconv.ParseBool(dbFldVal); err != nil {
			log.Warningf("Non Bool type for yang leaf-list item %v", dbFldVal)
		}
		var resBool bool = res.(bool)
		resPtr = &resBool
        case yang.Ybinary, yang.Ydecimal64, yang.Yenum, yang.Yidentityref, yang.Yint64, yang.Yuint64, yang.Ystring, yang.Yunion, yang.Yleafref:
                // TODO - handle the union type
                // Make sure to encode as string, expected by util_types.go: ytypes.yangToJSONType
                xfmrLogInfoAll("Yenum/Ystring/Yunion(having all members as strings) type for yangXpath %v", fldXpath)
                res = dbFldVal
		var resString string = res.(string)
		resPtr = &resString
	case yang.Yempty:
		logStr := fmt.Sprintf("Yang data type for xpath %v is Yempty.", fldXpath)
		log.Warning(logStr)
		err = errors.New(logStr)
        default:
		logStr := fmt.Sprintf("Unrecognized/Unhandled yang-data type(%v) for xpath %v.", fldXpath, yang.TypeKindToName[yngTerminalNdDtType])
                log.Warning(logStr)
                err = errors.New(logStr)
        }
	return res, resPtr, err
}

/*convert leaf-list in Db to leaf-list in yang*/
func processLfLstDbToYang(fieldXpath string, dbFldVal string, yngTerminalNdDtType yang.TypeKind) []interface{} {
	valLst := strings.Split(dbFldVal, ",")
	var resLst []interface{}

	xfmrLogInfoAll("xpath: %v, dbFldVal: %v", fieldXpath, dbFldVal)
	switch  yngTerminalNdDtType {
	case yang.Ybinary, yang.Ydecimal64, yang.Yenum, yang.Yidentityref, yang.Yint64, yang.Yuint64, yang.Ystring, yang.Yunion:
                // TODO - handle the union type.OC yang should have field xfmr.sonic-yang?
                // Make sure to encode as string, expected by util_types.go: ytypes.yangToJSONType:
		xfmrLogInfoAll("DB leaf-list and Yang leaf-list are of same data-type")
		for _, fldVal := range valLst {
			resLst = append(resLst, fldVal)
		}
	default:
		for _, fldVal := range valLst {
			resVal, _, err := DbToYangType(yngTerminalNdDtType, fieldXpath, fldVal)
			if err == nil {
				resLst = append(resLst, resVal)
			}
		}
	}
	return resLst
}

func sonicDbToYangTerminalNodeFill(field string, inParamsForGet xlateFromDbParams) {
	resField := field
	value := ""

	if inParamsForGet.dbDataMap != nil {
		tblInstFields, dbDataExists := (*inParamsForGet.dbDataMap)[inParamsForGet.curDb][inParamsForGet.tbl][inParamsForGet.tblKey]
		if dbDataExists {
			fieldVal, valueExists := tblInstFields.Field[field]
			if !valueExists {
				return
			}
			value = fieldVal
		} else {
			return
		}
	}

	if strings.HasSuffix(field, "@") {
		fldVals := strings.Split(field, "@")
		resField = fldVals[0]
	}
	fieldXpath := inParamsForGet.tbl + "/" + resField
	xDbSpecMapEntry, ok := xDbSpecMap[fieldXpath]
	if !ok {
		log.Warningf("No entry found in xDbSpecMap for xpath %v", fieldXpath)
		return
	}
	if xDbSpecMapEntry.dbEntry == nil {
		log.Warningf("Yang entry is nil in xDbSpecMap for xpath %v", fieldXpath)
		return
	}

	yangType := yangTypeGet(xDbSpecMapEntry.dbEntry)
	yngTerminalNdDtType := xDbSpecMapEntry.dbEntry.Type.Kind
	if yangType ==  YANG_LEAF_LIST {
		/* this should never happen but just adding for safetty */
		if !strings.HasSuffix(field, "@") {
			log.Warningf("Leaf-list in Sonic yang should also be a leaf-list in DB, its not for xpath %v", fieldXpath)
			return
		}
		resLst := processLfLstDbToYang(fieldXpath, value, yngTerminalNdDtType)
		inParamsForGet.resultMap[resField] = resLst
	} else { /* yangType is leaf - there are only 2 types of yang terminal node leaf and leaf-list */
		resVal, _, err := DbToYangType(yngTerminalNdDtType, fieldXpath, value)
		if err != nil {
			log.Warningf("Failure in converting Db value type to yang type for xpath", fieldXpath)
		} else {
			inParamsForGet.resultMap[resField] = resVal
		}
	}
}

func sonicDbToYangListFill(inParamsForGet xlateFromDbParams) []typeMapOfInterface {
	var mapSlice []typeMapOfInterface
	dbDataMap := inParamsForGet.dbDataMap
	table := inParamsForGet.tbl
	dbIdx := inParamsForGet.curDb
	xpath := inParamsForGet.xpath
	dbTblData := (*dbDataMap)[dbIdx][table]

	for keyStr := range dbTblData {
		curMap := make(map[string]interface{})
		linParamsForGet := formXlateFromDbParams(inParamsForGet.dbs[dbIdx], inParamsForGet.dbs, dbIdx, inParamsForGet.ygRoot, inParamsForGet.uri, inParamsForGet.requestUri, xpath, inParamsForGet.oper, table, keyStr, dbDataMap, inParamsForGet.txCache, curMap, inParamsForGet.validate)
		sonicDbToYangDataFill(linParamsForGet)
		curMap = linParamsForGet.resultMap
		dbDataMap = linParamsForGet.dbDataMap
		inParamsForGet.dbDataMap = dbDataMap
		dbSpecData, ok := xDbSpecMap[table]
		if ok && dbSpecData.keyName == nil {
			yangKeys := yangKeyFromEntryGet(xDbSpecMap[xpath].dbEntry)
			sonicKeyDataAdd(dbIdx, yangKeys, table, keyStr, curMap)
		}
		if len(curMap) > 0 {
			mapSlice = append(mapSlice, curMap)
		}
	}
	return mapSlice
}

func sonicDbToYangDataFill(inParamsForGet xlateFromDbParams) {
	xpath := inParamsForGet.xpath
	uri := inParamsForGet.uri
	table := inParamsForGet.tbl
	key := inParamsForGet.tblKey
	resultMap := inParamsForGet.resultMap
	dbDataMap := inParamsForGet.dbDataMap
	dbIdx := inParamsForGet.curDb
	yangNode, ok := xDbSpecMap[xpath]

	if ok  && yangNode.dbEntry != nil {
		xpathPrefix := table
		if len(table) > 0 { xpathPrefix += "/" }

		for yangChldName := range yangNode.dbEntry.Dir {
			chldXpath := xpathPrefix+yangChldName
			if xDbSpecMap[chldXpath] != nil && xDbSpecMap[chldXpath].dbEntry != nil {
				chldYangType := yangTypeGet(xDbSpecMap[chldXpath].dbEntry)

				if  chldYangType == YANG_LEAF || chldYangType == YANG_LEAF_LIST {
					xfmrLogInfoAll("tbl(%v), k(%v), yc(%v)", table, key, yangChldName)
					fldName := yangChldName
					if chldYangType == YANG_LEAF_LIST  {
						fldName = fldName + "@"
					}
				        curUri := inParamsForGet.uri + "/" + yangChldName
					linParamsForGet := formXlateFromDbParams(nil, inParamsForGet.dbs, dbIdx, inParamsForGet.ygRoot, curUri, inParamsForGet.requestUri, curUri, inParamsForGet.oper, table, key, dbDataMap, inParamsForGet.txCache, resultMap, inParamsForGet.validate)
                                        sonicDbToYangTerminalNodeFill(fldName, linParamsForGet)
					resultMap = linParamsForGet.resultMap
					inParamsForGet.resultMap = resultMap
				} else if chldYangType == YANG_CONTAINER {
					curMap := make(map[string]interface{})
					curUri := xpath + "/" + yangChldName
					// container can have a static key, so extract key for current container
					_, curKey, curTable := sonicXpathKeyExtract(curUri)
					// use table-name as xpath from now on
					d := inParamsForGet.dbs[xDbSpecMap[curTable].dbIndex]
					linParamsForGet := formXlateFromDbParams(d, inParamsForGet.dbs, xDbSpecMap[curTable].dbIndex, inParamsForGet.ygRoot, curUri, inParamsForGet.requestUri, curTable, inParamsForGet.oper, curTable, curKey, dbDataMap, inParamsForGet.txCache, curMap, inParamsForGet.validate)
					sonicDbToYangDataFill(linParamsForGet)
					curMap = linParamsForGet.resultMap
					dbDataMap = linParamsForGet.dbDataMap
					if len(curMap) > 0 {
						resultMap[yangChldName] = curMap
					} else {
						xfmrLogInfoAll("Empty container for xpath(%v)", curUri)
					}
					inParamsForGet.dbDataMap = linParamsForGet.dbDataMap
					inParamsForGet.resultMap = resultMap
				} else if chldYangType == YANG_LIST {
					pathList := strings.Split(uri, "/")
					// Skip the list entries if the uri has specific list query
					if len(pathList) > SONIC_TABLE_INDEX+1 && !strings.Contains(uri,yangChldName) {
						xfmrLogInfoAll("Skipping yangChldName: %v, pathList:%v, len:%v", yangChldName, pathList, len(pathList))
					} else {
						var mapSlice []typeMapOfInterface
						curUri := xpath + "/" + yangChldName
						inParamsForGet.uri = curUri
						inParamsForGet.xpath = curUri
						mapSlice = sonicDbToYangListFill(inParamsForGet)
						dbDataMap = inParamsForGet.dbDataMap
						if len(key) > 0 && len(mapSlice) == 1 {// Single instance query. Don't return array of maps
							for k, val := range mapSlice[0] {
								resultMap[k] = val
							}

						} else if len(mapSlice) > 0 {
							resultMap[yangChldName] = mapSlice
						} else {
							xfmrLogInfoAll("Empty list for xpath(%v)", curUri)
						}
						inParamsForGet.resultMap = resultMap
					}
				} else if chldYangType == YANG_CHOICE || chldYangType == YANG_CASE {
					curUri := table + "/" + yangChldName
					inParamsForGet.uri = curUri
					inParamsForGet.xpath = curUri
					inParamsForGet.curDb = xDbSpecMap[table].dbIndex
					sonicDbToYangDataFill(inParamsForGet)
					dbDataMap = inParamsForGet.dbDataMap
					resultMap = inParamsForGet.resultMap
				} else {
					xfmrLogInfoAll("Not handled case %v", chldXpath)
				}
			} else {
				xfmrLogInfoAll("Yang entry not found for %v", chldXpath)
			}
		}
	}
}

/* Traverse db map and create json for cvl yang */
func directDbToYangJsonCreate(inParamsForGet xlateFromDbParams) (string, bool, error) {
	var err error
	uri := inParamsForGet.uri
	dbDataMap := inParamsForGet.dbDataMap
	resultMap := inParamsForGet.resultMap
	xpath, key, table := sonicXpathKeyExtract(uri)
	inParamsForGet.xpath = xpath
	inParamsForGet.tbl = table
	inParamsForGet.tblKey = key

	if len(xpath) > 0 {
		var dbNode *dbInfo

		if len(table) > 0 {
			tokens:= strings.Split(xpath, "/")
			if tokens[SONIC_TABLE_INDEX] == table {
				fieldName := tokens[len(tokens)-1]
				dbSpecField := table + "/" + fieldName
				_, ok := xDbSpecMap[dbSpecField]
				if ok && (xDbSpecMap[dbSpecField].fieldType == YANG_LEAF || xDbSpecMap[dbSpecField].fieldType == YANG_LEAF_LIST) {
					dbNode = xDbSpecMap[dbSpecField]
					xpath = dbSpecField
					inParamsForGet.xpath = xpath
				} else {
					dbNode = xDbSpecMap[table]
				}
			}
		} else {
			dbNode = xDbSpecMap[xpath]
		}

		if dbNode != nil && dbNode.dbEntry != nil {
			cdb   := db.ConfigDB
			yangType := yangTypeGet(dbNode.dbEntry)
			if len(table) > 0 {
				cdb = xDbSpecMap[table].dbIndex
			}
			inParamsForGet.curDb = cdb

			if yangType == YANG_LEAF || yangType == YANG_LEAF_LIST {
				fldName := xDbSpecMap[xpath].dbEntry.Name
				if yangType == YANG_LEAF_LIST  {
					fldName = fldName + "@"
				}
				linParamsForGet := formXlateFromDbParams(nil, inParamsForGet.dbs, cdb, inParamsForGet.ygRoot, xpath, inParamsForGet.requestUri, uri, inParamsForGet.oper, table, key, dbDataMap, inParamsForGet.txCache, resultMap, inParamsForGet.validate)
				sonicDbToYangTerminalNodeFill(fldName, linParamsForGet)
				resultMap = linParamsForGet.resultMap
			} else if yangType == YANG_CONTAINER {
				if len(table) > 0 {
					xpath = table
					inParamsForGet.xpath = xpath
				}
				sonicDbToYangDataFill(inParamsForGet)
				resultMap = inParamsForGet.resultMap
			} else if yangType == YANG_LIST {
				mapSlice := sonicDbToYangListFill(inParamsForGet)
				if len(key) > 0 && len(mapSlice) == 1 {// Single instance query. Don't return array of maps
                                                for k, val := range mapSlice[0] {
                                                        resultMap[k] = val
                                                }

                                } else if len(mapSlice) > 0 {
					pathl := strings.Split(xpath, "/")
					lname := pathl[len(pathl) - 1]
					resultMap[lname] = mapSlice
				}
			}
		}
	}

	jsonMapData, _ := json.Marshal(resultMap)
	isEmptyPayload := isJsonDataEmpty(string(jsonMapData))
	jsonData := fmt.Sprintf("%v", string(jsonMapData))
        if isEmptyPayload {
		log.Warning("No data available")
        }
        return jsonData, isEmptyPayload, err
}

func tableNameAndKeyFromDbMapGet(dbDataMap map[string]map[string]db.Value) (string, string, error) {
    tableName := ""
    tableKey  := ""
    for tn, tblData := range dbDataMap {
        tableName = tn
        for kname := range tblData {
            tableKey = kname
        }
    }
    return tableName, tableKey, nil
}

func fillDbDataMapForTbl(uri string, xpath string, tblName string, tblKey string, cdb db.DBNum, dbs [db.MaxDB]*db.DB, dbTblKeyGetCache map[db.DBNum]map[string]map[string]bool) (map[db.DBNum]map[string]map[string]db.Value, error) {
	var err error
	dbresult  := make(RedisDbMap)
	dbresult[cdb] = make(map[string]map[string]db.Value)
	dbFormat := KeySpec{}
	dbFormat.Ts.Name = tblName
	dbFormat.DbNum = cdb
	if tblKey != "" {
		if tblSpecInfo, ok := xDbSpecMap[tblName]; ok && tblSpecInfo.hasXfmrFn {
			/* key from uri should be converted into redis-db key, to read data */
			tblKey, err = dbKeyValueXfmrHandler(CREATE, cdb, tblName, tblKey)
			if err != nil {
				log.Warningf("Value-xfmr for table(%v) & key(%v) didn't do conversion.", tblName, tblKey)
				return nil, err
			}
		}

		dbFormat.Key.Comp = append(dbFormat.Key.Comp, tblKey)
	}
	err = TraverseDb(dbs, dbFormat, &dbresult, nil, dbTblKeyGetCache)
	if err != nil {
		log.Warningf("TraverseDb() didn't fetch data for tbl(DB num) %v(%v) for xpath %v", tblName, cdb, xpath)
		return nil, err
	}
	if _, ok := dbresult[cdb]; !ok {
		logStr := fmt.Sprintf("TraverseDb() did not populate Db data for tbl(DB num) %v(%v) for xpath %v", tblName, cdb, xpath)
		err = fmt.Errorf("%v", logStr)
		return nil, err
	}
	return dbresult, err

}

// Assumption: All tables are from the same DB
func dbDataFromTblXfmrGet(tbl string, inParams XfmrParams, dbDataMap *map[db.DBNum]map[string]map[string]db.Value, dbTblKeyGetCache map[db.DBNum]map[string]map[string]bool) error {
    // skip the query if the table is already visited
    if _,ok := (*dbDataMap)[inParams.curDb][tbl]; ok {
       if len(inParams.key) > 0 {
          if  _,ok = (*dbDataMap)[inParams.curDb][tbl][inParams.key]; ok {
             return nil
          }
       } else {
          return nil
       }
    }
    xpath, _, _ := XfmrRemoveXPATHPredicates(inParams.uri)

	terminalNodeGet  := false
	qdbMapHasTblData := false
	qdbMapHasTblKeyData := false
	if !xYangSpecMap[xpath].hasNonTerminalNode  && len(inParams.key) > 0 {
		terminalNodeGet = true
	}
	if qdbMap, getOk := dbTblKeyGetCache[inParams.curDb]; getOk {
		if dbTblData, tblPresent := qdbMap[tbl]; tblPresent {
			qdbMapHasTblData = true
			if _, keyPresent := dbTblData[inParams.key]; keyPresent {
				qdbMapHasTblKeyData = true;
			}
		}
	}

	if !qdbMapHasTblData || (terminalNodeGet && qdbMapHasTblData && !qdbMapHasTblKeyData) {
		curDbDataMap, err := fillDbDataMapForTbl(inParams.uri, xpath, tbl, inParams.key, inParams.curDb, inParams.dbs, dbTblKeyGetCache)
		if err == nil {
			mapCopy((*dbDataMap)[inParams.curDb], curDbDataMap[inParams.curDb])
		}
	}
    return nil
}

func yangListDataFill(inParamsForGet xlateFromDbParams, isFirstCall bool) error {
	var tblList []string
	dbs := inParamsForGet.dbs
	ygRoot := inParamsForGet.ygRoot
	uri := inParamsForGet.uri
	requestUri := inParamsForGet.requestUri
	dbDataMap := inParamsForGet.dbDataMap
	txCache := inParamsForGet.txCache
	cdb := inParamsForGet.curDb
	resultMap := inParamsForGet.resultMap
	xpath := inParamsForGet.xpath
	tbl := inParamsForGet.tbl
	tblKey := inParamsForGet.tblKey


	_, ok := xYangSpecMap[xpath]
	if ok {
	if xYangSpecMap[xpath].xfmrTbl != nil {
		xfmrTblFunc := *xYangSpecMap[xpath].xfmrTbl
		if len(xfmrTblFunc) > 0 {
			inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, tblKey, dbDataMap, nil, nil, txCache)
			tblList, _   = xfmrTblHandlerFunc(xfmrTblFunc, inParams, inParamsForGet.xfmrDbTblKeyCache)
			inParamsForGet.dbDataMap = dbDataMap
			inParamsForGet.ygRoot = ygRoot
			if len(tblList) != 0 {
				for _, curTbl := range tblList {
					dbDataFromTblXfmrGet(curTbl, inParams, dbDataMap, inParamsForGet.dbTblKeyGetCache)
					inParamsForGet.dbDataMap = dbDataMap
					inParamsForGet.ygRoot = ygRoot
				}
			}
		}
		if tbl != "" {
			if !contains(tblList, tbl) {
				tblList = append(tblList, tbl)
			}
		}
	} else if tbl != "" && xYangSpecMap[xpath].xfmrTbl == nil {
		tblList = append(tblList, tbl)
	} else if tbl == "" && xYangSpecMap[xpath].xfmrTbl == nil {
		// Handling for case: Parent list is not associated with a tableName but has children containers/lists having tableNames.
		if tblKey != "" {
			var mapSlice []typeMapOfInterface
			instMap, err := yangListInstanceDataFill(inParamsForGet, isFirstCall)
			dbDataMap = inParamsForGet.dbDataMap
			if err != nil {
				log.Infof("Error(%v) returned for %v", err, uri)
			} else if ((instMap != nil)  && (len(instMap) > 0)) {
				mapSlice = append(mapSlice, instMap)
			}

			if len(mapSlice) > 0 {
				listInstanceGet := false
				// Check if it is a list instance level Get
				if ((strings.HasSuffix(uri, "]")) || (strings.HasSuffix(uri, "]/"))) {
					listInstanceGet = true
					for k, v := range mapSlice[0] {
						resultMap[k] = v
					}
				}
				if !listInstanceGet {
					resultMap[xYangSpecMap[xpath].yangEntry.Name] = mapSlice
				}
				inParamsForGet.resultMap = resultMap
			}
		}
	}
	}

	for _, tbl = range(tblList) {
		inParamsForGet.tbl = tbl

		tblData, ok := (*dbDataMap)[cdb][tbl]

		if ok {
			var mapSlice []typeMapOfInterface
			for dbKey := range tblData {
				inParamsForGet.tblKey = dbKey
				instMap, err := yangListInstanceDataFill(inParamsForGet, isFirstCall)
				dbDataMap = inParamsForGet.dbDataMap
				if err != nil {
					log.Infof("Error(%v) returned for %v", err, uri)
				} else if ((instMap != nil)  && (len(instMap) > 0)) {
					mapSlice = append(mapSlice, instMap)
				}
			}

			if len(mapSlice) > 0 {
				listInstanceGet := false
				/*Check if it is a list instance level Get*/
				if ((strings.HasSuffix(uri, "]")) || (strings.HasSuffix(uri, "]/"))) {
					listInstanceGet = true
					for k, v := range mapSlice[0] {
						resultMap[k] = v
					}
				}
				if !listInstanceGet {
					if _, specOk := xYangSpecMap[xpath]; specOk {
					if _, ok := resultMap[xYangSpecMap[xpath].yangEntry.Name]; ok {
						mlen := len(resultMap[xYangSpecMap[xpath].yangEntry.Name].([]typeMapOfInterface))
						for i := 0; i < mlen; i++ {
							mapSlice = append(mapSlice, resultMap[xYangSpecMap[xpath].yangEntry.Name].([]typeMapOfInterface)[i])
						}
					}
					resultMap[xYangSpecMap[xpath].yangEntry.Name] = mapSlice
					inParamsForGet.resultMap = resultMap
					}
				}
			} else {
				xfmrLogInfoAll("Empty slice for (\"%v\").\r\n", uri)
			}
		}
	}// end of tblList for

	return nil
}

func yangListInstanceDataFill(inParamsForGet xlateFromDbParams, isFirstCall bool) (typeMapOfInterface,error) {

	var err error
	curMap := make(map[string]interface{})
	err = nil
	dbs := inParamsForGet.dbs
	ygRoot := inParamsForGet.ygRoot
	uri := inParamsForGet.uri
	requestUri := inParamsForGet.requestUri
	dbDataMap := inParamsForGet.dbDataMap
	txCache := inParamsForGet.txCache
	cdb := inParamsForGet.curDb
	xpath := inParamsForGet.xpath
	tbl := inParamsForGet.tbl
	dbKey := inParamsForGet.tblKey

	curKeyMap, curUri, err := dbKeyToYangDataConvert(uri, requestUri, xpath, tbl, dbDataMap, dbKey, dbs[cdb].Opts.KeySeparator, txCache)
        if ((err != nil) || (curKeyMap == nil) || (len(curKeyMap) == 0)) {
                xfmrLogInfoAll("Skip filling list instance for uri %v since no yang  key found corresponding to db-key %v", uri, dbKey)
               return curMap, err
        }
	parentXpath := parentXpathGet(xpath)
	_, ok := xYangSpecMap[xpath]
	if ok && len(xYangSpecMap[xpath].xfmrFunc) > 0 {
		if isFirstCall || (!isFirstCall && (uri != requestUri) && ((len(xYangSpecMap[parentXpath].xfmrFunc) == 0) ||
			(len(xYangSpecMap[parentXpath].xfmrFunc) > 0 && (xYangSpecMap[parentXpath].xfmrFunc != xYangSpecMap[xpath].xfmrFunc)))) {
			xfmrLogInfoAll("Parent subtree already handled cur uri: %v", xpath)
			inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, curUri, requestUri, GET, dbKey, dbDataMap, nil, nil, txCache)
			err := xfmrHandlerFunc(inParams)
			inParamsForGet.ygRoot = ygRoot
			inParamsForGet.dbDataMap = dbDataMap
			if err != nil {
				xfmrLogInfoAll("Error returned by %v: %v", xYangSpecMap[xpath].xfmrFunc, err)
			}
		}
		if xYangSpecMap[xpath].hasChildSubTree {
			linParamsForGet := formXlateFromDbParams(dbs[cdb], dbs, cdb, ygRoot, curUri, requestUri, xpath, inParamsForGet.oper, tbl, dbKey, dbDataMap, inParamsForGet.txCache, curMap, inParamsForGet.validate)
			linParamsForGet.xfmrDbTblKeyCache = inParamsForGet.xfmrDbTblKeyCache
			linParamsForGet.dbTblKeyGetCache = inParamsForGet.dbTblKeyGetCache
			yangDataFill(linParamsForGet)
			curMap = linParamsForGet.resultMap
			dbDataMap = linParamsForGet.dbDataMap
			ygRoot = linParamsForGet.ygRoot
			inParamsForGet.dbDataMap = dbDataMap
			inParamsForGet.ygRoot = ygRoot
		}
	} else {
		xpathKeyExtRet, _ := xpathKeyExtract(dbs[cdb], ygRoot, GET, curUri, requestUri, nil, txCache, inParamsForGet.xfmrDbTblKeyCache)
		keyFromCurUri := xpathKeyExtRet.dbKey
		inParamsForGet.ygRoot = ygRoot
		if dbKey == keyFromCurUri || keyFromCurUri == "" {
			if dbKey == keyFromCurUri {
				for k, kv := range curKeyMap {
					curMap[k] = kv
				}
			}
			curXpath, _, _ := XfmrRemoveXPATHPredicates(curUri)
			linParamsForGet := formXlateFromDbParams(dbs[cdb], dbs, cdb, ygRoot, curUri, requestUri, curXpath, inParamsForGet.oper, tbl, dbKey, dbDataMap, inParamsForGet.txCache, curMap, inParamsForGet.validate)
			linParamsForGet.xfmrDbTblKeyCache = inParamsForGet.xfmrDbTblKeyCache
			linParamsForGet.dbTblKeyGetCache = inParamsForGet.dbTblKeyGetCache
			yangDataFill(linParamsForGet)
			curMap = linParamsForGet.resultMap
			dbDataMap = linParamsForGet.dbDataMap
			ygRoot = linParamsForGet.ygRoot
			inParamsForGet.dbDataMap = dbDataMap
			inParamsForGet.ygRoot = ygRoot
		}
	}
	return curMap, err
}

func terminalNodeProcess(inParamsForGet xlateFromDbParams) (map[string]interface{}, error) {
	xfmrLogInfoAll("Received xpath - %v, uri - %v, table - %v, table key - %v", inParamsForGet.xpath, inParamsForGet.uri, inParamsForGet.tbl, inParamsForGet.tblKey)
	var err error
	resFldValMap := make(map[string]interface{})
	xpath := inParamsForGet.xpath
	dbs := inParamsForGet.dbs
	ygRoot := inParamsForGet.ygRoot
	uri := inParamsForGet.uri
	tbl := inParamsForGet.tbl
	tblKey := inParamsForGet.tblKey
	requestUri := inParamsForGet.requestUri
	dbDataMap := inParamsForGet.dbDataMap
	txCache := inParamsForGet.txCache

	_, ok := xYangSpecMap[xpath]
	if !ok || xYangSpecMap[xpath].yangEntry == nil {
		logStr := fmt.Sprintf("No yang entry found for xpath %v.", xpath)
		err = fmt.Errorf("%v", logStr)
		return resFldValMap, err
	}

	cdb := xYangSpecMap[xpath].dbIndex
	if len(xYangSpecMap[xpath].xfmrField) > 0 {
		inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, tblKey, dbDataMap, nil, nil, txCache)
		fldValMap, err := leafXfmrHandlerFunc(inParams, xYangSpecMap[xpath].xfmrField)
		inParamsForGet.ygRoot = ygRoot
		inParamsForGet.dbDataMap = dbDataMap
		if err != nil {
			xfmrLogInfoAll("No data from field transformer for %v: %v.", uri, err)
			return resFldValMap, err
		}
		if (uri == requestUri) {
			yangType := yangTypeGet(xYangSpecMap[xpath].yangEntry)
			if len(fldValMap) == 0 {
				// field transformer returns empty map when no data in DB
				if ((yangType == YANG_LEAF) || ((yangType == YANG_LEAF_LIST) && ((strings.HasSuffix(uri, "]")) || (strings.HasSuffix(uri, "]/"))))) {
					log.Warningf("Field transformer returned empty data , uri  - %v", requestUri)
					err = tlerr.NotFoundError{Format:"Resource not found"}
					return resFldValMap, err
				}
			} else {
				if ((yangType == YANG_LEAF_LIST) && ((strings.HasSuffix(uri, "]")) || (strings.HasSuffix(uri, "]/")))) {
					return resFldValMap, nil
				}
			}
		}
		for lf, val := range fldValMap {
			resFldValMap[lf] = val
		}
	} else {
		dbFldName := xYangSpecMap[xpath].fieldName
		if dbFldName == XFMR_NONE_STRING {
			return resFldValMap, err
		}
		/* if there is no transformer extension/annotation then it means leaf-list in yang is also leaflist in db */
		if len(dbFldName) > 0  && !xYangSpecMap[xpath].isKey {
			yangType := yangTypeGet(xYangSpecMap[xpath].yangEntry)
			yngTerminalNdDtType := xYangSpecMap[xpath].yangEntry.Type.Kind
			if yangType ==  YANG_LEAF_LIST {
				dbFldName += "@"
				val, ok := (*dbDataMap)[cdb][tbl][tblKey].Field[dbFldName]
				leafLstInstGetReq := false

				ruriXpath, _, _ := XfmrRemoveXPATHPredicates(inParamsForGet.requestUri)
				rYangType := ""
				if rSpecInfo, rok := xYangSpecMap[ruriXpath]; rok {
					rYangType = yangTypeGet(rSpecInfo.yangEntry)
				}

				if ((strings.HasSuffix(requestUri, "]")) || (strings.HasSuffix(requestUri, "]/"))) && (rYangType == YANG_LEAF_LIST) {
					xfmrLogInfoAll("Request URI is leaf-list instance GET - %v", requestUri)
					leafLstInstGetReq = true
				}
				if ok {
					if leafLstInstGetReq {
						leafListInstVal, valErr := extractLeafListInstFromUri(requestUri)
						if valErr != nil {
							return resFldValMap, valErr
						}
						dbSpecField := tbl + "/" + strings.TrimSuffix(dbFldName, "@")
						dbSpecFieldInfo, dbSpecOk := xDbSpecMap[dbSpecField]
						if dbSpecOk && dbSpecFieldInfo.xfmrValue != nil {
							inParams := formXfmrDbInputRequest(CREATE, cdb, tbl, tblKey, dbFldName, leafListInstVal)
							retVal, valXfmrErr := valueXfmrHandler(inParams, *dbSpecFieldInfo.xfmrValue)
							if valXfmrErr != nil {
								log.Warningf("value-xfmr:fldpath(\"%v\") val(\"%v\"):err(\"%v\").", dbSpecField, leafListInstVal, valXfmrErr)
								return resFldValMap, valXfmrErr
							}
							log.Info("valueXfmrHandler() retuned ", retVal)
							leafListInstVal = retVal
						}
						if !leafListInstExists((*dbDataMap)[cdb][tbl][tblKey].Field[dbFldName], leafListInstVal) {
							log.Warningf("Queried leaf-list instance does not exists, uri  - %v, dbData - %v", requestUri, (*dbDataMap)[cdb][tbl][tblKey].Field[dbFldName])
							err = tlerr.NotFoundError{Format:"Resource not found"}
						}
						if err == nil {
							/* Since translib already fills in ygRoot with queried leaf-list instance, do not
							   fill in resFldValMap or else Unmarshall of payload(resFldValMap) into ygotTgt in
							   app layer will create duplicate instances in result.
							 */
							 log.Info("Queried leaf-list instance exists but Since translib already fills in ygRoot with queried leaf-list instance do not populate payload.")
						 }
						return resFldValMap, err
					} else {
						resLst := processLfLstDbToYang(xpath, val, yngTerminalNdDtType)
						resFldValMap[xYangSpecMap[xpath].yangEntry.Name] = resLst
					}
				} else {
					if leafLstInstGetReq {
						log.Warningf("Queried leaf-list does not exist in DB, uri  - %v", requestUri)
						err = tlerr.NotFoundError{Format:"Resource not found"}
					}
				}
			} else {
				val, ok := (*dbDataMap)[cdb][tbl][tblKey].Field[dbFldName]
				if ok {
					resVal, _, err := DbToYangType(yngTerminalNdDtType, xpath, val)
					if err != nil {
						log.Warning("Conversion of Db value type to yang type for field didn't happen.", xpath)
					} else {
						resFldValMap[xYangSpecMap[xpath].yangEntry.Name] = resVal
					}
				} else {
					xfmrLogInfoAll("Field value does not exist in DB for - %v" , uri)
					err = tlerr.NotFoundError{Format:"Resource not found"}
				}
			}
		}
	}
	return resFldValMap, err
}

func yangDataFill(inParamsForGet xlateFromDbParams) error {
	var err error
	validate := inParamsForGet.validate
	isValid := validate
	dbs := inParamsForGet.dbs
	ygRoot := inParamsForGet.ygRoot
	uri := inParamsForGet.uri
	requestUri := inParamsForGet.requestUri
	dbDataMap := inParamsForGet.dbDataMap
	txCache := inParamsForGet.txCache
	resultMap := inParamsForGet.resultMap
	xpath := inParamsForGet.xpath
	var chldUri string

	yangNode, ok := xYangSpecMap[xpath]

	if ok  && yangNode.yangEntry != nil {
		for yangChldName := range yangNode.yangEntry.Dir {
			chldXpath := xpath+"/"+yangChldName
			if xYangSpecMap[chldXpath] != nil && xYangSpecMap[chldXpath].nameWithMod != nil {
				chldUri   = uri+"/"+ *(xYangSpecMap[chldXpath].nameWithMod)
			} else {
				chldUri   = uri+"/"+yangChldName
			}
			inParamsForGet.xpath = chldXpath
			inParamsForGet.uri = chldUri
			if xYangSpecMap[chldXpath] != nil && xYangSpecMap[chldXpath].yangEntry != nil {
				cdb := xYangSpecMap[chldXpath].dbIndex
				inParamsForGet.curDb = cdb
				if len(xYangSpecMap[chldXpath].validateFunc) > 0 && !validate {
					xpathKeyExtRet, _ := xpathKeyExtract(dbs[cdb], ygRoot, GET, chldUri, requestUri, nil, txCache, inParamsForGet.xfmrDbTblKeyCache)
					inParamsForGet.ygRoot = ygRoot
					// TODO - handle non CONFIG-DB
					inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, chldUri, requestUri, GET, xpathKeyExtRet.dbKey, dbDataMap, nil, nil, txCache)
					res := validateHandlerFunc(inParams)
					if !res {
						continue
					} else {
						isValid = res
					}
					inParamsForGet.validate = isValid
					inParamsForGet.dbDataMap = dbDataMap
					inParamsForGet.ygRoot = ygRoot
				}
				chldYangType := xYangSpecMap[chldXpath].yangDataType
				if  chldYangType == YANG_LEAF || chldYangType == YANG_LEAF_LIST {
					if len(xYangSpecMap[xpath].xfmrFunc) > 0 {
						continue
					}
					fldValMap, err := terminalNodeProcess(inParamsForGet)
					dbDataMap = inParamsForGet.dbDataMap
					ygRoot = inParamsForGet.ygRoot
					if err != nil {
						xfmrLogInfoAll("Failed to get data(\"%v\").", chldUri)
					}
					for lf, val := range fldValMap {
						resultMap[lf] = val
					}
					inParamsForGet.resultMap = resultMap
				} else if chldYangType == YANG_CONTAINER {
					xpathKeyExtRet, _ := xpathKeyExtract(dbs[cdb], ygRoot, GET, chldUri, requestUri, nil, txCache, inParamsForGet.xfmrDbTblKeyCache)
					tblKey := xpathKeyExtRet.dbKey
					chtbl := xpathKeyExtRet.tableName
					inParamsForGet.ygRoot = ygRoot

					if _, ok := (*dbDataMap)[cdb][chtbl]; !ok && len(chtbl) > 0 {
						childDBKey := ""
						terminalNodeGet  := false
						qdbMapHasTblData := false
						qdbMapHasTblKeyData := false
						if !xYangSpecMap[chldXpath].hasNonTerminalNode {
							childDBKey      = tblKey
							terminalNodeGet = true
						}
						if qdbMap, getOk := inParamsForGet.dbTblKeyGetCache[cdb]; getOk {
							if dbTblData, tblPresent := qdbMap[chtbl]; tblPresent {
								qdbMapHasTblData = true
								if _, keyPresent := dbTblData[tblKey]; keyPresent {
									qdbMapHasTblKeyData = true;
								}
							}
						}

						if !qdbMapHasTblData || (terminalNodeGet && qdbMapHasTblData && !qdbMapHasTblKeyData) {
						curDbDataMap, err := fillDbDataMapForTbl(chldUri, chldXpath, chtbl, childDBKey, cdb, dbs, inParamsForGet.dbTblKeyGetCache)
						if err == nil {
							mapCopy((*dbDataMap)[cdb], curDbDataMap[cdb])
							inParamsForGet.dbDataMap = dbDataMap
						}
					    }
					}
					cname := xYangSpecMap[chldXpath].yangEntry.Name
					if xYangSpecMap[chldXpath].xfmrTbl != nil {
						xfmrTblFunc := *xYangSpecMap[chldXpath].xfmrTbl
						if len(xfmrTblFunc) > 0 {
							inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, chldUri, requestUri, GET, tblKey, dbDataMap, nil, nil, txCache)
							tblList, _ := xfmrTblHandlerFunc(xfmrTblFunc, inParams, inParamsForGet.xfmrDbTblKeyCache)
							inParamsForGet.dbDataMap = dbDataMap
							inParamsForGet.ygRoot = ygRoot
							if len(tblList) > 1 {
								log.Warningf("Table transformer returned more than one table for container %v", chldXpath)
							}
							if len(tblList) == 0 {
								continue
							}
							dbDataFromTblXfmrGet(tblList[0], inParams, dbDataMap, inParamsForGet.dbTblKeyGetCache)
							inParamsForGet.dbDataMap = dbDataMap
							inParamsForGet.ygRoot = ygRoot
							chtbl = tblList[0]
						}
					}
					if len(xYangSpecMap[chldXpath].xfmrFunc) > 0 {
						if (len(xYangSpecMap[xpath].xfmrFunc) == 0) ||
						(len(xYangSpecMap[xpath].xfmrFunc) > 0   &&
						(xYangSpecMap[xpath].xfmrFunc != xYangSpecMap[chldXpath].xfmrFunc)) {
							inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, chldUri, requestUri, GET, "", dbDataMap, nil, nil, txCache)
							err := xfmrHandlerFunc(inParams)
							inParamsForGet.dbDataMap = dbDataMap
							inParamsForGet.ygRoot = ygRoot
							if err != nil {
								xfmrLogInfoAll("Error returned by %v: %v", xYangSpecMap[xpath].xfmrFunc, err)
							}
						}
						if !xYangSpecMap[chldXpath].hasChildSubTree {
							continue
						}
					}
					cmap2 := make(map[string]interface{})
					linParamsForGet := formXlateFromDbParams(dbs[cdb], dbs, cdb, ygRoot, chldUri, requestUri, chldXpath, inParamsForGet.oper, chtbl, tblKey, dbDataMap, inParamsForGet.txCache, cmap2, inParamsForGet.validate)
					linParamsForGet.xfmrDbTblKeyCache = inParamsForGet.xfmrDbTblKeyCache
					linParamsForGet.dbTblKeyGetCache = inParamsForGet.dbTblKeyGetCache
					err  = yangDataFill(linParamsForGet)
					cmap2 = linParamsForGet.resultMap
					dbDataMap = linParamsForGet.dbDataMap
					ygRoot = linParamsForGet.ygRoot
					if err != nil && len(cmap2) == 0 {
						xfmrLogInfoAll("Empty container.(\"%v\").\r\n", chldUri)
					} else {
						if len(cmap2) > 0 {
							resultMap[cname] = cmap2
						}
						inParamsForGet.resultMap = resultMap
					}
					inParamsForGet.dbDataMap = dbDataMap
					inParamsForGet.ygRoot = ygRoot
				} else if chldYangType ==  YANG_LIST {
					xpathKeyExtRet, _ := xpathKeyExtract(dbs[cdb], ygRoot, GET, chldUri, requestUri, nil, txCache, inParamsForGet.xfmrDbTblKeyCache)
					inParamsForGet.ygRoot = ygRoot
					cdb = xYangSpecMap[chldXpath].dbIndex
					inParamsForGet.curDb = cdb
					if len(xYangSpecMap[chldXpath].xfmrFunc) > 0 {
						if (len(xYangSpecMap[xpath].xfmrFunc) == 0) ||
						(len(xYangSpecMap[xpath].xfmrFunc) > 0   &&
						(xYangSpecMap[xpath].xfmrFunc != xYangSpecMap[chldXpath].xfmrFunc)) {
							inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, chldUri, requestUri, GET, "", dbDataMap, nil, nil, txCache)
							err := xfmrHandlerFunc(inParams)
							if err != nil {
								xfmrLogInfoAll("Error returned by %v: %v", xYangSpecMap[chldXpath].xfmrFunc, err)
							}
							inParamsForGet.dbDataMap = dbDataMap
							inParamsForGet.ygRoot = ygRoot
						}
						if !xYangSpecMap[chldXpath].hasChildSubTree {
							continue
						}
					}
					ynode, ok := xYangSpecMap[chldXpath]
					lTblName := ""
					if ok && ynode.tableName != nil {
						lTblName = *ynode.tableName
					}
					if _, ok := (*dbDataMap)[cdb][lTblName]; !ok && len(lTblName) > 0 {
						curDbDataMap, err := fillDbDataMapForTbl(chldUri, chldXpath, lTblName, "", cdb, dbs, inParamsForGet.dbTblKeyGetCache)
						if err == nil {
							mapCopy((*dbDataMap)[cdb], curDbDataMap[cdb])
							inParamsForGet.dbDataMap = dbDataMap
						}
					}
					linParamsForGet := formXlateFromDbParams(dbs[cdb], dbs, cdb, ygRoot, chldUri, requestUri, chldXpath, inParamsForGet.oper, lTblName, xpathKeyExtRet.dbKey, dbDataMap, inParamsForGet.txCache, resultMap, inParamsForGet.validate)
					linParamsForGet.xfmrDbTblKeyCache = inParamsForGet.xfmrDbTblKeyCache
					linParamsForGet.dbTblKeyGetCache = inParamsForGet.dbTblKeyGetCache
					yangListDataFill(linParamsForGet, false)
					resultMap = linParamsForGet.resultMap
					dbDataMap = linParamsForGet.dbDataMap
					ygRoot = linParamsForGet.ygRoot
					inParamsForGet.dbDataMap = dbDataMap
					inParamsForGet.resultMap = resultMap
					inParamsForGet.ygRoot = ygRoot

				} else if chldYangType == "choice" || chldYangType == "case" {
					yangDataFill(inParamsForGet)
					resultMap = inParamsForGet.resultMap
					dbDataMap = inParamsForGet.dbDataMap
				} else {
					return err
				}
			}
		}
	}
	return err
}

/* Traverse linear db-map data and add to nested json data */
func dbDataToYangJsonCreate(inParamsForGet xlateFromDbParams) (string, bool, error) {
	var err error
	var fldSbtErr error // used only when direct query on leaf/leaf-list having subtree
	var fldErr error //used only when direct query on leaf/leaf-list having field transformer
	jsonData := "{}"
	resultMap := make(map[string]interface{})
        d := inParamsForGet.d
        dbs := inParamsForGet.dbs
        ygRoot := inParamsForGet.ygRoot
        uri := inParamsForGet.uri
        requestUri := inParamsForGet.requestUri
        dbDataMap := inParamsForGet.dbDataMap
        txCache := inParamsForGet.txCache
	cdb := inParamsForGet.curDb
	inParamsForGet.resultMap = resultMap

	if isSonicYang(uri) {
		return directDbToYangJsonCreate(inParamsForGet)
	} else {
		xpathKeyExtRet, _ := xpathKeyExtract(d, ygRoot, GET, uri, requestUri, nil, txCache, inParamsForGet.xfmrDbTblKeyCache)

		inParamsForGet.xpath = xpathKeyExtRet.xpath
		inParamsForGet.tbl = xpathKeyExtRet.tableName
		inParamsForGet.tblKey = xpathKeyExtRet.dbKey
		inParamsForGet.ygRoot = ygRoot
		yangNode, ok := xYangSpecMap[xpathKeyExtRet.xpath]
		if ok {
			/* Invoke pre-xfmr is present for the yang module */
			moduleName := "/" + strings.Split(uri, "/")[1]
			xfmrLogInfo("Module name for uri %s is %s", uri, moduleName)
			if modSpecInfo, specOk := xYangSpecMap[moduleName]; specOk && (len(modSpecInfo.xfmrPre) > 0) {
				inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, "", dbDataMap, nil, nil, txCache)
				err = preXfmrHandlerFunc(modSpecInfo.xfmrPre, inParams)
				xfmrLogInfo("Invoked pre transformer: %v, dbDataMap: %v ", modSpecInfo.xfmrPre, dbDataMap)
				if err != nil {
					log.Warningf("Pre-transformer: %v failed.(err:%v)", modSpecInfo.xfmrPre, err)
					return jsonData, true, err
				}
				inParamsForGet.dbDataMap = dbDataMap
				inParamsForGet.ygRoot    = ygRoot
			}

			yangType := yangTypeGet(yangNode.yangEntry)
			validateHandlerFlag := false
			tableXfmrFlag := false
			IsValidate := false
			if len(xYangSpecMap[xpathKeyExtRet.xpath].validateFunc) > 0 {
				inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, xpathKeyExtRet.dbKey, dbDataMap, nil, nil, txCache)
				res := validateHandlerFunc(inParams)
				inParamsForGet.dbDataMap = dbDataMap
				inParamsForGet.ygRoot = ygRoot
				if !res {
					validateHandlerFlag = true
					/* cannot immediately return from here since reXpath yangtype decides the return type */
				} else {
					IsValidate = res
				}
			}
			inParamsForGet.validate = IsValidate
			isList := false
			switch yangType {
			case YANG_LIST:
				isList = true
			case YANG_LEAF, YANG_LEAF_LIST, YANG_CONTAINER:
				isList = false
			default:
				xfmrLogInfo("Unknown yang object type for path %v", xpathKeyExtRet.xpath)
				isList = true //do not want non-list processing to happen
			}
			/*If yangtype is a list separate code path is to be taken in case of table transformer
			since that code path already handles the calling of table transformer and subsequent processing
			*/
			if (!validateHandlerFlag) && (!isList) {
				if xYangSpecMap[xpathKeyExtRet.xpath].xfmrTbl != nil {
					xfmrTblFunc := *xYangSpecMap[xpathKeyExtRet.xpath].xfmrTbl
					if len(xfmrTblFunc) > 0 {
						inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, xpathKeyExtRet.dbKey, dbDataMap, nil, nil, txCache)
						tblList, _ := xfmrTblHandlerFunc(xfmrTblFunc, inParams, inParamsForGet.xfmrDbTblKeyCache)
						inParamsForGet.dbDataMap = dbDataMap
						inParamsForGet.ygRoot = ygRoot
						if len(tblList) > 1 {
							log.Warningf("Table transformer returned more than one table for container %v", xpathKeyExtRet.xpath)
						}
						if len(tblList) == 0 {
							log.Warningf("Table transformer returned no table for conatiner %v", xpathKeyExtRet.xpath)
							tableXfmrFlag = true
						}
						if !tableXfmrFlag {
                                                      for _, tbl := range tblList {
                                                               dbDataFromTblXfmrGet(tbl, inParams, dbDataMap, inParamsForGet.dbTblKeyGetCache)
							       inParamsForGet.dbDataMap = dbDataMap
							       inParamsForGet.ygRoot = ygRoot
                                                      }

						}
					} else {
						log.Warningf("empty table transformer function name for xpath - %v", xpathKeyExtRet.xpath)
						tableXfmrFlag = true
					}
				}
			}

			for {
				done := true
				if yangType ==  YANG_LEAF || yangType == YANG_LEAF_LIST {
					yangName := xYangSpecMap[xpathKeyExtRet.xpath].yangEntry.Name
					if validateHandlerFlag || tableXfmrFlag {
						resultMap[yangName] = ""
						break
					}
					if len(xYangSpecMap[xpathKeyExtRet.xpath].xfmrFunc) > 0 {
						inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, "", dbDataMap, nil, nil, txCache)
						fldSbtErr = xfmrHandlerFunc(inParams)
						if fldSbtErr != nil {
							/*For request Uri pointing to leaf/leaf-list having subtree, error will be propagated
							  to handle check of leaf/leaf-list-instance existence in Db , which will be performed 
							  by subtree
							 */
							xfmrLogInfo("Error returned by %v: %v", xYangSpecMap[xpathKeyExtRet.xpath].xfmrFunc, err)
							inParamsForGet.ygRoot = ygRoot
							break
						}
						inParamsForGet.dbDataMap = dbDataMap
						inParamsForGet.ygRoot = ygRoot
					} else {
						tbl, key, _ := tableNameAndKeyFromDbMapGet((*dbDataMap)[cdb])
						inParamsForGet.tbl = tbl
						inParamsForGet.tblKey = key
						var fldValMap map[string]interface{}
						fldValMap, fldErr = terminalNodeProcess(inParamsForGet)
						if ((fldErr != nil) || (len(fldValMap) == 0)) {
							if fldErr == nil {
								if yangType == YANG_LEAF {
									xfmrLogInfo("Empty terminal node (\"%v\").", uri)
									fldErr = tlerr.NotFoundError{Format:"Resource Not found"}
								} else if ((yangType == YANG_LEAF_LIST) && ((strings.HasSuffix(uri, "]")) || (strings.HasSuffix(uri, "]/")))) {
									jsonMapData, _ := json.Marshal(resultMap)
									jsonData        = fmt.Sprintf("%v", string(jsonMapData))
									return jsonData, false, nil
								}
							}
						}
						resultMap = fldValMap
					}
					break

				} else if yangType == YANG_CONTAINER {
					cmap  := make(map[string]interface{})
					resultMap = cmap
					if validateHandlerFlag || tableXfmrFlag {
						break
					}
					if len(xYangSpecMap[xpathKeyExtRet.xpath].xfmrFunc) > 0 {
						inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, "", dbDataMap, nil, nil, txCache)
						err := xfmrHandlerFunc(inParams)
						if err != nil {
							xfmrLogInfo("Error returned by %v: %v", xYangSpecMap[xpathKeyExtRet.xpath].xfmrFunc, err)
							return jsonData, true, err
						}
						inParamsForGet.dbDataMap = dbDataMap
						inParamsForGet.ygRoot = ygRoot
						if !xYangSpecMap[xpathKeyExtRet.xpath].hasChildSubTree {
							break
						}
					}
					inParamsForGet.resultMap = make(map[string]interface{})
					err = yangDataFill(inParamsForGet)
					if err != nil {
						xfmrLogInfo("Empty container(\"%v\").\r\n", uri)
					}
					resultMap = inParamsForGet.resultMap
					break
				} else if yangType == YANG_LIST {
					isFirstCall := true
					if len(xYangSpecMap[xpathKeyExtRet.xpath].xfmrFunc) > 0 {
						inParams := formXfmrInputRequest(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, GET, "", dbDataMap, nil, nil, txCache)
						err := xfmrHandlerFunc(inParams)
						if err != nil {
							if (((strings.HasSuffix(uri, "]")) || (strings.HasSuffix(uri, "]/"))) && (uri == requestUri)) {
								// The error handling here is for the deferred resource check error being handled by the subtree for virtual table cases.
								log.Warningf("Subtree at list instance level returns error %v for  uri  - %v", err, uri)
								return jsonData, true, err

							} else {

								xfmrLogInfo("Error returned by %v: %v", xYangSpecMap[xpathKeyExtRet.xpath].xfmrFunc, err)
							}
						}
						isFirstCall = false
						inParamsForGet.dbDataMap = dbDataMap
						inParamsForGet.ygRoot = ygRoot
						if !xYangSpecMap[xpathKeyExtRet.xpath].hasChildSubTree {
							break
						}
					}
					inParamsForGet.resultMap = make(map[string]interface{})
					err = yangListDataFill(inParamsForGet, isFirstCall)
					if err != nil {
						xfmrLogInfo("yangListDataFill failed for list case(\"%v\").\r\n", uri)
					}
					resultMap = inParamsForGet.resultMap
					break
				} else {
					log.Warningf("Unknown yang object type for path %v", xpathKeyExtRet.xpath)
					break
				}
				if done {
					break
				}
			} //end of for
		}
	}

	jsonMapData, _ := json.Marshal(resultMap)
	isEmptyPayload := isJsonDataEmpty(string(jsonMapData))
	jsonData        = fmt.Sprintf("%v", string(jsonMapData))
	if fldSbtErr != nil {
		/*error should be propagated only when request Uri points to leaf/leaf-list-instance having subtree,
		  This is to handle check of leaf/leaf-list-instance existence in Db , which will be performed 
                  by subtree, and depending whether queried node exists or not subtree should return error
		*/
		return jsonData, isEmptyPayload, fldSbtErr
	}
	if fldErr != nil {
		/* error should be propagated only when request Uri points to leaf/leaf-list-instance and the data 
		   is not available(via field-xfmr or field name)
		*/
		return jsonData, isEmptyPayload, fldErr
	}

	return jsonData, isEmptyPayload, nil
}
