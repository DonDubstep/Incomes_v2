export const MONTHS_SHORT = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'];
export const MONTHS_FULL = ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'];

// Месяц как целое число (год * 12 + месяц) - удобно для арифметики
export const monthIndex = (key) => {
  const [y, m] = key.split('-').map(Number);
  return y * 12 + m - 1;
};
export const keyFromIndex = (i) => `${Math.floor(i / 12)}-${String((i % 12) + 1).padStart(2, '0')}`;
export const todayIndex = () => {
  const d = new Date();
  return d.getFullYear() * 12 + d.getMonth();
};
export const periodLabel = (idx) => `${MONTHS_FULL[idx % 12]} ${Math.floor(idx / 12)}`;

// "Август 2026" или "2026-08" -> "2026-08"
export function parsePeriod(str) {
  if (typeof str !== 'string') return null;
  if (/^\d{4}-\d{2}/.test(str)) return str.slice(0, 7);
  const [name, year] = str.trim().split(/\s+/);
  const m = MONTHS_FULL.findIndex((x) => x.toLowerCase() === name?.toLowerCase());
  return m < 0 || !year ? null : `${year}-${String(m + 1).padStart(2, '0')}`;
}

// rows: [{period, amount}] -> Map('2026-08' -> 210000)
export function toMap(rows) {
  const map = new Map();
  (rows ?? []).forEach((r) => {
    const key = parsePeriod(r.period);
    if (key) map.set(key, Number(r.amount));
  });
  return map;
}

// 12 месяцев, заканчивающихся на (latestIdx - offset * 12).
// fillZero: месяцы внутри диапазона данных без записей считаются нулём (доходы).
export function buildWindow(map, latestIdx, earliestIdx, offset, fillZero) {
  const endIdx = latestIdx - offset * 12;
  return Array.from({ length: 12 }, (_, i) => {
    const idx = endIdx - 11 + i;
    const mi = idx % 12;
    const year = Math.floor(idx / 12);
    const inRange = idx >= earliestIdx && idx <= latestIdx;
    return {
      key: keyFromIndex(idx),
      label: MONTHS_SHORT[mi],
      tipLabel: `${MONTHS_SHORT[mi]} ${year}`,
      fullLabel: `${MONTHS_FULL[mi]} ${year}`,
      year: i === 0 || mi === 0 ? year : null, // год под подписью - у первого столбца и у января
      value: map.get(keyFromIndex(idx)) ?? (fillZero && inRange ? 0 : null),
    };
  });
}

// Максимум, минимум и среднее по последним n месяцам (для графика общей суммы)
export function lastNStats(map, latestIdx, n) {
  const pts = [];
  for (let i = 0; i < n; i++) {
    const v = map.get(keyFromIndex(latestIdx - i));
    if (v != null) pts.push({ idx: latestIdx - i, value: v });
  }
  if (!pts.length) return null;
  return {
    max: pts.reduce((a, b) => (b.value > a.value ? b : a)),
    min: pts.reduce((a, b) => (b.value < a.value ? b : a)),
    mean: pts.reduce((s, p) => s + p.value, 0) / pts.length,
  };
}

export const percentChange = (cur, prev) =>
  prev > 0 && cur != null ? ((cur - prev) / prev) * 100 : null;

export const formatPercent = (p) =>
  `${p > 0 ? '+' : ''}${p.toLocaleString('ru-RU', { maximumFractionDigits: 1 })}%`;

export const formatMoney = (value, unit = '₽') =>
  `${value.toLocaleString('ru-RU', { maximumFractionDigits: 2 })} ${unit}`;

// 200000 -> "200K", 2500000 -> "2,5M"
export function formatAxis(v) {
  const abs = Math.abs(v);
  const fmt = (n) => n.toLocaleString('ru-RU', { maximumFractionDigits: 1 });
  if (abs >= 1e9) return `${fmt(v / 1e9)}B`;
  if (abs >= 1e6) return `${fmt(v / 1e6)}M`;
  if (abs >= 1e3) return `${fmt(v / 1e3)}K`;
  return fmt(v);
}

// Подбирает "красивую" верхнюю границу и деления оси (до 4 интервалов)
export function niceScale(max) {
  if (!(max > 0)) return { top: 1, ticks: [0] };
  const raw = max / 4;
  const mag = 10 ** Math.floor(Math.log10(raw));
  const step = [1, 2, 2.5, 5, 10].find((n) => raw / mag <= n) * mag;
  const count = Math.ceil(max / step - 1e-9);
  return {
    top: step * count,
    ticks: Array.from({ length: count + 1 }, (_, i) => Number((i * step).toPrecision(12))),
  };
}

// "USD" -> "$", "AED" -> "AED"; при неизвестном коде возвращаем как есть
export function currencySymbol(abbr) {
  try {
    const parts = new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: abbr,
      currencyDisplay: 'narrowSymbol',
    }).formatToParts(0);
    return parts.find((p) => p.type === 'currency')?.value ?? abbr;
  } catch {
    return abbr;
  }
}
