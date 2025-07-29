package core

import (
	ipUtils "blogx_server/utils/ip"
	"fmt"
	"net" // 导入net包以使用IP地址解析功能
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

// searcher 是一个全局的 ip2region 搜索器。
var searcher *xdb.Searcher

// InitIPDB 用于初始化IP地址数据库。
// 这个函数只在程序启动时调用一次。
func InitIPDB() {
	// 定义ip2region数据库文件的路径。这个文件包含了IP地址和地理位置的映射关系。
	// 注意：新版本的ip2region使用 .xdb 格式的数据库文件。
	var dbPath = "init/ip2region.xdb"
	// 使用 xdb.NewWithFileOnly 方法从指定路径加载IP数据库文件。
	// 这个方法会将整个数据库文件一次性加载到内存中，以实现最快的查询速度。
	// 加载成功后，会返回一个可用于查询的 searcher 实例。
	_searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		// 如果在加载数据库文件时发生任何错误（如文件不存在、文件损坏等），
		// 则记录一条致命错误日志，并立即终止整个程序的运行。
		// 这表明IP地址查询功能被视为本程序的核心关键服务，不可或缺。
		logrus.Fatalf("init ip2region.db error: %v", err)
		return
	}
	// 将成功创建的搜索器实例赋值给全局变量 searcher，以便在其他函数中调用。
	searcher = _searcher
}

// GetIpAddr 根据给定的IP地址字符串，查询并返回其地理位置。
// 返回的地址会经过格式化，变得更易读。
func GetIpAddr(ip string) (addr string) {
	// 第一步：将IP地址字符串解析为 net.IP 类型，为后续处理做准备。
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		// 如果传入的字符串无法被解析成一个有效的IP地址，
		// 则记录警告并返回，避免后续逻辑出错。
		logrus.Warnf("无法解析的IP地址字符串: %s", ip)
		return "异常地址"
	}

	// 第二步：优化处理，判断是不是内网IP地址。
	// 调用工具函数 HasLocalIP 进行快速判断。
	if ipUtils.HasLocalIP(parsedIP) {
		// 如果是内网IP（如 192.168.x.x, 127.0.0.1 等），
		// 则没有必要进行数据库查询，直接返回"内网"，提高处理效率。
		return "内网"
	}

	// 第三步：对于公网IP，使用预先加载好的 searcher 工具进行数据库查询。
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		// 如果查询过程中发生错误（例如，ip2region库本身无法处理这个IP），
		// 则记录一条警告日志，并返回一个表示地址异常的字符串。
		logrus.Warnf("查询ip地址归属地失败 %s", err)
		return "异常地址"
	}

	// 第四步：解析和格式化查询结果。
	// ip2region返回的 `region` 字符串格式为：国家|区域|省份|城市|运营商
	// 例如: "中国|0|广东省|深圳市|电信" (其中'0'表示该字段数据为空)
	_addrList := strings.Split(region, "|")
	if len(_addrList) != 5 {
		// 正常情况下，返回的列表应该有5个元素。
		// 如果不是，说明可能遇到了一个异常的IP数据或查询结果。
		logrus.Warnf("异常的ip地址 %s", ip)
		return "未知地址"
	}

	// 为了方便阅读，我们只提取 国家、省份、城市 这三个关键部分。
	// 根据 ip2region 的数据格式：
	// _addrList[0] 是 国家
	// _addrList[2] 是 省份
	// _addrList[3] 是 城市
	country := _addrList[0]
	province := _addrList[2]
	city := _addrList[3]

	// 第五步：根据获取到的地址信息，智能地组合成最合适的显示格式。
	// 优先显示最详细的“省·市”信息。
	if province != "0" && city != "0" {
		return fmt.Sprintf("%s·%s", province, city)
	}
	// 如果没有城市信息，但有省份信息，则显示“国家·省份”。
	if country != "0" && province != "0" {
		return fmt.Sprintf("%s·%s", country, province)
	}
	// 如果连省份信息都没有，则只显示国家。
	if country != "0" {
		return country
	}
	// 如果以上所有关键字段都为空，则返回数据库查询到的原始字符串。
	return region
}
