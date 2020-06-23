# modbustodb
a modbus server that can collect data from device to insert into influxdb
can load modbus device config to memory and collect data from devices via modbus protocol automatically

# api
config file is cfg/config.json 
1.:8081/data   show the datapoint in channel-device-datapoint configure hierarchy.
2.:8081/reload reload configure

