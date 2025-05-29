package selector

/*
builder 接口面临哪些问题需要解决？
	1. 当 grpc 服务器只有一个的时候，如何处理？
	2. 当 grpc 服务器有多个的时候，如何处理？
	3. 当 grpc 服务器负载不同的时候如何处理？或者说：当多个连接都处于 ready 状态的时候，应该如何选择？
	4. 当 连接失败的时候，如何处理？是否启动重试等等

grpc 为了解决这些问题，把链路分为不同的阶段：
	1. balancer 构建的阶段
	2. 子链路具体的连接阶段：一个 grpc 服务器地址对应一个连接，多个地址的时候就会有多个子连接
	3. 子连接的选择问题（picker接口完成）
	4. balancer 状态
	5. 链路创建，删除，更新

负载均衡相关的接口：
	1. Builder 接口：用于构建一个 balancer 接口实例  // !!!Important
	2. SubConn 接口：主要负责具体的连接
	3. Picker 接口：主要负责从众多的链接里，按照负载均衡算法选择一个连接，供客户端使用 // !!!!!!Important 重中之重
	4. Balancer 接口：主要负责更新 clientConn 状态，更新 subConn 状态 // !!!Important
	5. ClientConn 接口：主要负责链路的维护，包括创建一个子链路，删除一个子链路，更新 ClientConn 状态
*/