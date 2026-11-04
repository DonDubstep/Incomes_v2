import { useMemo, useState } from 'react';
import { formatAxis, niceScale } from '../../utils/statistics';
import styles from './BarChart.module.css';

// series: [{ key, label, value|null }]
function BarChart({ series, formatValue, loading = false, emptyText = 'Нет данных' }) {
  const [hovered, setHovered] = useState(null);

  const max = Math.max(0, ...series.map((p) => p.value ?? 0));
  const { top, ticks } = useMemo(() => niceScale(max), [max]);
  const isEmpty = max <= 0;

  // По умолчанию подсвечиваем месяц с максимальным значением
  const bestIndex = useMemo(() => {
    let idx = -1;
    series.forEach((p, i) => {
      if (p.value !== null && (idx === -1 || p.value > series[idx].value)) idx = i;
    });
    return idx;
  }, [series]);
  const activeIndex = hovered ?? bestIndex;

  const pct = (v) => (v / top) * 100;

  return (
    <div className={`${styles.chart} ${loading ? styles.loading : ''}`}>
      <div className={styles.yAxis}>
        {ticks.map((t) => (
          <span key={t} className={styles.yLabel} style={{ bottom: `${pct(t)}%` }}>{formatAxis(t)}</span>
        ))}
      </div>

      <div className={styles.plot}>
        {ticks.map((t) => (
          <div key={t} className={styles.gridLine} style={{ bottom: `${pct(t)}%` }} />
        ))}

        <div className={styles.bars} onMouseLeave={() => setHovered(null)}>
          {series.map((p, i) => {
            const active = i === activeIndex;
            const h = p.value === null ? 0 : pct(Math.max(0, p.value));
            const tipPos = i === 0 ? styles.tipFirst : i === series.length - 1 ? styles.tipLast : '';
            return (
              <button
                key={p.key}
                type="button"
                className={styles.col}
                onMouseEnter={() => setHovered(i)}
                onFocus={() => setHovered(i)}
                onBlur={() => setHovered(null)}
                aria-label={`${p.fullLabel}: ${p.value === null ? 'нет данных' : formatValue(p.value)}`}
              >
                {p.value > 0 && (
                  <span className={`${styles.bar} ${active ? styles.barActive : ''}`} style={{ height: `max(2px, ${h}%)` }} />
                )}
                {active && p.value !== null && !isEmpty && !loading && (
                  <span className={`${styles.tooltip} ${tipPos}`} style={{ bottom: `calc(${h}% + 12px)` }}>
                    <span className={styles.tipLabel}>{p.tipLabel ?? p.label}</span>
                    <span className={styles.tipValue}>{formatValue(p.value)}</span>
                  </span>
                )}
              </button>
            );
          })}
        </div>

        {(loading || isEmpty) && (
          <div className={styles.overlay}>{loading ? 'Загрузка...' : emptyText}</div>
        )}
      </div>

      <div className={styles.xAxis}>
        {series.map((p, i) => (
          <span key={p.key} className={i === activeIndex ? styles.xActive : ''}>
            {p.label}
            {p.year && <small className={styles.xYear}>{p.year}</small>}
          </span>
        ))}
      </div>
    </div>
  );
}

export default BarChart;
