package com.ug.datacube.common.parser;

import com.google.protobuf.Message;
import com.ug.datacube.common.flatten.FlattenedRTAPageviewImpression;
import com.ug.datacube.proto.pageview.PageviewService;

import java.util.ArrayList;
import java.util.List;

public class RTAPageviewImpressionParser extends Parser<FlattenedRTAPageviewImpression> {
    @Override
    public List<FlattenedRTAPageviewImpression> parse(Message message) {
        PageviewService.Pageview pageview = (PageviewService.Pageview) message;
        List<FlattenedRTAPageviewImpression> result = new ArrayList<>();

        String traceId = pageview.getTraceId();

        List<PageviewService.Imp> impList = pageview.getImpList();
        for (PageviewService.Imp imp : impList) {
            result.add(FlattenedRTAPageviewImpression.builder()
                    .traceId(traceId)
                    .impTime(0L)
                    .build());
        }
        return result;
    }
}
