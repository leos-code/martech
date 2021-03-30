package com.ug.datacube.common.parser;

import com.google.protobuf.Message;
import com.ug.datacube.common.flatten.FlattenedRTAPageviewConversion;
import com.ug.datacube.proto.pageview.PageviewService;

import java.util.ArrayList;
import java.util.List;

public class RTAPageviewConversionParser extends Parser<FlattenedRTAPageviewConversion> {
    @Override
    public List<FlattenedRTAPageviewConversion> parse(Message message) {
        PageviewService.Pageview pageview = (PageviewService.Pageview) message;
        List<FlattenedRTAPageviewConversion> result = new ArrayList<>();

        String traceId = pageview.getTraceId();
        for (PageviewService.Imp imp : pageview.getImpList()) {
            for (PageviewService.Click click : imp.getClickList()) {
                for (PageviewService.Conversion conversion : click.getConversionList()) {
                    result.add(FlattenedRTAPageviewConversion.builder()
                            .traceId(traceId)
                            .actionTime(conversion.getActionTime())
                            .actionType(conversion.getActionType())
                            .conversionCount(1)
                            .build());
                }
            }
        }
        return result;
    }
}
