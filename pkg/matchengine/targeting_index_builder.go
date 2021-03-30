package matchengine

import (
	"crypto/md5"
	"sort"

	"github.com/tencentad/martech/api/proto/targeting"
	pb "github.com/golang/protobuf/proto"
)

type ConHitSlice []*targeting.TokenIndex_ConjunctionHit

// Len
func (h ConHitSlice) Len() int {
	return len(h)
}

// Less
func (h ConHitSlice) Less(i, j int) bool {
	if h[i].ConjunctionId == h[j].ConjunctionId {
		return h[i].PredicateId > h[j].PredicateId
	} else {
		return h[i].ConjunctionId < h[j].ConjunctionId
	}
}

// Swap
func (h ConHitSlice) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

type AidSlice []*targeting.IdIndex

// Len
func (h AidSlice) Len() int {
	return len(h)
}

// Less
func (h AidSlice) Less(i, j int) bool {
	return h[i].TargetingId < h[j].TargetingId
}

// Swap
func (h AidSlice) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func calcConjFinger(conj *targeting.TargetingDNF_Conjunction) string {
	data, _ := pb.Marshal(conj)
	res := md5.Sum(data)
	return string(res[:])
}

// Interval 区间
type Interval struct {
	Begin uint64
	End   uint64
}

type TokenType string

const (
	TokenTypeNumber = "number"
	TokenTypeString = "string"
)

// Token 包含属性名和对应的取值区间
type Token struct {
	Type     TokenType
	Field    string
	Interval Interval
	String   string
}

// TargetingIndexBuilder 定向索引构造器
type TargetingIndexBuilder struct {
	targetingIndex *targeting.TargetingIndex

	tokenVec      []*Token
	fieldTokenMap map[string]*tokenMap // field => tokenMap

	conjIndexVec    []*targeting.ConjunctionIndex
	newConjIndexVec []*targeting.ConjunctionIndex
	conjFingerMap   map[string]uint32

	tokenIndexVec []*targeting.TokenIndex
	idIndexVec    []*targeting.IdIndex

	//tokenIndex map[string]*targeting.IntervalCoverIndex
	searchTokenIndex map[string]*targeting.SearchTokenIndex

	targetingVec []*targeting.Targeting

	conjIdRemapping []uint32
}

type tokenMap struct {
	intervalMap map[Interval]uint32
	stringMap   map[string]uint32
}

func newTokenMap() *tokenMap {
	return &tokenMap{
		intervalMap: make(map[Interval]uint32),
		stringMap:   make(map[string]uint32),
	}
}

// NewTargetingIndexBuilder
func NewTargetingIndexBuilder(index *targeting.TargetingIndex) *TargetingIndexBuilder {

	return &TargetingIndexBuilder{
		targetingIndex:   index,
		fieldTokenMap:    make(map[string]*tokenMap),
		conjFingerMap:    make(map[string]uint32),
		searchTokenIndex: make(map[string]*targeting.SearchTokenIndex),
	}
}

// GetToken 构造索引过程中，获取对应的token
func (builder *TargetingIndexBuilder) GetToken(token *Token) (uint32, *targeting.TokenIndex) {

	var tokenMap *tokenMap
	var ok bool
	field := token.Field
	if tokenMap, ok = builder.fieldTokenMap[field]; !ok {
		tokenMap = newTokenMap()
		builder.fieldTokenMap[field] = tokenMap
	}

	var tokenIndex *targeting.TokenIndex
	var tokenId uint32
	if token.Type == TokenTypeNumber {
		tokenId, ok = tokenMap.intervalMap[token.Interval]
	} else if token.Type == TokenTypeString {
		tokenId, ok = tokenMap.stringMap[token.String]
	}

	if !ok {
		tokenId = uint32(len(builder.tokenVec))
		builder.tokenVec = append(builder.tokenVec, token)
		if token.Type == TokenTypeNumber {
			tokenMap.intervalMap[token.Interval] = tokenId
		} else if token.Type == TokenTypeString {
			tokenMap.stringMap[token.String] = tokenId
		}
		tokenIndex = &targeting.TokenIndex{}
		builder.tokenIndexVec = append(builder.tokenIndexVec, tokenIndex)
	} else {
		tokenIndex = builder.tokenIndexVec[tokenId]
	}
	return tokenId, tokenIndex
}

// Build 根据定向表达式构造索引
func (builder *TargetingIndexBuilder) Build(targetingVec []*targeting.Targeting) {
	builder.targetingVec = targetingVec
	for i, target := range targetingVec {
		localId := uint32(i)
		builder.ProcessBe(target, localId)
	}
	builder.IdRemapping()
	builder.ReArrangeData()
	BuildTokenCoverIndex(builder.fieldTokenMap, builder.searchTokenIndex)
	builder.SaveData()
}

// ProcessBe 处理布尔表达式，生成DNF
func (builder *TargetingIndexBuilder) ProcessBe(target *targeting.Targeting, localId uint32) {
	var dnf targeting.TargetingDNF
	ConvertToTargetingDNF(target, &dnf)
	tokenSet := make(map[uint32]bool)
	for _, conj := range dnf.Conjunction {
		finger := calcConjFinger(conj)
		conjId, conjExist := builder.conjFingerMap[finger]
		var conjIndex *targeting.ConjunctionIndex
		if conjExist {
			conjIndex = builder.conjIndexVec[conjId]
		} else {
			conjId = uint32(len(builder.conjIndexVec))
			conjIndex = &targeting.ConjunctionIndex{}
			builder.conjFingerMap[finger] = conjId
			builder.conjIndexVec = append(builder.conjIndexVec, conjIndex)
		}

		var preVec, notPreVec []*targeting.Predicate
		for _, predicate := range conj.Predicate {
			if predicate.Not {
				notPreVec = append(notPreVec, predicate)
			} else {
				preVec = append(preVec, predicate)
			}
		}
		builder.ProcessConjunction(localId, conjId, conjExist, preVec, notPreVec, tokenSet)
		conjIndex.HitCount = uint32(len(preVec))
		conjIndex.LocalIds = append(conjIndex.LocalIds, localId)
	}
	builder.idIndexVec = append(builder.idIndexVec, &targeting.IdIndex{TargetingId: target.Id, LocalId: localId})
}

// IdRemapping 按照顺序重新生成ID
func (builder *TargetingIndexBuilder) IdRemapping() {
	var maxHitCount uint32 = 0
	for _, conjIndex := range builder.conjIndexVec {
		if conjIndex.HitCount > maxHitCount {
			maxHitCount = conjIndex.HitCount
		}
	}
	hitList := make([][]uint32, maxHitCount+1)
	for i, conjIndex := range builder.conjIndexVec {
		hitCount := conjIndex.HitCount
		hitList[hitCount] = append(hitList[hitCount], uint32(i))
	}
	var newId uint32 = 0
	builder.conjIdRemapping = make([]uint32, len(builder.conjIndexVec))
	for _, slice := range hitList {
		for _, oldId := range slice {
			builder.conjIdRemapping[oldId] = newId
			newId++
		}
	}
}

// SaveData 保存结果
func (builder *TargetingIndexBuilder) SaveData() {
	builder.targetingIndex.IdIndex = builder.idIndexVec
	builder.targetingIndex.TokenIndex = builder.tokenIndexVec
	builder.targetingIndex.ConjunctionIndex = builder.newConjIndexVec
	builder.targetingIndex.SearchTokenIndex = builder.searchTokenIndex
	builder.targetingIndex.Targeting = builder.targetingVec
}

// ReArrangeData 对索引数据进行重新排序
func (builder *TargetingIndexBuilder) ReArrangeData() {
	// conjIndex 处理
	builder.newConjIndexVec = make([]*targeting.ConjunctionIndex, len(builder.conjIndexVec))
	for i, conjIndex := range builder.conjIndexVec {
		newId := builder.conjIdRemapping[i]
		builder.newConjIndexVec[newId] = conjIndex
	}
	// tokenIndex
	for _, tokenIndex := range builder.tokenIndexVec {
		for _, conjHit := range tokenIndex.ConjunctionHit {
			oldId := conjHit.ConjunctionId
			conjHit.ConjunctionId = builder.conjIdRemapping[oldId]
		}
		sort.Sort(ConHitSlice(tokenIndex.ConjunctionHit))
	}
	sort.Sort(AidSlice(builder.idIndexVec))
}

// ProcessConjunction 处理合取
func (builder *TargetingIndexBuilder) ProcessConjunction(localId, conjId uint32, conjExist bool, preVec,
	notPreVec []*targeting.Predicate, tokenSet map[uint32]bool) {

	for p, predicate := range preVec {
		pid := uint32(p)
		for _, value := range predicate.Value {
			token := convertPredicateValue2Token(predicate.Field, value)
			tokenId, tokenIndex := builder.GetToken(token)
			adHit := targeting.TokenIndex_TargetingHit{LocalId: localId, IsId: true}
			if !conjExist {
				conjunctionHit := targeting.TokenIndex_ConjunctionHit{ConjunctionId: conjId, PredicateId: pid}
				tokenIndex.ConjunctionHit = append(tokenIndex.ConjunctionHit, &conjunctionHit)
			}
			if _, ok := tokenSet[tokenId]; !ok {
				tokenSet[tokenId] = true
				tokenIndex.TargetingHit = append(tokenIndex.TargetingHit, &adHit)
			}
		}
	}
	if conjExist {
		return
	}
	for p, predicate := range notPreVec {
		pid := uint32(p + len(preVec))
		for _, value := range predicate.Value {
			token := convertPredicateValue2Token(predicate.Field, value)
			_, tokenIndex := builder.GetToken(token)
			conjunctionHit := targeting.TokenIndex_ConjunctionHit{ConjunctionId: conjId, PredicateId: pid}
			tokenIndex.ConjunctionHit = append(tokenIndex.ConjunctionHit, &conjunctionHit)
		}
	}
}

func convertPredicateValue2Token(field string, value *targeting.Predicate_Value) *Token {
	token := &Token{
		Field: field,
	}
	if value.Type == targeting.Predicate_Value_ID {
		token.Interval = Interval{Begin: value.Id, End: value.Id + 1}
		token.Type = TokenTypeNumber
	} else if value.Type == targeting.Predicate_Value_RANGE {
		token.Interval = Interval{Begin: value.Range.Begin, End: value.Range.End}
		token.Type = TokenTypeNumber
	} else if value.Type == targeting.Predicate_Value_String {
		token.String = value.Str
		token.Type = TokenTypeString
	}

	return token
}
