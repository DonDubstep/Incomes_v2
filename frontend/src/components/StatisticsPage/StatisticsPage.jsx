import { useEffect, useMemo, useState } from 'react';
import Header from '../Header/Header';
import Select from '../Select/Select';
import BarChart from '../BarChart/BarChart';
import {
  MONTHS_SHORT, monthIndex, keyFromIndex, todayIndex, periodLabel, toMap, buildWindow,
  lastNStats, percentChange, formatPercent, formatMoney, currencySymbol,
} from '../../utils/statistics';
import styles from './StatisticsPage.module.css';

// Имя GET-параметра "N" из описания API (если на бэкенде другое - поменять здесь)
const PERIOD_PARAM = 'n';
const YEAR = 12;

const METRIC_OPTIONS = [
  { value: 'income', label: 'Доходы' },
  { value: 'total', label: 'Общая сумма' },
];
const RANGE_OPTIONS = [
  { value: 3, label: '3 месяца' },
  { value: 6, label: 'Полгода' },
  { value: 12, label: 'Год' },
];

const fetchJson = (url) =>
  fetch(url).then((res) => {
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return res.json();
  });

// Достаёт число из ответа вида {amount: 1}, [{amount: 1}] или просто 1
const pick = (r, field) => {
  const v = Array.isArray(r) ? r[0] : r;
  const x = v !== null && typeof v === 'object' ? v[field] : v;
  return x == null ? null : Number(x);
};

const monthLabel = (idx) => `${MONTHS_SHORT[idx % 12]} ${Math.floor(idx / 12)}`;

function StatisticsPage() {
  const [metric, setMetric] = useState('income');
  const [currency, setCurrency] = useState('all'); // 'all' | id валюты
  const [currencies, setCurrencies] = useState([]);
  const [range, setRange] = useState(12);
  const [offset, setOffset] = useState(0); // 0 - последние 12 месяцев, 1 - 12 до них, ...
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [yearSum, setYearSum] = useState(null);
  const [totalData, setTotalData] = useState(null);
  const [rangeStats, setRangeStats] = useState(null);

  const isIncome = metric === 'income';

  useEffect(() => {
    fetchJson('/api/14')
      .then((list) => setCurrencies(list ?? []))
      .catch((err) => console.error('Ошибка при загрузке валют', err));
  }, []);

  // Все месяцы выбранного показателя
  useEffect(() => {
    let cancelled = false;
    const url = isIncome ? '/api/11' : `/api/13?currency=${currency === 'all' ? 0 : currency}`;
    setLoading(true);
    setError(null);
    setRows([]);
    setOffset(0);
    fetchJson(url)
      .then((r) => !cancelled && setRows(r ?? []))
      .catch((err) => {
        console.error('Ошибка при загрузке', err);
        if (!cancelled) setError('Не удалось загрузить данные');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => { cancelled = true; };
  }, [isIncome, currency]);

  // Доходы за год (общая карточка)
  useEffect(() => {
    fetchJson('/api/10')
      .then(setYearSum)
      .catch((err) => console.error('Ошибка при загрузке суммы за год', err));
  }, []);

  // Общая сумма на текущий момент
  useEffect(() => {
    fetchJson('/api/5')
      .then(setTotalData)
      .catch((err) => console.error('Ошибка при загрузке суммы за год', err));
  }, []);


  // Максимум / минимум / среднее доходов за выбранный период
  useEffect(() => {
    if (!isIncome) return undefined;
    let cancelled = false;
    fetchJson(`/api/12?${PERIOD_PARAM}=${range}`)
      .then((r) => !cancelled && setRangeStats({
        max: pick(r, 'max_income'), min: pick(r, 'min_income'), mean: pick(r, 'mean_income'),
      }))
      .catch((err) => { console.error(err); if (!cancelled) setRangeStats(null); });
    return () => { cancelled = true; };
  }, [isIncome, range]);

  const model = useMemo(() => {
    const map = toMap(rows);
    const idxs = [...map.keys()].map(monthIndex);
    const latestIdx = idxs.length ? Math.max(...idxs) : todayIndex();
    const earliestIdx = idxs.length ? Math.min(...idxs) : latestIdx;
    const startIdx = latestIdx - offset * 12 - 11;
    return {
      map,
      latestIdx,
      series: buildWindow(map, latestIdx, earliestIdx, offset, isIncome),
      canBack: startIdx > earliestIdx,
      rangeLabel: `${monthLabel(startIdx)} – ${monthLabel(startIdx + 11)}`,
    };
  }, [rows, offset, isIncome]);

  const { map, latestIdx, series } = model;

  const selectedCurrency = currencies.find((c) => String(c.id) === currency);
  const unit = !isIncome && selectedCurrency
    ? currencySymbol(selectedCurrency.abbr || selectedCurrency.name)
    : '₽';
  const money = (v) => (v == null ? '—' : formatMoney(v, unit));

  // Главное число и сравнение с годом ранее
  let mainValue;
  let prevValue;
  if (isIncome) {
    mainValue = yearSum;
    let sum = 0;
    let found = false;
    for (let i = 12; i < 24; i++) {
      const v = map.get(keyFromIndex(latestIdx - i));
      if (v != null) { sum += v; found = true; }
    }
    prevValue = found ? sum : null;
  } else {
    mainValue = totalData;    
    let sum = 0;
    let found = false;
    for (let i = 12; i < 24; i++) {
      const v = map.get(keyFromIndex(latestIdx - i));
      if (v != null) { sum += v; found = true; }
    }
    prevValue = found ? sum : null;
  }
  const delta = percentChange(mainValue, prevValue);

  // Три показателя: для доходов - с бэкенда (/12), для общей суммы - считаем из графика
  const totalStats = isIncome ? null : lastNStats(map, latestIdx, range);
  const cards = isIncome
    ? [
        { title: 'Максимальный доход', value: rangeStats?.max },
        { title: 'Минимальный доход', value: rangeStats?.min },
        { title: 'Средний доход в месяц', value: rangeStats?.mean },
      ]
    : [
        { title: 'Максимальная сумма', value: totalStats?.max.value, sub: totalStats && periodLabel(totalStats.max.idx) },
        { title: 'Минимальная сумма', value: totalStats?.min.value, sub: totalStats && periodLabel(totalStats.min.idx) },
        { title: 'Средняя сумма за месяц', value: totalStats?.mean },
      ];

  return (
    <>
      <Header title="Статистика">
        <Select value={metric} onChange={setMetric} options={METRIC_OPTIONS} ariaLabel="Показатель" />
      </Header>

      <section className={styles.card}>
        <div className={styles.cardHeader}>
          <div>
            <div className={styles.label}>
              {isIncome ? 'Доходы за последний год' : 'Общая сумма на текущий момент'}
            </div>
            <div className={styles.bigValue}>{money(mainValue)}</div>
            {delta !== null && (
              <div className={styles.delta}>
                <span className={delta >= 0 ? styles.up : styles.down}>
                  {delta >= 0 ? '↑' : '↓'} {formatPercent(delta)}
                </span>{' '}
                по сравнению с предыдущим годом
              </div>
            )}
          </div>

          {!isIncome && (
            <Select
              value={currency}
              onChange={setCurrency}
              ariaLabel="Валюта"
              options={[
                { value: 'all', label: 'Все валюты (в ₽)' },
                ...currencies.map((c) => ({ value: String(c.id), label: c.name })),
              ]}
            />
          )}
        </div>

        {error && <div className={styles.error}>{error}</div>}

        <div className={styles.nav}>
          <button
            type="button"
            className={styles.navBtn}
            onClick={() => setOffset((o) => o + 1)}
            disabled={!model.canBack || loading}
            aria-label="Назад на 12 месяцев"
          >‹</button>
          <span className={styles.navLabel}>{model.rangeLabel}</span>
          <button
            type="button"
            className={styles.navBtn}
            onClick={() => setOffset((o) => o - 1)}
            disabled={offset === 0 || loading}
            aria-label="Вперёд на 12 месяцев"
          >›</button>
        </div>

        <BarChart series={series} formatValue={money} loading={loading} emptyText="Нет данных за период" />
      </section>

      <div className={styles.rangeBar}>
        <span className={styles.label}>Показатели за последние</span>
        <div className={styles.segmented} role="group" aria-label="Период показателей">
          {RANGE_OPTIONS.map((o) => (
            <button
              key={o.value}
              type="button"
              className={`${styles.segBtn} ${range === o.value ? styles.segActive : ''}`}
              onClick={() => setRange(o.value)}
              aria-pressed={range === o.value}
            >{o.label}</button>
          ))}
        </div>
      </div>

      <div className={styles.summary}>
        {cards.map((c) => (
          <div key={c.title} className={styles.summaryCard}>
            <div className={styles.label}>{c.title}</div>
            <div className={styles.summaryTitle}>{money(c.value)}</div>
            {c.sub && <div className={styles.summaryValue}>{c.sub}</div>}
          </div>
        ))}
      </div>
    </>
  );
}

export default StatisticsPage;
