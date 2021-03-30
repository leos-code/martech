package com.ug.datacube.common.parser;

import com.google.protobuf.Message;
import com.ug.datacube.common.flatten.FlattenedRecord;

import java.util.List;

public abstract class Parser<T extends FlattenedRecord> {

    public abstract List<T> parse(Message message);
}
