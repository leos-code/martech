package com.ug.datacube.common.parser;

import com.google.protobuf.Message;
import com.ug.datacube.common.flatten.FlattenedRTAPageviewRequest;
import com.ug.datacube.proto.pageview.PageviewService;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

/**
 * 这个类是对 rtaPageview 的request进行打平的parser
 */
public class RTAPageviewRequestParser extends Parser<FlattenedRTAPageviewRequest> implements Serializable {
    @Override
    public List<FlattenedRTAPageviewRequest> parse(Message message) {
        PageviewService.Pageview pageview = (PageviewService.Pageview) message;
        List<FlattenedRTAPageviewRequest> result = new ArrayList<>();

        String traceId = pageview.getTraceId();
        long processTime = pageview.getProcessTime();
        String fromPlatform = pageview.getFromPlatform();
        PageviewService.Device device = pageview.getDevice();
        String imei = device.getImei();
        String idfa = device.getIdfa();
        String oaid = device.getOaid();
        int innerCostMs = pageview.getInnerCostMs();

        // 这是作为请求的数据，过程中忽略其他的需要打平的数据
        result.add(FlattenedRTAPageviewRequest.builder()
                .traceId(traceId)
                .processTime(processTime)
                .fromPlatform(fromPlatform)
                .imei(imei)
                .idfa(idfa)
                .oaid(oaid)
                .innerCostMs((long) innerCostMs)
                .requestCount(1)
                .build()
        );
        return result;
    }
}
