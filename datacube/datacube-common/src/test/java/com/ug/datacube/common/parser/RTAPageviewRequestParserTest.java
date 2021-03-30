package com.ug.datacube.common.parser;

import com.ug.datacube.proto.pageview.PageviewService;
import org.junit.Test;

public class RTAPageviewRequestParserTest {

    @Test
    public void test() {
        PageviewService.Pageview pageview = PageviewService.Pageview.getDefaultInstance();
        RTAPageviewRequestParser parser = new RTAPageviewRequestParser();
        System.out.println(parser.parse(pageview));
    }

}