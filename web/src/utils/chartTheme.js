/*
 * OneHost 图表主题（ECharts）
 * ------------------------------------------------------------
 * 从 theme.scss 的 CSS 变量读取品牌色 / 文本色 / 边框色，生成亮色与暗色两套 ECharts 默认配置，
 * 并提供 useChartTheme() 在 pinia theme store 切换（html.dark）时自动重新渲染图表。
 *
 * 用法：
 *   import { withChartTheme, useChartTheme, areaGradient, getChartPalette } from '@/utils/chartTheme'
 *   chart.setOption(withChartTheme({ xAxis: {...}, yAxis: {...}, series: [...] }), true)
 *   useChartTheme(() => renderChart(lastData))   // 主题切换时重绘
 */
import { watch, nextTick, onBeforeUnmount, getCurrentInstance } from 'vue'
import { useThemeStore } from '@/pinia/modules/theme'

// CSS 变量读取失败时（SSR / 单元测试 / 变量缺失）使用的回退值，与 theme.scss 保持一致
const FALLBACK = {
  light: {
    '--primary-color': '#4f46e5',
    '--brand-accent': '#0ea5e9',
    '--success-color': '#10b981',
    '--warning-color': '#f59e0b',
    '--brand-accent-2': '#8b5cf6',
    '--error-color': '#ef4444',
    '--text-color-primary': '#0f172a',
    '--text-color-secondary': '#475569',
    '--text-color-tertiary': '#94a3b8',
    '--border-color': '#e5e7ef',
    '--border-color-light': '#eef0f5',
    '--surface-elevated': '#ffffff',
    '--card-bg-solid': '#ffffff'
  },
  dark: {
    '--primary-color': '#6366f1',
    '--brand-accent': '#38bdf8',
    '--success-color': '#34d399',
    '--warning-color': '#fbbf24',
    '--brand-accent-2': '#a78bfa',
    '--error-color': '#f87171',
    '--text-color-primary': '#e8ebf4',
    '--text-color-secondary': '#a3acc2',
    '--text-color-tertiary': '#737d96',
    '--border-color': 'rgba(148, 163, 184, 0.14)',
    '--border-color-light': 'rgba(148, 163, 184, 0.09)',
    '--surface-elevated': '#161a27',
    '--card-bg-solid': '#11141f'
  }
}

const FONT_FAMILY =
  "Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif"

/** 当前文档是否处于暗色模式 */
export const isDarkMode = () =>
  typeof document !== 'undefined' && document.documentElement.classList.contains('dark')

/** 读取 :root / html.dark 上的 CSS 变量 */
export const readCssVar = (name, dark = isDarkMode()) => {
  const fallback = (dark ? FALLBACK.dark : FALLBACK.light)[name] || ''
  if (typeof window === 'undefined' || typeof document === 'undefined' || !window.getComputedStyle) {
    return fallback
  }
  const value = window.getComputedStyle(document.documentElement).getPropertyValue(name)
  return (value && value.trim()) || fallback
}

/** 将 #rrggbb / #rgb 转为带透明度的 rgba；其他格式（rgb/rgba）尽量原样处理 */
export const withAlpha = (color, alpha) => {
  if (!color) return `rgba(99, 102, 241, ${alpha})`
  const c = color.trim()
  if (c.startsWith('#')) {
    let hex = c.slice(1)
    if (hex.length === 3) hex = hex.split('').map((ch) => ch + ch).join('')
    const num = parseInt(hex.slice(0, 6), 16)
    if (Number.isNaN(num)) return c
    return `rgba(${(num >> 16) & 255}, ${(num >> 8) & 255}, ${num & 255}, ${alpha})`
  }
  const m = c.match(/rgba?\(([^)]+)\)/i)
  if (m) {
    const [r, g, b] = m[1].split(',').map((s) => s.trim())
    return `rgba(${r}, ${g}, ${b}, ${alpha})`
  }
  return c
}

/** 系列色板：品牌靛蓝 → 天青 → 翠绿 → 琥珀 → 紫罗兰 → 玫红 */
export const getChartPalette = (dark = isDarkMode()) => [
  readCssVar('--primary-color', dark),
  readCssVar('--brand-accent', dark),
  readCssVar('--success-color', dark),
  readCssVar('--warning-color', dark),
  readCssVar('--brand-accent-2', dark),
  readCssVar('--error-color', dark)
]

/** 主题令牌（供组件需要单独取色时使用） */
export const getChartTokens = (dark = isDarkMode()) => ({
  dark,
  palette: getChartPalette(dark),
  primary: readCssVar('--primary-color', dark),
  accent: readCssVar('--brand-accent', dark),
  success: readCssVar('--success-color', dark),
  warning: readCssVar('--warning-color', dark),
  danger: readCssVar('--error-color', dark),
  violet: readCssVar('--brand-accent-2', dark),
  text: readCssVar('--text-color-primary', dark),
  textSecondary: readCssVar('--text-color-secondary', dark),
  textMuted: readCssVar('--text-color-tertiary', dark),
  border: readCssVar('--border-color', dark),
  borderLight: readCssVar('--border-color-light', dark),
  surface: readCssVar('--surface-elevated', dark),
  card: readCssVar('--card-bg-solid', dark)
})

/** 垂直线性渐变（面积图填充），返回 ECharts 可直接识别的纯对象，无需引入 echarts */
export const areaGradient = (color, from = 0.28, to = 0.02) => ({
  type: 'linear',
  x: 0,
  y: 0,
  x2: 0,
  y2: 1,
  colorStops: [
    { offset: 0, color: withAlpha(color, from) },
    { offset: 1, color: withAlpha(color, to) }
  ]
})

const axisStyle = (tk, kind) => {
  const base = {
    axisLine: { show: kind === 'category', lineStyle: { color: tk.border } },
    axisTick: { show: false },
    axisLabel: { color: tk.textMuted, fontSize: 11, fontFamily: FONT_FAMILY },
    nameTextStyle: { color: tk.textMuted, fontSize: 11, fontFamily: FONT_FAMILY },
    splitLine: {
      show: kind !== 'category',
      lineStyle: { color: tk.borderLight, type: 'dashed' }
    }
  }
  return base
}

/** 亮/暗两套 ECharts 默认配置 */
export const buildChartTheme = (dark = isDarkMode()) => {
  const tk = getChartTokens(dark)
  return {
    tokens: tk,
    base: {
      color: tk.palette,
      backgroundColor: 'transparent',
      textStyle: { color: tk.textSecondary, fontFamily: FONT_FAMILY },
      animationDuration: 600,
      animationEasing: 'cubicOut',
      grid: { top: 40, right: 16, bottom: 12, left: 12, containLabel: true },
      legend: {
        icon: 'roundRect',
        itemWidth: 10,
        itemHeight: 10,
        itemGap: 16,
        textStyle: { color: tk.textSecondary, fontSize: 12, fontFamily: FONT_FAMILY },
        inactiveColor: dark ? '#3b4256' : '#cbd5e1'
      },
      tooltip: {
        backgroundColor: tk.surface,
        borderColor: tk.border,
        borderWidth: 1,
        padding: [8, 12],
        textStyle: { color: tk.text, fontSize: 12, fontFamily: FONT_FAMILY },
        extraCssText: `border-radius: 10px; box-shadow: ${
          dark ? '0 12px 28px -8px rgba(0,0,0,0.6)' : '0 12px 28px -8px rgba(15,23,42,0.18)'
        };`,
        axisPointer: {
          lineStyle: { color: withAlpha(tk.primary, 0.45), type: 'dashed' },
          crossStyle: { color: withAlpha(tk.primary, 0.45) },
          shadowStyle: { color: withAlpha(tk.primary, 0.06) },
          label: { backgroundColor: tk.primary, color: '#fff', fontSize: 11 }
        }
      }
    },
    categoryAxis: axisStyle(tk, 'category'),
    valueAxis: axisStyle(tk, 'value')
  }
}

const isPlainObject = (v) => v !== null && typeof v === 'object' && !Array.isArray(v) && !(v instanceof Date)

/** 深合并：override 优先，数组整体替换，函数原样保留 */
export const deepMerge = (base, override) => {
  if (!isPlainObject(base)) return override === undefined ? base : override
  if (!isPlainObject(override)) return override === undefined ? { ...base } : override
  const out = { ...base }
  Object.keys(override).forEach((key) => {
    const ov = override[key]
    if (ov === undefined) return
    out[key] = isPlainObject(ov) && isPlainObject(base[key]) ? deepMerge(base[key], ov) : ov
  })
  return out
}

const themeAxis = (axis, theme) => {
  if (!axis) return axis
  const apply = (a) => deepMerge(a.type === 'value' || a.type === 'log' ? theme.valueAxis : theme.categoryAxis, a)
  return Array.isArray(axis) ? axis.map(apply) : apply(axis)
}

/**
 * 把组件自己的 option 与主题默认值合并。
 * 组件中显式设置的属性优先；坐标轴按 category/value 分别套用样式。
 */
export const withChartTheme = (option = {}, dark = isDarkMode()) => {
  const theme = buildChartTheme(dark)
  const merged = deepMerge(theme.base, option)
  if (option.xAxis) merged.xAxis = themeAxis(option.xAxis, theme)
  if (option.yAxis) merged.yAxis = themeAxis(option.yAxis, theme)
  if (option.tooltip === undefined) delete merged.tooltip
  if (option.legend === undefined) delete merged.legend
  if (option.grid === undefined && !option.xAxis && !option.yAxis) delete merged.grid
  return merged
}

/**
 * 主题切换时调用 rerender（例如重新 setOption）。
 * 监听 pinia theme store 的 isDark；组件卸载时自动停止。
 * 返回 stop 函数，可手动取消。
 */
export const useChartTheme = (rerender) => {
  let stop = () => {}
  try {
    const themeStore = useThemeStore()
    stop = watch(
      () => themeStore.isDark,
      async () => {
        // 等待 html.dark 类及 CSS 变量生效后再读取颜色
        await nextTick()
        if (typeof window !== 'undefined' && window.requestAnimationFrame) {
          window.requestAnimationFrame(() => rerender())
        } else {
          rerender()
        }
      },
      { flush: 'post' }
    )
  } catch (e) {
    // pinia 未就绪（极少见）时静默降级：图表仍可正常渲染，只是不随主题切换
    console.warn('[chartTheme] theme store unavailable:', e)
  }
  if (getCurrentInstance()) onBeforeUnmount(() => stop())
  return stop
}
