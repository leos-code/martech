package com.ug.datacube.common.factories.datasource;

import com.google.protobuf.Message;
import com.ug.datacube.proto.tool.ProtoMessageSchema;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.streaming.connectors.kafka.FlinkKafkaConsumer;

import java.util.Properties;

public class UgFlinkKafkaConsumer<T extends Message> extends FlinkKafkaConsumer<T> {
    private final ProtoMessageSchema protoMessageSchema;

    public UgFlinkKafkaConsumer(String topic, DeserializationSchema<T> valueDeserializer, Properties props) {
        super(topic, valueDeserializer, props);
        this.protoMessageSchema = (ProtoMessageSchema) valueDeserializer;
    }

    @Override
    public void open(Configuration configuration) throws Exception {
        super.open(configuration);
        this.protoMessageSchema.initMetrics(getRuntimeContext());
        this.protoMessageSchema.initInvokeFunc();
    }
}
