# onehost 线上回归测试清单（BUG_LIST）

> **用途**：替代 WorkBuddy「未发送」消息暂存队列，持久化记录测试发现的 bug，避免客户端刷新 / 清缓存导致内容丢失。
> **用法**：每次测试发现新 bug，直接告诉我「记一下：xxxx」，我会追加到第四节并标注状态；修复后更新状态。
> **最后更新**：2026-08-11

---

## 测试登录凭据
- 后台地址：`https://onec.ypvps.com/#/login`
- 用户名：`admin`
- 密码：（由 owner 保管，不入库）  ← **当前实测可用**（2026-08-11 调 `/api/v1/auth/login` 返回 200 + token）
- ⚠️ 注意：旧文档《线上全面测试报告》《线上复查记录》里写的 旧密码 **已过时**。密码曾被并发会话多次改动，登录后建议改回你自设的密码。

---

## 一、历史 5 个待回归 Bug（均在 2026-08-10 复查中修复，待你重新验证）

| # | 问题 | 现象 | 状态 |
|---|------|------|------|
| 1 | 登录加载失败 / 弹回登录页 | 登录成功却提示「加载失败」并弹回登录页 | ✅ 已修复 — 待回归 |
| 2 | 新增产品保存 400 | 新增产品报 400 参数错误 | ✅ 已修复 — 待回归 |
| 3 | 站点配置保存无反应 | 点保存没反应 | ✅ 已修复 — 待回归（v0.03 还修了 `Save()` 覆盖 + 首页动态化） |
| 4 | 工单分类英文 | 工单「分类」显示英文未汉化 | ✅ 已修复 — 待回归 |
| 5 | 易支付 | 配置保存后刷新丢失；用户充值/下单报 400 | ✅ 已修复 — 待回归（注意：低于最低充值额会被正常拒绝，非 bug） |

---

## 二、已修复参考（回归时对照用）
- `system_configs` 唯一索引缺陷 → SMTP 等配置「已填却报未配置」（v0.02）
- 首页动态化 / 推荐产品 / 站点配置 `Save()` 覆盖（v0.03，本地提交 `60b562c`）
- 登录 502：系部署重启窗口 / 并发会话改坏端口所致，**非代码缺陷**

---

## 三、当前已知风险
- 并发会话（IP `117.187.124.154`）曾反复改密码 / 端口 / 库导致 502，**建议关闭多余会话**，避免再次互相覆盖。
- GitHub 推送卡在 token 失效，本地最新 `v0.05`（`3c6a055`）未同步到 `qdmz/onehost`，但**生产已部署且健康**（v0.03/v0.04/v0.05 均已部署）。

---

## 四、新发现 Bug（实时追加）

| # | 问题 | 现象 | 状态 |
|---|------|------|------|
| 6 | 首页栏目开关失效（核心根因） | 站点配置里打开「显示平台 / 赞助方 / 推荐产品」等开关并保存，前端首页始终不显示对应栏目（最初「站点配置设置项没被前端动态显示」的真正原因） | ✅ 已修复（v0.05） |

> **根因**：`server/model/product/product.go` 中 `ShowPlatformsSection` / `ShowSponsorsSection` / `ShowRecommendedSection` 三个字段**未显式指定 gorm 列名**，gorm 默认按字段名生成列 `show_platforms_section` / `show_sponsors_section` / `show_recommended_section`；但迁移加的列、前端/接口用的键都是 `show_platforms` / `show_sponsors` / `show_recommended`。结果：**更新写进 `show_recommended` 列，读取却从 `show_recommended_section` 列取**，导致开关永远 False、怎么点都不生效。`show_nav`（字段名 `ShowNav`→列名 `show_nav`）与 `recommended_title`（→`recommended_title`）恰好一致，所以那两个能存。

## 五、历史回归修复记录
- **2026-08-11 站点配置保存 500（v0.03 引入的回归）**：后端 `UpdateSiteConfigFields` 用 JSON 字段名（`custom_footer`/`custom_header`/`show_yipay`）直接当 DB 列名做 UPDATE，而真实列名是 `footer_html`/`header_html`/`show_yi_pay`，导致 `Unknown column` 500。修复：在 `server/service/product/site.go` 增加 `jsonKeyToColumn` 映射，更新前转成真实列名。线上已部署，v0.04。
- **2026-08-11 首页栏目开关失效（模型列名不匹配）**：如上「Bug6」。修复：给三个字段加 `gorm:"column:show_platforms;..."` 等显式列名。线上已部署并实测开关可落库，`deploy/verify.py` 验证 **12/12 全 PASS**。版本 v0.05（本地已提交 `3c6a055`，GitHub 推送待新 token）。
- **2026-08-11 登录「加载失败」根治**：根因为浏览器残留旧 `index.html` 缓存指向已删除的 chunk（每次部署 JS hash 变化）。修复：`web/src/router/guards.js` 的 `router.onError` 在 chunk 加载失败时自动 `location.reload()` 一次；并已在 nginx 对 `index.html` 加 `Cache-Control: no-cache`。
