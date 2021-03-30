package com.ug.datacube.common.parser;

import com.google.protobuf.Message;
import com.ug.datacube.common.flatten.FlattenedRTAPageviewStrategy;
import com.ug.datacube.proto.pageview.PageviewService;
import com.ug.datacube.proto.rta.RTAStrategyService;

import java.util.ArrayList;
import java.util.List;

public class RTAPageviewStrategyParser extends Parser<FlattenedRTAPageviewStrategy> {
    @Override
    public List<FlattenedRTAPageviewStrategy> parse(Message message) {
        PageviewService.Pageview pageview = (PageviewService.Pageview) message;
        List<FlattenedRTAPageviewStrategy> result = new ArrayList<>();

        String traceId = pageview.getTraceId();
        long processTime = pageview.getProcessTime();

        List<RTAStrategyService.RTAStrategy> hitStrategyList = pageview.getHitStrategyList();
        for (RTAStrategyService.RTAStrategy rtaStrategy : hitStrategyList) {
            result.add(FlattenedRTAPageviewStrategy.builder()
                    .traceId(traceId)
                    .processTime(processTime)
                    .rtaStrategyId(rtaStrategy.getId())
                    .rtaValidPlatform(rtaStrategy.getPlatform())
                    .strategyIds(rtaStrategy.getStrategyIdList())
                    .hitStrategyCount(1)
                    .build());
        }
        return result;
    }
}
