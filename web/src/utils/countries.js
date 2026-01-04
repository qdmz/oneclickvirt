// 国家/地区数据
export const countries = [
  // 亚洲 - 东亚
  { code: 'CN', name: '中国', flag: '🇨🇳', region: '东亚' },
  { code: 'HK', name: '中国香港', flag: '🇭🇰', region: '东亚' },
  { code: 'TW', name: '中国台湾', flag: '🇨🇳', region: '东亚' },
  { code: 'MO', name: '中国澳门', flag: '🇲🇴', region: '东亚' },
  { code: 'JP', name: '日本', flag: '🇯🇵', region: '东亚' },
  { code: 'KR', name: '韩国', flag: '🇰🇷', region: '东亚' },
  { code: 'KP', name: '朝鲜', flag: '🇰🇵', region: '东亚' },
  { code: 'MN', name: '蒙古', flag: '🇲🇳', region: '东亚' },
  
  // 亚洲 - 东南亚
  { code: 'SG', name: '新加坡', flag: '🇸🇬', region: '东南亚' },
  { code: 'MY', name: '马来西亚', flag: '🇲🇾', region: '东南亚' },
  { code: 'TH', name: '泰国', flag: '🇹🇭', region: '东南亚' },
  { code: 'VN', name: '越南', flag: '🇻🇳', region: '东南亚' },
  { code: 'ID', name: '印度尼西亚', flag: '🇮🇩', region: '东南亚' },
  { code: 'PH', name: '菲律宾', flag: '🇵🇭', region: '东南亚' },
  { code: 'BN', name: '文莱', flag: '🇧🇳', region: '东南亚' },
  { code: 'LA', name: '老挝', flag: '🇱🇦', region: '东南亚' },
  { code: 'KH', name: '柬埔寨', flag: '🇰🇭', region: '东南亚' },
  { code: 'MM', name: '缅甸', flag: '🇲🇲', region: '东南亚' },
  { code: 'TL', name: '东帝汶', flag: '🇹🇱', region: '东南亚' },
  
  // 亚洲 - 南亚
  { code: 'IN', name: '印度', flag: '🇮🇳', region: '南亚' },
  { code: 'PK', name: '巴基斯坦', flag: '🇵🇰', region: '南亚' },
  { code: 'BD', name: '孟加拉国', flag: '🇧🇩', region: '南亚' },
  { code: 'LK', name: '斯里兰卡', flag: '🇱🇰', region: '南亚' },
  { code: 'NP', name: '尼泊尔', flag: '🇳🇵', region: '南亚' },
  { code: 'BT', name: '不丹', flag: '🇧🇹', region: '南亚' },
  { code: 'MV', name: '马尔代夫', flag: '🇲🇻', region: '南亚' },
  { code: 'AF', name: '阿富汗', flag: '🇦🇫', region: '南亚' },
  
  // 亚洲 - 中亚
  { code: 'KZ', name: '哈萨克斯坦', flag: '🇰🇿', region: '中亚' },
  { code: 'UZ', name: '乌兹别克斯坦', flag: '🇺🇿', region: '中亚' },
  { code: 'TJ', name: '塔吉克斯坦', flag: '🇹🇯', region: '中亚' },
  { code: 'KG', name: '吉尔吉斯斯坦', flag: '🇰🇬', region: '中亚' },
  { code: 'TM', name: '土库曼斯坦', flag: '🇹🇲', region: '中亚' },
  
  // 亚洲 - 西亚/中东
  { code: 'TR', name: '土耳其', flag: '🇹🇷', region: '西亚/中东' },
  { code: 'IR', name: '伊朗', flag: '🇮🇷', region: '西亚/中东' },
  { code: 'IQ', name: '伊拉克', flag: '🇮🇶', region: '西亚/中东' },
  { code: 'SY', name: '叙利亚', flag: '🇸🇾', region: '西亚/中东' },
  { code: 'JO', name: '约旦', flag: '🇯🇴', region: '西亚/中东' },
  { code: 'LB', name: '黎巴嫩', flag: '🇱🇧', region: '西亚/中东' },
  { code: 'IL', name: '以色列', flag: '🇮🇱', region: '西亚/中东' },
  { code: 'SA', name: '沙特阿拉伯', flag: '🇸🇦', region: '西亚/中东' },
  { code: 'AE', name: '阿联酋', flag: '🇦🇪', region: '西亚/中东' },
  { code: 'QA', name: '卡塔尔', flag: '🇶🇦', region: '西亚/中东' },
  { code: 'KW', name: '科威特', flag: '🇰🇼', region: '西亚/中东' },
  { code: 'BH', name: '巴林', flag: '🇧🇭', region: '西亚/中东' },
  { code: 'OM', name: '阿曼', flag: '🇴🇲', region: '西亚/中东' },
  { code: 'YE', name: '也门', flag: '🇾🇪', region: '西亚/中东' },
  { code: 'GE', name: '格鲁吉亚', flag: '🇬🇪', region: '西亚/中东' },
  { code: 'AM', name: '亚美尼亚', flag: '🇦🇲', region: '西亚/中东' },
  { code: 'AZ', name: '阿塞拜疆', flag: '🇦🇿', region: '西亚/中东' },
  { code: 'CY', name: '塞浦路斯', flag: '🇨🇾', region: '西亚/中东' },
  
  // 欧洲 - 西欧
  { code: 'GB', name: '英国', flag: '🇬🇧', region: '西欧' },
  { code: 'IE', name: '爱尔兰', flag: '🇮🇪', region: '西欧' },
  { code: 'FR', name: '法国', flag: '🇫🇷', region: '西欧' },
  { code: 'DE', name: '德国', flag: '🇩🇪', region: '西欧' },
  { code: 'NL', name: '荷兰', flag: '🇳🇱', region: '西欧' },
  { code: 'BE', name: '比利时', flag: '🇧🇪', region: '西欧' },
  { code: 'LU', name: '卢森堡', flag: '🇱🇺', region: '西欧' },
  { code: 'CH', name: '瑞士', flag: '🇨🇭', region: '西欧' },
  { code: 'AT', name: '奥地利', flag: '🇦🇹', region: '西欧' },
  { code: 'LI', name: '列支敦士登', flag: '🇱🇮', region: '西欧' },
  { code: 'MC', name: '摩纳哥', flag: '🇲🇨', region: '西欧' },
  
  // 欧洲 - 北欧
  { code: 'SE', name: '瑞典', flag: '🇸🇪', region: '北欧' },
  { code: 'NO', name: '挪威', flag: '🇳🇴', region: '北欧' },
  { code: 'FI', name: '芬兰', flag: '🇫🇮', region: '北欧' },
  { code: 'DK', name: '丹麦', flag: '🇩🇰', region: '北欧' },
  { code: 'IS', name: '冰岛', flag: '🇮🇸', region: '北欧' },
  
  // 欧洲 - 南欧
  { code: 'IT', name: '意大利', flag: '🇮🇹', region: '南欧' },
  { code: 'ES', name: '西班牙', flag: '🇪🇸', region: '南欧' },
  { code: 'PT', name: '葡萄牙', flag: '🇵🇹', region: '南欧' },
  { code: 'GR', name: '希腊', flag: '🇬🇷', region: '南欧' },
  { code: 'MT', name: '马耳他', flag: '🇲🇹', region: '南欧' },
  { code: 'SM', name: '圣马力诺', flag: '🇸🇲', region: '南欧' },
  { code: 'VA', name: '梵蒂冈', flag: '🇻🇦', region: '南欧' },
  { code: 'AD', name: '安道尔', flag: '🇦🇩', region: '南欧' },
  { code: 'AL', name: '阿尔巴尼亚', flag: '🇦🇱', region: '南欧' },
  { code: 'RS', name: '塞尔维亚', flag: '🇷🇸', region: '南欧' },
  { code: 'ME', name: '黑山', flag: '🇲🇪', region: '南欧' },
  { code: 'BA', name: '波黑', flag: '🇧🇦', region: '南欧' },
  { code: 'HR', name: '克罗地亚', flag: '🇭🇷', region: '南欧' },
  { code: 'SI', name: '斯洛文尼亚', flag: '🇸🇮', region: '南欧' },
  { code: 'MK', name: '北马其顿', flag: '🇲🇰', region: '南欧' },
  
  // 欧洲 - 中欧
  { code: 'PL', name: '波兰', flag: '🇵🇱', region: '中欧' },
  { code: 'CZ', name: '捷克', flag: '🇨🇿', region: '中欧' },
  { code: 'SK', name: '斯洛伐克', flag: '🇸🇰', region: '中欧' },
  { code: 'HU', name: '匈牙利', flag: '🇭🇺', region: '中欧' },
  
  // 欧洲 - 东欧
  { code: 'RU', name: '俄罗斯', flag: '🇷🇺', region: '东欧' },
  { code: 'UA', name: '乌克兰', flag: '🇺🇦', region: '东欧' },
  { code: 'BY', name: '白俄罗斯', flag: '🇧🇾', region: '东欧' },
  { code: 'MD', name: '摩尔多瓦', flag: '🇲🇩', region: '东欧' },
  { code: 'RO', name: '罗马尼亚', flag: '🇷🇴', region: '东欧' },
  { code: 'BG', name: '保加利亚', flag: '🇧🇬', region: '东欧' },
  { code: 'EE', name: '爱沙尼亚', flag: '🇪🇪', region: '东欧' },
  { code: 'LV', name: '拉脱维亚', flag: '🇱🇻', region: '东欧' },
  { code: 'LT', name: '立陶宛', flag: '🇱🇹', region: '东欧' },
  
  // 北美洲
  { code: 'US', name: '美国', flag: '🇺🇸', region: '北美洲' },
  { code: 'CA', name: '加拿大', flag: '🇨🇦', region: '北美洲' },
  { code: 'MX', name: '墨西哥', flag: '🇲🇽', region: '北美洲' },
  { code: 'GT', name: '危地马拉', flag: '🇬🇹', region: '北美洲' },
  { code: 'BZ', name: '伯利兹', flag: '🇧🇿', region: '北美洲' },
  { code: 'SV', name: '萨尔瓦多', flag: '🇸🇻', region: '北美洲' },
  { code: 'HN', name: '洪都拉斯', flag: '🇭🇳', region: '北美洲' },
  { code: 'NI', name: '尼加拉瓜', flag: '🇳🇮', region: '北美洲' },
  { code: 'CR', name: '哥斯达黎加', flag: '🇨🇷', region: '北美洲' },
  { code: 'PA', name: '巴拿马', flag: '🇵🇦', region: '北美洲' },
  { code: 'CU', name: '古巴', flag: '🇨🇺', region: '加勒比海' },
  { code: 'JM', name: '牙买加', flag: '🇯🇲', region: '加勒比海' },
  { code: 'HT', name: '海地', flag: '🇭🇹', region: '加勒比海' },
  { code: 'DO', name: '多米尼加', flag: '🇩🇴', region: '加勒比海' },
  { code: 'BS', name: '巴哈马', flag: '🇧🇸', region: '加勒比海' },
  { code: 'BB', name: '巴巴多斯', flag: '🇧🇧', region: '加勒比海' },
  { code: 'TT', name: '特立尼达和多巴哥', flag: '��', region: '加勒比海' },
  
  // 南美洲
  { code: 'BR', name: '巴西', flag: '🇧🇷', region: '南美洲' },
  { code: 'AR', name: '阿根廷', flag: '🇦🇷', region: '南美洲' },
  { code: 'CL', name: '智利', flag: '🇨🇱', region: '南美洲' },
  { code: 'CO', name: '哥伦比亚', flag: '🇨🇴', region: '南美洲' },
  { code: 'PE', name: '秘鲁', flag: '🇵🇪', region: '南美洲' },
  { code: 'VE', name: '委内瑞拉', flag: '🇻🇪', region: '南美洲' },
  { code: 'EC', name: '厄瓜多尔', flag: '🇪🇨', region: '南美洲' },
  { code: 'BO', name: '玻利维亚', flag: '🇧🇴', region: '南美洲' },
  { code: 'PY', name: '巴拉圭', flag: '🇵🇾', region: '南美洲' },
  { code: 'UY', name: '乌拉圭', flag: '🇺🇾', region: '南美洲' },
  { code: 'GY', name: '圭亚那', flag: '🇬🇾', region: '南美洲' },
  { code: 'SR', name: '苏里南', flag: '🇸🇷', region: '南美洲' },
  { code: 'GF', name: '法属圭亚那', flag: '🇬🇫', region: '南美洲' },
  
  // 非洲 - 北非
  { code: 'EG', name: '埃及', flag: '🇪🇬', region: '北非' },
  { code: 'LY', name: '利比亚', flag: '🇱🇾', region: '北非' },
  { code: 'TN', name: '突尼斯', flag: '🇹🇳', region: '北非' },
  { code: 'DZ', name: '阿尔及利亚', flag: '🇩🇿', region: '北非' },
  { code: 'MA', name: '摩洛哥', flag: '🇲🇦', region: '北非' },
  { code: 'SD', name: '苏丹', flag: '🇸🇩', region: '北非' },
  { code: 'SS', name: '南苏丹', flag: '🇸🇸', region: '北非' },
  
  // 非洲 - 西非
  { code: 'NG', name: '尼日利亚', flag: '🇳🇬', region: '西非' },
  { code: 'GH', name: '加纳', flag: '🇬🇭', region: '西非' },
  { code: 'CI', name: '科特迪瓦', flag: '🇨🇮', region: '西非' },
  { code: 'SN', name: '塞内加尔', flag: '🇸🇳', region: '西非' },
  { code: 'ML', name: '马里', flag: '🇲🇱', region: '西非' },
  { code: 'BF', name: '布基纳法索', flag: '🇧🇫', region: '西非' },
  { code: 'NE', name: '尼日尔', flag: '🇳🇪', region: '西非' },
  { code: 'GN', name: '几内亚', flag: '🇬🇳', region: '西非' },
  { code: 'SL', name: '塞拉利昂', flag: '🇸🇱', region: '西非' },
  { code: 'LR', name: '利比里亚', flag: '🇱🇷', region: '西非' },
  { code: 'TG', name: '多哥', flag: '🇹🇬', region: '西非' },
  { code: 'BJ', name: '贝宁', flag: '🇧🇯', region: '西非' },
  
  // 非洲 - 中非
  { code: 'TD', name: '乍得', flag: '🇹🇩', region: '中非' },
  { code: 'CF', name: '中非共和国', flag: '🇨🇫', region: '中非' },
  { code: 'CM', name: '喀麦隆', flag: '🇨🇲', region: '中非' },
  { code: 'GQ', name: '赤道几内亚', flag: '🇬🇶', region: '中非' },
  { code: 'GA', name: '加蓬', flag: '🇬🇦', region: '中非' },
  { code: 'CG', name: '刚果(布)', flag: '🇨🇬', region: '中非' },
  { code: 'CD', name: '刚果(金)', flag: '🇨🇩', region: '中非' },
  { code: 'AO', name: '安哥拉', flag: '🇦🇴', region: '中非' },
  
  // 非洲 - 东非
  { code: 'ET', name: '埃塞俄比亚', flag: '🇪🇹', region: '东非' },
  { code: 'KE', name: '肯尼亚', flag: '🇰🇪', region: '东非' },
  { code: 'UG', name: '乌干达', flag: '🇺🇬', region: '东非' },
  { code: 'TZ', name: '坦桑尼亚', flag: '🇹🇿', region: '东非' },
  { code: 'RW', name: '卢旺达', flag: '🇷🇼', region: '东非' },
  { code: 'BI', name: '布隆迪', flag: '🇧🇮', region: '东非' },
  { code: 'SO', name: '索马里', flag: '🇸🇴', region: '东非' },
  { code: 'DJ', name: '吉布提', flag: '🇩🇯', region: '东非' },
  { code: 'ER', name: '厄立特里亚', flag: '🇪🇷', region: '东非' },
  { code: 'MG', name: '马达加斯加', flag: '🇲🇬', region: '东非' },
  { code: 'MU', name: '毛里求斯', flag: '🇲🇺', region: '东非' },
  { code: 'SC', name: '塞舌尔', flag: '🇸🇨', region: '东非' },
  { code: 'KM', name: '科摩罗', flag: '🇰🇲', region: '东非' },
  
  // 非洲 - 南非
  { code: 'ZA', name: '南非', flag: '🇿🇦', region: '南部非洲' },
  { code: 'ZW', name: '津巴布韦', flag: '🇿🇼', region: '南部非洲' },
  { code: 'ZM', name: '赞比亚', flag: '🇿🇲', region: '南部非洲' },
  { code: 'MW', name: '马拉维', flag: '🇲🇼', region: '南部非洲' },
  { code: 'MZ', name: '莫桑比克', flag: '🇲🇿', region: '南部非洲' },
  { code: 'BW', name: '博茨瓦纳', flag: '🇧🇼', region: '南部非洲' },
  { code: 'NA', name: '纳米比亚', flag: '🇳🇦', region: '南部非洲' },
  { code: 'LS', name: '莱索托', flag: '🇱🇸', region: '南部非洲' },
  { code: 'SZ', name: '斯威士兰', flag: '🇸🇿', region: '南部非洲' },
  
  // 大洋洲
  { code: 'AU', name: '澳大利亚', flag: '🇦🇺', region: '大洋洲' },
  { code: 'NZ', name: '新西兰', flag: '🇳🇿', region: '大洋洲' },
  { code: 'FJ', name: '斐济', flag: '🇫🇯', region: '大洋洲' },
  { code: 'PG', name: '巴布亚新几内亚', flag: '🇵🇬', region: '大洋洲' },
  { code: 'SB', name: '所罗门群岛', flag: '🇸🇧', region: '大洋洲' },
  { code: 'VU', name: '瓦努阿图', flag: '🇻🇺', region: '大洋洲' },
  { code: 'NC', name: '新喀里多尼亚', flag: '🇳🇨', region: '大洋洲' },
  { code: 'PF', name: '法属波利尼西亚', flag: '🇵🇫', region: '大洋洲' },
  { code: 'WS', name: '萨摩亚', flag: '🇼🇸', region: '大洋洲' },
  { code: 'TO', name: '汤加', flag: '🇹🇴', region: '大洋洲' },
  { code: 'CK', name: '库克群岛', flag: '🇨🇰', region: '大洋洲' },
  { code: 'NU', name: '纽埃', flag: '🇳🇺', region: '大洋洲' },
  { code: 'PW', name: '帕劳', flag: '🇵🇼', region: '大洋洲' },
  { code: 'FM', name: '密克罗尼西亚', flag: '🇫🇲', region: '大洋洲' },
  { code: 'MH', name: '马绍尔群岛', flag: '🇲🇭', region: '大洋洲' },
  { code: 'KI', name: '基里巴斯', flag: '🇰🇮', region: '大洋洲' },
  { code: 'NR', name: '瑙鲁', flag: '🇳🇷', region: '大洋洲' },
  { code: 'TV', name: '图瓦卢', flag: '🇹🇻', region: '大洋洲' }
]

// 根据国家代码获取国旗
export function getFlagEmoji(countryCode) {
  const country = countries.find(c => c.code === countryCode)
  return country ? country.flag : '🌍'
}

// 根据国家代码获取国家信息
export function getCountryInfo(countryCode) {
  return countries.find(c => c.code === countryCode)
}

// 根据国家名称获取国家信息
export function getCountryByName(countryName) {
  return countries.find(c => c.name === countryName)
}

// 按地区分组
export function getCountriesByRegion() {
  const grouped = {}
  countries.forEach(country => {
    if (!grouped[country.region]) {
      grouped[country.region] = []
    }
    grouped[country.region].push(country)
  })
  return grouped
}
