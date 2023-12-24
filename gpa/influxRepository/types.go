package influxRepository

// Measurement 태그 지정용으로 사용.
// example
//
//		type Cpu struct{
//		    Measurement `influxdb:"measurement:cpu"`
//		    Usage       float `influxdb:"field:usage"`
//	     ...
//		}
type Measurement struct{}
