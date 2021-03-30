package com.ug.datacube.common.factories.datasource;

import com.google.protobuf.Message;
import com.typesafe.config.Config;
import com.ug.datacube.common.vars.ConfigVars;
import com.ug.datacube.common.vars.KafkaVars;
import com.ug.datacube.proto.tool.ProtoMessageSchema;
import org.apache.commons.lang3.ClassUtils;
import org.apache.flink.streaming.api.functions.source.RichParallelSourceFunction;

import java.util.Properties;

public class KafkaStreamClient extends StreamClient {
    public KafkaStreamClient(Config config) {
        topic = config.getString(ConfigVars.TOPIC);
        pbClassName = config.getString(ConfigVars.PBCLASSNAME);
        // todo 这里现在只设置了基础的两个配置，还需要调研之后添加
        properties = new Properties();
        properties.setProperty(
                KafkaVars.BOOTSTRAPSERVICES, config.getString(ConfigVars.BOOTSTRAPSERVICES));
        properties.setProperty(KafkaVars.GROUPID, config.getString(ConfigVars.GROUPID));
    }

    /**
     * 这个方法负责返回一个kafka的数据源
     *
     * @return
     */
    @Override
    public RichParallelSourceFunction<Message> getSourceFunction() {
        try {
            Class pbClass = ClassUtils.getClass(pbClassName);
            ProtoMessageSchema schema = new ProtoMessageSchema(pbClass);
            // todo 弄清楚这里的 setStartFromGroupOffsets 是否可能会产生不符合预期的效果
            return new UgFlinkKafkaConsumer<>(topic, schema, properties).setStartFromGroupOffsets();
        } catch (ClassNotFoundException e) {
            e.printStackTrace();
        }
        return null;
    }
}
