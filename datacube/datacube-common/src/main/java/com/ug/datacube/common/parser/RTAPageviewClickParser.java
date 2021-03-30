package com.ug.datacube.common.parser;

import com.google.protobuf.Message;
import com.ug.datacube.common.flatten.FlattenedRTAPageviewClick;
import com.ug.datacube.proto.pageview.PageviewService;

import java.util.ArrayList;
import java.util.List;

public class RTAPageviewClickParser extends Parser<FlattenedRTAPageviewClick> {
    @Override
    public List<FlattenedRTAPageviewClick> parse(Message message) {
        PageviewService.Pageview pageview = (PageviewService.Pageview) message;
        List<FlattenedRTAPageviewClick> result = new ArrayList<>();

        String traceId = pageview.getTraceId();
        for (PageviewService.Imp imp : pageview.getImpList()) {
            for (PageviewService.Click click : imp.getClickList()) {
                result.add(FlattenedRTAPageviewClick.builder()
                        .clickTime(click.getClickTime())
                        .clickCount(1)
                        .build());
            }
        }
        return null;
    }
}
