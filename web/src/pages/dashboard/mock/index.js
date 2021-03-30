const Mock = require('mockjs');

// UserInfo
// import { userInfo } from './data/user';
// Mock.mock(/dashboard\/user\/info/, 'get', userInfo);

// UserList
// import { userList } from './data/user';
// Mock.mock(/dashboard\/user\/list/, 'get', userList);

// RoleList
// import { roleList } from './data/role';
// Mock.mock(/dashboard\/role/, 'get', roleList);

// Authority
import { casbinMenuAdmin } from './data/casbin'
// Mock.mock(/dashboard\/user\/authority/, 'get', casbinMenuAdmin);
import { casbinEmptyData } from './data/casbin';
// Mock.mock(/dashboard\/user\/authority/, 'get', casbinEmptyData);


// Material
import {materialEdit, materialList, materialUpload, materialSubmitAudit} from "./data/material";
Mock.mock(/material\/list/, 'get', materialList);
Mock.mock(/material\/file\/upload/, 'post', materialUpload)
Mock.mock(/material\/edit/, 'post', materialEdit);
Mock.mock(/material\/audit\/submit/, 'post', materialSubmitAudit);

// ObjectList
import { objectList, AdvertiserObjectList } from './data/object';
// Mock.mock(/dashboard\/organization\/object/, 'get', objectList);
// Mock.mock(/dashboard\/organization\/object/, 'get', AdvertiserObjectList);

// Advertiser
import { AdvertiserList, AuthorityList } from './data/advertiser';
// Mock.mock(/advertiser\/authorize/, 'get', AuthorityList);
// Mock.mock(/advertiser/, 'get', AdvertiserList);
// Rta
import  { targetingList, targetingEdit} from "./data/targeting";

Mock.mock(/targeting\/list/, 'get', targetingList)
Mock.mock(/targeting\/edit/, 'post', targetingEdit)

// Schema
import  {targetingSchema} from "./data/schema";
Mock.mock(/schema/, 'get', targetingSchema)

